// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package tf5muxserver_test

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-go/tfprotov5"
	"github.com/hashicorp/terraform-plugin-go/tfprotov5/tf5server"
	"github.com/hashicorp/terraform-plugin-mux/tf5muxserver"
)

func ExampleNewMuxServer() {
	ctx := context.Background()
	providers := []func() tfprotov5.ProviderServer{
		// Example terraform-plugin-sdk ProviderServer function
		// sdkprovider.New("version")().GRPCProvider,
		//
		// Example terraform-plugin-go ProviderServer function
		// goprovider.Provider(),
	}

	// requests will be routed to whichever server advertises support for
	// them in the GetSchema response. Only one server may advertise
	// support for any given resource, data source, or the provider or
	// provider_meta schemas. An error will be returned if more than one
	// server claims support.
	muxServer, err := tf5muxserver.NewMuxServer(ctx, providers...)

	if err != nil {
		log.Fatalln(err.Error())
	}

	// Use the result to start a muxed provider
	err = tf5server.Serve("registry.terraform.io/namespace/example", muxServer.ProviderServer)

	if err != nil {
		log.Fatalln(err.Error())
	}
}
