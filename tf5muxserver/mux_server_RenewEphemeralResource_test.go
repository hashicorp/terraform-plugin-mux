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

func TestMuxServerRenewEphemeralResource(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	testServer1 := &tf5testserver.TestServer{
		GetProviderSchemaResponse: &tfprotov5.GetProviderSchemaResponse{
			EphemeralResourceSchemas: map[string]*tfprotov5.Schema{
				"test_ephemeral_resource_server1": {},
			},
		},
	}
	testServer2 := &tf5testserver.TestServer{
		GetProviderSchemaResponse: &tfprotov5.GetProviderSchemaResponse{
			EphemeralResourceSchemas: map[string]*tfprotov5.Schema{
				"test_ephemeral_resource_server2": {},
			},
		},
	}
	servers := []func() tfprotov5.ProviderServer{testServer1.ProviderServer, testServer2.ProviderServer}
	muxServer, err := tf5muxserver.NewMuxServer(ctx, servers...)

	if err != nil {
		t.Fatalf("unexpected error setting up factory: %s", err)
	}

	_, err = muxServer.ProviderServer().RenewEphemeralResource(ctx, &tfprotov5.RenewEphemeralResourceRequest{
		TypeName: "test_ephemeral_resource_server1",
	})

	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	if !testServer1.RenewEphemeralResourceCalled["test_ephemeral_resource_server1"] {
		t.Errorf("expected test_ephemeral_resource_server1 RenewEphemeralResource to be called on server1")
	}

	if testServer2.RenewEphemeralResourceCalled["test_ephemeral_resource_server1"] {
		t.Errorf("unexpected test_ephemeral_resource_server1 RenewEphemeralResource called on server2")
	}

	_, err = muxServer.ProviderServer().RenewEphemeralResource(ctx, &tfprotov5.RenewEphemeralResourceRequest{
		TypeName: "test_ephemeral_resource_server2",
	})

	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	if testServer1.RenewEphemeralResourceCalled["test_ephemeral_resource_server2"] {
		t.Errorf("unexpected test_ephemeral_resource_server2 RenewEphemeralResource called on server1")
	}

	if !testServer2.RenewEphemeralResourceCalled["test_ephemeral_resource_server2"] {
		t.Errorf("expected test_ephemeral_resource_server2 RenewEphemeralResource to be called on server2")
	}
}
