package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type APIClient struct {
	BaseURL    *url.URL
	HTTPClient *http.Client
	Token      string
}

type APIError struct {
	StatusCode int
	Message    string
}

func (e *APIError) Error() string {
	return fmt.Sprintf("Vscale API returned HTTP %d: %s", e.StatusCode, e.Message)
}

type serverResponse struct {
	CTID            int64  `json:"ctid"`
	Name            string `json:"name"`
	MadeFrom        string `json:"made_from"`
	Rplan           string `json:"rplan"`
	Location        string `json:"location"`
	Status          string `json:"status"`
	PublicAddresses struct {
		Address string `json:"address"`
	} `json:"public_address"`
}

type createServerRequest struct {
	MakeFrom string  `json:"make_from"`
	Rplan    string  `json:"rplan"`
	Name     string  `json:"name"`
	Location string  `json:"location"`
	DoStart  bool    `json:"do_start"`
	Keys     []int64 `json:"keys,omitempty"`
}

func NewAPIClient(token, rawURL string) (*APIClient, error) {
	parsed, err := url.Parse(strings.TrimRight(rawURL, "/") + "/")
	if err != nil {
		return nil, err
	}
	if parsed.Scheme != "https" && parsed.Scheme != "http" {
		return nil, fmt.Errorf("API URL must use http or https")
	}
	return &APIClient{BaseURL: parsed, HTTPClient: http.DefaultClient, Token: token}, nil
}

func (c *APIClient) CreateServer(ctx context.Context, request createServerRequest) (serverResponse, error) {
	var server serverResponse
	if err := c.do(ctx, http.MethodPost, "scalets", request, &server); err != nil {
		return server, err
	}
	return c.waitForServer(ctx, server.CTID)
}

func (c *APIClient) GetServer(ctx context.Context, id int64) (serverResponse, error) {
	var server serverResponse
	err := c.do(ctx, http.MethodGet, "scalets/"+strconv.FormatInt(id, 10), nil, &server)
	return server, err
}

func (c *APIClient) DeleteServer(ctx context.Context, id int64) error {
	if err := c.do(ctx, http.MethodDelete, "scalets/"+strconv.FormatInt(id, 10), nil, nil); err != nil {
		var apiErr *APIError
		if errors.As(err, &apiErr) && apiErr.StatusCode == http.StatusNotFound {
			return nil
		}
		return err
	}

	deadline := time.Now().Add(5 * time.Minute)
	for time.Now().Before(deadline) {
		_, err := c.GetServer(ctx, id)
		if err != nil {
			var apiErr *APIError
			if errors.As(err, &apiErr) && apiErr.StatusCode == http.StatusNotFound {
				return nil
			}
		}
		if err := sleepContext(ctx, 2*time.Second); err != nil {
			return err
		}
	}
	return fmt.Errorf("timed out waiting for server %d deletion", id)
}

func (c *APIClient) waitForServer(ctx context.Context, id int64) (serverResponse, error) {
	deadline := time.Now().Add(10 * time.Minute)
	var server serverResponse
	for time.Now().Before(deadline) {
		var err error
		server, err = c.GetServer(ctx, id)
		if err == nil && server.Status != "creating" && server.Status != "installing" {
			return server, nil
		}
		if err != nil {
			return server, err
		}
		if err := sleepContext(ctx, 3*time.Second); err != nil {
			return server, err
		}
	}
	return server, fmt.Errorf("timed out waiting for server %d creation", id)
}

func (c *APIClient) do(ctx context.Context, method, path string, requestBody, responseBody any) error {
	var body io.Reader
	if requestBody != nil {
		encoded, err := json.Marshal(requestBody)
		if err != nil {
			return err
		}
		body = bytes.NewReader(encoded)
	}
	req, err := http.NewRequestWithContext(ctx, method, c.BaseURL.ResolveReference(&url.URL{Path: path}).String(), body)
	if err != nil {
		return err
	}
	req.Header.Set("X-Token", c.Token)
	req.Header.Set("Accept", "application/json")
	if requestBody != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	response, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}
	if res.StatusCode < 200 || res.StatusCode >= 300 {
		message := strings.TrimSpace(string(response))
		if message == "" {
			message = res.Status
		}
		return &APIError{StatusCode: res.StatusCode, Message: message}
	}
	if responseBody != nil && len(response) > 0 {
		if err := json.Unmarshal(response, responseBody); err != nil {
			return fmt.Errorf("decode Vscale response: %w", err)
		}
	}
	return nil
}

func sleepContext(ctx context.Context, duration time.Duration) error {
	timer := time.NewTimer(duration)
	defer timer.Stop()
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-timer.C:
		return nil
	}
}
