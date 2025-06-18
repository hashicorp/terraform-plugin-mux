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

func TestMuxServerListResource(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	testServer1 := &tf5testserver.TestServer{
		GetProviderSchemaResponse: &tfprotov5.GetProviderSchemaResponse{
			ListResourceSchemas: map[string]*tfprotov5.Schema{
				"test_list_resource_server1": {},
			},
		},
	}
	testServer2 := &tf5testserver.TestServer{
		GetProviderSchemaResponse: &tfprotov5.GetProviderSchemaResponse{
			ListResourceSchemas: map[string]*tfprotov5.Schema{
				"test_list_resource_server2": {},
			},
		},
	}
	servers := []func() tfprotov5.ProviderServer{testServer1.ProviderServer, testServer2.ProviderServer}
	muxServer, err := tf5muxserver.NewMuxServer(ctx, servers...)

	if err != nil {
		t.Fatalf("unexpected error setting up factory: %s", err)
	}

	//nolint:staticcheck // Intentionally verifying interface implementation
	listResourceServer, ok := muxServer.ProviderServer().(tfprotov5.ProviderServerWithListResource)
	if !ok {
		t.Fatal("muxServer should implement tfprotov5.ProviderServerWithListResource")
	}

	_, err = listResourceServer.ListResource(ctx, &tfprotov5.ListResourceRequest{
		TypeName: "test_list_resource_server1",
	})

	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	if !testServer1.ListResourceCalled["test_list_resource_server1"] {
		t.Errorf("expected test_list_resource_server1 ListResource to be called on server1")
	}

	if testServer2.ListResourceCalled["test_list_resource_server1"] {
		t.Errorf("unexpected test_list_resource_server1 ListResource called on server2")
	}

	_, err = listResourceServer.ListResource(ctx, &tfprotov5.ListResourceRequest{
		TypeName: "test_list_resource_server2",
	})

	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	if testServer1.ListResourceCalled["test_list_resource_server2"] {
		t.Errorf("unexpected test_list_resource_server2 ListResource called on server1")
	}

	if !testServer2.ListResourceCalled["test_list_resource_server2"] {
		t.Errorf("expected test_list_resource_server2 ListResource to be called on server2")
	}
}
