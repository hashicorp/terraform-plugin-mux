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

func TestMuxServerCallFunction(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	testServer1 := &tf6testserver.TestServer{
		GetProviderSchemaResponse: &tfprotov6.GetProviderSchemaResponse{
			Functions: map[string]*tfprotov6.Function{
				"test_function1": {},
			},
		},
	}
	testServer2 := &tf6testserver.TestServer{
		GetProviderSchemaResponse: &tfprotov6.GetProviderSchemaResponse{
			Functions: map[string]*tfprotov6.Function{
				"test_function2": {},
			},
		},
	}

	servers := []func() tfprotov6.ProviderServer{testServer1.ProviderServer, testServer2.ProviderServer}
	muxServer, err := tf6muxserver.NewMuxServer(ctx, servers...)

	if err != nil {
		t.Fatalf("unexpected error setting up factory: %s", err)
	}

	// Reference: https://github.com/hashicorp/terraform-plugin-mux/issues/210
	functionServer, ok := muxServer.ProviderServer().(tfprotov6.FunctionServer)

	if !ok {
		t.Fatal("muxServer should implement tfprotov6.FunctionServer")
	}

	// _, err = muxServer.ProviderServer().CallFunction(ctx, &tfprotov6.CallFunctionRequest{
	_, err = functionServer.CallFunction(ctx, &tfprotov6.CallFunctionRequest{
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

	// _, err = muxServer.ProviderServer().CallFunction(ctx, &tfprotov6.CallFunctionRequest{
	_, err = functionServer.CallFunction(ctx, &tfprotov6.CallFunctionRequest{
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
