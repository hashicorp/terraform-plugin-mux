// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package tf6muxserver_test

import (
	"context"
	"strings"
	"sync"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-go/tftypes"

	"github.com/hashicorp/terraform-plugin-mux/internal/tf6testserver"
	"github.com/hashicorp/terraform-plugin-mux/tf6muxserver"
)

func TestMuxServerGetDataSourceServer_GetMetadata(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	testServer1 := &tf6testserver.TestServer{
		GetMetadataResponse: &tfprotov6.GetMetadataResponse{
			DataSources: []tfprotov6.DataSourceMetadata{
				{
					TypeName: "test_datasource_server1",
				},
			},
		},
	}
	testServer2 := &tf6testserver.TestServer{
		GetMetadataResponse: &tfprotov6.GetMetadataResponse{
			DataSources: []tfprotov6.DataSourceMetadata{
				{
					TypeName: "test_datasource_server2",
				},
			},
		},
	}

	servers := []func() tfprotov6.ProviderServer{testServer1.ProviderServer, testServer2.ProviderServer}
	muxServer, err := tf6muxserver.NewMuxServer(ctx, servers...)

	if err != nil {
		t.Fatalf("unexpected error setting up factory: %s", err)
	}

	// When GetProviderSchemaOptional is enabled, the secondary provider
	// instances will receive non-GetProviderSchema RPCs such as
	// ValidateDataResourceConfig which will cause getDataSourceServer to perform
	// server discovery. This testing also simulates concurrent operations from
	// Terraform to verify the mutex does not deadlock.
	var wg sync.WaitGroup

	terraformOp := func() {
		defer wg.Done()

		_, _ = muxServer.ProviderServer().ValidateDataResourceConfig(ctx, &tfprotov6.ValidateDataResourceConfigRequest{
			TypeName: "test_datasource_server1",
		})
	}

	wg.Add(2)
	go terraformOp()
	go terraformOp()

	wg.Wait()

	if !testServer1.ValidateDataResourceConfigCalled["test_datasource_server1"] {
		t.Errorf("expected test_datasource_server1 ValidateDataResourceConfig to be called on server1")
	}
}

func TestMuxServerGetDataSourceServer_GetMetadata_Duplicate(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	testServer1 := &tf6testserver.TestServer{
		GetMetadataResponse: &tfprotov6.GetMetadataResponse{
			DataSources: []tfprotov6.DataSourceMetadata{
				{
					TypeName: "test_datasource_server", // intentionally duplicated
				},
			},
		},
	}
	testServer2 := &tf6testserver.TestServer{
		GetMetadataResponse: &tfprotov6.GetMetadataResponse{
			DataSources: []tfprotov6.DataSourceMetadata{
				{
					TypeName: "test_datasource_server", // intentionally duplicated
				},
			},
		},
	}

	servers := []func() tfprotov6.ProviderServer{testServer1.ProviderServer, testServer2.ProviderServer}
	muxServer, err := tf6muxserver.NewMuxServer(ctx, servers...)

	if err != nil {
		t.Fatalf("unexpected error setting up factory: %s", err)
	}

	// When GetProviderSchemaOptional is enabled, the secondary provider
	// instances will receive non-GetProviderSchema RPCs such as
	// ValidateDataResourceConfig which will cause getDataSourceServer to perform
	// server discovery. This testing also simulates concurrent operations from
	// Terraform to verify the mutex does not deadlock.
	var wg sync.WaitGroup

	expectedDiags := []*tfprotov6.Diagnostic{
		{
			Severity: tfprotov6.DiagnosticSeverityError,
			Summary:  "Invalid Provider Server Combination",
			Detail: "The combined provider has multiple implementations of the same data source type across underlying providers. " +
				"Data source types must be implemented by only one underlying provider. " +
				"This is always an issue in the provider implementation and should be reported to the provider developers.\n\n" +
				"Duplicate data source type: test_datasource_server",
		},
	}

	terraformOp := func() {
		defer wg.Done()

		resp, _ := muxServer.ProviderServer().ValidateDataResourceConfig(ctx, &tfprotov6.ValidateDataResourceConfigRequest{
			TypeName: "test_datasource_server",
		})

		if diff := cmp.Diff(resp.Diagnostics, expectedDiags); diff != "" {
			t.Errorf("unexpected diagnostics difference: %s", diff)
		}
	}

	wg.Add(2)
	go terraformOp()
	go terraformOp()

	wg.Wait()

	if testServer1.ValidateDataResourceConfigCalled["test_datasource_server"] {
		t.Errorf("unexpected test_datasource_server ValidateDataResourceConfig called on server1")
	}

	if testServer2.ValidateDataResourceConfigCalled["test_datasource_server"] {
		t.Errorf("unexpected test_datasource_server ValidateDataResourceConfig called on server2")
	}
}

func TestMuxServerGetDataSourceServer_GetMetadata_Partial(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	testServer1 := &tf6testserver.TestServer{
		GetMetadataResponse: &tfprotov6.GetMetadataResponse{
			DataSources: []tfprotov6.DataSourceMetadata{
				{
					TypeName: "test_datasource_server1",
				},
			},
		},
	}
	testServer2 := &tf6testserver.TestServer{
		GetProviderSchemaResponse: &tfprotov6.GetProviderSchemaResponse{
			DataSourceSchemas: map[string]*tfprotov6.Schema{
				"test_datasource_server2": {},
			},
		},
	}

	servers := []func() tfprotov6.ProviderServer{testServer1.ProviderServer, testServer2.ProviderServer}
	muxServer, err := tf6muxserver.NewMuxServer(ctx, servers...)

	if err != nil {
		t.Fatalf("unexpected error setting up factory: %s", err)
	}

	// When GetProviderSchemaOptional is enabled, the secondary provider
	// instances will receive non-GetProviderSchema RPCs such as
	// ValidateDataResourceConfig which will cause getDataSourceServer to perform
	// server discovery. This testing also simulates concurrent operations from
	// Terraform to verify the mutex does not deadlock.
	var wg sync.WaitGroup

	terraformOp := func() {
		defer wg.Done()

		_, _ = muxServer.ProviderServer().ValidateDataResourceConfig(ctx, &tfprotov6.ValidateDataResourceConfigRequest{
			TypeName: "test_datasource_server1",
		})
	}

	wg.Add(2)
	go terraformOp()
	go terraformOp()

	wg.Wait()

	if !testServer1.ValidateDataResourceConfigCalled["test_datasource_server1"] {
		t.Errorf("expected test_datasource_server1 ValidateDataResourceConfig to be called on server1")
	}
}

