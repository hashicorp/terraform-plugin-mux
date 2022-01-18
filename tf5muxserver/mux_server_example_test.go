package tfmux

import (
	"context"
	"log"
	"os"

	"github.com/hashicorp/terraform-plugin-go/tfprotov5"
)

func ExampleNewSchemaServerFactory_v2protocol() {
	ctx := context.Background()

	// the ProviderServer from SDKv2
	// usually this is the Provider function
	var sdkv2 func() tfprotov5.ProviderServer

	// the ProviderServer from the new protocol package
	var protocolServer func() tfprotov5.ProviderServer

	// requests will be routed to whichever server advertises support for
	// them in the GetSchema response. Only one server may advertise
	// support for any given resource, data source, or the provider or
	// provider_meta schemas. An error will be returned if more than one
	// server claims support.
	_, err := NewSchemaServerFactory(ctx, sdkv2, protocolServer)
	if err != nil {
		log.Println(err.Error())
		os.Exit(1)
	}

	// use the result when instantiating the terraform-plugin-sdk.plugin.Serve
	/*
		plugin.Serve(&plugin.ServeOpts{
			GRPCProviderFunc: plugin.GRPCProviderFunc(factory),
		})
	*/
}
