// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package tf6muxserver_test

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-mux/internal/tf6testserver"
	"github.com/hashicorp/terraform-plugin-mux/tf6muxserver"
)

func TestMuxServerConfigureProvider(t *testing.T) {
	t.Parallel()

	testServers := [5]*tf6testserver.TestServer{{}, {}, {}, {}, {}}

	servers := []func() tfprotov6.ProviderServer{
		testServers[0].ProviderServer,
		testServers[1].ProviderServer,
		testServers[2].ProviderServer,
		testServers[3].ProviderServer,
		testServers[4].ProviderServer,
	}

	muxServer, err := tf6muxserver.NewMuxServer(context.Background(), servers...)

	if err != nil {
		t.Fatalf("error setting up muxer: %s", err)
	}

	_, err = muxServer.ProviderServer().ConfigureProvider(context.Background(), &tfprotov6.ConfigureProviderRequest{})

	if err != nil {
		t.Fatalf("error calling ConfigureProvider: %s", err)
	}

	for num, testServer := range testServers {
		if !testServer.ConfigureProviderCalled {
			t.Errorf("configure not called on server%d", num+1)
		}
	}
}