func TestMuxServerGetDataSourceServer_GetProviderSchema(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	testServer1 := &tf6testserver.TestServer{
		GetProviderSchemaResponse: &tfprotov6.GetProviderSchemaResponse{
			DataSourceSchemas: map[string]*tfprotov6.Schema{
				"test_datasource_server1": {},
			},
		},
	}
	testServer2 := &tf6testserver.TestServer{
		GetProviderSchemaResponse: &tfprotov6.GetProviderSchemaResponse{
			DataSourceSchemas: map[string]*tfprotov6.Schema{
				"test_datasource_server2": {},
			},
		},
	}

	servers := []func() tfprotov6.ProviderServer{testServer1.ProviderServer, testServer2.ProviderServer}
	muxServer, err := tf6muxserver.NewMuxServer(ctx, servers...)

	if err != nil {
		t.Fatalf("unexpected error setting up factory: %s", err)
	}

	// When GetProviderSchemaOptional is enabled, the secondary provider
	// instances will receive non-GetProviderSchema RPCs such as
	// ValidateDataResourceConfig which will cause getDataSourceServer to perform
	// server discovery. This testing also simulates concurrent operations from
	// Terraform to verify the mutex does not deadlock.
	var wg sync.WaitGroup

	terraformOp := func() {
		defer wg.Done()

		_, _ = muxServer.ProviderServer().ValidateDataResourceConfig(ctx, &tfprotov6.ValidateDataResourceConfigRequest{
			TypeName: "test_datasource_server1",
		})
	}

	wg.Add(2)
	go terraformOp()
	go terraformOp()

	wg.Wait()

	if !testServer1.ValidateDataResourceConfigCalled["test_datasource_server1"] {
		t.Errorf("expected test_datasource_server1 ValidateDataResourceConfig to be called on server1")
	}
}

func TestMuxServerGetDataSourceServer_GetProviderSchema_Duplicate(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	testServer1 := &tf6testserver.TestServer{
		GetProviderSchemaResponse: &tfprotov6.GetProviderSchemaResponse{
			DataSourceSchemas: map[string]*tfprotov6.Schema{
				"test_datasource_server": {}, // intentionally duplicated
			},
		},
	}
	testServer2 := &tf6testserver.TestServer{
		GetProviderSchemaResponse: &tfprotov6.GetProviderSchemaResponse{
			DataSourceSchemas: map[string]*tfprotov6.Schema{
				"test_datasource_server": {}, // intentionally duplicated
			},
		},
	}

	servers := []func() tfprotov6.ProviderServer{testServer1.ProviderServer, testServer2.ProviderServer}
	muxServer, err := tf6muxserver.NewMuxServer(ctx, servers...)

	if err != nil {
		t.Fatalf("unexpected error setting up factory: %s", err)
	}

	// When GetProviderSchemaOptional is enabled, the secondary provider
	// instances will receive non-GetProviderSchema RPCs such as
	// ValidateDataResourceConfig which will cause getDataSourceServer to perform
	// server discovery. This testing also simulates concurrent operations from
	// Terraform to verify the mutex does not deadlock.
	var wg sync.WaitGroup

	expectedDiags := []*tfprotov6.Diagnostic{
		{
			Severity: tfprotov6.DiagnosticSeverityError,
			Summary:  "Invalid Provider Server Combination",
			Detail: "The combined provider has multiple implementations of the same data source type across underlying providers. " +
				"Data source types must be implemented by only one underlying provider. " +
				"This is always an issue in the provider implementation and should be reported to the provider developers.\n\n" +
				"Duplicate data source type: test_datasource_server",
		},
	}

	terraformOp := func() {
		defer wg.Done()

		resp, _ := muxServer.ProviderServer().ValidateDataResourceConfig(ctx, &tfprotov6.ValidateDataResourceConfigRequest{
			TypeName: "test_datasource_server",
		})

		if diff := cmp.Diff(resp.Diagnostics, expectedDiags); diff != "" {
			t.Errorf("unexpected diagnostics difference: %s", diff)
		}
	}

	wg.Add(2)
	go terraformOp()
	go terraformOp()

	wg.Wait()

	if testServer1.ValidateDataResourceConfigCalled["test_datasource_server"] {
		t.Errorf("unexpected test_datasource_server ValidateDataResourceConfig called on server1")
	}

	if testServer2.ValidateDataResourceConfigCalled["test_datasource_server"] {
		t.Errorf("unexpected test_datasource_server ValidateDataResourceConfig called on server2")
	}
}

func TestMuxServerGetDataSourceServer_Missing(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	testServer1 := &tf6testserver.TestServer{
		GetMetadataResponse: &tfprotov6.GetMetadataResponse{
			DataSources: []tfprotov6.DataSourceMetadata{
				{
					TypeName: "test_datasource_server1",
				},
			},
		},
	}
	testServer2 := &tf6testserver.TestServer{
		GetMetadataResponse: &tfprotov6.GetMetadataResponse{
			DataSources: []tfprotov6.DataSourceMetadata{
				{
					TypeName: "test_datasource_server2",
				},
			},
		},
	}

	servers := []func() tfprotov6.ProviderServer{testServer1.ProviderServer, testServer2.ProviderServer}
	muxServer, err := tf6muxserver.NewMuxServer(ctx, servers...)

	if err != nil {
		t.Fatalf("unexpected error setting up factory: %s", err)
	}

	// When GetProviderSchemaOptional is enabled, the secondary provider
	// instances will receive non-GetProviderSchema RPCs such as
	// ValidateDataResourceConfig which will cause getDataSourceServer to perform
	// server discovery. This testing also simulates concurrent operations from
	// Terraform to verify the mutex does not deadlock.
	var wg sync.WaitGroup

	expectedDiags := []*tfprotov6.Diagnostic{
		{
			Severity: tfprotov6.DiagnosticSeverityError,
			Summary:  "Data Source Not Implemented",
			Detail: "The combined provider does not implement the requested data source type. " +
				"This is always an issue in the provider implementation and should be reported to the provider developers.\n\n" +
				"Missing data source type: test_datasource_nonexistent",
		},
	}

	terraformOp := func() {
		defer wg.Done()

		resp, _ := muxServer.ProviderServer().ValidateDataResourceConfig(ctx, &tfprotov6.ValidateDataResourceConfigRequest{
			TypeName: "test_datasource_nonexistent",
		})

		if diff := cmp.Diff(resp.Diagnostics, expectedDiags); diff != "" {
			t.Errorf("unexpected diagnostics difference: %s", diff)
		}
	}

	wg.Add(2)
	go terraformOp()
	go terraformOp()

	wg.Wait()

	if testServer1.ValidateDataResourceConfigCalled["test_datasource_nonexistent"] {
		t.Errorf("unexpected test_datasource_nonexistent ValidateDataResourceConfig called on server1")
	}

	if testServer2.ValidateDataResourceConfigCalled["test_datasource_nonexistent"] {
		t.Errorf("unexpected test_datasource_nonexistent ValidateDataResourceConfig called on server2")
	}
}

func TestMuxServerGetFunctionServer_GetProviderSchema(t *testing.T) {
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

	// When GetProviderSchemaOptional is enabled, the secondary provider
	// instances will receive non-GetProviderSchema RPCs such as
	// CallFunction which will cause getFunctionServer to perform
	// server discovery. This testing also simulates concurrent operations from
	// Terraform to verify the mutex does not deadlock.
	var wg sync.WaitGroup

	terraformOp := func() {
		defer wg.Done()

		_, _ = muxServer.ProviderServer().CallFunction(ctx, &tfprotov6.CallFunctionRequest{
			Name: "test_function1",
		})
	}

	wg.Add(2)
	go terraformOp()
	go terraformOp()

	wg.Wait()

	if !testServer1.CallFunctionCalled["test_function1"] {
		t.Errorf("expected test_function1 CallFunction to be called on server1")
	}
}

