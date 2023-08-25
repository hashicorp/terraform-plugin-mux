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
		expectedDataSources        []tfprotov6.DataSourceMetadata
		expectedDiagnostics        []*tfprotov6.Diagnostic
		expectedResources          []tfprotov6.ResourceMetadata
		expectedServerCapabilities *tfprotov6.ServerCapabilities
	}{
		"combined": {
			servers: []func() tfprotov6.ProviderServer{
				(&tf6testserver.TestServer{
					GetMetadataResponse: &tfprotov6.GetMetadataResponse{
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
					},
				}).ProviderServer,
				(&tf6testserver.TestServer{
					GetMetadataResponse: &tfprotov6.GetMetadataResponse{
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
					},
				}).ProviderServer,
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
			expectedServerCapabilities: &tfprotov6.ServerCapabilities{
				GetProviderSchemaOptional: true,
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
			expectedResources: []tfprotov6.ResourceMetadata{},
			expectedServerCapabilities: &tfprotov6.ServerCapabilities{
				GetProviderSchemaOptional: true,
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
			expectedResources: []tfprotov6.ResourceMetadata{
				{
					TypeName: "test_foo",
				},
			},
			expectedServerCapabilities: &tfprotov6.ServerCapabilities{
				GetProviderSchemaOptional: true,
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
			expectedDataSources: []tfprotov6.DataSourceMetadata{},
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
			expectedDataSources: []tfprotov6.DataSourceMetadata{},
			expectedDiagnostics: []*tfprotov6.Diagnostic{
				{
					Severity: tfprotov6.DiagnosticSeverityError,
					Summary:  "test error summary",
					Detail:   "test error details",
				},
			},
			expectedResources: []tfprotov6.ResourceMetadata{},
			expectedServerCapabilities: &tfprotov6.ServerCapabilities{
				GetProviderSchemaOptional: true,
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
			expectedResources: []tfprotov6.ResourceMetadata{},
			expectedServerCapabilities: &tfprotov6.ServerCapabilities{
				GetProviderSchemaOptional: true,
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
			expectedDataSources: []tfprotov6.DataSourceMetadata{},
			expectedDiagnostics: []*tfprotov6.Diagnostic{
				{
					Severity: tfprotov6.DiagnosticSeverityWarning,
					Summary:  "test warning summary",
					Detail:   "test warning details",
				},
			},
			expectedResources: []tfprotov6.ResourceMetadata{},
			expectedServerCapabilities: &tfprotov6.ServerCapabilities{
				GetProviderSchemaOptional: true,
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
			expectedResources: []tfprotov6.ResourceMetadata{},
			expectedServerCapabilities: &tfprotov6.ServerCapabilities{
				GetProviderSchemaOptional: true,
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
			expectedResources: []tfprotov6.ResourceMetadata{},
			expectedServerCapabilities: &tfprotov6.ServerCapabilities{
				GetProviderSchemaOptional: true,
				PlanDestroy:               true,
			},
		},
	}

	for name, testCase := range testCases {
		name, testCase := name, testCase

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

			if diff := cmp.Diff(resp.DataSources, testCase.expectedDataSources); diff != "" {
				t.Errorf("data sources didn't match expectations: %s", diff)
			}

			if diff := cmp.Diff(resp.Diagnostics, testCase.expectedDiagnostics); diff != "" {
				t.Errorf("diagnostics didn't match expectations: %s", diff)
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
