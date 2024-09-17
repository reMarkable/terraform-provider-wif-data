// Copyright (c) HashiCorp, Inc.

package main

import (
	"context"
	"flag"
	"log"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"

	"github.com/remarkable/terraform-wif-data-provider/internal/data"
)

var version string = "dev"

func main() {
	var debug bool

	flag.BoolVar(&debug, "debug", false, "set to true to run the provider with support for debuggers like delve")
	flag.Parse()
	opts := providerserver.ServeOpts{
		Address: "registry.terraform.io/remarkable/terraform-wif-data-provider",
		Debug:   debug,
	}
	err := providerserver.Serve(context.Background(), data.New(version), opts)
	if err != nil {
		log.Fatal(err.Error())
	}
}