func TestMuxServerGetFunctionServer_GetProviderSchema_Duplicate(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	testServer1 := &tf6testserver.TestServer{
		GetProviderSchemaResponse: &tfprotov6.GetProviderSchemaResponse{
			Functions: map[string]*tfprotov6.Function{
				"test_function": {}, // intentionally duplicated
			},
		},
	}
	testServer2 := &tf6testserver.TestServer{
		GetProviderSchemaResponse: &tfprotov6.GetProviderSchemaResponse{
			Functions: map[string]*tfprotov6.Function{
				"test_function": {}, // intentionally duplicated
			},
		},
	}

	servers := []func() tfprotov6.ProviderServer{testServer1.ProviderServer, testServer2.ProviderServer}
	muxServer, err := tf6muxserver.NewMuxServer(ctx, servers...)

	if err != nil {
		t.Fatalf("unexpected error setting up factory: %s", err)
	}

	// When GetProviderSchemaOptional is enabled, the secondary provider
	// instances will receive non-GetProviderSchema RPCs such as
	// CallFunction which will cause getFunctionServer to perform
	// server discovery. This testing also simulates concurrent operations from
	// Terraform to verify the mutex does not deadlock.
	var wg sync.WaitGroup

	expectedError := &tfprotov6.FunctionError{
		Text: "Invalid Provider Server Combination: The combined provider has multiple implementations of the same function name across underlying providers. " +
			"Functions must be implemented by only one underlying provider. " +
			"This is always an issue in the provider implementation and should be reported to the provider developers.\n\n" +
			"Duplicate function: test_function",
	}

	terraformOp := func() {
		defer wg.Done()

		resp, _ := muxServer.ProviderServer().CallFunction(ctx, &tfprotov6.CallFunctionRequest{
			Name: "test_function",
		})

		if diff := cmp.Diff(resp.Error, expectedError); diff != "" {
			t.Errorf("unexpected diagnostics difference: %s", diff)
		}
	}

	wg.Add(2)
	go terraformOp()
	go terraformOp()

	wg.Wait()

	if testServer1.CallFunctionCalled["test_function"] {
		t.Errorf("unexpected test_function CallFunction called on server1")
	}

	if testServer2.CallFunctionCalled["test_function"] {
		t.Errorf("unexpected test_function CallFunction called on server2")
	}
}

func TestMuxServerGetFunctionServer_GetMetadata(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	testServer1 := &tf6testserver.TestServer{
		GetMetadataResponse: &tfprotov6.GetMetadataResponse{
			Functions: []tfprotov6.FunctionMetadata{
				{
					Name: "test_function1",
				},
			},
		},
	}
	testServer2 := &tf6testserver.TestServer{
		GetMetadataResponse: &tfprotov6.GetMetadataResponse{
			Functions: []tfprotov6.FunctionMetadata{
				{
					Name: "test_function2",
				},
			},
		},
	}

	servers := []func() tfprotov6.ProviderServer{testServer1.ProviderServer, testServer2.ProviderServer}
	muxServer, err := tf6muxserver.NewMuxServer(ctx, servers...)

	if err != nil {
		t.Fatalf("unexpected error setting up factory: %s", err)
	}

	// When GetProviderSchemaOptional is enabled, the secondary provider
	// instances will receive non-GetProviderSchema RPCs such as
	// CallFunction which will cause getFunctionServer to perform
	// server discovery. This testing also simulates concurrent operations from
	// Terraform to verify the mutex does not deadlock.
	var wg sync.WaitGroup

	terraformOp := func() {
		defer wg.Done()

		_, _ = muxServer.ProviderServer().CallFunction(ctx, &tfprotov6.CallFunctionRequest{
			Name: "test_function1",
		})
	}

	wg.Add(2)
	go terraformOp()
	go terraformOp()

	wg.Wait()

	if !testServer1.CallFunctionCalled["test_function1"] {
		t.Errorf("expected test_function1 CallFunction to be called on server1")
	}
}

func TestMuxServerGetFunctionServer_GetMetadata_Duplicate(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	testServer1 := &tf6testserver.TestServer{
		GetMetadataResponse: &tfprotov6.GetMetadataResponse{
			Functions: []tfprotov6.FunctionMetadata{
				{
					Name: "test_function", // intentionally duplicated
				},
			},
		},
	}
	testServer2 := &tf6testserver.TestServer{
		GetMetadataResponse: &tfprotov6.GetMetadataResponse{
			Functions: []tfprotov6.FunctionMetadata{
				{
					Name: "test_function", // intentionally duplicated
				},
			},
		},
	}

	servers := []func() tfprotov6.ProviderServer{testServer1.ProviderServer, testServer2.ProviderServer}
	muxServer, err := tf6muxserver.NewMuxServer(ctx, servers...)

	if err != nil {
		t.Fatalf("unexpected error setting up factory: %s", err)
	}

	// When GetProviderSchemaOptional is enabled, the secondary provider
	// instances will receive non-GetProviderSchema RPCs such as
	// CallFunction which will cause getFunctionServer to perform
	// server discovery. This testing also simulates concurrent operations from
	// Terraform to verify the mutex does not deadlock.
	var wg sync.WaitGroup

	expectedError := &tfprotov6.FunctionError{
		Text: "Invalid Provider Server Combination: The combined provider has multiple implementations of the same function name across underlying providers. " +
			"Functions must be implemented by only one underlying provider. " +
			"This is always an issue in the provider implementation and should be reported to the provider developers.\n\n" +
			"Duplicate function: test_function",
	}

	terraformOp := func() {
		defer wg.Done()

		resp, _ := muxServer.ProviderServer().CallFunction(ctx, &tfprotov6.CallFunctionRequest{
			Name: "test_function",
		})

		if diff := cmp.Diff(resp.Error, expectedError); diff != "" {
			t.Errorf("unexpected diagnostics difference: %s", diff)
		}
	}

	wg.Add(2)
	go terraformOp()
	go terraformOp()

	wg.Wait()

	if testServer1.CallFunctionCalled["test_function"] {
		t.Errorf("unexpected test_function CallFunction called on server1")
	}

	if testServer2.CallFunctionCalled["test_function"] {
		t.Errorf("unexpected test_function CallFunction called on server2")
	}
}

