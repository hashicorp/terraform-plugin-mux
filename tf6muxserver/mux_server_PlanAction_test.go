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

func TestMuxServerPlanAction(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	testServer1 := &tf6testserver.TestServer{
		GetProviderSchemaResponse: &tfprotov6.GetProviderSchemaResponse{
			ActionSchemas: map[string]*tfprotov6.ActionSchema{
				"test_action_server1": {},
			},
		},
	}
	testServer2 := &tf6testserver.TestServer{
		GetProviderSchemaResponse: &tfprotov6.GetProviderSchemaResponse{
			ActionSchemas: map[string]*tfprotov6.ActionSchema{
				"test_action_server2": {},
			},
		},
	}
	servers := []func() tfprotov6.ProviderServer{testServer1.ProviderServer, testServer2.ProviderServer}
	muxServer, err := tf6muxserver.NewMuxServer(ctx, servers...)

	if err != nil {
		t.Fatalf("unexpected error setting up factory: %s", err)
	}

	//nolint:staticcheck // Intentionally verifying interface implementation
	actionServer, ok := muxServer.ProviderServer().(tfprotov6.ProviderServerWithActions)
	if !ok {
		t.Fatal("muxServer should implement tfprotov6.ProviderServerWithActions")
	}

	_, err = actionServer.PlanAction(ctx, &tfprotov6.PlanActionRequest{
		ActionType: "test_action_server1",
	})

	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	if !testServer1.PlanActionCalled["test_action_server1"] {
		t.Errorf("expected test_action_server1 PlanAction to be called on server1")
	}

	if testServer2.PlanActionCalled["test_action_server1"] {
		t.Errorf("unexpected test_action_server1 PlanAction called on server2")
	}

	_, err = actionServer.PlanAction(ctx, &tfprotov6.PlanActionRequest{
		ActionType: "test_action_server2",
	})

	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	if testServer1.PlanActionCalled["test_action_server2"] {
		t.Errorf("unexpected test_action_server2 PlanAction called on server1")
	}

	if !testServer2.PlanActionCalled["test_action_server2"] {
		t.Errorf("expected test_action_server2 PlanAction to be called on server2")
	}
}
