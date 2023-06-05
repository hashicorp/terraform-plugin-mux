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

func TestMuxServerValidateResourceTypeConfig(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	testServer1 := &tf5testserver.TestServer{
		ResourceSchemas: map[string]*tfprotov5.Schema{
			"test_resource_server1": {},
		},
	}
	testServer2 := &tf5testserver.TestServer{
		ResourceSchemas: map[string]*tfprotov5.Schema{
			"test_resource_server2": {},
		},
	}
	servers := []func() tfprotov5.ProviderServer{testServer1.ProviderServer, testServer2.ProviderServer}
	muxServer, err := tf5muxserver.NewMuxServer(ctx, servers...)

	if err != nil {
		t.Fatalf("unexpected error setting up factory: %s", err)
	}

	_, err = muxServer.ProviderServer().ValidateResourceTypeConfig(ctx, &tfprotov5.ValidateResourceTypeConfigRequest{
		TypeName: "test_resource_server1",
	})

	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	if !testServer1.ValidateResourceTypeConfigCalled["test_resource_server1"] {
		t.Errorf("expected test_resource_server1 ValidateResourceTypeConfig to be called on server1")
	}

	if testServer2.ValidateResourceTypeConfigCalled["test_resource_server1"] {
		t.Errorf("unexpected test_resource_server1 ValidateResourceTypeConfig called on server2")
	}

	_, err = muxServer.ProviderServer().ValidateResourceTypeConfig(ctx, &tfprotov5.ValidateResourceTypeConfigRequest{
		TypeName: "test_resource_server2",
	})

	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	if testServer1.ValidateResourceTypeConfigCalled["test_resource_server2"] {
		t.Errorf("unexpected test_resource_server2 ValidateResourceTypeConfig called on server1")
	}

	if !testServer2.ValidateResourceTypeConfigCalled["test_resource_server2"] {
		t.Errorf("expected test_resource_server2 ValidateResourceTypeConfig to be called on server2")
	}
}
