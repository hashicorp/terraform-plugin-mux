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

func TestMuxServerValidateDataResourceConfig(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	testServer1 := &tf6testserver.TestServer{
		DataSourceSchemas: map[string]*tfprotov6.Schema{
			"test_data_source_server1": {},
		},
	}
	testServer2 := &tf6testserver.TestServer{
		DataSourceSchemas: map[string]*tfprotov6.Schema{
			"test_data_source_server2": {},
		},
	}

	servers := []func() tfprotov6.ProviderServer{testServer1.ProviderServer, testServer2.ProviderServer}
	muxServer, err := tf6muxserver.NewMuxServer(ctx, servers...)

	if err != nil {
		t.Fatalf("unexpected error setting up factory: %s", err)
	}

	_, err = muxServer.ProviderServer().ValidateDataResourceConfig(ctx, &tfprotov6.ValidateDataResourceConfigRequest{
		TypeName: "test_data_source_server1",
	})

	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	if !testServer1.ValidateDataResourceConfigCalled["test_data_source_server1"] {
		t.Errorf("expected test_data_source_server1 ValidateDataResourceConfig to be called on server1")
	}

	if testServer2.ValidateDataResourceConfigCalled["test_data_source_server1"] {
		t.Errorf("unexpected test_data_source_server1 ValidateDataResourceConfig called on server2")
	}

	_, err = muxServer.ProviderServer().ValidateDataResourceConfig(ctx, &tfprotov6.ValidateDataResourceConfigRequest{
		TypeName: "test_data_source_server2",
	})

	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	if testServer1.ValidateDataResourceConfigCalled["test_data_source_server2"] {
		t.Errorf("unexpected test_data_source_server2 ValidateDataResourceConfig called on server1")
	}

	if !testServer2.ValidateDataResourceConfigCalled["test_data_source_server2"] {
		t.Errorf("expected test_data_source_server2 ValidateDataResourceConfig to be called on server2")
	}
}
