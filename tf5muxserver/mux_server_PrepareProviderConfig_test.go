// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package tf5muxserver_test

import (
	"context"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/terraform-plugin-go/tfprotov5"
	"github.com/hashicorp/terraform-plugin-go/tftypes"

	"github.com/hashicorp/terraform-plugin-mux/internal/tf5testserver"
	"github.com/hashicorp/terraform-plugin-mux/tf5muxserver"
)

func TestMuxServerPrepareProviderConfig(t *testing.T) {
	t.Parallel()

	config, err := tfprotov5.NewDynamicValue(tftypes.Object{
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

	config2, err := tfprotov5.NewDynamicValue(tftypes.Object{
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

	configSchema := tfprotov5.Schema{
		Block: &tfprotov5.SchemaBlock{
			Attributes: []*tfprotov5.SchemaAttribute{
				{
					Name: "hello",
					Type: tftypes.String,
				},
			},
		},
	}

	testCases := map[string]struct {
		testServers      [3]*tf5testserver.TestServer
		expectedError    error
		expectedResponse *tfprotov5.PrepareProviderConfigResponse
	}{
		"error-once": {
			testServers: [3]*tf5testserver.TestServer{
				{
					PrepareProviderConfigResponse: &tfprotov5.PrepareProviderConfigResponse{
						Diagnostics: []*tfprotov5.Diagnostic{
							{
								Severity: tfprotov5.DiagnosticSeverityError,
								Summary:  "test error summary",
								Detail:   "test error details",
							},
						},
					},
				},
				{},
				{},
			},
			expectedResponse: &tfprotov5.PrepareProviderConfigResponse{
				Diagnostics: []*tfprotov5.Diagnostic{
					{
						Severity: tfprotov5.DiagnosticSeverityError,
						Summary:  "test error summary",
						Detail:   "test error details",
					},
				},
				PreparedConfig: &config,
			},
		},
		"error-multiple": {
			testServers: [3]*tf5testserver.TestServer{
				{
					PrepareProviderConfigResponse: &tfprotov5.PrepareProviderConfigResponse{
						Diagnostics: []*tfprotov5.Diagnostic{
							{
								Severity: tfprotov5.DiagnosticSeverityError,
								Summary:  "test error summary",
								Detail:   "test error details",
							},
						},
					},
				},
				{},
				{
					PrepareProviderConfigResponse: &tfprotov5.PrepareProviderConfigResponse{
						Diagnostics: []*tfprotov5.Diagnostic{
							{
								Severity: tfprotov5.DiagnosticSeverityError,
								Summary:  "test error summary",
								Detail:   "test error details",
							},
						},
					},
				},
			},
			expectedResponse: &tfprotov5.PrepareProviderConfigResponse{
				Diagnostics: []*tfprotov5.Diagnostic{
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
				PreparedConfig: &config,
			},
		},
		"warning-once": {
			testServers: [3]*tf5testserver.TestServer{
				{
					PrepareProviderConfigResponse: &tfprotov5.PrepareProviderConfigResponse{
						Diagnostics: []*tfprotov5.Diagnostic{
							{
								Severity: tfprotov5.DiagnosticSeverityWarning,
								Summary:  "test warning summary",
								Detail:   "test warning details",
							},
						},
					},
				},
				{},
				{},
			},
			expectedResponse: &tfprotov5.PrepareProviderConfigResponse{
				Diagnostics: []*tfprotov5.Diagnostic{
					{
						Severity: tfprotov5.DiagnosticSeverityWarning,
						Summary:  "test warning summary",
						Detail:   "test warning details",
					},
				},
				PreparedConfig: &config,
			},
		},
		"warning-multiple": {
			testServers: [3]*tf5testserver.TestServer{
				{
					PrepareProviderConfigResponse: &tfprotov5.PrepareProviderConfigResponse{
						Diagnostics: []*tfprotov5.Diagnostic{
							{
								Severity: tfprotov5.DiagnosticSeverityWarning,
								Summary:  "test warning summary",
								Detail:   "test warning details",
							},
						},
					},
				},
				{},
				{
					PrepareProviderConfigResponse: &tfprotov5.PrepareProviderConfigResponse{
						Diagnostics: []*tfprotov5.Diagnostic{
							{
								Severity: tfprotov5.DiagnosticSeverityWarning,
								Summary:  "test warning summary",
								Detail:   "test warning details",
							},
						},
					},
				},
			},
			expectedResponse: &tfprotov5.PrepareProviderConfigResponse{
				Diagnostics: []*tfprotov5.Diagnostic{
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
				PreparedConfig: &config,
			},
		},
		"warning-then-error": {
			testServers: [3]*tf5testserver.TestServer{
				{
					PrepareProviderConfigResponse: &tfprotov5.PrepareProviderConfigResponse{
						Diagnostics: []*tfprotov5.Diagnostic{
							{
								Severity: tfprotov5.DiagnosticSeverityWarning,
								Summary:  "test warning summary",
								Detail:   "test warning details",
							},
						},
					},
				},
				{},
				{

					PrepareProviderConfigResponse: &tfprotov5.PrepareProviderConfigResponse{
						Diagnostics: []*tfprotov5.Diagnostic{
							{
								Severity: tfprotov5.DiagnosticSeverityError,
								Summary:  "test error summary",
								Detail:   "test error details",
							},
						},
					},
				},
			},
			expectedResponse: &tfprotov5.PrepareProviderConfigResponse{
				Diagnostics: []*tfprotov5.Diagnostic{
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
				PreparedConfig: &config,
			},
		},
		"no-response": {
			testServers: [3]*tf5testserver.TestServer{
				{},
				{},
				{},
			},
			expectedResponse: &tfprotov5.PrepareProviderConfigResponse{
				PreparedConfig: &config,
			},
		},
		"PreparedConfig-once": {
			testServers: [3]*tf5testserver.TestServer{
				{
					GetProviderSchemaResponse: &tfprotov5.GetProviderSchemaResponse{
						Provider: &configSchema,
					},
					PrepareProviderConfigResponse: &tfprotov5.PrepareProviderConfigResponse{
						PreparedConfig: &config,
					},
				},
				{},
				{},
			},
			expectedResponse: &tfprotov5.PrepareProviderConfigResponse{
				PreparedConfig: &config,
			},
		},
		"PreparedConfig-once-and-error": {
			testServers: [3]*tf5testserver.TestServer{
				{
					GetProviderSchemaResponse: &tfprotov5.GetProviderSchemaResponse{
						Provider: &configSchema,
					},
					PrepareProviderConfigResponse: &tfprotov5.PrepareProviderConfigResponse{
						PreparedConfig: &config,
					},
				},
				{},
				{
					GetProviderSchemaResponse: &tfprotov5.GetProviderSchemaResponse{
						Provider: &configSchema,
					},
					PrepareProviderConfigResponse: &tfprotov5.PrepareProviderConfigResponse{
						Diagnostics: []*tfprotov5.Diagnostic{
							{
								Severity: tfprotov5.DiagnosticSeverityError,
								Summary:  "test error summary",
								Detail:   "test error details",
							},
						},
						PreparedConfig: &config,
					},
				},
			},
			expectedResponse: &tfprotov5.PrepareProviderConfigResponse{
				Diagnostics: []*tfprotov5.Diagnostic{
					{
						Severity: tfprotov5.DiagnosticSeverityError,
						Summary:  "test error summary",
						Detail:   "test error details",
					},
				},
				PreparedConfig: &config,
			},
		},
		"PreparedConfig-once-and-warning": {
			testServers: [3]*tf5testserver.TestServer{
				{
					GetProviderSchemaResponse: &tfprotov5.GetProviderSchemaResponse{
						Provider: &configSchema,
					},
					PrepareProviderConfigResponse: &tfprotov5.PrepareProviderConfigResponse{
						PreparedConfig: &config,
					},
				},
				{},
				{
					PrepareProviderConfigResponse: &tfprotov5.PrepareProviderConfigResponse{
						Diagnostics: []*tfprotov5.Diagnostic{
							{
								Severity: tfprotov5.DiagnosticSeverityWarning,
								Summary:  "test warning summary",
								Detail:   "test warning details",
							},
						},
					},
				},
			},
			expectedResponse: &tfprotov5.PrepareProviderConfigResponse{
				Diagnostics: []*tfprotov5.Diagnostic{
					{
						Severity: tfprotov5.DiagnosticSeverityWarning,
						Summary:  "test warning summary",
						Detail:   "test warning details",
					},
				},
				PreparedConfig: &config,
			},
		},
		"PreparedConfig-multiple-different": {
			testServers: [3]*tf5testserver.TestServer{
				{
					GetProviderSchemaResponse: &tfprotov5.GetProviderSchemaResponse{
						Provider: &configSchema,
					},
					PrepareProviderConfigResponse: &tfprotov5.PrepareProviderConfigResponse{
						PreparedConfig: &config,
					},
				},
				{},
				{
					GetProviderSchemaResponse: &tfprotov5.GetProviderSchemaResponse{
						Provider: &configSchema,
					},
					PrepareProviderConfigResponse: &tfprotov5.PrepareProviderConfigResponse{
						PreparedConfig: &config2, // intentionally ignored
					},
				},
			},
			expectedResponse: &tfprotov5.PrepareProviderConfigResponse{
				PreparedConfig: &config,
			},
		},
		"PreparedConfig-multiple-equal": {
			testServers: [3]*tf5testserver.TestServer{
				{
					GetProviderSchemaResponse: &tfprotov5.GetProviderSchemaResponse{
						Provider: &configSchema,
					},
					PrepareProviderConfigResponse: &tfprotov5.PrepareProviderConfigResponse{
						PreparedConfig: &config,
					},
				},
				{},
				{
					GetProviderSchemaResponse: &tfprotov5.GetProviderSchemaResponse{
						Provider: &configSchema,
					},
					PrepareProviderConfigResponse: &tfprotov5.PrepareProviderConfigResponse{
						PreparedConfig: &config,
					},
				},
			},
			expectedResponse: &tfprotov5.PrepareProviderConfigResponse{
				PreparedConfig: &config,
			},
		},
	}

	for name, testCase := range testCases {

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			ctx := context.Background()
			servers := []func() tfprotov5.ProviderServer{
				testCase.testServers[0].ProviderServer,
				testCase.testServers[1].ProviderServer,
				testCase.testServers[2].ProviderServer,
			}

			muxServer, err := tf5muxserver.NewMuxServer(ctx, servers...)

			if err != nil {
				t.Fatalf("error setting up muxer: %s", err)
			}

			got, err := muxServer.ProviderServer().PrepareProviderConfig(ctx, &tfprotov5.PrepareProviderConfigRequest{
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
				if !testServer.PrepareProviderConfigCalled {
					t.Errorf("PrepareProviderConfig not called on server%d", num+1)
				}
			}
		})
	}
}
