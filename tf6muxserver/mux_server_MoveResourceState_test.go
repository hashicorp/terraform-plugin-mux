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

func TestMuxServerMoveResourceState(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	testServer1 := &tf6testserver.TestServer{
		GetProviderSchemaResponse: &tfprotov6.GetProviderSchemaResponse{
			ResourceSchemas: map[string]*tfprotov6.Schema{
				"test_resource1": {},
			},
		},
	}
	testServer2 := &tf6testserver.TestServer{
		GetProviderSchemaResponse: &tfprotov6.GetProviderSchemaResponse{
			ResourceSchemas: map[string]*tfprotov6.Schema{
				"test_resource2": {},
			},
		},
	}

	servers := []func() tfprotov6.ProviderServer{testServer1.ProviderServer, testServer2.ProviderServer}
	muxServer, err := tf6muxserver.NewMuxServer(ctx, servers...)

	if err != nil {
		t.Fatalf("unexpected error setting up factory: %s", err)
	}

	// Reference: https://github.com/hashicorp/terraform-plugin-mux/issues/219
	//nolint:staticcheck // Intentionally verifying interface implementation
	resourceServer, ok := muxServer.ProviderServer().(tfprotov6.ResourceServerWithMoveResourceState)

	if !ok {
		t.Fatal("muxServer should implement tfprotov6.ResourceServerWithMoveResourceState")
	}

	// _, err = muxServer.ProviderServer().MoveResourceState(ctx, &tfprotov6.MoveResourceStateRequest{
	_, err = resourceServer.MoveResourceState(ctx, &tfprotov6.MoveResourceStateRequest{
		TargetTypeName: "test_resource1",
	})

	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	if !testServer1.MoveResourceStateCalled["test_resource1"] {
		t.Errorf("expected test_resource1 MoveResourceState to be called on server1")
	}

	if testServer2.MoveResourceStateCalled["test_resource1"] {
		t.Errorf("unexpected test_resource1 MoveResourceState called on server2")
	}

	// _, err = muxServer.ProviderServer().MoveResourceState(ctx, &tfprotov6.MoveResourceStateRequest{
	_, err = resourceServer.MoveResourceState(ctx, &tfprotov6.MoveResourceStateRequest{
		TargetTypeName: "test_resource2",
	})

	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	if testServer1.MoveResourceStateCalled["test_resource2"] {
		t.Errorf("unexpected test_resource2 MoveResourceState called on server1")
	}

	if !testServer2.MoveResourceStateCalled["test_resource2"] {
		t.Errorf("expected test_resource2 MoveResourceState to be called on server2")
	}
}
