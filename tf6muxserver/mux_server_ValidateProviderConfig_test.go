// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package tf6muxserver_test

import (
	"context"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-go/tftypes"

	"github.com/hashicorp/terraform-plugin-mux/internal/tf6testserver"
	"github.com/hashicorp/terraform-plugin-mux/tf6muxserver"
)

func TestMuxServerValidateProviderConfig(t *testing.T) {
	t.Parallel()

	config, err := tfprotov6.NewDynamicValue(tftypes.Object{
		AttributeTypes: map[string]tftypes.Type{
			"hello": tftypes.String,
		},
	}, tftypes.NewValue(tftypes.Object{
		AttributeTypes: map[string]tftypes.Type{
			"hello": tftypes.String,
		},
	}, map[string]tftypes.Value{
		"hello": tftypes.NewValue(tftypes.String, "world"),
	}))

	if err != nil {
		t.Fatalf("error constructing config: %s", err)
	}

	config2, err := tfprotov6.NewDynamicValue(tftypes.Object{
		AttributeTypes: map[string]tftypes.Type{
			"hello": tftypes.String,
		},
	}, tftypes.NewValue(tftypes.Object{
		AttributeTypes: map[string]tftypes.Type{
			"hello": tftypes.String,
		},
	}, map[string]tftypes.Value{
		"hello": tftypes.NewValue(tftypes.String, "goodbye"),
	}))

	if err != nil {
		t.Fatalf("error constructing config: %s", err)
	}

	configSchema := tfprotov6.Schema{
		Block: &tfprotov6.SchemaBlock{
			Attributes: []*tfprotov6.SchemaAttribute{
				{
					Name: "hello",
					Type: tftypes.String,
				},
			},
		},
	}

	testCases := map[string]struct {
		testServers      [3]*tf6testserver.TestServer
		expectedError    error
		expectedResponse *tfprotov6.ValidateProviderConfigResponse
	}{
		"error-once": {
			testServers: [3]*tf6testserver.TestServer{
				{
					ValidateProviderConfigResponse: &tfprotov6.ValidateProviderConfigResponse{
						Diagnostics: []*tfprotov6.Diagnostic{
							{
								Severity: tfprotov6.DiagnosticSeverityError,
								Summary:  "test error summary",
								Detail:   "test error details",
							},
						},
					},
				},
				{},
				{},
			},
			expectedResponse: &tfprotov6.ValidateProviderConfigResponse{
				Diagnostics: []*tfprotov6.Diagnostic{
					{
						Severity: tfprotov6.DiagnosticSeverityError,
						Summary:  "test error summary",
						Detail:   "test error details",
					},
				},
				PreparedConfig: &config,
			},
		},
		"error-multiple": {
			testServers: [3]*tf6testserver.TestServer{
				{
					ValidateProviderConfigResponse: &tfprotov6.ValidateProviderConfigResponse{
						Diagnostics: []*tfprotov6.Diagnostic{
							{
								Severity: tfprotov6.DiagnosticSeverityError,
								Summary:  "test error summary",
								Detail:   "test error details",
							},
						},
					},
				},
				{},
				{
					ValidateProviderConfigResponse: &tfprotov6.ValidateProviderConfigResponse{
						Diagnostics: []*tfprotov6.Diagnostic{
							{
								Severity: tfprotov6.DiagnosticSeverityError,
								Summary:  "test error summary",
								Detail:   "test error details",
							},
						},
					},
				},
			},
			expectedResponse: &tfprotov6.ValidateProviderConfigResponse{
				Diagnostics: []*tfprotov6.Diagnostic{
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
				PreparedConfig: &config,
			},
		},
		"warning-once": {
			testServers: [3]*tf6testserver.TestServer{
				{
					ValidateProviderConfigResponse: &tfprotov6.ValidateProviderConfigResponse{
						Diagnostics: []*tfprotov6.Diagnostic{
							{
								Severity: tfprotov6.DiagnosticSeverityWarning,
								Summary:  "test warning summary",
								Detail:   "test warning details",
							},
						},
					},
				},
				{},
				{},
			},
			expectedResponse: &tfprotov6.ValidateProviderConfigResponse{
				Diagnostics: []*tfprotov6.Diagnostic{
					{
						Severity: tfprotov6.DiagnosticSeverityWarning,
						Summary:  "test warning summary",
						Detail:   "test warning details",
					},
				},
				PreparedConfig: &config,
			},
		},
		"warning-multiple": {
			testServers: [3]*tf6testserver.TestServer{
				{
					ValidateProviderConfigResponse: &tfprotov6.ValidateProviderConfigResponse{
						Diagnostics: []*tfprotov6.Diagnostic{
							{
								Severity: tfprotov6.DiagnosticSeverityWarning,
								Summary:  "test warning summary",
								Detail:   "test warning details",
							},
						},
					},
				},
				{},
				{
					ValidateProviderConfigResponse: &tfprotov6.ValidateProviderConfigResponse{
						Diagnostics: []*tfprotov6.Diagnostic{
							{
								Severity: tfprotov6.DiagnosticSeverityWarning,
								Summary:  "test warning summary",
								Detail:   "test warning details",
							},
						},
					},
				},
			},
			expectedResponse: &tfprotov6.ValidateProviderConfigResponse{
				Diagnostics: []*tfprotov6.Diagnostic{
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
				PreparedConfig: &config,
			},
		},
		"warning-then-error": {
			testServers: [3]*tf6testserver.TestServer{
				{
					ValidateProviderConfigResponse: &tfprotov6.ValidateProviderConfigResponse{
						Diagnostics: []*tfprotov6.Diagnostic{
							{
								Severity: tfprotov6.DiagnosticSeverityWarning,
								Summary:  "test warning summary",
								Detail:   "test warning details",
							},
						},
					},
				},
				{},
				{
					ValidateProviderConfigResponse: &tfprotov6.ValidateProviderConfigResponse{
						Diagnostics: []*tfprotov6.Diagnostic{
							{
								Severity: tfprotov6.DiagnosticSeverityError,
								Summary:  "test error summary",
								Detail:   "test error details",
							},
						},
					},
				},
			},
			expectedResponse: &tfprotov6.ValidateProviderConfigResponse{
				Diagnostics: []*tfprotov6.Diagnostic{
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
				PreparedConfig: &config,
			},
		},
		"no-response": {
			testServers: [3]*tf6testserver.TestServer{
				{},
				{},
				{},
			},
			expectedResponse: &tfprotov6.ValidateProviderConfigResponse{
				PreparedConfig: &config,
			},
		},
		"PreparedConfig-once": {
			testServers: [3]*tf6testserver.TestServer{
				{
					GetProviderSchemaResponse: &tfprotov6.GetProviderSchemaResponse{
						Provider: &configSchema,
					},
					ValidateProviderConfigResponse: &tfprotov6.ValidateProviderConfigResponse{
						PreparedConfig: &config,
					},
				},
				{},
				{},
			},
			expectedResponse: &tfprotov6.ValidateProviderConfigResponse{
				PreparedConfig: &config,
			},
		},
		"PreparedConfig-once-and-error": {
			testServers: [3]*tf6testserver.TestServer{
				{
					GetProviderSchemaResponse: &tfprotov6.GetProviderSchemaResponse{
						Provider: &configSchema,
					},
					ValidateProviderConfigResponse: &tfprotov6.ValidateProviderConfigResponse{
						PreparedConfig: &config,
					},
				},
				{
					GetProviderSchemaResponse: &tfprotov6.GetProviderSchemaResponse{
						Provider: &configSchema,
					},
				},
				{
					GetProviderSchemaResponse: &tfprotov6.GetProviderSchemaResponse{
						Provider: &configSchema,
					},
					ValidateProviderConfigResponse: &tfprotov6.ValidateProviderConfigResponse{
						Diagnostics: []*tfprotov6.Diagnostic{
							{
								Severity: tfprotov6.DiagnosticSeverityError,
								Summary:  "test error summary",
								Detail:   "test error details",
							},
						},
					},
				},
			},
			expectedResponse: &tfprotov6.ValidateProviderConfigResponse{
				Diagnostics: []*tfprotov6.Diagnostic{
					{
						Severity: tfprotov6.DiagnosticSeverityError,
						Summary:  "test error summary",
						Detail:   "test error details",
					},
				},
				PreparedConfig: &config,
			},
		},
		"PreparedConfig-once-and-warning": {
			testServers: [3]*tf6testserver.TestServer{
				{
					GetProviderSchemaResponse: &tfprotov6.GetProviderSchemaResponse{
						Provider: &configSchema,
					},
					ValidateProviderConfigResponse: &tfprotov6.ValidateProviderConfigResponse{
						PreparedConfig: &config,
					},
				},
				{
					GetProviderSchemaResponse: &tfprotov6.GetProviderSchemaResponse{
						Provider: &configSchema,
					},
				},
				{
					GetProviderSchemaResponse: &tfprotov6.GetProviderSchemaResponse{
						Provider: &configSchema,
					},
					ValidateProviderConfigResponse: &tfprotov6.ValidateProviderConfigResponse{
						Diagnostics: []*tfprotov6.Diagnostic{
							{
								Severity: tfprotov6.DiagnosticSeverityWarning,
								Summary:  "test warning summary",
								Detail:   "test warning details",
							},
						},
					},
				},
			},
			expectedResponse: &tfprotov6.ValidateProviderConfigResponse{
				Diagnostics: []*tfprotov6.Diagnostic{
					{
						Severity: tfprotov6.DiagnosticSeverityWarning,
						Summary:  "test warning summary",
						Detail:   "test warning details",
					},
				},
				PreparedConfig: &config,
			},
		},
		"PreparedConfig-multiple-different": {
			testServers: [3]*tf6testserver.TestServer{
				{
					GetProviderSchemaResponse: &tfprotov6.GetProviderSchemaResponse{
						Provider: &configSchema,
					},
					ValidateProviderConfigResponse: &tfprotov6.ValidateProviderConfigResponse{
						PreparedConfig: &config,
					},
				},
				{
					GetProviderSchemaResponse: &tfprotov6.GetProviderSchemaResponse{
						Provider: &configSchema,
					},
				},
				{
					GetProviderSchemaResponse: &tfprotov6.GetProviderSchemaResponse{
						Provider: &configSchema,
					},
					ValidateProviderConfigResponse: &tfprotov6.ValidateProviderConfigResponse{
						PreparedConfig: &config2, // intentionally ignored
					},
				},
			},
			expectedResponse: &tfprotov6.ValidateProviderConfigResponse{
				PreparedConfig: &config,
			},
		},
		"PreparedConfig-multiple-equal": {
			testServers: [3]*tf6testserver.TestServer{
				{
					GetProviderSchemaResponse: &tfprotov6.GetProviderSchemaResponse{
						Provider: &configSchema,
					},
					ValidateProviderConfigResponse: &tfprotov6.ValidateProviderConfigResponse{
						PreparedConfig: &config,
					},
				},
				{
					GetProviderSchemaResponse: &tfprotov6.GetProviderSchemaResponse{
						Provider: &configSchema,
					},
				},
				{
					GetProviderSchemaResponse: &tfprotov6.GetProviderSchemaResponse{
						Provider: &configSchema,
					},
					ValidateProviderConfigResponse: &tfprotov6.ValidateProviderConfigResponse{
						PreparedConfig: &config,
					},
				},
			},
			expectedResponse: &tfprotov6.ValidateProviderConfigResponse{
				PreparedConfig: &config,
			},
		},
	}

	for name, testCase := range testCases {

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			ctx := context.Background()
			servers := []func() tfprotov6.ProviderServer{
				testCase.testServers[0].ProviderServer,
				testCase.testServers[1].ProviderServer,
				testCase.testServers[2].ProviderServer,
			}

			muxServer, err := tf6muxserver.NewMuxServer(ctx, servers...)

			if err != nil {
				t.Fatalf("error setting up muxer: %s", err)
			}

			got, err := muxServer.ProviderServer().ValidateProviderConfig(ctx, &tfprotov6.ValidateProviderConfigRequest{
				Config: &config,
			})

			if err != nil {
				if testCase.expectedError == nil {
					t.Fatalf("wanted no error, got error: %s", err)
				}

				if !strings.Contains(err.Error(), testCase.expectedError.Error()) {
					t.Fatalf("wanted error %q, got error: %s", testCase.expectedError.Error(), err.Error())
				}
			}

			if err == nil && testCase.expectedError != nil {
				t.Fatalf("got no error, wanted err: %s", testCase.expectedError)
			}

			if !cmp.Equal(got, testCase.expectedResponse) {
				t.Errorf("unexpected response: %s", cmp.Diff(got, testCase.expectedResponse))
			}

			for num, testServer := range testCase.testServers {
				if !testServer.ValidateProviderConfigCalled {
					t.Errorf("ValidateProviderConfig not called on server%d", num+1)
				}
			}
		})
	}
}
