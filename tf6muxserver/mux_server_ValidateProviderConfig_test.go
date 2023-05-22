// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package tf6muxserver_test

import (
	"context"
	"fmt"
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
			},
		},
		"no-response": {
			testServers: [3]*tf6testserver.TestServer{
				{},
				{},
				{},
			},
		},
		"PreparedConfig-once": {
			testServers: [3]*tf6testserver.TestServer{
				{
					ProviderSchema: &configSchema,
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
					ProviderSchema: &configSchema,
					ValidateProviderConfigResponse: &tfprotov6.ValidateProviderConfigResponse{
						PreparedConfig: &config,
					},
				},
				{
					ProviderSchema: &configSchema,
				},
				{
					ProviderSchema: &configSchema,
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
					ProviderSchema: &configSchema,
					ValidateProviderConfigResponse: &tfprotov6.ValidateProviderConfigResponse{
						PreparedConfig: &config,
					},
				},
				{
					ProviderSchema: &configSchema,
				},
				{
					ProviderSchema: &configSchema,
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
					ProviderSchema: &configSchema,
					ValidateProviderConfigResponse: &tfprotov6.ValidateProviderConfigResponse{
						PreparedConfig: &config,
					},
				},
				{
					ProviderSchema: &configSchema,
				},
				{
					ProviderSchema: &configSchema,
					ValidateProviderConfigResponse: &tfprotov6.ValidateProviderConfigResponse{
						PreparedConfig: &config2,
					},
				},
			},
			expectedError: fmt.Errorf("got different PrepareProviderConfig PreparedConfig response from multiple servers, not sure which to use"),
		},
		"PreparedConfig-multiple-equal": {
			testServers: [3]*tf6testserver.TestServer{
				{
					ProviderSchema: &configSchema,
					ValidateProviderConfigResponse: &tfprotov6.ValidateProviderConfigResponse{
						PreparedConfig: &config,
					},
				},
				{
					ProviderSchema: &configSchema,
				},
				{
					ProviderSchema: &configSchema,
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
		name, testCase := name, testCase

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			servers := []func() tfprotov6.ProviderServer{
				testCase.testServers[0].ProviderServer,
				testCase.testServers[1].ProviderServer,
				testCase.testServers[2].ProviderServer,
			}

			muxServer, err := tf6muxserver.NewMuxServer(context.Background(), servers...)

			if err != nil {
				t.Fatalf("error setting up muxer: %s", err)
			}

			got, err := muxServer.ProviderServer().ValidateProviderConfig(context.Background(), &tfprotov6.ValidateProviderConfigRequest{
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
