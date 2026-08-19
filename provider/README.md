# Custom Vscale Provider

This provider uses the Terraform Plugin Framework and the Vscale REST API.

The exported resource is `vscale_server`.

## API mapping

- `POST /v1/scalets` creates a server.
- `GET /v1/scalets/{ctid}` reads a server.
- `DELETE /v1/scalets/{ctid}` deletes a server.
- Authentication uses the `X-Token` header.

The provider waits for asynchronous server creation and deletion by polling the API. All server arguments are replacement-only because the Vscale API does not provide a safe generic update operation for them.

## Build

```text
go test ./...
go build -o bin/terraform-provider-vscale_v0.1.6.exe .
```

The root `Makefile` builds this binary into a Terraform filesystem mirror and sets `TF_CLI_CONFIG_FILE` before running Terraform. This avoids downloading the custom provider from the public Registry.
