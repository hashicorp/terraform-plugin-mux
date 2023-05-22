// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package tf6muxserver_test

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6/tf6server"
	"github.com/hashicorp/terraform-plugin-mux/tf6muxserver"
)

func ExampleNewMuxServer() {
	ctx := context.Background()
	providers := []func() tfprotov6.ProviderServer{
		// Example terraform-plugin-framework ProviderServer function
		// func() tfprotov6.ProviderServer {
		//   return tfsdk.NewProtocol6Server(frameworkprovider.New("version")())
		// },
		//
		// Example terraform-plugin-go ProviderServer function
		// goprovider.Provider(),
	}

	// requests will be routed to whichever server advertises support for
	// them in the GetSchema response. Only one server may advertise
	// support for any given resource, data source, or the provider or
	// provider_meta schemas. An error will be returned if more than one
	// server claims support.
	muxServer, err := tf6muxserver.NewMuxServer(ctx, providers...)

	if err != nil {
		log.Fatalln(err.Error())
	}

	// Use the result to start a muxed provider
	err = tf6server.Serve("registry.terraform.io/namespace/example", muxServer.ProviderServer)

	if err != nil {
		log.Fatalln(err.Error())
	}
}
