// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package tf5muxserver_test

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-go/tfprotov5"
	"github.com/hashicorp/terraform-plugin-mux/internal/tf5testserver"
	"github.com/hashicorp/terraform-plugin-mux/tf5muxserver"
)

func TestMuxServerConfigureProvider(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	testServers := [5]*tf5testserver.TestServer{{}, {}, {}, {}, {}}

	servers := []func() tfprotov5.ProviderServer{
		testServers[0].ProviderServer,
		testServers[1].ProviderServer,
		testServers[2].ProviderServer,
		testServers[3].ProviderServer,
		testServers[4].ProviderServer,
	}

	muxServer, err := tf5muxserver.NewMuxServer(ctx, servers...)

	if err != nil {
		t.Fatalf("error setting up muxer: %s", err)
	}

	// Required to populate routers
	_, err = muxServer.GetProviderSchema(ctx, &tfprotov5.GetProviderSchemaRequest{})

	if err != nil {
		t.Fatalf("unexpected error calling GetProviderSchema: %s", err)
	}

	_, err = muxServer.ProviderServer().ConfigureProvider(ctx, &tfprotov5.ConfigureProviderRequest{})

	if err != nil {
		t.Fatalf("error calling ConfigureProvider: %s", err)
	}

	for num, testServer := range testServers {
		if !testServer.ConfigureProviderCalled {
			t.Errorf("configure not called on server%d", num+1)
		}
	}
}