func TestMuxServerGetFunctionServer_GetMetadata_Partial(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	testServer1 := &tf6testserver.TestServer{
		GetMetadataResponse: &tfprotov6.GetMetadataResponse{
			Functions: []tfprotov6.FunctionMetadata{
				{
					Name: "test_function1",
				},
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

	// When GetProviderSchemaOptional is enabled, the secondary provider
	// instances will receive non-GetProviderSchema RPCs such as
	// CallFunction which will cause getFunctionServer to perform
	// server discovery. This testing also simulates concurrent operations from
	// Terraform to verify the mutex does not deadlock.
	var wg sync.WaitGroup

	terraformOp := func() {
		defer wg.Done()

		_, _ = muxServer.ProviderServer().CallFunction(ctx, &tfprotov6.CallFunctionRequest{
			Name: "test_function1",
		})
	}

	wg.Add(2)
	go terraformOp()
	go terraformOp()

	wg.Wait()

	if !testServer1.CallFunctionCalled["test_function1"] {
		t.Errorf("expected test_function1 CallFunction to be called on server1")
	}
}

func TestMuxServerGetFunctionServer_Missing(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	testServer1 := &tf6testserver.TestServer{
		GetMetadataResponse: &tfprotov6.GetMetadataResponse{
			Functions: []tfprotov6.FunctionMetadata{
				{
					Name: "test_function1",
				},
			},
		},
	}
	testServer2 := &tf6testserver.TestServer{
		GetMetadataResponse: &tfprotov6.GetMetadataResponse{
			Functions: []tfprotov6.FunctionMetadata{
				{
					Name: "test_function2",
				},
			},
		},
	}

	servers := []func() tfprotov6.ProviderServer{testServer1.ProviderServer, testServer2.ProviderServer}
	muxServer, err := tf6muxserver.NewMuxServer(ctx, servers...)

	if err != nil {
		t.Fatalf("unexpected error setting up factory: %s", err)
	}

	// When GetProviderSchemaOptional is enabled, the secondary provider
	// instances will receive non-GetProviderSchema RPCs such as
	// CallFunction which will cause getFunctionServer to perform
	// server discovery. This testing also simulates concurrent operations from
	// Terraform to verify the mutex does not deadlock.
	var wg sync.WaitGroup

	expectedError := &tfprotov6.FunctionError{
		Text: "Function Not Implemented: The combined provider does not implement the requested function. " +
			"This is always an issue in the provider implementation and should be reported to the provider developers.\n\n" +
			"Missing function: test_function_nonexistent",
	}

	terraformOp := func() {
		defer wg.Done()

		resp, _ := muxServer.ProviderServer().CallFunction(ctx, &tfprotov6.CallFunctionRequest{
			Name: "test_function_nonexistent",
		})

		if diff := cmp.Diff(resp.Error, expectedError); diff != "" {
			t.Errorf("unexpected diagnostics difference: %s", diff)
		}
	}

	wg.Add(2)
	go terraformOp()
	go terraformOp()

	wg.Wait()

	if testServer1.CallFunctionCalled["test_function_nonexistent"] {
		t.Errorf("unexpected test_function_nonexistent CallFunction called on server1")
	}

	if testServer2.CallFunctionCalled["test_function_nonexistent"] {
		t.Errorf("unexpected test_function_nonexistent CallFunction called on server2")
	}
}

func TestMuxServerGetResourceServer_GetProviderSchema(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	testServer1 := &tf6testserver.TestServer{
		GetProviderSchemaResponse: &tfprotov6.GetProviderSchemaResponse{
			ResourceSchemas: map[string]*tfprotov6.Schema{
				"test_resource_server1": {},
			},
		},
	}
	testServer2 := &tf6testserver.TestServer{
		GetProviderSchemaResponse: &tfprotov6.GetProviderSchemaResponse{
			ResourceSchemas: map[string]*tfprotov6.Schema{
				"test_resource_server2": {},
			},
		},
	}

	servers := []func() tfprotov6.ProviderServer{testServer1.ProviderServer, testServer2.ProviderServer}
	muxServer, err := tf6muxserver.NewMuxServer(ctx, servers...)

	if err != nil {
		t.Fatalf("unexpected error setting up factory: %s", err)
	}

	// When GetProviderSchemaOptional is enabled, the secondary provider
	// instances will receive non-GetProviderSchema RPCs such as
	// ValidateResourceConfig which will cause getResourceServer to perform
	// server discovery. This testing also simulates concurrent operations from
	// Terraform to verify the mutex does not deadlock.
	var wg sync.WaitGroup

	terraformOp := func() {
		defer wg.Done()

		_, _ = muxServer.ProviderServer().ValidateResourceConfig(ctx, &tfprotov6.ValidateResourceConfigRequest{
			TypeName: "test_resource_server1",
		})
	}

	wg.Add(2)
	go terraformOp()
	go terraformOp()

	wg.Wait()

	if !testServer1.ValidateResourceConfigCalled["test_resource_server1"] {
		t.Errorf("expected test_resource_server1 ValidateResourceConfig to be called on server1")
	}
}

func TestMuxServerGetResourceServer_GetProviderSchema_Duplicate(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	testServer1 := &tf6testserver.TestServer{
		GetProviderSchemaResponse: &tfprotov6.GetProviderSchemaResponse{
			ResourceSchemas: map[string]*tfprotov6.Schema{
				"test_resource_server": {}, // intentionally duplicated
			},
		},
	}
	testServer2 := &tf6testserver.TestServer{
		GetProviderSchemaResponse: &tfprotov6.GetProviderSchemaResponse{
			ResourceSchemas: map[string]*tfprotov6.Schema{
				"test_resource_server": {}, // intentionally duplicated
			},
		},
	}

	servers := []func() tfprotov6.ProviderServer{testServer1.ProviderServer, testServer2.ProviderServer}
	muxServer, err := tf6muxserver.NewMuxServer(ctx, servers...)

	if err != nil {
		t.Fatalf("unexpected error setting up factory: %s", err)
	}

	// When GetProviderSchemaOptional is enabled, the secondary provider
	// instances will receive non-GetProviderSchema RPCs such as
	// ValidateDataResourceConfig which will cause getResourceServer to perform
	// server discovery. This testing also simulates concurrent operations from
	// Terraform to verify the mutex does not deadlock.
	var wg sync.WaitGroup

	expectedDiags := []*tfprotov6.Diagnostic{
		{
			Severity: tfprotov6.DiagnosticSeverityError,
			Summary:  "Invalid Provider Server Combination",
			Detail: "The combined provider has multiple implementations of the same resource type across underlying providers. " +
				"Resource types must be implemented by only one underlying provider. " +
				"This is always an issue in the provider implementation and should be reported to the provider developers.\n\n" +
				"Duplicate resource type: test_resource_server",
		},
	}

	terraformOp := func() {
		defer wg.Done()

		resp, _ := muxServer.ProviderServer().ValidateResourceConfig(ctx, &tfprotov6.ValidateResourceConfigRequest{
			TypeName: "test_resource_server",
		})

		if diff := cmp.Diff(resp.Diagnostics, expectedDiags); diff != "" {
			t.Errorf("unexpected diagnostics difference: %s", diff)
		}
	}

	wg.Add(2)
	go terraformOp()
	go terraformOp()

	wg.Wait()

	if testServer1.ValidateResourceConfigCalled["test_resource_server"] {
		t.Errorf("unexpected test_resource_server ValidateDataResourceConfig called on server1")
	}

	if testServer2.ValidateResourceConfigCalled["test_resource_server"] {
		t.Errorf("unexpected test_resource_server ValidateDataResourceConfig called on server2")
	}
}

func TestMuxServerGetResourceServer_GetMetadata(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	testServer1 := &tf6testserver.TestServer{
		GetMetadataResponse: &tfprotov6.GetMetadataResponse{
			Resources: []tfprotov6.ResourceMetadata{
				{
					TypeName: "test_resource_server1",
				},
			},
		},
	}
	testServer2 := &tf6testserver.TestServer{
		GetMetadataResponse: &tfprotov6.GetMetadataResponse{
			Resources: []tfprotov6.ResourceMetadata{
				{
					TypeName: "test_resource_server2",
				},
			},
		},
	}

	servers := []func() tfprotov6.ProviderServer{testServer1.ProviderServer, testServer2.ProviderServer}
	muxServer, err := tf6muxserver.NewMuxServer(ctx, servers...)

	if err != nil {
		t.Fatalf("unexpected error setting up factory: %s", err)
	}

	// When GetProviderSchemaOptional is enabled, the secondary provider
	// instances will receive non-GetProviderSchema RPCs such as
	// ValidateResourceConfig which will cause getResourceServer to perform
	// server discovery. This testing also simulates concurrent operations from
	// Terraform to verify the mutex does not deadlock.
	var wg sync.WaitGroup

	terraformOp := func() {
		defer wg.Done()

		_, _ = muxServer.ProviderServer().ValidateResourceConfig(ctx, &tfprotov6.ValidateResourceConfigRequest{
			TypeName: "test_resource_server1",
		})
	}

	wg.Add(2)
	go terraformOp()
	go terraformOp()

	wg.Wait()

	if !testServer1.ValidateResourceConfigCalled["test_resource_server1"] {
		t.Errorf("expected test_resource_server1 ValidateResourceConfig to be called on server1")
	}
}

func TestMuxServerGetResourceServer_GetMetadata_Duplicate(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	testServer1 := &tf6testserver.TestServer{
		GetMetadataResponse: &tfprotov6.GetMetadataResponse{
			Resources: []tfprotov6.ResourceMetadata{
				{
					TypeName: "test_resource_server", // intentionally duplicated
				},
			},
		},
	}
	testServer2 := &tf6testserver.TestServer{
		GetMetadataResponse: &tfprotov6.GetMetadataResponse{
			Resources: []tfprotov6.ResourceMetadata{
				{
					TypeName: "test_resource_server", // intentionally duplicated
				},
			},
		},
	}

	servers := []func() tfprotov6.ProviderServer{testServer1.ProviderServer, testServer2.ProviderServer}
	muxServer, err := tf6muxserver.NewMuxServer(ctx, servers...)

	if err != nil {
		t.Fatalf("unexpected error setting up factory: %s", err)
	}

	// When GetProviderSchemaOptional is enabled, the secondary provider
	// instances will receive non-GetProviderSchema RPCs such as
	// ValidateDataResourceConfig which will cause getResourceServer to perform
	// server discovery. This testing also simulates concurrent operations from
	// Terraform to verify the mutex does not deadlock.
	var wg sync.WaitGroup

	expectedDiags := []*tfprotov6.Diagnostic{
		{
			Severity: tfprotov6.DiagnosticSeverityError,
			Summary:  "Invalid Provider Server Combination",
			Detail: "The combined provider has multiple implementations of the same resource type across underlying providers. " +
				"Resource types must be implemented by only one underlying provider. " +
				"This is always an issue in the provider implementation and should be reported to the provider developers.\n\n" +
				"Duplicate resource type: test_resource_server",
		},
	}

	terraformOp := func() {
		defer wg.Done()

		resp, _ := muxServer.ProviderServer().ValidateResourceConfig(ctx, &tfprotov6.ValidateResourceConfigRequest{
			TypeName: "test_resource_server",
		})

		if diff := cmp.Diff(resp.Diagnostics, expectedDiags); diff != "" {
			t.Errorf("unexpected diagnostics difference: %s", diff)
		}
	}

	wg.Add(2)
	go terraformOp()
	go terraformOp()

	wg.Wait()

	if testServer1.ValidateResourceConfigCalled["test_resource_server"] {
		t.Errorf("unexpected test_resource_server ValidateDataResourceConfig called on server1")
	}

	if testServer2.ValidateResourceConfigCalled["test_resource_server"] {
		t.Errorf("unexpected test_resource_server ValidateDataResourceConfig called on server2")
	}
}

func TestMuxServerGetResourceServer_GetMetadata_Partial(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	testServer1 := &tf6testserver.TestServer{
		GetMetadataResponse: &tfprotov6.GetMetadataResponse{
			Resources: []tfprotov6.ResourceMetadata{
				{
					TypeName: "test_resource_server1",
				},
			},
		},
	}
	testServer2 := &tf6testserver.TestServer{
		GetProviderSchemaResponse: &tfprotov6.GetProviderSchemaResponse{
			ResourceSchemas: map[string]*tfprotov6.Schema{
				"test_resource_server2": {},
			},
		},
	}

	servers := []func() tfprotov6.ProviderServer{testServer1.ProviderServer, testServer2.ProviderServer}
	muxServer, err := tf6muxserver.NewMuxServer(ctx, servers...)

	if err != nil {
		t.Fatalf("unexpected error setting up factory: %s", err)
	}

	// When GetProviderSchemaOptional is enabled, the secondary provider
	// instances will receive non-GetProviderSchema RPCs such as
	// ValidateResourceConfig which will cause getResourceServer to perform
	// server discovery. This testing also simulates concurrent operations from
	// Terraform to verify the mutex does not deadlock.
	var wg sync.WaitGroup

	terraformOp := func() {
		defer wg.Done()

		_, _ = muxServer.ProviderServer().ValidateResourceConfig(ctx, &tfprotov6.ValidateResourceConfigRequest{
			TypeName: "test_resource_server1",
		})
	}

	wg.Add(2)
	go terraformOp()
	go terraformOp()

	wg.Wait()

	if !testServer1.ValidateResourceConfigCalled["test_resource_server1"] {
		t.Errorf("expected test_resource_server1 ValidateResourceConfig to be called on server1")
	}
}

func TestMuxServerGetResourceServer_Missing(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	testServer1 := &tf6testserver.TestServer{
		GetMetadataResponse: &tfprotov6.GetMetadataResponse{
			Resources: []tfprotov6.ResourceMetadata{
				{
					TypeName: "test_resource_server1",
				},
			},
		},
	}
	testServer2 := &tf6testserver.TestServer{
		GetMetadataResponse: &tfprotov6.GetMetadataResponse{
			Resources: []tfprotov6.ResourceMetadata{
				{
					TypeName: "test_resource_server2",
				},
			},
		},
	}

	servers := []func() tfprotov6.ProviderServer{testServer1.ProviderServer, testServer2.ProviderServer}
	muxServer, err := tf6muxserver.NewMuxServer(ctx, servers...)

	if err != nil {
		t.Fatalf("unexpected error setting up factory: %s", err)
	}

	// When GetProviderSchemaOptional is enabled, the secondary provider
	// instances will receive non-GetProviderSchema RPCs such as
	// ValidateResourceConfig which will cause getResourceServer to perform
	// server discovery. This testing also simulates concurrent operations from
	// Terraform to verify the mutex does not deadlock.
	var wg sync.WaitGroup

	expectedDiags := []*tfprotov6.Diagnostic{
		{
			Severity: tfprotov6.DiagnosticSeverityError,
			Summary:  "Resource Not Implemented",
			Detail: "The combined provider does not implement the requested resource type. " +
				"This is always an issue in the provider implementation and should be reported to the provider developers.\n\n" +
				"Missing resource type: test_resource_nonexistent",
		},
	}

	terraformOp := func() {
		defer wg.Done()

		resp, _ := muxServer.ProviderServer().ValidateResourceConfig(ctx, &tfprotov6.ValidateResourceConfigRequest{
			TypeName: "test_resource_nonexistent",
		})

		if diff := cmp.Diff(resp.Diagnostics, expectedDiags); diff != "" {
			t.Errorf("unexpected diagnostics difference: %s", diff)
		}
	}

	wg.Add(2)
	go terraformOp()
	go terraformOp()

	wg.Wait()

	if testServer1.ValidateResourceConfigCalled["test_resource_nonexistent"] {
		t.Errorf("unexpected test_resource_nonexistent ValidateDataResourceConfig called on server1")
	}

	if testServer2.ValidateResourceConfigCalled["test_resource_nonexistent"] {
		t.Errorf("unexpected test_resource_nonexistent ValidateDataResourceConfig called on server2")
	}
}

func TestNewMuxServer(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		servers       []func() tfprotov6.ProviderServer
		expectedError error
	}{
		"duplicate-data-source": {
			servers: []func() tfprotov6.ProviderServer{
				(&tf6testserver.TestServer{
					GetProviderSchemaResponse: &tfprotov6.GetProviderSchemaResponse{
						DataSourceSchemas: map[string]*tfprotov6.Schema{
							"test_foo": {},
						},
					},
				}).ProviderServer,
				(&tf6testserver.TestServer{
					GetProviderSchemaResponse: &tfprotov6.GetProviderSchemaResponse{
						DataSourceSchemas: map[string]*tfprotov6.Schema{
							"test_foo": {},
						},
					},
				}).ProviderServer,
			},
			expectedError: nil, // deferred to GetProviderSchema
		},
		"duplicate-resource": {
			servers: []func() tfprotov6.ProviderServer{
				(&tf6testserver.TestServer{
					GetProviderSchemaResponse: &tfprotov6.GetProviderSchemaResponse{
						ResourceSchemas: map[string]*tfprotov6.Schema{
							"test_foo": {},
						},
					},
				}).ProviderServer,
				(&tf6testserver.TestServer{
					GetProviderSchemaResponse: &tfprotov6.GetProviderSchemaResponse{
						ResourceSchemas: map[string]*tfprotov6.Schema{
							"test_foo": {},
						},
					},
				}).ProviderServer,
			},
			expectedError: nil, // deferred to GetProviderSchema
		},
		"provider-mismatch": {
			servers: []func() tfprotov6.ProviderServer{
				(&tf6testserver.TestServer{
					GetProviderSchemaResponse: &tfprotov6.GetProviderSchemaResponse{
						Provider: &tfprotov6.Schema{
							Version: 1,
							Block: &tfprotov6.SchemaBlock{
								Version: 1,
								Attributes: []*tfprotov6.SchemaAttribute{
									{
										Name:            "account_id",
										Type:            tftypes.String,
										Required:        true,
										Description:     "the account ID to make requests for",
										DescriptionKind: tfprotov6.StringKindPlain,
									},
								},
								BlockTypes: []*tfprotov6.SchemaNestedBlock{
									{
										TypeName: "feature",
										Nesting:  tfprotov6.SchemaNestedBlockNestingModeList,
										Block: &tfprotov6.SchemaBlock{
											Version:         1,
											Description:     "features to enable on the provider",
											DescriptionKind: tfprotov6.StringKindPlain,
											Attributes: []*tfprotov6.SchemaAttribute{
												{
													Name:            "feature_id",
													Type:            tftypes.Number,
													Required:        true,
													Description:     "The ID of the feature",
													DescriptionKind: tfprotov6.StringKindPlain,
												},
												{
													Name:            "enabled",
													Type:            tftypes.Bool,
													Required:        true,
													Description:     "whether the feature is enabled",
													DescriptionKind: tfprotov6.StringKindPlain,
												},
											},
										},
									},
								},
							},
						},
					},
				}).ProviderServer,
				(&tf6testserver.TestServer{
					GetProviderSchemaResponse: &tfprotov6.GetProviderSchemaResponse{
						Provider: &tfprotov6.Schema{
							Version: 1,
							Block: &tfprotov6.SchemaBlock{
								Version: 1,
								Attributes: []*tfprotov6.SchemaAttribute{
									{
										Name:            "account_id",
										Type:            tftypes.String,
										Required:        true,
										Description:     "the account ID to make requests for",
										DescriptionKind: tfprotov6.StringKindPlain,
									},
								},
								BlockTypes: []*tfprotov6.SchemaNestedBlock{
									{
										TypeName: "feature",
										Nesting:  tfprotov6.SchemaNestedBlockNestingModeList,
										Block: &tfprotov6.SchemaBlock{
											Version:         1,
											Description:     "features to enable on the provider",
											DescriptionKind: tfprotov6.StringKindPlain,
											Attributes: []*tfprotov6.SchemaAttribute{
												{
													Name:            "feature_id",
													Type:            tftypes.Number,
													Required:        true,
													Description:     "The ID of the feature",
													DescriptionKind: tfprotov6.StringKindPlain,
												},
												{
													Name:            "enabled",
													Type:            tftypes.Bool,
													Optional:        true,
													Description:     "whether the feature is enabled",
													DescriptionKind: tfprotov6.StringKindPlain,
												},
											},
										},
									},
								},
							},
						},
					},
				}).ProviderServer,
			},
			expectedError: nil, // deferred to GetProviderSchema
		},
		"provider-ordering": {
			servers: []func() tfprotov6.ProviderServer{
				(&tf6testserver.TestServer{
					GetProviderSchemaResponse: &tfprotov6.GetProviderSchemaResponse{
						Provider: &tfprotov6.Schema{
							Version: 1,
							Block: &tfprotov6.SchemaBlock{
								Version: 1,
								Attributes: []*tfprotov6.SchemaAttribute{
									{
										Name:            "account_id",
										Type:            tftypes.String,
										Required:        true,
										Description:     "the account ID to make requests for",
										DescriptionKind: tfprotov6.StringKindPlain,
									},
									{
										Name:            "secret",
										Type:            tftypes.String,
										Required:        true,
										Description:     "the secret to authenticate with",
										DescriptionKind: tfprotov6.StringKindPlain,
									},
								},
								BlockTypes: []*tfprotov6.SchemaNestedBlock{
									{
										TypeName: "other_feature",
										Nesting:  tfprotov6.SchemaNestedBlockNestingModeList,
										Block: &tfprotov6.SchemaBlock{
											Version:         1,
											Description:     "features to enable on the provider",
											DescriptionKind: tfprotov6.StringKindPlain,
											Attributes: []*tfprotov6.SchemaAttribute{
												{
													Name:            "enabled",
													Type:            tftypes.Bool,
													Required:        true,
													Description:     "whether the feature is enabled",
													DescriptionKind: tfprotov6.StringKindPlain,
												},
												{
													Name:            "feature_id",
													Type:            tftypes.Number,
													Required:        true,
													Description:     "The ID of the feature",
													DescriptionKind: tfprotov6.StringKindPlain,
												},
											},
										},
									},
									{
										TypeName: "feature",
										Nesting:  tfprotov6.SchemaNestedBlockNestingModeList,
										Block: &tfprotov6.SchemaBlock{
											Version:         1,
											Description:     "features to enable on the provider",
											DescriptionKind: tfprotov6.StringKindPlain,
											Attributes: []*tfprotov6.SchemaAttribute{
												{
													Name:            "feature_id",
													Type:            tftypes.Number,
													Required:        true,
													Description:     "The ID of the feature",
													DescriptionKind: tfprotov6.StringKindPlain,
												},
												{
													Name:            "enabled",
													Type:            tftypes.Bool,
													Required:        true,
													Description:     "whether the feature is enabled",
													DescriptionKind: tfprotov6.StringKindPlain,
												},
											},
										},
									},
								},
							},
						},
					},
				}).ProviderServer,
				(&tf6testserver.TestServer{
					GetProviderSchemaResponse: &tfprotov6.GetProviderSchemaResponse{
						Provider: &tfprotov6.Schema{
							Version: 1,
							Block: &tfprotov6.SchemaBlock{
								Version: 1,
								Attributes: []*tfprotov6.SchemaAttribute{
									{
										Name:            "secret",
										Type:            tftypes.String,
										Required:        true,
										Description:     "the secret to authenticate with",
										DescriptionKind: tfprotov6.StringKindPlain,
									},
									{
										Name:            "account_id",
										Type:            tftypes.String,
										Required:        true,
										Description:     "the account ID to make requests for",
										DescriptionKind: tfprotov6.StringKindPlain,
									},
								},
								BlockTypes: []*tfprotov6.SchemaNestedBlock{
									{
										TypeName: "feature",
										Nesting:  tfprotov6.SchemaNestedBlockNestingModeList,
										Block: &tfprotov6.SchemaBlock{
											Version:         1,
											Description:     "features to enable on the provider",
											DescriptionKind: tfprotov6.StringKindPlain,
											Attributes: []*tfprotov6.SchemaAttribute{
												{
													Name:            "enabled",
													Type:            tftypes.Bool,
													Required:        true,
													Description:     "whether the feature is enabled",
													DescriptionKind: tfprotov6.StringKindPlain,
												},
												{
													Name:            "feature_id",
													Type:            tftypes.Number,
													Required:        true,
													Description:     "The ID of the feature",
													DescriptionKind: tfprotov6.StringKindPlain,
												},
											},
										},
									},
									{
										TypeName: "other_feature",
										Nesting:  tfprotov6.SchemaNestedBlockNestingModeList,
										Block: &tfprotov6.SchemaBlock{
											Version:         1,
											Description:     "features to enable on the provider",
											DescriptionKind: tfprotov6.StringKindPlain,
											Attributes: []*tfprotov6.SchemaAttribute{
												{
													Name:            "enabled",
													Type:            tftypes.Bool,
													Required:        true,
													Description:     "whether the feature is enabled",
													DescriptionKind: tfprotov6.StringKindPlain,
												},
												{
													Name:            "feature_id",
													Type:            tftypes.Number,
													Required:        true,
													Description:     "The ID of the feature",
													DescriptionKind: tfprotov6.StringKindPlain,
												},
											},
										},
									},
								},
							},
						},
					},
				}).ProviderServer,
			},
		},
		"provider-meta-mismatch": {
			servers: []func() tfprotov6.ProviderServer{
				(&tf6testserver.TestServer{
					GetProviderSchemaResponse: &tfprotov6.GetProviderSchemaResponse{
						ProviderMeta: &tfprotov6.Schema{
							Version: 1,
							Block: &tfprotov6.SchemaBlock{
								Version: 1,
								Attributes: []*tfprotov6.SchemaAttribute{
									{
										Name:            "account_id",
										Type:            tftypes.String,
										Required:        true,
										Description:     "the account ID to make requests for",
										DescriptionKind: tfprotov6.StringKindPlain,
									},
								},
								BlockTypes: []*tfprotov6.SchemaNestedBlock{
									{
										TypeName: "feature",
										Nesting:  tfprotov6.SchemaNestedBlockNestingModeList,
										Block: &tfprotov6.SchemaBlock{
											Version:         1,
											Description:     "features to enable on the provider",
											DescriptionKind: tfprotov6.StringKindPlain,
											Attributes: []*tfprotov6.SchemaAttribute{
												{
													Name:            "feature_id",
													Type:            tftypes.Number,
													Required:        true,
													Description:     "The ID of the feature",
													DescriptionKind: tfprotov6.StringKindPlain,
												},
												{
													Name:            "enabled",
													Type:            tftypes.Bool,
													Required:        true,
													Description:     "whether the feature is enabled",
													DescriptionKind: tfprotov6.StringKindPlain,
												},
											},
										},
									},
								},
							},
						},
					},
				}).ProviderServer,
				(&tf6testserver.TestServer{
					GetProviderSchemaResponse: &tfprotov6.GetProviderSchemaResponse{
						ProviderMeta: &tfprotov6.Schema{
							Version: 1,
							Block: &tfprotov6.SchemaBlock{
								Version: 1,
								Attributes: []*tfprotov6.SchemaAttribute{
									{
										Name:            "account_id",
										Type:            tftypes.String,
										Required:        true,
										Description:     "the account ID to make requests for",
										DescriptionKind: tfprotov6.StringKindPlain,
									},
								},
								BlockTypes: []*tfprotov6.SchemaNestedBlock{
									{
										TypeName: "feature",
										Nesting:  tfprotov6.SchemaNestedBlockNestingModeList,
										Block: &tfprotov6.SchemaBlock{
											Version:         1,
											Description:     "features to enable on the provider",
											DescriptionKind: tfprotov6.StringKindPlain,
											Attributes: []*tfprotov6.SchemaAttribute{
												{
													Name:            "feature_id",
													Type:            tftypes.Number,
													Required:        true,
													Description:     "The ID of the feature",
													DescriptionKind: tfprotov6.StringKindPlain,
												},
												{
													Name:            "enabled",
													Type:            tftypes.Bool,
													Optional:        true,
													Description:     "whether the feature is enabled",
													DescriptionKind: tfprotov6.StringKindPlain,
												},
											},
										},
									},
								},
							},
						},
					},
				}).ProviderServer,
			},
			expectedError: nil, // deferred to GetProviderSchema
		},
		"provider-meta-ordering": {
			servers: []func() tfprotov6.ProviderServer{
				(&tf6testserver.TestServer{
					GetProviderSchemaResponse: &tfprotov6.GetProviderSchemaResponse{
						ProviderMeta: &tfprotov6.Schema{
							Version: 1,
							Block: &tfprotov6.SchemaBlock{
								Version: 1,
								Attributes: []*tfprotov6.SchemaAttribute{
									{
										Name:            "account_id",
										Type:            tftypes.String,
										Required:        true,
										Description:     "the account ID to make requests for",
										DescriptionKind: tfprotov6.StringKindPlain,
									},
									{
										Name:            "secret",
										Type:            tftypes.String,
										Required:        true,
										Description:     "the secret to authenticate with",
										DescriptionKind: tfprotov6.StringKindPlain,
									},
								},
								BlockTypes: []*tfprotov6.SchemaNestedBlock{
									{
										TypeName: "feature",
										Nesting:  tfprotov6.SchemaNestedBlockNestingModeList,
										Block: &tfprotov6.SchemaBlock{
											Version:         1,
											Description:     "features to enable on the provider",
											DescriptionKind: tfprotov6.StringKindPlain,
											Attributes: []*tfprotov6.SchemaAttribute{
												{
													Name:            "feature_id",
													Type:            tftypes.Number,
													Required:        true,
													Description:     "The ID of the feature",
													DescriptionKind: tfprotov6.StringKindPlain,
												},
												{
													Name:            "enabled",
													Type:            tftypes.Bool,
													Required:        true,
													Description:     "whether the feature is enabled",
													DescriptionKind: tfprotov6.StringKindPlain,
												},
											},
										},
									},
									{
										TypeName: "other_feature",
										Nesting:  tfprotov6.SchemaNestedBlockNestingModeList,
										Block: &tfprotov6.SchemaBlock{
											Version:         1,
											Description:     "features to enable on the provider",
											DescriptionKind: tfprotov6.StringKindPlain,
											Attributes: []*tfprotov6.SchemaAttribute{
												{
													Name:            "enabled",
													Type:            tftypes.Bool,
													Required:        true,
													Description:     "whether the feature is enabled",
													DescriptionKind: tfprotov6.StringKindPlain,
												},
												{
													Name:            "feature_id",
													Type:            tftypes.Number,
													Required:        true,
													Description:     "The ID of the feature",
													DescriptionKind: tfprotov6.StringKindPlain,
												},
											},
										},
									},
								},
							},
						},
					},
				}).ProviderServer,
				(&tf6testserver.TestServer{
					GetProviderSchemaResponse: &tfprotov6.GetProviderSchemaResponse{
						ProviderMeta: &tfprotov6.Schema{
							Version: 1,
							Block: &tfprotov6.SchemaBlock{
								Version: 1,
								Attributes: []*tfprotov6.SchemaAttribute{
									{
										Name:            "secret",
										Type:            tftypes.String,
										Required:        true,
										Description:     "the secret to authenticate with",
										DescriptionKind: tfprotov6.StringKindPlain,
									},
									{
										Name:            "account_id",
										Type:            tftypes.String,
										Required:        true,
										Description:     "the account ID to make requests for",
										DescriptionKind: tfprotov6.StringKindPlain,
									},
								},
								BlockTypes: []*tfprotov6.SchemaNestedBlock{
									{
										TypeName: "other_feature",
										Nesting:  tfprotov6.SchemaNestedBlockNestingModeList,
										Block: &tfprotov6.SchemaBlock{
											Version:         1,
											Description:     "features to enable on the provider",
											DescriptionKind: tfprotov6.StringKindPlain,
											Attributes: []*tfprotov6.SchemaAttribute{
												{
													Name:            "feature_id",
													Type:            tftypes.Number,
													Required:        true,
													Description:     "The ID of the feature",
													DescriptionKind: tfprotov6.StringKindPlain,
												},
												{
													Name:            "enabled",
													Type:            tftypes.Bool,
													Required:        true,
													Description:     "whether the feature is enabled",
													DescriptionKind: tfprotov6.StringKindPlain,
												},
											},
										},
									},
									{
										TypeName: "feature",
										Nesting:  tfprotov6.SchemaNestedBlockNestingModeList,
										Block: &tfprotov6.SchemaBlock{
											Version:         1,
											Description:     "features to enable on the provider",
											DescriptionKind: tfprotov6.StringKindPlain,
											Attributes: []*tfprotov6.SchemaAttribute{
												{
													Name:            "enabled",
													Type:            tftypes.Bool,
													Required:        true,
													Description:     "whether the feature is enabled",
													DescriptionKind: tfprotov6.StringKindPlain,
												},
												{
													Name:            "feature_id",
													Type:            tftypes.Number,
													Required:        true,
													Description:     "The ID of the feature",
													DescriptionKind: tfprotov6.StringKindPlain,
												},
											},
										},
									},
								},
							},
						},
					},
				}).ProviderServer,
			},
		},
		"nested block MinItems and MaxItems diff are ignored": {
			servers: []func() tfprotov6.ProviderServer{
				(&tf6testserver.TestServer{
					GetProviderSchemaResponse: &tfprotov6.GetProviderSchemaResponse{
						Provider: &tfprotov6.Schema{
							Version: 1,
							Block: &tfprotov6.SchemaBlock{
								Version: 1,
								BlockTypes: []*tfprotov6.SchemaNestedBlock{
									{
										TypeName: "feature",
										Nesting:  tfprotov6.SchemaNestedBlockNestingModeList,
										Block: &tfprotov6.SchemaBlock{
											Version:         1,
											Description:     "features to enable on the provider",
											DescriptionKind: tfprotov6.StringKindPlain,
											Attributes:      []*tfprotov6.SchemaAttribute{},
										},
										MinItems: 1,
										MaxItems: 2,
									},
								},
							},
						},
					},
				}).ProviderServer,
				(&tf6testserver.TestServer{
					GetProviderSchemaResponse: &tfprotov6.GetProviderSchemaResponse{
						Provider: &tfprotov6.Schema{
							Version: 1,
							Block: &tfprotov6.SchemaBlock{
								Version: 1,
								BlockTypes: []*tfprotov6.SchemaNestedBlock{
									{
										TypeName: "feature",
										Nesting:  tfprotov6.SchemaNestedBlockNestingModeList,
										Block: &tfprotov6.SchemaBlock{
											Version:         1,
											Description:     "features to enable on the provider",
											DescriptionKind: tfprotov6.StringKindPlain,
											Attributes:      []*tfprotov6.SchemaAttribute{},
										},
										MinItems: 0,
										MaxItems: 0,
									},
								},
							},
						},
					},
				}).ProviderServer,
			},
		},
		"server-capabilities": {
			servers: []func() tfprotov6.ProviderServer{
				(&tf6testserver.TestServer{
					GetProviderSchemaResponse: &tfprotov6.GetProviderSchemaResponse{
						ResourceSchemas: map[string]*tfprotov6.Schema{
							"test_with_server_capabilities": {},
						},
						ServerCapabilities: &tfprotov6.ServerCapabilities{
							PlanDestroy: true,
						},
					},
				}).ProviderServer,
				(&tf6testserver.TestServer{
					GetProviderSchemaResponse: &tfprotov6.GetProviderSchemaResponse{
						ResourceSchemas: map[string]*tfprotov6.Schema{
							"test_without_server_capabilities": {},
						},
					},
				}).ProviderServer,
			},
			expectedError: nil,
		},
	}

	for name, testCase := range testCases {

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			_, err := tf6muxserver.NewMuxServer(context.Background(), testCase.servers...)

			if err != nil {
				if testCase.expectedError == nil {
					t.Fatalf("unexpected error: %s", err)
				}

				if !strings.Contains(err.Error(), testCase.expectedError.Error()) {
					t.Fatalf("expected error %q, got: %s", testCase.expectedError, err)
				}
			}

			if err == nil && testCase.expectedError != nil {
				t.Fatalf("expected error: %s", testCase.expectedError)
			}
		})
	}
}
