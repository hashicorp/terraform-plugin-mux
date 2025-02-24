// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package tf6muxserver_test

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-go/tftypes"

	"github.com/hashicorp/terraform-plugin-mux/internal/tf6testserver"
	"github.com/hashicorp/terraform-plugin-mux/tf6muxserver"
)

func TestMuxServerGetFunctions(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		servers  []func() tfprotov6.ProviderServer
		expected *tfprotov6.GetFunctionsResponse
	}{
		"combined": {
			servers: []func() tfprotov6.ProviderServer{
				(&tf6testserver.TestServer{
					GetFunctionsResponse: &tfprotov6.GetFunctionsResponse{
						Functions: map[string]*tfprotov6.Function{
							"test_function1": {
								Return: &tfprotov6.FunctionReturn{
									Type: tftypes.String,
								},
							},
						},
					},
				}).ProviderServer,
				(&tf6testserver.TestServer{
					GetFunctionsResponse: &tfprotov6.GetFunctionsResponse{
						Functions: map[string]*tfprotov6.Function{
							"test_function2": {
								Return: &tfprotov6.FunctionReturn{
									Type: tftypes.String,
								},
							},
							"test_function3": {
								Return: &tfprotov6.FunctionReturn{
									Type: tftypes.String,
								},
							},
						},
					},
				}).ProviderServer,
			},
			expected: &tfprotov6.GetFunctionsResponse{
				Functions: map[string]*tfprotov6.Function{
					"test_function1": {
						Return: &tfprotov6.FunctionReturn{
							Type: tftypes.String,
						},
					},
					"test_function2": {
						Return: &tfprotov6.FunctionReturn{
							Type: tftypes.String,
						},
					},
					"test_function3": {
						Return: &tfprotov6.FunctionReturn{
							Type: tftypes.String,
						},
					},
				},
			},
		},
		"duplicate-function": {
			servers: []func() tfprotov6.ProviderServer{
				(&tf6testserver.TestServer{
					GetFunctionsResponse: &tfprotov6.GetFunctionsResponse{
						Functions: map[string]*tfprotov6.Function{
							"test_function": {
								Return: &tfprotov6.FunctionReturn{
									Type: tftypes.String,
								},
							},
						},
					},
				}).ProviderServer,
				(&tf6testserver.TestServer{
					GetFunctionsResponse: &tfprotov6.GetFunctionsResponse{
						Functions: map[string]*tfprotov6.Function{
							"test_function": {
								Return: &tfprotov6.FunctionReturn{
									Type: tftypes.String,
								},
							},
						},
					},
				}).ProviderServer,
			},
			expected: &tfprotov6.GetFunctionsResponse{
				Diagnostics: []*tfprotov6.Diagnostic{
					{
						Severity: tfprotov6.DiagnosticSeverityError,
						Summary:  "Invalid Provider Server Combination",
						Detail: "The combined provider has multiple implementations of the same function name across underlying providers. " +
							"Functions must be implemented by only one underlying provider. " +
							"This is always an issue in the provider implementation and should be reported to the provider developers.\n\n" +
							"Duplicate function: test_function",
					},
				},
				Functions: map[string]*tfprotov6.Function{
					"test_function": {
						Return: &tfprotov6.FunctionReturn{
							Type: tftypes.String,
						},
					},
				},
			},
		},
		"error-once": {
			servers: []func() tfprotov6.ProviderServer{
				(&tf6testserver.TestServer{
					GetFunctionsResponse: &tfprotov6.GetFunctionsResponse{
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
			expected: &tfprotov6.GetFunctionsResponse{
				Diagnostics: []*tfprotov6.Diagnostic{
					{
						Severity: tfprotov6.DiagnosticSeverityError,
						Summary:  "test error summary",
						Detail:   "test error details",
					},
				},
				Functions: map[string]*tfprotov6.Function{},
			},
		},
		"error-multiple": {
			servers: []func() tfprotov6.ProviderServer{
				(&tf6testserver.TestServer{
					GetFunctionsResponse: &tfprotov6.GetFunctionsResponse{
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
					GetFunctionsResponse: &tfprotov6.GetFunctionsResponse{
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
			expected: &tfprotov6.GetFunctionsResponse{
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
				Functions: map[string]*tfprotov6.Function{},
			},
		},
		"warning-once": {
			servers: []func() tfprotov6.ProviderServer{
				(&tf6testserver.TestServer{
					GetFunctionsResponse: &tfprotov6.GetFunctionsResponse{
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
			expected: &tfprotov6.GetFunctionsResponse{
				Diagnostics: []*tfprotov6.Diagnostic{
					{
						Severity: tfprotov6.DiagnosticSeverityWarning,
						Summary:  "test warning summary",
						Detail:   "test warning details",
					},
				},
				Functions: map[string]*tfprotov6.Function{},
			},
		},
		"warning-multiple": {
			servers: []func() tfprotov6.ProviderServer{
				(&tf6testserver.TestServer{
					GetFunctionsResponse: &tfprotov6.GetFunctionsResponse{
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
					GetFunctionsResponse: &tfprotov6.GetFunctionsResponse{
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
			expected: &tfprotov6.GetFunctionsResponse{
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
				Functions: map[string]*tfprotov6.Function{},
			},
		},
		"warning-then-error": {
			servers: []func() tfprotov6.ProviderServer{
				(&tf6testserver.TestServer{
					GetFunctionsResponse: &tfprotov6.GetFunctionsResponse{
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
					GetFunctionsResponse: &tfprotov6.GetFunctionsResponse{
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
			expected: &tfprotov6.GetFunctionsResponse{
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
				Functions: map[string]*tfprotov6.Function{},
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

			resp, err := muxServer.ProviderServer().GetFunctions(context.Background(), &tfprotov6.GetFunctionsRequest{})

			if err != nil {
				t.Fatalf("unexpected error: %s", err)
			}

			if diff := cmp.Diff(resp, testCase.expected); diff != "" {
				t.Errorf("unexpected difference: %s", diff)
			}
		})
	}
}
