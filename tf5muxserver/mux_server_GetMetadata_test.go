// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package tf5muxserver_test

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/terraform-plugin-go/tfprotov5"

	"github.com/hashicorp/terraform-plugin-mux/internal/tf5testserver"
	"github.com/hashicorp/terraform-plugin-mux/tf5muxserver"
)

func TestMuxServerGetMetadata(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		servers                    []func() tfprotov5.ProviderServer
		expectedDataSources        []tfprotov5.DataSourceMetadata
		expectedDiagnostics        []*tfprotov5.Diagnostic
		expectedFunctions          []tfprotov5.FunctionMetadata
		expectedResources          []tfprotov5.ResourceMetadata
		expectedServerCapabilities *tfprotov5.ServerCapabilities
	}{
		"combined": {
			servers: []func() tfprotov5.ProviderServer{
				(&tf5testserver.TestServer{
					GetMetadataResponse: &tfprotov5.GetMetadataResponse{
						Resources: []tfprotov5.ResourceMetadata{
							{
								TypeName: "test_foo",
							},
							{
								TypeName: "test_bar",
							},
						},
						DataSources: []tfprotov5.DataSourceMetadata{
							{
								TypeName: "test_foo",
							},
						},
						Functions: []tfprotov5.FunctionMetadata{
							{
								Name: "test_function1",
							},
						},
					},
				}).ProviderServer,
				(&tf5testserver.TestServer{
					GetMetadataResponse: &tfprotov5.GetMetadataResponse{
						Resources: []tfprotov5.ResourceMetadata{
							{
								TypeName: "test_quux",
							},
						},
						DataSources: []tfprotov5.DataSourceMetadata{
							{
								TypeName: "test_bar",
							},
							{
								TypeName: "test_quux",
							},
						},
						Functions: []tfprotov5.FunctionMetadata{
							{
								Name: "test_function2",
							},
							{
								Name: "test_function3",
							},
						},
					},
				}).ProviderServer,
			},
			expectedResources: []tfprotov5.ResourceMetadata{
				{
					TypeName: "test_foo",
				},
				{
					TypeName: "test_bar",
				},
				{
					TypeName: "test_quux",
				},
			},
			expectedDataSources: []tfprotov5.DataSourceMetadata{
				{
					TypeName: "test_foo",
				},
				{
					TypeName: "test_bar",
				},
				{
					TypeName: "test_quux",
				},
			},
			expectedFunctions: []tfprotov5.FunctionMetadata{
				{
					Name: "test_function1",
				},
				{
					Name: "test_function2",
				},
				{
					Name: "test_function3",
				},
			},
			expectedServerCapabilities: &tfprotov5.ServerCapabilities{
				GetProviderSchemaOptional: true,
				MoveResourceState:         true,
				PlanDestroy:               true,
			},
		},
		"duplicate-data-source-type": {
			servers: []func() tfprotov5.ProviderServer{
				(&tf5testserver.TestServer{
					GetMetadataResponse: &tfprotov5.GetMetadataResponse{
						DataSources: []tfprotov5.DataSourceMetadata{
							{
								TypeName: "test_foo",
							},
						},
					},
				}).ProviderServer,
				(&tf5testserver.TestServer{
					GetMetadataResponse: &tfprotov5.GetMetadataResponse{
						DataSources: []tfprotov5.DataSourceMetadata{
							{
								TypeName: "test_foo",
							},
						},
					},
				}).ProviderServer,
			},
			expectedDataSources: []tfprotov5.DataSourceMetadata{
				{
					TypeName: "test_foo",
				},
			},
			expectedDiagnostics: []*tfprotov5.Diagnostic{
				{
					Severity: tfprotov5.DiagnosticSeverityError,
					Summary:  "Invalid Provider Server Combination",
					Detail: "The combined provider has multiple implementations of the same data source type across underlying providers. " +
						"Data source types must be implemented by only one underlying provider. " +
						"This is always an issue in the provider implementation and should be reported to the provider developers.\n\n" +
						"Duplicate data source type: test_foo",
				},
			},
			expectedFunctions: []tfprotov5.FunctionMetadata{},
			expectedResources: []tfprotov5.ResourceMetadata{},
			expectedServerCapabilities: &tfprotov5.ServerCapabilities{
				GetProviderSchemaOptional: true,
				MoveResourceState:         true,
				PlanDestroy:               true,
			},
		},
		"duplicate-function": {
			servers: []func() tfprotov5.ProviderServer{
				(&tf5testserver.TestServer{
					GetMetadataResponse: &tfprotov5.GetMetadataResponse{
						Functions: []tfprotov5.FunctionMetadata{
							{
								Name: "test_function",
							},
						},
					},
				}).ProviderServer,
				(&tf5testserver.TestServer{
					GetMetadataResponse: &tfprotov5.GetMetadataResponse{
						Functions: []tfprotov5.FunctionMetadata{
							{
								Name: "test_function",
							},
						},
					},
				}).ProviderServer,
			},
			expectedDataSources: []tfprotov5.DataSourceMetadata{},
			expectedDiagnostics: []*tfprotov5.Diagnostic{
				{
					Severity: tfprotov5.DiagnosticSeverityError,
					Summary:  "Invalid Provider Server Combination",
					Detail: "The combined provider has multiple implementations of the same function name across underlying providers. " +
						"Functions must be implemented by only one underlying provider. " +
						"This is always an issue in the provider implementation and should be reported to the provider developers.\n\n" +
						"Duplicate function: test_function",
				},
			},
			expectedFunctions: []tfprotov5.FunctionMetadata{
				{
					Name: "test_function",
				},
			},
			expectedResources: []tfprotov5.ResourceMetadata{},
			expectedServerCapabilities: &tfprotov5.ServerCapabilities{
				GetProviderSchemaOptional: true,
				MoveResourceState:         true,
				PlanDestroy:               true,
			},
		},
		"duplicate-resource-type": {
			servers: []func() tfprotov5.ProviderServer{
				(&tf5testserver.TestServer{
					GetMetadataResponse: &tfprotov5.GetMetadataResponse{
						Resources: []tfprotov5.ResourceMetadata{
							{
								TypeName: "test_foo",
							},
						},
					},
				}).ProviderServer,
				(&tf5testserver.TestServer{
					GetMetadataResponse: &tfprotov5.GetMetadataResponse{
						Resources: []tfprotov5.ResourceMetadata{
							{
								TypeName: "test_foo",
							},
						},
					},
				}).ProviderServer,
			},
			expectedDataSources: []tfprotov5.DataSourceMetadata{},
			expectedDiagnostics: []*tfprotov5.Diagnostic{
				{
					Severity: tfprotov5.DiagnosticSeverityError,
					Summary:  "Invalid Provider Server Combination",
					Detail: "The combined provider has multiple implementations of the same resource type across underlying providers. " +
						"Resource types must be implemented by only one underlying provider. " +
						"This is always an issue in the provider implementation and should be reported to the provider developers.\n\n" +
						"Duplicate resource type: test_foo",
				},
			},
			expectedFunctions: []tfprotov5.FunctionMetadata{},
			expectedResources: []tfprotov5.ResourceMetadata{
				{
					TypeName: "test_foo",
				},
			},
			expectedServerCapabilities: &tfprotov5.ServerCapabilities{
				GetProviderSchemaOptional: true,
				MoveResourceState:         true,
				PlanDestroy:               true,
			},
		},
		"server-capabilities": {
			servers: []func() tfprotov5.ProviderServer{
				(&tf5testserver.TestServer{
					GetMetadataResponse: &tfprotov5.GetMetadataResponse{
						Resources: []tfprotov5.ResourceMetadata{
							{
								TypeName: "test_with_server_capabilities",
							},
						},
						ServerCapabilities: &tfprotov5.ServerCapabilities{
							GetProviderSchemaOptional: true,
							MoveResourceState:         true,
							PlanDestroy:               true,
						},
					},
				}).ProviderServer,
				(&tf5testserver.TestServer{
					GetMetadataResponse: &tfprotov5.GetMetadataResponse{
						Resources: []tfprotov5.ResourceMetadata{
							{
								TypeName: "test_without_server_capabilities",
							},
						},
					},
				}).ProviderServer,
			},
			expectedDataSources: []tfprotov5.DataSourceMetadata{},
			expectedFunctions:   []tfprotov5.FunctionMetadata{},
			expectedResources: []tfprotov5.ResourceMetadata{
				{
					TypeName: "test_with_server_capabilities",
				},
				{
					TypeName: "test_without_server_capabilities",
				},
			},
			expectedServerCapabilities: &tfprotov5.ServerCapabilities{
				GetProviderSchemaOptional: true,
				MoveResourceState:         true,
				PlanDestroy:               true,
			},
		},
		"error-once": {
			servers: []func() tfprotov5.ProviderServer{
				(&tf5testserver.TestServer{
					GetMetadataResponse: &tfprotov5.GetMetadataResponse{
						Diagnostics: []*tfprotov5.Diagnostic{
							{
								Severity: tfprotov5.DiagnosticSeverityError,
								Summary:  "test error summary",
								Detail:   "test error details",
							},
						},
					},
				}).ProviderServer,
				(&tf5testserver.TestServer{}).ProviderServer,
				(&tf5testserver.TestServer{}).ProviderServer,
			},
			expectedDataSources: []tfprotov5.DataSourceMetadata{},
			expectedDiagnostics: []*tfprotov5.Diagnostic{
				{
					Severity: tfprotov5.DiagnosticSeverityError,
					Summary:  "test error summary",
					Detail:   "test error details",
				},
			},
			expectedFunctions: []tfprotov5.FunctionMetadata{},
			expectedResources: []tfprotov5.ResourceMetadata{},
			expectedServerCapabilities: &tfprotov5.ServerCapabilities{
				GetProviderSchemaOptional: true,
				MoveResourceState:         true,
				PlanDestroy:               true,
			},
		},
		"error-multiple": {
			servers: []func() tfprotov5.ProviderServer{
				(&tf5testserver.TestServer{
					GetMetadataResponse: &tfprotov5.GetMetadataResponse{
						Diagnostics: []*tfprotov5.Diagnostic{
							{
								Severity: tfprotov5.DiagnosticSeverityError,
								Summary:  "test error summary",
								Detail:   "test error details",
							},
						},
					},
				}).ProviderServer,
				(&tf5testserver.TestServer{}).ProviderServer,
				(&tf5testserver.TestServer{
					GetMetadataResponse: &tfprotov5.GetMetadataResponse{
						Diagnostics: []*tfprotov5.Diagnostic{
							{
								Severity: tfprotov5.DiagnosticSeverityError,
								Summary:  "test error summary",
								Detail:   "test error details",
							},
						},
					},
				}).ProviderServer,
			},
			expectedDataSources: []tfprotov5.DataSourceMetadata{},
			expectedDiagnostics: []*tfprotov5.Diagnostic{
				{
					Severity: tfprotov5.DiagnosticSeverityError,
					Summary:  "test error summary",
					Detail:   "test error details",
				},
				{
					Severity: tfprotov5.DiagnosticSeverityError,
					Summary:  "test error summary",
					Detail:   "test error details",
				},
			},
			expectedFunctions: []tfprotov5.FunctionMetadata{},
			expectedResources: []tfprotov5.ResourceMetadata{},
			expectedServerCapabilities: &tfprotov5.ServerCapabilities{
				GetProviderSchemaOptional: true,
				MoveResourceState:         true,
				PlanDestroy:               true,
			},
		},
		"warning-once": {
			servers: []func() tfprotov5.ProviderServer{
				(&tf5testserver.TestServer{
					GetMetadataResponse: &tfprotov5.GetMetadataResponse{
						Diagnostics: []*tfprotov5.Diagnostic{
							{
								Severity: tfprotov5.DiagnosticSeverityWarning,
								Summary:  "test warning summary",
								Detail:   "test warning details",
							},
						},
					},
				}).ProviderServer,
				(&tf5testserver.TestServer{}).ProviderServer,
				(&tf5testserver.TestServer{}).ProviderServer,
			},
			expectedDataSources: []tfprotov5.DataSourceMetadata{},
			expectedDiagnostics: []*tfprotov5.Diagnostic{
				{
					Severity: tfprotov5.DiagnosticSeverityWarning,
					Summary:  "test warning summary",
					Detail:   "test warning details",
				},
			},
			expectedFunctions: []tfprotov5.FunctionMetadata{},
			expectedResources: []tfprotov5.ResourceMetadata{},
			expectedServerCapabilities: &tfprotov5.ServerCapabilities{
				GetProviderSchemaOptional: true,
				MoveResourceState:         true,
				PlanDestroy:               true,
			},
		},
		"warning-multiple": {
			servers: []func() tfprotov5.ProviderServer{
				(&tf5testserver.TestServer{
					GetMetadataResponse: &tfprotov5.GetMetadataResponse{
						Diagnostics: []*tfprotov5.Diagnostic{
							{
								Severity: tfprotov5.DiagnosticSeverityWarning,
								Summary:  "test warning summary",
								Detail:   "test warning details",
							},
						},
					},
				}).ProviderServer,
				(&tf5testserver.TestServer{}).ProviderServer,
				(&tf5testserver.TestServer{
					GetMetadataResponse: &tfprotov5.GetMetadataResponse{
						Diagnostics: []*tfprotov5.Diagnostic{
							{
								Severity: tfprotov5.DiagnosticSeverityWarning,
								Summary:  "test warning summary",
								Detail:   "test warning details",
							},
						},
					},
				}).ProviderServer,
			},
			expectedDataSources: []tfprotov5.DataSourceMetadata{},
			expectedDiagnostics: []*tfprotov5.Diagnostic{
				{
					Severity: tfprotov5.DiagnosticSeverityWarning,
					Summary:  "test warning summary",
					Detail:   "test warning details",
				},
				{
					Severity: tfprotov5.DiagnosticSeverityWarning,
					Summary:  "test warning summary",
					Detail:   "test warning details",
				},
			},
			expectedFunctions: []tfprotov5.FunctionMetadata{},
			expectedResources: []tfprotov5.ResourceMetadata{},
			expectedServerCapabilities: &tfprotov5.ServerCapabilities{
				GetProviderSchemaOptional: true,
				MoveResourceState:         true,
				PlanDestroy:               true,
			},
		},
		"warning-then-error": {
			servers: []func() tfprotov5.ProviderServer{
				(&tf5testserver.TestServer{
					GetMetadataResponse: &tfprotov5.GetMetadataResponse{
						Diagnostics: []*tfprotov5.Diagnostic{
							{
								Severity: tfprotov5.DiagnosticSeverityWarning,
								Summary:  "test warning summary",
								Detail:   "test warning details",
							},
						},
					},
				}).ProviderServer,
				(&tf5testserver.TestServer{}).ProviderServer,
				(&tf5testserver.TestServer{
					GetMetadataResponse: &tfprotov5.GetMetadataResponse{
						Diagnostics: []*tfprotov5.Diagnostic{
							{
								Severity: tfprotov5.DiagnosticSeverityError,
								Summary:  "test error summary",
								Detail:   "test error details",
							},
						},
					},
				}).ProviderServer,
			},
			expectedDataSources: []tfprotov5.DataSourceMetadata{},
			expectedDiagnostics: []*tfprotov5.Diagnostic{
				{
					Severity: tfprotov5.DiagnosticSeverityWarning,
					Summary:  "test warning summary",
					Detail:   "test warning details",
				},
				{
					Severity: tfprotov5.DiagnosticSeverityError,
					Summary:  "test error summary",
					Detail:   "test error details",
				},
			},
			expectedFunctions: []tfprotov5.FunctionMetadata{},
			expectedResources: []tfprotov5.ResourceMetadata{},
			expectedServerCapabilities: &tfprotov5.ServerCapabilities{
				GetProviderSchemaOptional: true,
				MoveResourceState:         true,
				PlanDestroy:               true,
			},
		},
	}

	for name, testCase := range testCases {
		name, testCase := name, testCase

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			muxServer, err := tf5muxserver.NewMuxServer(context.Background(), testCase.servers...)

			if err != nil {
				t.Fatalf("unexpected error: %s", err)
			}

			resp, err := muxServer.ProviderServer().GetMetadata(context.Background(), &tfprotov5.GetMetadataRequest{})

			if err != nil {
				t.Fatalf("unexpected error: %s", err)
			}

			if diff := cmp.Diff(resp.DataSources, testCase.expectedDataSources); diff != "" {
				t.Errorf("data sources didn't match expectations: %s", diff)
			}

			if diff := cmp.Diff(resp.Diagnostics, testCase.expectedDiagnostics); diff != "" {
				t.Errorf("diagnostics didn't match expectations: %s", diff)
			}

			if diff := cmp.Diff(resp.Functions, testCase.expectedFunctions); diff != "" {
				t.Errorf("functions didn't match expectations: %s", diff)
			}

			if diff := cmp.Diff(resp.Resources, testCase.expectedResources); diff != "" {
				t.Errorf("resources didn't match expectations: %s", diff)
			}

			if diff := cmp.Diff(resp.ServerCapabilities, testCase.expectedServerCapabilities); diff != "" {
				t.Errorf("server capabilities didn't match expectations: %s", diff)
			}
		})
	}
}
