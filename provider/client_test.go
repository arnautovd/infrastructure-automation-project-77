package main

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestCreateServerUsesVscaleAPIContract(t *testing.T) {
	var seenToken string
	var seenRequest createServerRequest
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		seenToken = req.Header.Get("X-Token")
		if req.Method == http.MethodPost {
			if err := json.NewDecoder(req.Body).Decode(&seenRequest); err != nil {
				t.Fatalf("decode request: %v", err)
			}
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"ctid":42}`))
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"ctid":42,"name":"web-01","made_from":"debian_12_64","rplan":"small","location":"msk0","status":"active","public_address":{"address":"203.0.113.10"}}`))
	}))
	defer server.Close()

	client, err := NewAPIClient("test-token", server.URL+"/v1/")
	if err != nil {
		t.Fatalf("new client: %v", err)
	}
	client.HTTPClient = server.Client()
	got, err := client.CreateServer(context.Background(), createServerRequest{
		Name: "web-01", MakeFrom: "debian_12_64", Rplan: "small", Location: "msk0", DoStart: true, Keys: []int64{7},
	})
	if err != nil {
		t.Fatalf("create server: %v", err)
	}
	if got.CTID != 42 || got.PublicAddresses.Address != "203.0.113.10" {
		t.Fatalf("unexpected server: %#v", got)
	}
	if seenToken != "test-token" {
		t.Fatalf("unexpected token header: %q", seenToken)
	}
	if seenRequest.Name != "web-01" || len(seenRequest.Keys) != 1 || seenRequest.Keys[0] != 7 {
		t.Fatalf("unexpected create request: %#v", seenRequest)
	}
}

func TestStringSetToIDs(t *testing.T) {
	value := types.SetValueMust(types.StringType, []attr.Value{
		types.StringValue("10"),
		types.StringValue("20"),
	})
	got, err := stringSetToIDs(context.Background(), value)
	if err != nil {
		t.Fatalf("convert IDs: %v", err)
	}
	if len(got) != 2 || got[0] != 10 || got[1] != 20 {
		t.Fatalf("unexpected IDs: %#v", got)
	}
}
