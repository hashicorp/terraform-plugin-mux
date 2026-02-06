// Copyright IBM Corp. 2020, 2026
// SPDX-License-Identifier: MPL-2.0

package tf6muxserver_test

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-go/tfprotov6"

	"github.com/hashicorp/terraform-plugin-mux/internal/tf6testserver"
	"github.com/hashicorp/terraform-plugin-mux/tf6muxserver"
)

func TestMuxServerStateStore_ConfigureStateStore(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	testServer1 := &tf6testserver.TestServer{
		GetMetadataResponse: &tfprotov6.GetMetadataResponse{
			StateStores: []tfprotov6.StateStoreMetadata{
				{
					TypeName: "test_statestore_server1",
				},
			},
		},
	}
	testServer2 := &tf6testserver.TestServer{
		GetMetadataResponse: &tfprotov6.GetMetadataResponse{
			StateStores: []tfprotov6.StateStoreMetadata{
				{
					TypeName: "test_statestore_server2",
				},
			},
		},
	}

	servers := []func() tfprotov6.ProviderServer{testServer1.ProviderServer, testServer2.ProviderServer}
	muxServer, err := tf6muxserver.NewMuxServer(ctx, servers...)

	if err != nil {
		t.Fatalf("unexpected error setting up mux server: %s", err)
	}

	_, err = muxServer.ProviderServer().(tfprotov6.StateStoreServer).ConfigureStateStore(ctx, &tfprotov6.ConfigureStateStoreRequest{
		TypeName: "test_statestore_server1",
	})

	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	if !testServer1.ConfigureStateStoreCalled["test_statestore_server1"] {
		t.Errorf("expected test_statestore_server1 ConfigureStateStore to be called on server1")
	}

	if testServer2.ConfigureStateStoreCalled["test_statestore_server1"] {
		t.Errorf("unexpected test_statestore_server1 ConfigureStateStore called on server2")
	}

	_, err = muxServer.ProviderServer().(tfprotov6.StateStoreServer).ConfigureStateStore(ctx, &tfprotov6.ConfigureStateStoreRequest{
		TypeName: "test_statestore_server2",
	})

	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	if !testServer2.ConfigureStateStoreCalled["test_statestore_server2"] {
		t.Errorf("expected test_statestore_server2 ConfigureStateStore to be called on server2")
	}

	if testServer1.ConfigureStateStoreCalled["test_statestore_server2"] {
		t.Errorf("unexpected test_statestore_server2 ConfigureStateStore called on server1")
	}
}
