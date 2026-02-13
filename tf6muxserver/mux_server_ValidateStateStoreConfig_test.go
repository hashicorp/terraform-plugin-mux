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

func TestMuxServerValidateStateStoreConfig(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	testServer1 := &tf6testserver.TestServer{
		GetProviderSchemaResponse: &tfprotov6.GetProviderSchemaResponse{
			StateStoreSchemas: map[string]*tfprotov6.Schema{
				"test_state_store_server1": {},
			},
		},
	}
	testServer2 := &tf6testserver.TestServer{
		GetProviderSchemaResponse: &tfprotov6.GetProviderSchemaResponse{
			StateStoreSchemas: map[string]*tfprotov6.Schema{
				"test_state_store_server2": {},
			},
		},
	}
	servers := []func() tfprotov6.ProviderServer{testServer1.ProviderServer, testServer2.ProviderServer}
	muxServer, err := tf6muxserver.NewMuxServer(ctx, servers...)

	if err != nil {
		t.Fatalf("unexpected error setting up factory: %s", err)
	}

	//nolint:staticcheck // Intentionally verifying interface implementation
	stateStoreServer, ok := muxServer.ProviderServer().(tfprotov6.StateStoreServer)
	if !ok {
		t.Fatal("muxServer should implement tfprotov6.StateStoreServer")
	}

	_, err = stateStoreServer.ValidateStateStoreConfig(ctx, &tfprotov6.ValidateStateStoreConfigRequest{
		TypeName: "test_state_store_server1",
	})

	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	if !testServer1.ValidateStateStoreConfigCalled["test_state_store_server1"] {
		t.Errorf("expected test_state_store_server1 ValidateStateStoreConfig to be called on server1")
	}

	if testServer2.ValidateStateStoreConfigCalled["test_state_store_server1"] {
		t.Errorf("unexpected test_state_store_server1 ValidateStateStoreConfig called on server2")
	}

	_, err = stateStoreServer.ValidateStateStoreConfig(ctx, &tfprotov6.ValidateStateStoreConfigRequest{
		TypeName: "test_state_store_server2",
	})

	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	if testServer1.ValidateStateStoreConfigCalled["test_state_store_server2"] {
		t.Errorf("unexpected test_state_store_server2 ValidateStateStoreConfig called on server1")
	}

	if !testServer2.ValidateStateStoreConfigCalled["test_state_store_server2"] {
		t.Errorf("expected test_state_store_server2 ValidateStateStoreConfig to be called on server2")
	}
}
