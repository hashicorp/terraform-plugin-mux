// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package tf5muxserver_test

import (
	"context"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/terraform-plugin-go/tfprotov5"

	"github.com/hashicorp/terraform-plugin-mux/internal/tf5testserver"
	"github.com/hashicorp/terraform-plugin-mux/tf5muxserver"
)

func TestMuxServerStopProvider(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	testServers := [5]*tf5testserver.TestServer{
		{
			GetProviderSchemaResponse: &tfprotov5.GetProviderSchemaResponse{},
		},
		{
			GetProviderSchemaResponse: &tfprotov5.GetProviderSchemaResponse{},
			StopProviderResponse: &tfprotov5.StopProviderResponse{
				Error: "error in server2",
			},
		},
		{
			GetProviderSchemaResponse: &tfprotov5.GetProviderSchemaResponse{},
		},
		{
			GetProviderSchemaResponse: &tfprotov5.GetProviderSchemaResponse{},
			StopProviderResponse: &tfprotov5.StopProviderResponse{
				Error: "error in server4",
			},
		},
		{
			GetProviderSchemaResponse: &tfprotov5.GetProviderSchemaResponse{},
		},
	}

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

	resp, err := muxServer.ProviderServer().StopProvider(ctx, &tfprotov5.StopProviderRequest{})

	if err != nil {
		t.Fatalf("error calling StopProvider: %s", err)
	}

	expectedRespError := strings.Join([]string{"error in server2", "error in server4"}, "\n")

	if diff := cmp.Diff(resp.Error, expectedRespError); diff != "" {
		t.Errorf("unexpected response Error difference: %s", diff)
	}

	for num, testServer := range testServers {
		if !testServer.StopProviderCalled {
			t.Errorf("StopProvider not called on server%d", num+1)
		}
	}
}
