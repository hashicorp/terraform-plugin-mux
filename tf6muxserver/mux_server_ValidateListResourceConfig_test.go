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

func TestMuxServerValidateListResourceConfig(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	testServer1 := &tf6testserver.TestServer{
		GetProviderSchemaResponse: &tfprotov6.GetProviderSchemaResponse{
			ListResourceSchemas: map[string]*tfprotov6.Schema{
				"test_list_resource_server1": {},
			},
		},
	}
	testServer2 := &tf6testserver.TestServer{
		GetProviderSchemaResponse: &tfprotov6.GetProviderSchemaResponse{
			ListResourceSchemas: map[string]*tfprotov6.Schema{
				"test_list_resource_server2": {},
			},
		},
	}
	servers := []func() tfprotov6.ProviderServer{testServer1.ProviderServer, testServer2.ProviderServer}
	muxServer, err := tf6muxserver.NewMuxServer(ctx, servers...)

	if err != nil {
		t.Fatalf("unexpected error setting up factory: %s", err)
	}

	//nolint:staticcheck // Intentionally verifying interface implementation
	listResourceServer, ok := muxServer.ProviderServer().(tfprotov6.ProviderServerWithListResource)
	if !ok {
		t.Fatal("muxServer should implement tfprotov6.ProviderServerWithListResource")
	}

	_, err = listResourceServer.ValidateListResourceConfig(ctx, &tfprotov6.ValidateListResourceConfigRequest{
		TypeName: "test_list_resource_server1",
	})

	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	if !testServer1.ValidateListResourceConfigCalled["test_list_resource_server1"] {
		t.Errorf("expected test_list_resource_server1 ValidateListResourceConfig to be called on server1")
	}

	if testServer2.ValidateListResourceConfigCalled["test_list_resource_server1"] {
		t.Errorf("unexpected test_list_resource_server1 ValidateListResourceConfig called on server2")
	}

	_, err = listResourceServer.ValidateListResourceConfig(ctx, &tfprotov6.ValidateListResourceConfigRequest{
		TypeName: "test_list_resource_server2",
	})

	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	if testServer1.ValidateListResourceConfigCalled["test_list_resource_server2"] {
		t.Errorf("unexpected test_list_resource_server2 ValidateListResourceConfig called on server1")
	}

	if !testServer2.ValidateListResourceConfigCalled["test_list_resource_server2"] {
		t.Errorf("expected test_list_resource_server2 ValidateListResourceConfig to be called on server2")
	}
}
