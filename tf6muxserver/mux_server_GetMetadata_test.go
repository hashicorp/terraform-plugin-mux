// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package tf6muxserver_test

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"

	"github.com/hashicorp/terraform-plugin-mux/internal/tf6testserver"
	"github.com/hashicorp/terraform-plugin-mux/tf6muxserver"
)

func TestMuxServerGetMetadata(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		servers                    []func() tfprotov6.ProviderServer
		expectedActions            []tfprotov6.ActionMetadata
		expectedDataSources        []tfprotov6.DataSourceMetadata
		expectedDiagnostics        []*tfprotov6.Diagnostic
		expectedEphemeralResources []tfprotov6.EphemeralResourceMetadata
		expectedListResources      []tfprotov6.ListResourceMetadata
		expectedFunctions          []tfprotov6.FunctionMetadata
		expectedResources          []tfprotov6.ResourceMetadata
		expectedServerCapabilities *tfprotov6.ServerCapabilities
	}{
		"combined": {
			servers: []func() tfprotov6.ProviderServer{
				(&tf6testserver.TestServer{
					GetMetadataResponse: &tfprotov6.GetMetadataResponse{
						Actions: []tfprotov6.ActionMetadata{
							{
								TypeName: "test_foo",
							},
							{
								TypeName: "test_bar",
							},
						},
						Resources: []tfprotov6.ResourceMetadata{
							{
								TypeName: "test_foo",
							},
							{
								TypeName: "test_bar",
							},
						},
						DataSources: []tfprotov6.DataSourceMetadata{
							{
								TypeName: "test_foo",
							},
						},
						Functions: []tfprotov6.FunctionMetadata{
							{
								Name: "test_function1",
							},
						},
						EphemeralResources: []tfprotov6.EphemeralResourceMetadata{
							{
								TypeName: "test_foo",
							},
							{
								TypeName: "test_bar",
							},
						},
						ListResources: []tfprotov6.ListResourceMetadata{
							{
								TypeName: "test_foo",
							},
							{
								TypeName: "test_bar",
							},
						},
					},
				}).ProviderServer,
				(&tf6testserver.TestServer{
					GetMetadataResponse: &tfprotov6.GetMetadataResponse{
						Actions: []tfprotov6.ActionMetadata{
							{
								TypeName: "test_quux",
							},
						},
						Resources: []tfprotov6.ResourceMetadata{
							{
								TypeName: "test_quux",
							},
						},
						DataSources: []tfprotov6.DataSourceMetadata{
							{
								TypeName: "test_bar",
							},
							{
								TypeName: "test_quux",
							},
						},
						Functions: []tfprotov6.FunctionMetadata{
							{
								Name: "test_function2",
							},
							{
								Name: "test_function3",
							},
						},
						EphemeralResources: []tfprotov6.EphemeralResourceMetadata{
							{
								TypeName: "test_quux",
							},
						},
						ListResources: []tfprotov6.ListResourceMetadata{
							{
								TypeName: "test_quux",
							},
						},
					},
				}).ProviderServer,
			},
			expectedActions: []tfprotov6.ActionMetadata{
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
			expectedResources: []tfprotov6.ResourceMetadata{
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
			expectedDataSources: []tfprotov6.DataSourceMetadata{
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
			expectedFunctions: []tfprotov6.FunctionMetadata{
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
			expectedEphemeralResources: []tfprotov6.EphemeralResourceMetadata{
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
			expectedListResources: []tfprotov6.ListResourceMetadata{
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
			expectedServerCapabilities: &tfprotov6.ServerCapabilities{
				GetProviderSchemaOptional: true,
				MoveResourceState:         true,
				PlanDestroy:               true,
			},
		},
		"duplicate-action": {
			servers: []func() tfprotov6.ProviderServer{
				(&tf6testserver.TestServer{
					GetMetadataResponse: &tfprotov6.GetMetadataResponse{
						Actions: []tfprotov6.ActionMetadata{
							{
								TypeName: "test_foo",
							},
						},
					},
				}).ProviderServer,
				(&tf6testserver.TestServer{
					GetMetadataResponse: &tfprotov6.GetMetadataResponse{
						Actions: []tfprotov6.ActionMetadata{
							{
								TypeName: "test_foo",
							},
						},
					},
				}).ProviderServer,
			},
			expectedActions: []tfprotov6.ActionMetadata{
				{
					TypeName: "test_foo",
				},
			},
			expectedDataSources: []tfprotov6.DataSourceMetadata{},
			expectedDiagnostics: []*tfprotov6.Diagnostic{
				{
					Severity: tfprotov6.DiagnosticSeverityError,
					Summary:  "Invalid Provider Server Combination",
					Detail: "The combined provider has multiple implementations of the same action type across underlying providers. " +
						"Actions must be implemented by only one underlying provider. " +
						"This is always an issue in the provider implementation and should be reported to the provider developers.\n\n" +
						"Duplicate action: test_foo",
				},
			},
			expectedEphemeralResources: []tfprotov6.EphemeralResourceMetadata{},
			expectedListResources:      []tfprotov6.ListResourceMetadata{},
			expectedFunctions:          []tfprotov6.FunctionMetadata{},
			expectedResources:          []tfprotov6.ResourceMetadata{},
			expectedServerCapabilities: &tfprotov6.ServerCapabilities{
				GetProviderSchemaOptional: true,
				MoveResourceState:         true,
				PlanDestroy:               true,
			},
		},
		"duplicate-data-source-type": {
			servers: []func() tfprotov6.ProviderServer{
				(&tf6testserver.TestServer{
					GetMetadataResponse: &tfprotov6.GetMetadataResponse{
						DataSources: []tfprotov6.DataSourceMetadata{
							{
								TypeName: "test_foo",
							},
						},
					},
				}).ProviderServer,
				(&tf6testserver.TestServer{
					GetMetadataResponse: &tfprotov6.GetMetadataResponse{
						DataSources: []tfprotov6.DataSourceMetadata{
							{
								TypeName: "test_foo",
							},
						},
					},
				}).ProviderServer,
			},
			expectedActions: []tfprotov6.ActionMetadata{},
			expectedDataSources: []tfprotov6.DataSourceMetadata{
				{
					TypeName: "test_foo",
				},
			},
			expectedDiagnostics: []*tfprotov6.Diagnostic{
				{
					Severity: tfprotov6.DiagnosticSeverityError,
					Summary:  "Invalid Provider Server Combination",
					Detail: "The combined provider has multiple implementations of the same data source type across underlying providers. " +
						"Data source types must be implemented by only one underlying provider. " +
						"This is always an issue in the provider implementation and should be reported to the provider developers.\n\n" +
						"Duplicate data source type: test_foo",
				},
			},
			expectedEphemeralResources: []tfprotov6.EphemeralResourceMetadata{},
			expectedListResources:      []tfprotov6.ListResourceMetadata{},
			expectedFunctions:          []tfprotov6.FunctionMetadata{},
			expectedResources:          []tfprotov6.ResourceMetadata{},
			expectedServerCapabilities: &tfprotov6.ServerCapabilities{
				GetProviderSchemaOptional: true,
				MoveResourceState:         true,
				PlanDestroy:               true,
			},
		},
		"duplicate-ephemeral-resource-type": {
			servers: []func() tfprotov6.ProviderServer{
				(&tf6testserver.TestServer{
					GetMetadataResponse: &tfprotov6.GetMetadataResponse{
						EphemeralResources: []tfprotov6.EphemeralResourceMetadata{
							{
								TypeName: "test_foo",
							},
						},
					},
				}).ProviderServer,
				(&tf6testserver.TestServer{
					GetMetadataResponse: &tfprotov6.GetMetadataResponse{
						EphemeralResources: []tfprotov6.EphemeralResourceMetadata{
							{
								TypeName: "test_foo",
							},
						},
					},
				}).ProviderServer,
			},
			expectedActions:     []tfprotov6.ActionMetadata{},
			expectedDataSources: []tfprotov6.DataSourceMetadata{},
			expectedDiagnostics: []*tfprotov6.Diagnostic{
				{
					Severity: tfprotov6.DiagnosticSeverityError,
					Summary:  "Invalid Provider Server Combination",
					Detail: "The combined provider has multiple implementations of the same ephemeral resource type across underlying providers. " +
						"Ephemeral resource types must be implemented by only one underlying provider. " +
						"This is always an issue in the provider implementation and should be reported to the provider developers.\n\n" +
						"Duplicate ephemeral resource type: test_foo",
				},
			},
			expectedEphemeralResources: []tfprotov6.EphemeralResourceMetadata{
				{
					TypeName: "test_foo",
				},
			},
			expectedListResources: []tfprotov6.ListResourceMetadata{},
			expectedFunctions:     []tfprotov6.FunctionMetadata{},
			expectedResources:     []tfprotov6.ResourceMetadata{},
			expectedServerCapabilities: &tfprotov6.ServerCapabilities{
				GetProviderSchemaOptional: true,
				MoveResourceState:         true,
				PlanDestroy:               true,
			},
		},
		"duplicate-list-resource-type": {
			servers: []func() tfprotov6.ProviderServer{
				(&tf6testserver.TestServer{
					GetMetadataResponse: &tfprotov6.GetMetadataResponse{
						ListResources: []tfprotov6.ListResourceMetadata{
							{
								TypeName: "test_foo",
							},
						},
					},
				}).ProviderServer,
				(&tf6testserver.TestServer{
					GetMetadataResponse: &tfprotov6.GetMetadataResponse{
						ListResources: []tfprotov6.ListResourceMetadata{
							{
								TypeName: "test_foo",
							},
						},
					},
				}).ProviderServer,
			},
			expectedActions:     []tfprotov6.ActionMetadata{},
			expectedDataSources: []tfprotov6.DataSourceMetadata{},
			expectedDiagnostics: []*tfprotov6.Diagnostic{
				{
					Severity: tfprotov6.DiagnosticSeverityError,
					Summary:  "Invalid Provider Server Combination",
					Detail: "The combined provider has multiple implementations of the same list resource type across underlying providers. " +
						"List resource types must be implemented by only one underlying provider. " +
						"This is always an issue in the provider implementation and should be reported to the provider developers.\n\n" +
						"Duplicate list resource type: test_foo",
				},
			},
			expectedListResources: []tfprotov6.ListResourceMetadata{
				{
					TypeName: "test_foo",
				},
			},
			expectedEphemeralResources: []tfprotov6.EphemeralResourceMetadata{},
			expectedFunctions:          []tfprotov6.FunctionMetadata{},
			expectedResources:          []tfprotov6.ResourceMetadata{},
			expectedServerCapabilities: &tfprotov6.ServerCapabilities{
				GetProviderSchemaOptional: true,
				MoveResourceState:         true,
				PlanDestroy:               true,
			},
		},
		"duplicate-function": {
			servers: []func() tfprotov6.ProviderServer{
				(&tf6testserver.TestServer{
					GetMetadataResponse: &tfprotov6.GetMetadataResponse{
						Functions: []tfprotov6.FunctionMetadata{
							{
								Name: "test_function",
							},
						},
					},
				}).ProviderServer,
				(&tf6testserver.TestServer{
					GetMetadataResponse: &tfprotov6.GetMetadataResponse{
						Functions: []tfprotov6.FunctionMetadata{
							{
								Name: "test_function",
							},
						},
					},
				}).ProviderServer,
			},
			expectedActions:     []tfprotov6.ActionMetadata{},
			expectedDataSources: []tfprotov6.DataSourceMetadata{},
			expectedDiagnostics: []*tfprotov6.Diagnostic{
				{
					Severity: tfprotov6.DiagnosticSeverityError,
					Summary:  "Invalid Provider Server Combination",
					Detail: "The combined provider has multiple implementations of the same function name across underlying providers. " +
						"Functions must be implemented by only one underlying provider. " +
						"This is always an issue in the provider implementation and should be reported to the provider developers.\n\n" +
						"Duplicate function: test_function",
				},
			},
			expectedEphemeralResources: []tfprotov6.EphemeralResourceMetadata{},
			expectedListResources:      []tfprotov6.ListResourceMetadata{},
			expectedFunctions: []tfprotov6.FunctionMetadata{
				{
					Name: "test_function",
				},
			},
			expectedResources: []tfprotov6.ResourceMetadata{},
			expectedServerCapabilities: &tfprotov6.ServerCapabilities{
				GetProviderSchemaOptional: true,
				MoveResourceState:         true,
				PlanDestroy:               true,
			},
		},
		"duplicate-resource-type": {
			servers: []func() tfprotov6.ProviderServer{
				(&tf6testserver.TestServer{
					GetMetadataResponse: &tfprotov6.GetMetadataResponse{
						Resources: []tfprotov6.ResourceMetadata{
							{
								TypeName: "test_foo",
							},
						},
					},
				}).ProviderServer,
				(&tf6testserver.TestServer{
					GetMetadataResponse: &tfprotov6.GetMetadataResponse{
						Resources: []tfprotov6.ResourceMetadata{
							{
								TypeName: "test_foo",
							},
						},
					},
				}).ProviderServer,
			},
			expectedActions:     []tfprotov6.ActionMetadata{},
			expectedDataSources: []tfprotov6.DataSourceMetadata{},
			expectedDiagnostics: []*tfprotov6.Diagnostic{
				{
					Severity: tfprotov6.DiagnosticSeverityError,
					Summary:  "Invalid Provider Server Combination",
					Detail: "The combined provider has multiple implementations of the same resource type across underlying providers. " +
						"Resource types must be implemented by only one underlying provider. " +
						"This is always an issue in the provider implementation and should be reported to the provider developers.\n\n" +
						"Duplicate resource type: test_foo",
				},
			},
			expectedEphemeralResources: []tfprotov6.EphemeralResourceMetadata{},
			expectedListResources:      []tfprotov6.ListResourceMetadata{},
			expectedFunctions:          []tfprotov6.FunctionMetadata{},
			expectedResources: []tfprotov6.ResourceMetadata{
				{
					TypeName: "test_foo",
				},
			},
			expectedServerCapabilities: &tfprotov6.ServerCapabilities{
				GetProviderSchemaOptional: true,
				MoveResourceState:         true,
				PlanDestroy:               true,
			},
		},
		"server-capabilities": {
			servers: []func() tfprotov6.ProviderServer{
				(&tf6testserver.TestServer{
					GetMetadataResponse: &tfprotov6.GetMetadataResponse{
						Resources: []tfprotov6.ResourceMetadata{
							{
								TypeName: "test_with_server_capabilities",
							},
						},
						ServerCapabilities: &tfprotov6.ServerCapabilities{
							GetProviderSchemaOptional: true,
							MoveResourceState:         true,
							PlanDestroy:               true,
						},
					},
				}).ProviderServer,
				(&tf6testserver.TestServer{
					GetMetadataResponse: &tfprotov6.GetMetadataResponse{
						Resources: []tfprotov6.ResourceMetadata{
							{
								TypeName: "test_without_server_capabilities",
							},
						},
					},
				}).ProviderServer,
			},
			expectedActions:            []tfprotov6.ActionMetadata{},
			expectedDataSources:        []tfprotov6.DataSourceMetadata{},
			expectedEphemeralResources: []tfprotov6.EphemeralResourceMetadata{},
			expectedListResources:      []tfprotov6.ListResourceMetadata{},
			expectedFunctions:          []tfprotov6.FunctionMetadata{},
			expectedResources: []tfprotov6.ResourceMetadata{
				{
					TypeName: "test_with_server_capabilities",
				},
				{
					TypeName: "test_without_server_capabilities",
				},
			},
			expectedServerCapabilities: &tfprotov6.ServerCapabilities{
				GetProviderSchemaOptional: true,
				MoveResourceState:         true,
				PlanDestroy:               true,
			},
		},
		"error-once": {
			servers: []func() tfprotov6.ProviderServer{
				(&tf6testserver.TestServer{
					GetMetadataResponse: &tfprotov6.GetMetadataResponse{
						Diagnostics: []*tfprotov6.Diagnostic{
							{
								Severity: tfprotov6.DiagnosticSeverityError,
								Summary:  "test error summary",
								Detail:   "test error details",
							},
						},
					},
				}).ProviderServer,
				(&tf6testserver.TestServer{}).ProviderServer,
				(&tf6testserver.TestServer{}).ProviderServer,
			},
			expectedActions:     []tfprotov6.ActionMetadata{},
			expectedDataSources: []tfprotov6.DataSourceMetadata{},
			expectedDiagnostics: []*tfprotov6.Diagnostic{
				{
					Severity: tfprotov6.DiagnosticSeverityError,
					Summary:  "test error summary",
					Detail:   "test error details",
				},
			},
			expectedEphemeralResources: []tfprotov6.EphemeralResourceMetadata{},
			expectedListResources:      []tfprotov6.ListResourceMetadata{},
			expectedFunctions:          []tfprotov6.FunctionMetadata{},
			expectedResources:          []tfprotov6.ResourceMetadata{},
			expectedServerCapabilities: &tfprotov6.ServerCapabilities{
				GetProviderSchemaOptional: true,
				MoveResourceState:         true,
				PlanDestroy:               true,
			},
		},
		"error-multiple": {
			servers: []func() tfprotov6.ProviderServer{
				(&tf6testserver.TestServer{
					GetMetadataResponse: &tfprotov6.GetMetadataResponse{
						Diagnostics: []*tfprotov6.Diagnostic{
							{
								Severity: tfprotov6.DiagnosticSeverityError,
								Summary:  "test error summary",
								Detail:   "test error details",
							},
						},
					},
				}).ProviderServer,
				(&tf6testserver.TestServer{}).ProviderServer,
				(&tf6testserver.TestServer{
					GetMetadataResponse: &tfprotov6.GetMetadataResponse{
						Diagnostics: []*tfprotov6.Diagnostic{
							{
								Severity: tfprotov6.DiagnosticSeverityError,
								Summary:  "test error summary",
								Detail:   "test error details",
							},
						},
					},
				}).ProviderServer,
			},
			expectedActions:     []tfprotov6.ActionMetadata{},
			expectedDataSources: []tfprotov6.DataSourceMetadata{},
			expectedDiagnostics: []*tfprotov6.Diagnostic{
				{
					Severity: tfprotov6.DiagnosticSeverityError,
					Summary:  "test error summary",
					Detail:   "test error details",
				},
				{
					Severity: tfprotov6.DiagnosticSeverityError,
					Summary:  "test error summary",
					Detail:   "test error details",
				},
			},
			expectedEphemeralResources: []tfprotov6.EphemeralResourceMetadata{},
			expectedListResources:      []tfprotov6.ListResourceMetadata{},
			expectedFunctions:          []tfprotov6.FunctionMetadata{},
			expectedResources:          []tfprotov6.ResourceMetadata{},
			expectedServerCapabilities: &tfprotov6.ServerCapabilities{
				GetProviderSchemaOptional: true,
				MoveResourceState:         true,
				PlanDestroy:               true,
			},
		},
		"warning-once": {
			servers: []func() tfprotov6.ProviderServer{
				(&tf6testserver.TestServer{
					GetMetadataResponse: &tfprotov6.GetMetadataResponse{
						Diagnostics: []*tfprotov6.Diagnostic{
							{
								Severity: tfprotov6.DiagnosticSeverityWarning,
								Summary:  "test warning summary",
								Detail:   "test warning details",
							},
						},
					},
				}).ProviderServer,
				(&tf6testserver.TestServer{}).ProviderServer,
				(&tf6testserver.TestServer{}).ProviderServer,
			},
			expectedActions:     []tfprotov6.ActionMetadata{},
			expectedDataSources: []tfprotov6.DataSourceMetadata{},
			expectedDiagnostics: []*tfprotov6.Diagnostic{
				{
					Severity: tfprotov6.DiagnosticSeverityWarning,
					Summary:  "test warning summary",
					Detail:   "test warning details",
				},
			},
			expectedEphemeralResources: []tfprotov6.EphemeralResourceMetadata{},
			expectedListResources:      []tfprotov6.ListResourceMetadata{},
			expectedFunctions:          []tfprotov6.FunctionMetadata{},
			expectedResources:          []tfprotov6.ResourceMetadata{},
			expectedServerCapabilities: &tfprotov6.ServerCapabilities{
				GetProviderSchemaOptional: true,
				MoveResourceState:         true,
				PlanDestroy:               true,
			},
		},
		"warning-multiple": {
			servers: []func() tfprotov6.ProviderServer{
				(&tf6testserver.TestServer{
					GetMetadataResponse: &tfprotov6.GetMetadataResponse{
						Diagnostics: []*tfprotov6.Diagnostic{
							{
								Severity: tfprotov6.DiagnosticSeverityWarning,
								Summary:  "test warning summary",
								Detail:   "test warning details",
							},
						},
					},
				}).ProviderServer,
				(&tf6testserver.TestServer{}).ProviderServer,
				(&tf6testserver.TestServer{
					GetMetadataResponse: &tfprotov6.GetMetadataResponse{
						Diagnostics: []*tfprotov6.Diagnostic{
							{
								Severity: tfprotov6.DiagnosticSeverityWarning,
								Summary:  "test warning summary",
								Detail:   "test warning details",
							},
						},
					},
				}).ProviderServer,
			},
			expectedActions:     []tfprotov6.ActionMetadata{},
			expectedDataSources: []tfprotov6.DataSourceMetadata{},
			expectedDiagnostics: []*tfprotov6.Diagnostic{
				{
					Severity: tfprotov6.DiagnosticSeverityWarning,
					Summary:  "test warning summary",
					Detail:   "test warning details",
				},
				{
					Severity: tfprotov6.DiagnosticSeverityWarning,
					Summary:  "test warning summary",
					Detail:   "test warning details",
				},
			},
			expectedEphemeralResources: []tfprotov6.EphemeralResourceMetadata{},
			expectedListResources:      []tfprotov6.ListResourceMetadata{},
			expectedFunctions:          []tfprotov6.FunctionMetadata{},
			expectedResources:          []tfprotov6.ResourceMetadata{},
			expectedServerCapabilities: &tfprotov6.ServerCapabilities{
				GetProviderSchemaOptional: true,
				MoveResourceState:         true,
				PlanDestroy:               true,
			},
		},
		"warning-then-error": {
			servers: []func() tfprotov6.ProviderServer{
				(&tf6testserver.TestServer{
					GetMetadataResponse: &tfprotov6.GetMetadataResponse{
						Diagnostics: []*tfprotov6.Diagnostic{
							{
								Severity: tfprotov6.DiagnosticSeverityWarning,
								Summary:  "test warning summary",
								Detail:   "test warning details",
							},
						},
					},
				}).ProviderServer,
				(&tf6testserver.TestServer{}).ProviderServer,
				(&tf6testserver.TestServer{
					GetMetadataResponse: &tfprotov6.GetMetadataResponse{
						Diagnostics: []*tfprotov6.Diagnostic{
							{
								Severity: tfprotov6.DiagnosticSeverityError,
								Summary:  "test error summary",
								Detail:   "test error details",
							},
						},
					},
				}).ProviderServer,
			},
			expectedActions:     []tfprotov6.ActionMetadata{},
			expectedDataSources: []tfprotov6.DataSourceMetadata{},
			expectedDiagnostics: []*tfprotov6.Diagnostic{
				{
					Severity: tfprotov6.DiagnosticSeverityWarning,
					Summary:  "test warning summary",
					Detail:   "test warning details",
				},
				{
					Severity: tfprotov6.DiagnosticSeverityError,
					Summary:  "test error summary",
					Detail:   "test error details",
				},
			},
			expectedEphemeralResources: []tfprotov6.EphemeralResourceMetadata{},
			expectedListResources:      []tfprotov6.ListResourceMetadata{},
			expectedFunctions:          []tfprotov6.FunctionMetadata{},
			expectedResources:          []tfprotov6.ResourceMetadata{},
			expectedServerCapabilities: &tfprotov6.ServerCapabilities{
				GetProviderSchemaOptional: true,
				MoveResourceState:         true,
				PlanDestroy:               true,
			},
		},
	}

	for name, testCase := range testCases {

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			muxServer, err := tf6muxserver.NewMuxServer(context.Background(), testCase.servers...)

			if err != nil {
				t.Fatalf("unexpected error: %s", err)
			}

			resp, err := muxServer.ProviderServer().GetMetadata(context.Background(), &tfprotov6.GetMetadataRequest{})

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

			if diff := cmp.Diff(resp.Functions, testCase.expectedFunctions); diff != "" {
				t.Errorf("functions didn't match expectations: %s", diff)
			}

			if diff := cmp.Diff(resp.EphemeralResources, testCase.expectedEphemeralResources); diff != "" {
				t.Errorf("ephemeral resources didn't match expectations: %s", diff)
			}

			if diff := cmp.Diff(resp.ListResources, testCase.expectedListResources); diff != "" {
				t.Errorf("list resources didn't match expectations: %s", diff)
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
