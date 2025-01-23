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

func TestMuxServerOpenEphemeralResource(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	testServer1 := &tf6testserver.TestServer{
		GetProviderSchemaResponse: &tfprotov6.GetProviderSchemaResponse{
			EphemeralResourceSchemas: map[string]*tfprotov6.Schema{
				"test_ephemeral_resource_server1": {},
			},
		},
	}
	testServer2 := &tf6testserver.TestServer{
		GetProviderSchemaResponse: &tfprotov6.GetProviderSchemaResponse{
			EphemeralResourceSchemas: map[string]*tfprotov6.Schema{
				"test_ephemeral_resource_server2": {},
			},
		},
	}
	servers := []func() tfprotov6.ProviderServer{testServer1.ProviderServer, testServer2.ProviderServer}
	muxServer, err := tf6muxserver.NewMuxServer(ctx, servers...)

	if err != nil {
		t.Fatalf("unexpected error setting up factory: %s", err)
	}

	_, err = muxServer.ProviderServer().OpenEphemeralResource(ctx, &tfprotov6.OpenEphemeralResourceRequest{
		TypeName: "test_ephemeral_resource_server1",
	})

	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	if !testServer1.OpenEphemeralResourceCalled["test_ephemeral_resource_server1"] {
		t.Errorf("expected test_ephemeral_resource_server1 OpenEphemeralResource to be called on server1")
	}

	if testServer2.OpenEphemeralResourceCalled["test_ephemeral_resource_server1"] {
		t.Errorf("unexpected test_ephemeral_resource_server1 OpenEphemeralResource called on server2")
	}

	_, err = muxServer.ProviderServer().OpenEphemeralResource(ctx, &tfprotov6.OpenEphemeralResourceRequest{
		TypeName: "test_ephemeral_resource_server2",
	})

	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	if testServer1.OpenEphemeralResourceCalled["test_ephemeral_resource_server2"] {
		t.Errorf("unexpected test_ephemeral_resource_server2 OpenEphemeralResource called on server1")
	}

	if !testServer2.OpenEphemeralResourceCalled["test_ephemeral_resource_server2"] {
		t.Errorf("expected test_ephemeral_resource_server2 OpenEphemeralResource to be called on server2")
	}
}
