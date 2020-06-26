package tfmux

import (
	"context"
	"log"
	"os"

	"github.com/hashicorp/terraform-plugin-mux/internal/tfplugin5"
)

func ExampleNewSchemaServerFactory_v2protocol() {
	ctx := context.Background()

	// the ProviderServer from SDKv2
	// usually this is the Provider function
	var sdkv2 func() tfplugin5.ProviderServer

	// the ProviderServer from the new protocol package
	var protocolServer func() tfplugin5.ProviderServer

	// requests will be routed to whichever server advertises support for
	// them in the GetSchema response. Only one server may advertise
	// support for any given resource, data source, or the provider or
	// provider_meta schemas. An error will be returned if more than one
	// server claims support.
	factory, err := NewSchemaServerFactory(ctx, sdkv2, protocolServer)
	if err != nil {
		log.Println(err.Error())
		os.Exit(1)
	}

	// the server exported by the protocol package doesn't need to support
	// PrepareProviderConfig, as that RPC call is only useful for the
	// legacy SDK. But the mux package doesn't and can't know if we're
	// using a legacy SDK or not, or if we're using multiple legacy SDKs
	// which should handle that call, so we just need to tell it explicitly
	// which server we want handling that RPC call.
	//
	// We use `0` to denote the first serverin the list. 1 would be
	// second, and so on. We use the position because our servers are
	// actually server factories, meaning the value with the
	// PrepareProviderConfig method on it is returned at runtime, and isn't
	// available to us now. Here we're saying to use the
	// PrepareProviderConfig method of the server returned from the first
	// factory we passed in, sdkv2. This preserves the lifecycle providers
	// have come to expect from PrepareProviderConfig.
	//
	// Not setting anything will yield a muxer that treats this RPC call as
	// a no-op.
	factory.PrepareProviderConfigServer = 0

	// use the result when instantiating the terraform-plugin-sdk.plugin.Serve
	/*
		plugin.Serve(&plugin.ServeOpts{
			GRPCProviderFunc: plugin.GRPCProviderFunc(factory),
		})
	*/
}
