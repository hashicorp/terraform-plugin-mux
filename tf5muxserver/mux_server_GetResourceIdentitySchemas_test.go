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

func TestMuxServerGetResourceIdentitySchema(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		servers                 []func() tfprotov5.ProviderServer
		expectedIdentitySchemas map[string]*tfprotov5.ResourceIdentitySchema
		expectedDiagnostics     []*tfprotov5.Diagnostic
	}{
		"combined": {
			servers: []func() tfprotov5.ProviderServer{
				(&tf5testserver.TestServer{
					GetResourceIdentitySchemasResponse: &tfprotov5.GetResourceIdentitySchemasResponse{
						IdentitySchemas: map[string]*tfprotov5.ResourceIdentitySchema{
							"test_resource_identity_foo": {
								Version: 1,
								IdentityAttributes: []*tfprotov5.ResourceIdentitySchemaAttribute{
									{
										Name:              "req",
										RequiredForImport: true,
										Description:       "this one's required",
									},
									{
										Name:              "opt",
										OptionalForImport: true,
										Description:       "this one's optional",
									},
								},
							},
							"test_resource_identity_bar": {
								Version: 1,
								IdentityAttributes: []*tfprotov5.ResourceIdentitySchemaAttribute{
									{
										Name:              "req",
										RequiredForImport: true,
										Description:       "this one's required",
									},
									{
										Name:              "opt",
										OptionalForImport: true,
										Description:       "this one's optional",
									},
								},
							},
						},
					},
				}).ProviderServer,
				(&tf5testserver.TestServer{
					GetResourceIdentitySchemasResponse: &tfprotov5.GetResourceIdentitySchemasResponse{
						IdentitySchemas: map[string]*tfprotov5.ResourceIdentitySchema{
							"test_resource_identity_foobar": {
								Version: 1,
								IdentityAttributes: []*tfprotov5.ResourceIdentitySchemaAttribute{
									{
										Name:              "req",
										RequiredForImport: true,
										Description:       "this one's required",
									},
									{
										Name:              "opt",
										OptionalForImport: true,
										Description:       "this one's optional",
									},
								},
							},
						},
					},
				}).ProviderServer,
			},
			expectedIdentitySchemas: map[string]*tfprotov5.ResourceIdentitySchema{
				"test_resource_identity_foo": {
					Version: 1,
					IdentityAttributes: []*tfprotov5.ResourceIdentitySchemaAttribute{
						{
							Name:              "req",
							RequiredForImport: true,
							Description:       "this one's required",
						},
						{
							Name:              "opt",
							OptionalForImport: true,
							Description:       "this one's optional",
						},
					},
				},
				"test_resource_identity_bar": {
					Version: 1,
					IdentityAttributes: []*tfprotov5.ResourceIdentitySchemaAttribute{
						{
							Name:              "req",
							RequiredForImport: true,
							Description:       "this one's required",
						},
						{
							Name:              "opt",
							OptionalForImport: true,
							Description:       "this one's optional",
						},
					},
				},
				"test_resource_identity_foobar": {
					Version: 1,
					IdentityAttributes: []*tfprotov5.ResourceIdentitySchemaAttribute{
						{
							Name:              "req",
							RequiredForImport: true,
							Description:       "this one's required",
						},
						{
							Name:              "opt",
							OptionalForImport: true,
							Description:       "this one's optional",
						},
					},
				},
			},
			expectedDiagnostics: []*tfprotov5.Diagnostic{},
		},
		"duplicate-identity-schema-type": {
			servers: []func() tfprotov5.ProviderServer{
				(&tf5testserver.TestServer{
					GetResourceIdentitySchemasResponse: &tfprotov5.GetResourceIdentitySchemasResponse{
						IdentitySchemas: map[string]*tfprotov5.ResourceIdentitySchema{
							"test_foo": {},
						},
					},
				}).ProviderServer,
				(&tf5testserver.TestServer{
					GetResourceIdentitySchemasResponse: &tfprotov5.GetResourceIdentitySchemasResponse{
						IdentitySchemas: map[string]*tfprotov5.ResourceIdentitySchema{
							"test_foo": {},
						},
					},
				}).ProviderServer,
			},
			expectedIdentitySchemas: map[string]*tfprotov5.ResourceIdentitySchema{
				"test_foo": {},
			},
			expectedDiagnostics: []*tfprotov5.Diagnostic{
				{
					Severity: tfprotov5.DiagnosticSeverityError,
					Summary:  "Invalid Provider Server Combination",
					Detail: "The combined provider has multiple implementations of the same resource identity across underlying providers. " +
						"Resource identity types must be implemented by only one underlying provider. " +
						"This is always an issue in the provider implementation and should be reported to the provider developers.\n\n" +
						"Duplicate identity type for resource: test_foo",
				},
			},
		},
		"error-once": {
			servers: []func() tfprotov5.ProviderServer{
				(&tf5testserver.TestServer{
					GetResourceIdentitySchemasResponse: &tfprotov5.GetResourceIdentitySchemasResponse{
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
			expectedIdentitySchemas: map[string]*tfprotov5.ResourceIdentitySchema{},
			expectedDiagnostics: []*tfprotov5.Diagnostic{
				{
					Severity: tfprotov5.DiagnosticSeverityError,
					Summary:  "test error summary",
					Detail:   "test error details",
				},
			},
		},
		"error-multiple": {
			servers: []func() tfprotov5.ProviderServer{
				(&tf5testserver.TestServer{
					GetResourceIdentitySchemasResponse: &tfprotov5.GetResourceIdentitySchemasResponse{
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
					GetResourceIdentitySchemasResponse: &tfprotov5.GetResourceIdentitySchemasResponse{
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
			expectedIdentitySchemas: map[string]*tfprotov5.ResourceIdentitySchema{},
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
		},
		"warning-once": {
			servers: []func() tfprotov5.ProviderServer{
				(&tf5testserver.TestServer{
					GetResourceIdentitySchemasResponse: &tfprotov5.GetResourceIdentitySchemasResponse{
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
			expectedIdentitySchemas: map[string]*tfprotov5.ResourceIdentitySchema{},
			expectedDiagnostics: []*tfprotov5.Diagnostic{
				{
					Severity: tfprotov5.DiagnosticSeverityWarning,
					Summary:  "test warning summary",
					Detail:   "test warning details",
				},
			},
		},
		"warning-multiple": {
			servers: []func() tfprotov5.ProviderServer{
				(&tf5testserver.TestServer{
					GetResourceIdentitySchemasResponse: &tfprotov5.GetResourceIdentitySchemasResponse{
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
					GetResourceIdentitySchemasResponse: &tfprotov5.GetResourceIdentitySchemasResponse{
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
			expectedIdentitySchemas: map[string]*tfprotov5.ResourceIdentitySchema{},
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
		},
		"warning-then-error": {
			servers: []func() tfprotov5.ProviderServer{
				(&tf5testserver.TestServer{
					GetResourceIdentitySchemasResponse: &tfprotov5.GetResourceIdentitySchemasResponse{
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
					GetResourceIdentitySchemasResponse: &tfprotov5.GetResourceIdentitySchemasResponse{
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
			expectedIdentitySchemas: map[string]*tfprotov5.ResourceIdentitySchema{},
		},
	}

	for name, testCase := range testCases {

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			muxServer, err := tf5muxserver.NewMuxServer(context.Background(), testCase.servers...)

			if err != nil {
				t.Fatalf("unexpected error: %s", err)
			}

			resp, err := muxServer.ProviderServer().GetResourceIdentitySchemas(context.Background(), &tfprotov5.GetResourceIdentitySchemasRequest{})

			if err != nil {
				t.Fatalf("unexpected error: %s", err)
			}

			if diff := cmp.Diff(resp.Diagnostics, testCase.expectedDiagnostics); diff != "" {
				t.Errorf("diagnostics didn't match expectations: %s", diff)
			}

			if diff := cmp.Diff(resp.IdentitySchemas, testCase.expectedIdentitySchemas); diff != "" {
				t.Errorf("identity schemas didn't match expectations: %s", diff)
			}
		})
	}
}
