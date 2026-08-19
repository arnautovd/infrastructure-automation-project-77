package main

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
)

var version = "dev"

func main() {
	err := providerserver.Serve(context.Background(), NewProvider, providerserver.ServeOpts{
		Address: "registry.terraform.io/arnautovd/vscale",
		Debug:   false,
	})
	if err != nil {
		log.Fatal(err)
	}
}
