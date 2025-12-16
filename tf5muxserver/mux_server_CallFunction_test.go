// Copyright IBM Corp. 2020, 2025
// SPDX-License-Identifier: MPL-2.0

package tf5muxserver_test

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-go/tfprotov5"

	"github.com/hashicorp/terraform-plugin-mux/internal/tf5testserver"
	"github.com/hashicorp/terraform-plugin-mux/tf5muxserver"
)

func TestMuxServerCallFunction(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	testServer1 := &tf5testserver.TestServer{
		GetProviderSchemaResponse: &tfprotov5.GetProviderSchemaResponse{
			Functions: map[string]*tfprotov5.Function{
				"test_function1": {},
			},
		},
	}
	testServer2 := &tf5testserver.TestServer{
		GetProviderSchemaResponse: &tfprotov5.GetProviderSchemaResponse{
			Functions: map[string]*tfprotov5.Function{
				"test_function2": {},
			},
		},
	}

	servers := []func() tfprotov5.ProviderServer{testServer1.ProviderServer, testServer2.ProviderServer}
	muxServer, err := tf5muxserver.NewMuxServer(ctx, servers...)

	if err != nil {
		t.Fatalf("unexpected error setting up factory: %s", err)
	}

	_, err = muxServer.ProviderServer().CallFunction(ctx, &tfprotov5.CallFunctionRequest{
		Name: "test_function1",
	})

	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	if !testServer1.CallFunctionCalled["test_function1"] {
		t.Errorf("expected test_function1 CallFunction to be called on server1")
	}

	if testServer2.CallFunctionCalled["test_function1"] {
		t.Errorf("unexpected test_function1 CallFunction called on server2")
	}

	_, err = muxServer.ProviderServer().CallFunction(ctx, &tfprotov5.CallFunctionRequest{
		Name: "test_function2",
	})

	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	if testServer1.CallFunctionCalled["test_function2"] {
		t.Errorf("unexpected test_function2 CallFunction called on server1")
	}

	if !testServer2.CallFunctionCalled["test_function2"] {
		t.Errorf("expected test_function2 CallFunction to be called on server2")
	}
}
