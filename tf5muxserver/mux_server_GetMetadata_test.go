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
		expectedActions            []tfprotov5.ActionMetadata
		expectedDataSources        []tfprotov5.DataSourceMetadata
		expectedDiagnostics        []*tfprotov5.Diagnostic
		expectedEphemeralResources []tfprotov5.EphemeralResourceMetadata
		expectedListResources      []tfprotov5.ListResourceMetadata
		expectedFunctions          []tfprotov5.FunctionMetadata
		expectedResources          []tfprotov5.ResourceMetadata
		expectedServerCapabilities *tfprotov5.ServerCapabilities
	}{
		"combined": {
			servers: []func() tfprotov5.ProviderServer{
				(&tf5testserver.TestServer{
					GetMetadataResponse: &tfprotov5.GetMetadataResponse{
						Actions: []tfprotov5.ActionMetadata{
							{
								TypeName: "test_foo",
							},
							{
								TypeName: "test_bar",
							},
						},
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
						EphemeralResources: []tfprotov5.EphemeralResourceMetadata{
							{
								TypeName: "test_foo",
							},
							{
								TypeName: "test_bar",
							},
						},
						ListResources: []tfprotov5.ListResourceMetadata{
							{
								TypeName: "test_foo",
							},
							{
								TypeName: "test_bar",
							},
						},
					},
				}).ProviderServer,
				(&tf5testserver.TestServer{
					GetMetadataResponse: &tfprotov5.GetMetadataResponse{
						Actions: []tfprotov5.ActionMetadata{
							{
								TypeName: "test_quux",
							},
						},
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
						EphemeralResources: []tfprotov5.EphemeralResourceMetadata{
							{
								TypeName: "test_quux",
							},
						},
						ListResources: []tfprotov5.ListResourceMetadata{
							{
								TypeName: "test_quux",
							},
						},
					},
				}).ProviderServer,
			},
			expectedActions: []tfprotov5.ActionMetadata{
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
			expectedEphemeralResources: []tfprotov5.EphemeralResourceMetadata{
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
			expectedListResources: []tfprotov5.ListResourceMetadata{
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
			expectedServerCapabilities: &tfprotov5.ServerCapabilities{
				GetProviderSchemaOptional: true,
				MoveResourceState:         true,
				PlanDestroy:               true,
			},
		},
		"duplicate-action": {
			servers: []func() tfprotov5.ProviderServer{
				(&tf5testserver.TestServer{
					GetMetadataResponse: &tfprotov5.GetMetadataResponse{
						Actions: []tfprotov5.ActionMetadata{
							{
								TypeName: "test_foo",
							},
						},
					},
				}).ProviderServer,
				(&tf5testserver.TestServer{
					GetMetadataResponse: &tfprotov5.GetMetadataResponse{
						Actions: []tfprotov5.ActionMetadata{
							{
								TypeName: "test_foo",
							},
						},
					},
				}).ProviderServer,
			},
			expectedActions: []tfprotov5.ActionMetadata{
				{
					TypeName: "test_foo",
				},
			},
			expectedDataSources: []tfprotov5.DataSourceMetadata{},
			expectedDiagnostics: []*tfprotov5.Diagnostic{
				{
					Severity: tfprotov5.DiagnosticSeverityError,
					Summary:  "Invalid Provider Server Combination",
					Detail: "The combined provider has multiple implementations of the same action type across underlying providers. " +
						"Actions must be implemented by only one underlying provider. " +
						"This is always an issue in the provider implementation and should be reported to the provider developers.\n\n" +
						"Duplicate action: test_foo",
				},
			},
			expectedEphemeralResources: []tfprotov5.EphemeralResourceMetadata{},
			expectedListResources:      []tfprotov5.ListResourceMetadata{},
			expectedFunctions:          []tfprotov5.FunctionMetadata{},
			expectedResources:          []tfprotov5.ResourceMetadata{},
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
			expectedActions: []tfprotov5.ActionMetadata{},
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
			expectedEphemeralResources: []tfprotov5.EphemeralResourceMetadata{},
			expectedListResources:      []tfprotov5.ListResourceMetadata{},
			expectedFunctions:          []tfprotov5.FunctionMetadata{},
			expectedResources:          []tfprotov5.ResourceMetadata{},
			expectedServerCapabilities: &tfprotov5.ServerCapabilities{
				GetProviderSchemaOptional: true,
				MoveResourceState:         true,
				PlanDestroy:               true,
			},
		},
		"duplicate-ephemeral-resource-type": {
			servers: []func() tfprotov5.ProviderServer{
				(&tf5testserver.TestServer{
					GetMetadataResponse: &tfprotov5.GetMetadataResponse{
						EphemeralResources: []tfprotov5.EphemeralResourceMetadata{
							{
								TypeName: "test_foo",
							},
						},
					},
				}).ProviderServer,
				(&tf5testserver.TestServer{
					GetMetadataResponse: &tfprotov5.GetMetadataResponse{
						EphemeralResources: []tfprotov5.EphemeralResourceMetadata{
							{
								TypeName: "test_foo",
							},
						},
					},
				}).ProviderServer,
			},
			expectedActions:     []tfprotov5.ActionMetadata{},
			expectedDataSources: []tfprotov5.DataSourceMetadata{},
			expectedDiagnostics: []*tfprotov5.Diagnostic{
				{
					Severity: tfprotov5.DiagnosticSeverityError,
					Summary:  "Invalid Provider Server Combination",
					Detail: "The combined provider has multiple implementations of the same ephemeral resource type across underlying providers. " +
						"Ephemeral resource types must be implemented by only one underlying provider. " +
						"This is always an issue in the provider implementation and should be reported to the provider developers.\n\n" +
						"Duplicate ephemeral resource type: test_foo",
				},
			},
			expectedEphemeralResources: []tfprotov5.EphemeralResourceMetadata{
				{
					TypeName: "test_foo",
				},
			},
			expectedListResources: []tfprotov5.ListResourceMetadata{},
			expectedFunctions:     []tfprotov5.FunctionMetadata{},
			expectedResources:     []tfprotov5.ResourceMetadata{},
			expectedServerCapabilities: &tfprotov5.ServerCapabilities{
				GetProviderSchemaOptional: true,
				MoveResourceState:         true,
				PlanDestroy:               true,
			},
		},
		"duplicate-list-resource-type": {
			servers: []func() tfprotov5.ProviderServer{
				(&tf5testserver.TestServer{
					GetMetadataResponse: &tfprotov5.GetMetadataResponse{
						ListResources: []tfprotov5.ListResourceMetadata{
							{
								TypeName: "test_foo",
							},
						},
					},
				}).ProviderServer,
				(&tf5testserver.TestServer{
					GetMetadataResponse: &tfprotov5.GetMetadataResponse{
						ListResources: []tfprotov5.ListResourceMetadata{
							{
								TypeName: "test_foo",
							},
						},
					},
				}).ProviderServer,
			},
			expectedActions:     []tfprotov5.ActionMetadata{},
			expectedDataSources: []tfprotov5.DataSourceMetadata{},
			expectedDiagnostics: []*tfprotov5.Diagnostic{
				{
					Severity: tfprotov5.DiagnosticSeverityError,
					Summary:  "Invalid Provider Server Combination",
					Detail: "The combined provider has multiple implementations of the same list resource type across underlying providers. " +
						"List resource types must be implemented by only one underlying provider. " +
						"This is always an issue in the provider implementation and should be reported to the provider developers.\n\n" +
						"Duplicate list resource type: test_foo",
				},
			},
			expectedEphemeralResources: []tfprotov5.EphemeralResourceMetadata{},
			expectedListResources: []tfprotov5.ListResourceMetadata{
				{
					TypeName: "test_foo",
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
			expectedActions:     []tfprotov5.ActionMetadata{},
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
			expectedEphemeralResources: []tfprotov5.EphemeralResourceMetadata{},
			expectedListResources:      []tfprotov5.ListResourceMetadata{},
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
			expectedActions:     []tfprotov5.ActionMetadata{},
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
			expectedEphemeralResources: []tfprotov5.EphemeralResourceMetadata{},
			expectedListResources:      []tfprotov5.ListResourceMetadata{},
			expectedFunctions:          []tfprotov5.FunctionMetadata{},
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
			expectedActions:            []tfprotov5.ActionMetadata{},
			expectedDataSources:        []tfprotov5.DataSourceMetadata{},
			expectedEphemeralResources: []tfprotov5.EphemeralResourceMetadata{},
			expectedListResources:      []tfprotov5.ListResourceMetadata{},
			expectedFunctions:          []tfprotov5.FunctionMetadata{},
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
			expectedActions:     []tfprotov5.ActionMetadata{},
			expectedDataSources: []tfprotov5.DataSourceMetadata{},
			expectedDiagnostics: []*tfprotov5.Diagnostic{
				{
					Severity: tfprotov5.DiagnosticSeverityError,
					Summary:  "test error summary",
					Detail:   "test error details",
				},
			},
			expectedEphemeralResources: []tfprotov5.EphemeralResourceMetadata{},
			expectedListResources:      []tfprotov5.ListResourceMetadata{},
			expectedFunctions:          []tfprotov5.FunctionMetadata{},
			expectedResources:          []tfprotov5.ResourceMetadata{},
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
			expectedActions:     []tfprotov5.ActionMetadata{},
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
			expectedEphemeralResources: []tfprotov5.EphemeralResourceMetadata{},
			expectedListResources:      []tfprotov5.ListResourceMetadata{},
			expectedFunctions:          []tfprotov5.FunctionMetadata{},
			expectedResources:          []tfprotov5.ResourceMetadata{},
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
			expectedActions:     []tfprotov5.ActionMetadata{},
			expectedDataSources: []tfprotov5.DataSourceMetadata{},
			expectedDiagnostics: []*tfprotov5.Diagnostic{
				{
					Severity: tfprotov5.DiagnosticSeverityWarning,
					Summary:  "test warning summary",
					Detail:   "test warning details",
				},
			},
			expectedEphemeralResources: []tfprotov5.EphemeralResourceMetadata{},
			expectedListResources:      []tfprotov5.ListResourceMetadata{},
			expectedFunctions:          []tfprotov5.FunctionMetadata{},
			expectedResources:          []tfprotov5.ResourceMetadata{},
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
			expectedActions:     []tfprotov5.ActionMetadata{},
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
			expectedEphemeralResources: []tfprotov5.EphemeralResourceMetadata{},
			expectedListResources:      []tfprotov5.ListResourceMetadata{},
			expectedFunctions:          []tfprotov5.FunctionMetadata{},
			expectedResources:          []tfprotov5.ResourceMetadata{},
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
			expectedActions:     []tfprotov5.ActionMetadata{},
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
			expectedEphemeralResources: []tfprotov5.EphemeralResourceMetadata{},
			expectedListResources:      []tfprotov5.ListResourceMetadata{},
			expectedFunctions:          []tfprotov5.FunctionMetadata{},
			expectedResources:          []tfprotov5.ResourceMetadata{},
			expectedServerCapabilities: &tfprotov5.ServerCapabilities{
				GetProviderSchemaOptional: true,
				MoveResourceState:         true,
				PlanDestroy:               true,
			},
		},
	}

	for name, testCase := range testCases {

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

			if diff := cmp.Diff(resp.Actions, testCase.expectedActions); diff != "" {
				t.Errorf("actions didn't match expectations: %s", diff)
			}

			if diff := cmp.Diff(resp.DataSources, testCase.expectedDataSources); diff != "" {
				t.Errorf("data sources didn't match expectations: %s", diff)
			}

			if diff := cmp.Diff(resp.Diagnostics, testCase.expectedDiagnostics); diff != "" {
				t.Errorf("diagnostics didn't match expectations: %s", diff)
			}

			if diff := cmp.Diff(resp.EphemeralResources, testCase.expectedEphemeralResources); diff != "" {
				t.Errorf("ephemeral resources didn't match expectations: %s", diff)
			}

			if diff := cmp.Diff(resp.ListResources, testCase.expectedListResources); diff != "" {
				t.Errorf("list resources didn't match expectations: %s", diff)
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
