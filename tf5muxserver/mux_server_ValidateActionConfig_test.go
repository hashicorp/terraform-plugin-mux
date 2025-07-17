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

func TestMuxServerValidateActionConfig(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	testServer1 := &tf5testserver.TestServer{
		GetProviderSchemaResponse: &tfprotov5.GetProviderSchemaResponse{
			ActionSchemas: map[string]*tfprotov5.ActionSchema{
				"test_action_server1": {},
			},
		},
	}
	testServer2 := &tf5testserver.TestServer{
		GetProviderSchemaResponse: &tfprotov5.GetProviderSchemaResponse{
			ActionSchemas: map[string]*tfprotov5.ActionSchema{
				"test_action_server2": {},
			},
		},
	}
	servers := []func() tfprotov5.ProviderServer{testServer1.ProviderServer, testServer2.ProviderServer}
	muxServer, err := tf5muxserver.NewMuxServer(ctx, servers...)

	if err != nil {
		t.Fatalf("unexpected error setting up factory: %s", err)
	}

	//nolint:staticcheck // Intentionally verifying interface implementation
	actionServer, ok := muxServer.ProviderServer().(tfprotov5.ProviderServerWithActions)
	if !ok {
		t.Fatal("muxServer should implement tfprotov5.ProviderServerWithActions")
	}

	_, err = actionServer.ValidateActionConfig(ctx, &tfprotov5.ValidateActionConfigRequest{
		ActionType: "test_action_server1",
	})

	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	if !testServer1.ValidateActionConfigCalled["test_action_server1"] {
		t.Errorf("expected test_action_server1 ValidateActionConfig to be called on server1")
	}

	if testServer2.ValidateActionConfigCalled["test_action_server1"] {
		t.Errorf("unexpected test_action_server1 ValidateActionConfig called on server2")
	}

	_, err = actionServer.ValidateActionConfig(ctx, &tfprotov5.ValidateActionConfigRequest{
		ActionType: "test_action_server2",
	})

	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	if testServer1.ValidateActionConfigCalled["test_action_server2"] {
		t.Errorf("unexpected test_action_server2 ValidateActionConfig called on server1")
	}

	if !testServer2.ValidateActionConfigCalled["test_action_server2"] {
		t.Errorf("expected test_action_server2 ValidateActionConfig to be called on server2")
	}
}
