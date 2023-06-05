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

func TestMuxServerStopProvider(t *testing.T) {
	t.Parallel()

	testServers := [5]*tf5testserver.TestServer{{}, {StopProviderError: "error in server2"}, {}, {StopProviderError: "error in server4"}, {}}

	servers := []func() tfprotov5.ProviderServer{
		testServers[0].ProviderServer,
		testServers[1].ProviderServer,
		testServers[2].ProviderServer,
		testServers[3].ProviderServer,
		testServers[4].ProviderServer,
	}

	muxServer, err := tf5muxserver.NewMuxServer(context.Background(), servers...)

	if err != nil {
		t.Fatalf("error setting up muxer: %s", err)
	}

	_, err = muxServer.ProviderServer().StopProvider(context.Background(), &tfprotov5.StopProviderRequest{})

	if err != nil {
		t.Fatalf("error calling StopProvider: %s", err)
	}

	for num, testServer := range testServers {
		if !testServer.StopProviderCalled {
			t.Errorf("StopProvider not called on server%d", num+1)
		}
	}
}
