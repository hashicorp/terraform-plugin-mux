// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package tf5muxserver_test

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/terraform-plugin-go/tfprotov5"
	"github.com/hashicorp/terraform-plugin-go/tftypes"

	"github.com/hashicorp/terraform-plugin-mux/internal/tf5testserver"
	"github.com/hashicorp/terraform-plugin-mux/tf5muxserver"
)

func TestMuxServerGetProviderSchema(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		servers                    []func() tfprotov5.ProviderServer
		expectedDataSourceSchemas  map[string]*tfprotov5.Schema
		expectedDiagnostics        []*tfprotov5.Diagnostic
		expectedFunctions          map[string]*tfprotov5.Function
		expectedProviderSchema     *tfprotov5.Schema
		expectedProviderMetaSchema *tfprotov5.Schema
		expectedResourceSchemas    map[string]*tfprotov5.Schema
		expectedServerCapabilities *tfprotov5.ServerCapabilities
	}{
		"combined": {
			servers: []func() tfprotov5.ProviderServer{
				(&tf5testserver.TestServer{
					GetProviderSchemaResponse: &tfprotov5.GetProviderSchemaResponse{
						Provider: &tfprotov5.Schema{
							Version: 1,
							Block: &tfprotov5.SchemaBlock{
								Version: 1,
								Attributes: []*tfprotov5.SchemaAttribute{
									{
										Name:            "account_id",
										Type:            tftypes.String,
										Required:        true,
										Description:     "the account ID to make requests for",
										DescriptionKind: tfprotov5.StringKindPlain,
									},
								},
								BlockTypes: []*tfprotov5.SchemaNestedBlock{
									{
										TypeName: "feature",
										Nesting:  tfprotov5.SchemaNestedBlockNestingModeList,
										Block: &tfprotov5.SchemaBlock{
											Version:         1,
											Description:     "features to enable on the provider",
											DescriptionKind: tfprotov5.StringKindPlain,
											Attributes: []*tfprotov5.SchemaAttribute{
												{
													Name:            "feature_id",
													Type:            tftypes.Number,
													Required:        true,
													Description:     "The ID of the feature",
													DescriptionKind: tfprotov5.StringKindPlain,
												},
												{
													Name:            "enabled",
													Type:            tftypes.Bool,
													Optional:        true,
													Description:     "whether the feature is enabled",
													DescriptionKind: tfprotov5.StringKindPlain,
												},
											},
										},
									},
								},
							},
						},
						ProviderMeta: &tfprotov5.Schema{
							Version: 4,
							Block: &tfprotov5.SchemaBlock{
								Version: 4,
								Attributes: []*tfprotov5.SchemaAttribute{
									{
										Name:            "module_id",
										Type:            tftypes.String,
										Optional:        true,
										Description:     "a unique identifier for the module",
										DescriptionKind: tfprotov5.StringKindPlain,
									},
								},
							},
						},
						ResourceSchemas: map[string]*tfprotov5.Schema{
							"test_foo": {
								Version: 1,
								Block: &tfprotov5.SchemaBlock{
									Version: 1,
									Attributes: []*tfprotov5.SchemaAttribute{
										{
											Name:            "airspeed_velocity",
											Type:            tftypes.Number,
											Required:        true,
											Description:     "the airspeed velocity of a swallow",
											DescriptionKind: tfprotov5.StringKindPlain,
										},
									},
								},
							},
							"test_bar": {
								Version: 1,
								Block: &tfprotov5.SchemaBlock{
									Version: 1,
									Attributes: []*tfprotov5.SchemaAttribute{
										{
											Name:            "name",
											Type:            tftypes.String,
											Optional:        true,
											Description:     "your name",
											DescriptionKind: tfprotov5.StringKindPlain,
										},
										{
											Name:            "color",
											Type:            tftypes.String,
											Optional:        true,
											Description:     "your favorite color",
											DescriptionKind: tfprotov5.StringKindPlain,
										},
									},
								},
							},
						},
						DataSourceSchemas: map[string]*tfprotov5.Schema{
							"test_foo": {
								Version: 1,
								Block: &tfprotov5.SchemaBlock{
									Version: 1,
									Attributes: []*tfprotov5.SchemaAttribute{
										{
											Name:            "current_time",
											Type:            tftypes.String,
											Computed:        true,
											Description:     "the current time in RFC 3339 format",
											DescriptionKind: tfprotov5.StringKindPlain,
										},
									},
								},
							},
						},
						Functions: map[string]*tfprotov5.Function{
							"test_function1": {
								Return: &tfprotov5.FunctionReturn{
									Type: tftypes.String,
								},
							},
						},
					},
				}).ProviderServer,
				(&tf5testserver.TestServer{
					GetProviderSchemaResponse: &tfprotov5.GetProviderSchemaResponse{
						Provider: &tfprotov5.Schema{
							Version: 1,
							Block: &tfprotov5.SchemaBlock{
								Version: 1,
								Attributes: []*tfprotov5.SchemaAttribute{
									{
										Name:            "account_id",
										Type:            tftypes.String,
										Required:        true,
										Description:     "the account ID to make requests for",
										DescriptionKind: tfprotov5.StringKindPlain,
									},
								},
								BlockTypes: []*tfprotov5.SchemaNestedBlock{
									{
										TypeName: "feature",
										Nesting:  tfprotov5.SchemaNestedBlockNestingModeList,
										Block: &tfprotov5.SchemaBlock{
											Version:         1,
											Description:     "features to enable on the provider",
											DescriptionKind: tfprotov5.StringKindPlain,
											Attributes: []*tfprotov5.SchemaAttribute{
												{
													Name:            "feature_id",
													Type:            tftypes.Number,
													Required:        true,
													Description:     "The ID of the feature",
													DescriptionKind: tfprotov5.StringKindPlain,
												},
												{
													Name:            "enabled",
													Type:            tftypes.Bool,
													Optional:        true,
													Description:     "whether the feature is enabled",
													DescriptionKind: tfprotov5.StringKindPlain,
												},
											},
										},
									},
								},
							},
						},
						ProviderMeta: &tfprotov5.Schema{
							Version: 4,
							Block: &tfprotov5.SchemaBlock{
								Version: 4,
								Attributes: []*tfprotov5.SchemaAttribute{
									{
										Name:            "module_id",
										Type:            tftypes.String,
										Optional:        true,
										Description:     "a unique identifier for the module",
										DescriptionKind: tfprotov5.StringKindPlain,
									},
								},
							},
						},
						ResourceSchemas: map[string]*tfprotov5.Schema{
							"test_quux": {
								Version: 1,
								Block: &tfprotov5.SchemaBlock{
									Version: 1,
									Attributes: []*tfprotov5.SchemaAttribute{
										{
											Name:            "a",
											Type:            tftypes.String,
											Required:        true,
											Description:     "the account ID to make requests for",
											DescriptionKind: tfprotov5.StringKindPlain,
										},
										{
											Name:     "b",
											Type:     tftypes.String,
											Required: true,
										},
									},
								},
							},
						},
						DataSourceSchemas: map[string]*tfprotov5.Schema{
							"test_bar": {
								Version: 1,
								Block: &tfprotov5.SchemaBlock{
									Version: 1,
									Attributes: []*tfprotov5.SchemaAttribute{
										{
											Name:            "a",
											Type:            tftypes.Number,
											Computed:        true,
											Description:     "some field that's set by the provider",
											DescriptionKind: tfprotov5.StringKindMarkdown,
										},
									},
								},
							},
							"test_quux": {
								Version: 1,
								Block: &tfprotov5.SchemaBlock{
									Version: 1,
									Attributes: []*tfprotov5.SchemaAttribute{
										{
											Name:            "abc",
											Type:            tftypes.Number,
											Computed:        true,
											Description:     "some other field that's set by the provider",
											DescriptionKind: tfprotov5.StringKindMarkdown,
										},
									},
								},
							},
						},
						Functions: map[string]*tfprotov5.Function{
							"test_function2": {
								Return: &tfprotov5.FunctionReturn{
									Type: tftypes.String,
								},
							},
							"test_function3": {
								Return: &tfprotov5.FunctionReturn{
									Type: tftypes.String,
								},
							},
						},
					},
				}).ProviderServer,
			},
			expectedProviderSchema: &tfprotov5.Schema{
				Version: 1,
				Block: &tfprotov5.SchemaBlock{
					Version: 1,
					Attributes: []*tfprotov5.SchemaAttribute{
						{
							Name:            "account_id",
							Type:            tftypes.String,
							Required:        true,
							Description:     "the account ID to make requests for",
							DescriptionKind: tfprotov5.StringKindPlain,
						},
					},
					BlockTypes: []*tfprotov5.SchemaNestedBlock{
						{
							TypeName: "feature",
							Nesting:  tfprotov5.SchemaNestedBlockNestingModeList,
							Block: &tfprotov5.SchemaBlock{
								Version:         1,
								Description:     "features to enable on the provider",
								DescriptionKind: tfprotov5.StringKindPlain,
								Attributes: []*tfprotov5.SchemaAttribute{
									{
										Name:            "feature_id",
										Type:            tftypes.Number,
										Required:        true,
										Description:     "The ID of the feature",
										DescriptionKind: tfprotov5.StringKindPlain,
									},
									{
										Name:            "enabled",
										Type:            tftypes.Bool,
										Optional:        true,
										Description:     "whether the feature is enabled",
										DescriptionKind: tfprotov5.StringKindPlain,
									},
								},
							},
						},
					},
				},
			},
			expectedProviderMetaSchema: &tfprotov5.Schema{
				Version: 4,
				Block: &tfprotov5.SchemaBlock{
					Version: 4,
					Attributes: []*tfprotov5.SchemaAttribute{
						{
							Name:            "module_id",
							Type:            tftypes.String,
							Optional:        true,
							Description:     "a unique identifier for the module",
							DescriptionKind: tfprotov5.StringKindPlain,
						},
					},
				},
			},
			expectedResourceSchemas: map[string]*tfprotov5.Schema{
				"test_foo": {
					Version: 1,
					Block: &tfprotov5.SchemaBlock{
						Version: 1,
						Attributes: []*tfprotov5.SchemaAttribute{
							{
								Name:            "airspeed_velocity",
								Type:            tftypes.Number,
								Required:        true,
								Description:     "the airspeed velocity of a swallow",
								DescriptionKind: tfprotov5.StringKindPlain,
							},
						},
					},
				},
				"test_bar": {
					Version: 1,
					Block: &tfprotov5.SchemaBlock{
						Version: 1,
						Attributes: []*tfprotov5.SchemaAttribute{
							{
								Name:            "name",
								Type:            tftypes.String,
								Optional:        true,
								Description:     "your name",
								DescriptionKind: tfprotov5.StringKindPlain,
							},
							{
								Name:            "color",
								Type:            tftypes.String,
								Optional:        true,
								Description:     "your favorite color",
								DescriptionKind: tfprotov5.StringKindPlain,
							},
						},
					},
				},
				"test_quux": {
					Version: 1,
					Block: &tfprotov5.SchemaBlock{
						Version: 1,
						Attributes: []*tfprotov5.SchemaAttribute{
							{
								Name:            "a",
								Type:            tftypes.String,
								Required:        true,
								Description:     "the account ID to make requests for",
								DescriptionKind: tfprotov5.StringKindPlain,
							},
							{
								Name:     "b",
								Type:     tftypes.String,
								Required: true,
							},
						},
					},
				},
			},
			expectedDataSourceSchemas: map[string]*tfprotov5.Schema{
				"test_foo": {
					Version: 1,
					Block: &tfprotov5.SchemaBlock{
						Version: 1,
						Attributes: []*tfprotov5.SchemaAttribute{
							{
								Name:            "current_time",
								Type:            tftypes.String,
								Computed:        true,
								Description:     "the current time in RFC 3339 format",
								DescriptionKind: tfprotov5.StringKindPlain,
							},
						},
					},
				},
				"test_bar": {
					Version: 1,
					Block: &tfprotov5.SchemaBlock{
						Version: 1,
						Attributes: []*tfprotov5.SchemaAttribute{
							{
								Name:            "a",
								Type:            tftypes.Number,
								Computed:        true,
								Description:     "some field that's set by the provider",
								DescriptionKind: tfprotov5.StringKindMarkdown,
							},
						},
					},
				},
				"test_quux": {
					Version: 1,
					Block: &tfprotov5.SchemaBlock{
						Version: 1,
						Attributes: []*tfprotov5.SchemaAttribute{
							{
								Name:            "abc",
								Type:            tftypes.Number,
								Computed:        true,
								Description:     "some other field that's set by the provider",
								DescriptionKind: tfprotov5.StringKindMarkdown,
							},
						},
					},
				},
			},
			expectedFunctions: map[string]*tfprotov5.Function{
				"test_function1": {
					Return: &tfprotov5.FunctionReturn{
						Type: tftypes.String,
					},
				},
				"test_function2": {
					Return: &tfprotov5.FunctionReturn{
						Type: tftypes.String,
					},
				},
				"test_function3": {
					Return: &tfprotov5.FunctionReturn{
						Type: tftypes.String,
					},
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
					GetProviderSchemaResponse: &tfprotov5.GetProviderSchemaResponse{
						DataSourceSchemas: map[string]*tfprotov5.Schema{
							"test_foo": {},
						},
					},
				}).ProviderServer,
				(&tf5testserver.TestServer{
					GetProviderSchemaResponse: &tfprotov5.GetProviderSchemaResponse{
						DataSourceSchemas: map[string]*tfprotov5.Schema{
							"test_foo": {},
						},
					},
				}).ProviderServer,
			},
			expectedDataSourceSchemas: map[string]*tfprotov5.Schema{
				"test_foo": {},
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
			expectedFunctions:       map[string]*tfprotov5.Function{},
			expectedResourceSchemas: map[string]*tfprotov5.Schema{},
			expectedServerCapabilities: &tfprotov5.ServerCapabilities{
				GetProviderSchemaOptional: true,
				MoveResourceState:         true,
				PlanDestroy:               true,
			},
		},
		"duplicate-function": {
			servers: []func() tfprotov5.ProviderServer{
				(&tf5testserver.TestServer{
					GetProviderSchemaResponse: &tfprotov5.GetProviderSchemaResponse{
						Functions: map[string]*tfprotov5.Function{
							"test_function": {},
						},
					},
				}).ProviderServer,
				(&tf5testserver.TestServer{
					GetProviderSchemaResponse: &tfprotov5.GetProviderSchemaResponse{
						Functions: map[string]*tfprotov5.Function{
							"test_function": {},
						},
					},
				}).ProviderServer,
			},
			expectedDataSourceSchemas: map[string]*tfprotov5.Schema{},
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
			expectedFunctions: map[string]*tfprotov5.Function{
				"test_function": {},
			},
			expectedResourceSchemas: map[string]*tfprotov5.Schema{},
			expectedServerCapabilities: &tfprotov5.ServerCapabilities{
				GetProviderSchemaOptional: true,
				MoveResourceState:         true,
				PlanDestroy:               true,
			},
		},
		"duplicate-resource-type": {
			servers: []func() tfprotov5.ProviderServer{
				(&tf5testserver.TestServer{
					GetProviderSchemaResponse: &tfprotov5.GetProviderSchemaResponse{
						ResourceSchemas: map[string]*tfprotov5.Schema{
							"test_foo": {},
						},
					},
				}).ProviderServer,
				(&tf5testserver.TestServer{
					GetProviderSchemaResponse: &tfprotov5.GetProviderSchemaResponse{
						ResourceSchemas: map[string]*tfprotov5.Schema{
							"test_foo": {},
						},
					},
				}).ProviderServer,
			},
			expectedDataSourceSchemas: map[string]*tfprotov5.Schema{},
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
			expectedFunctions: map[string]*tfprotov5.Function{},
			expectedResourceSchemas: map[string]*tfprotov5.Schema{
				"test_foo": {},
			},
			expectedServerCapabilities: &tfprotov5.ServerCapabilities{
				GetProviderSchemaOptional: true,
				MoveResourceState:         true,
				PlanDestroy:               true,
			},
		},
		"provider-mismatch": {
			servers: []func() tfprotov5.ProviderServer{
				(&tf5testserver.TestServer{
					GetProviderSchemaResponse: &tfprotov5.GetProviderSchemaResponse{
						Provider: &tfprotov5.Schema{
							Block: &tfprotov5.SchemaBlock{
								Attributes: []*tfprotov5.SchemaAttribute{
									{
										Name:     "testattribute1",
										Type:     tftypes.String,
										Required: true,
									},
								},
							},
						},
					},
				}).ProviderServer,
				(&tf5testserver.TestServer{
					GetProviderSchemaResponse: &tfprotov5.GetProviderSchemaResponse{
						Provider: &tfprotov5.Schema{
							Block: &tfprotov5.SchemaBlock{
								Attributes: []*tfprotov5.SchemaAttribute{
									{
										Name:     "testattribute2",
										Type:     tftypes.String,
										Required: true,
									},
								},
							},
						},
					},
				}).ProviderServer,
			},
			expectedDataSourceSchemas: map[string]*tfprotov5.Schema{},
			expectedDiagnostics: []*tfprotov5.Diagnostic{
				{
					Severity: tfprotov5.DiagnosticSeverityError,
					Summary:  "Invalid Provider Server Combination",
					Detail: "The combined provider has differing provider schema implementations across providers. " +
						"Provider schemas must be identical across providers. " +
						"This is always an issue in the provider implementation and should be reported to the provider developers.\n\n" +
						"Provider schema difference: " + cmp.Diff(
						&tfprotov5.Schema{
							Block: &tfprotov5.SchemaBlock{
								Attributes: []*tfprotov5.SchemaAttribute{
									{
										Name:     "testattribute2",
										Type:     tftypes.String,
										Required: true,
									},
								},
							},
						},
						&tfprotov5.Schema{
							Block: &tfprotov5.SchemaBlock{
								Attributes: []*tfprotov5.SchemaAttribute{
									{
										Name:     "testattribute1",
										Type:     tftypes.String,
										Required: true,
									},
								},
							},
						},
					),
				},
			},
			expectedFunctions: map[string]*tfprotov5.Function{},
			expectedProviderSchema: &tfprotov5.Schema{
				Block: &tfprotov5.SchemaBlock{
					Attributes: []*tfprotov5.SchemaAttribute{
						{
							Name:     "testattribute1",
							Type:     tftypes.String,
							Required: true,
						},
					},
				},
			},
			expectedResourceSchemas: map[string]*tfprotov5.Schema{},
			expectedServerCapabilities: &tfprotov5.ServerCapabilities{
				GetProviderSchemaOptional: true,
				MoveResourceState:         true,
				PlanDestroy:               true,
			},
		},
		"provider-meta-mismatch": {
			servers: []func() tfprotov5.ProviderServer{
				(&tf5testserver.TestServer{
					GetProviderSchemaResponse: &tfprotov5.GetProviderSchemaResponse{
						ProviderMeta: &tfprotov5.Schema{
							Block: &tfprotov5.SchemaBlock{
								Attributes: []*tfprotov5.SchemaAttribute{
									{
										Name:     "testattribute1",
										Type:     tftypes.String,
										Required: true,
									},
								},
							},
						},
					},
				}).ProviderServer,
				(&tf5testserver.TestServer{
					GetProviderSchemaResponse: &tfprotov5.GetProviderSchemaResponse{
						ProviderMeta: &tfprotov5.Schema{
							Block: &tfprotov5.SchemaBlock{
								Attributes: []*tfprotov5.SchemaAttribute{
									{
										Name:     "testattribute2",
										Type:     tftypes.String,
										Required: true,
									},
								},
							},
						},
					},
				}).ProviderServer,
			},
			expectedDataSourceSchemas: map[string]*tfprotov5.Schema{},
			expectedDiagnostics: []*tfprotov5.Diagnostic{
				{
					Severity: tfprotov5.DiagnosticSeverityError,
					Summary:  "Invalid Provider Server Combination",
					Detail: "The combined provider has differing provider meta schema implementations across providers. " +
						"Provider meta schemas must be identical across providers. " +
						"This is always an issue in the provider implementation and should be reported to the provider developers.\n\n" +
						"Provider meta schema difference: " + cmp.Diff(
						&tfprotov5.Schema{
							Block: &tfprotov5.SchemaBlock{
								Attributes: []*tfprotov5.SchemaAttribute{
									{
										Name:     "testattribute2",
										Type:     tftypes.String,
										Required: true,
									},
								},
							},
						},
						&tfprotov5.Schema{
							Block: &tfprotov5.SchemaBlock{
								Attributes: []*tfprotov5.SchemaAttribute{
									{
										Name:     "testattribute1",
										Type:     tftypes.String,
										Required: true,
									},
								},
							},
						},
					),
				},
			},
			expectedFunctions: map[string]*tfprotov5.Function{},
			expectedProviderMetaSchema: &tfprotov5.Schema{
				Block: &tfprotov5.SchemaBlock{
					Attributes: []*tfprotov5.SchemaAttribute{
						{
							Name:     "testattribute1",
							Type:     tftypes.String,
							Required: true,
						},
					},
				},
			},
			expectedResourceSchemas: map[string]*tfprotov5.Schema{},
			expectedServerCapabilities: &tfprotov5.ServerCapabilities{
				GetProviderSchemaOptional: true,
				MoveResourceState:         true,
				PlanDestroy:               true,
			},
		},
		"server-capabilities": {
			servers: []func() tfprotov5.ProviderServer{
				(&tf5testserver.TestServer{
					GetProviderSchemaResponse: &tfprotov5.GetProviderSchemaResponse{
						ResourceSchemas: map[string]*tfprotov5.Schema{
							"test_with_server_capabilities": {},
						},
						ServerCapabilities: &tfprotov5.ServerCapabilities{
							GetProviderSchemaOptional: true,
							PlanDestroy:               true,
						},
					},
				}).ProviderServer,
				(&tf5testserver.TestServer{
					GetProviderSchemaResponse: &tfprotov5.GetProviderSchemaResponse{
						ResourceSchemas: map[string]*tfprotov5.Schema{
							"test_without_server_capabilities": {},
						},
					},
				}).ProviderServer,
			},
			expectedDataSourceSchemas: map[string]*tfprotov5.Schema{},
			expectedFunctions:         map[string]*tfprotov5.Function{},
			expectedResourceSchemas: map[string]*tfprotov5.Schema{
				"test_with_server_capabilities":    {},
				"test_without_server_capabilities": {},
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
					GetProviderSchemaResponse: &tfprotov5.GetProviderSchemaResponse{
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
			expectedDataSourceSchemas: map[string]*tfprotov5.Schema{},
			expectedDiagnostics: []*tfprotov5.Diagnostic{
				{
					Severity: tfprotov5.DiagnosticSeverityError,
					Summary:  "test error summary",
					Detail:   "test error details",
				},
			},
			expectedFunctions:       map[string]*tfprotov5.Function{},
			expectedResourceSchemas: map[string]*tfprotov5.Schema{},
			expectedServerCapabilities: &tfprotov5.ServerCapabilities{
				GetProviderSchemaOptional: true,
				MoveResourceState:         true,
				PlanDestroy:               true,
			},
		},
		"error-multiple": {
			servers: []func() tfprotov5.ProviderServer{
				(&tf5testserver.TestServer{
					GetProviderSchemaResponse: &tfprotov5.GetProviderSchemaResponse{
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
					GetProviderSchemaResponse: &tfprotov5.GetProviderSchemaResponse{
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
			expectedDataSourceSchemas: map[string]*tfprotov5.Schema{},
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
			expectedFunctions:       map[string]*tfprotov5.Function{},
			expectedResourceSchemas: map[string]*tfprotov5.Schema{},
			expectedServerCapabilities: &tfprotov5.ServerCapabilities{
				GetProviderSchemaOptional: true,
				MoveResourceState:         true,
				PlanDestroy:               true,
			},
		},
		"warning-once": {
			servers: []func() tfprotov5.ProviderServer{
				(&tf5testserver.TestServer{
					GetProviderSchemaResponse: &tfprotov5.GetProviderSchemaResponse{
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
			expectedDataSourceSchemas: map[string]*tfprotov5.Schema{},
			expectedDiagnostics: []*tfprotov5.Diagnostic{
				{
					Severity: tfprotov5.DiagnosticSeverityWarning,
					Summary:  "test warning summary",
					Detail:   "test warning details",
				},
			},
			expectedFunctions:       map[string]*tfprotov5.Function{},
			expectedResourceSchemas: map[string]*tfprotov5.Schema{},
			expectedServerCapabilities: &tfprotov5.ServerCapabilities{
				GetProviderSchemaOptional: true,
				MoveResourceState:         true,
				PlanDestroy:               true,
			},
		},
		"warning-multiple": {
			servers: []func() tfprotov5.ProviderServer{
				(&tf5testserver.TestServer{
					GetProviderSchemaResponse: &tfprotov5.GetProviderSchemaResponse{
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
					GetProviderSchemaResponse: &tfprotov5.GetProviderSchemaResponse{
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
			expectedDataSourceSchemas: map[string]*tfprotov5.Schema{},
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
			expectedFunctions:       map[string]*tfprotov5.Function{},
			expectedResourceSchemas: map[string]*tfprotov5.Schema{},
			expectedServerCapabilities: &tfprotov5.ServerCapabilities{
				GetProviderSchemaOptional: true,
				MoveResourceState:         true,
				PlanDestroy:               true,
			},
		},
		"warning-then-error": {
			servers: []func() tfprotov5.ProviderServer{
				(&tf5testserver.TestServer{
					GetProviderSchemaResponse: &tfprotov5.GetProviderSchemaResponse{
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
					GetProviderSchemaResponse: &tfprotov5.GetProviderSchemaResponse{
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
			expectedDataSourceSchemas: map[string]*tfprotov5.Schema{},
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
			expectedFunctions:       map[string]*tfprotov5.Function{},
			expectedResourceSchemas: map[string]*tfprotov5.Schema{},
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

			resp, err := muxServer.ProviderServer().GetProviderSchema(context.Background(), &tfprotov5.GetProviderSchemaRequest{})

			if err != nil {
				t.Fatalf("unexpected error: %s", err)
			}

			if diff := cmp.Diff(resp.DataSourceSchemas, testCase.expectedDataSourceSchemas); diff != "" {
				t.Errorf("data source schemas didn't match expectations: %s", diff)
			}

			if diff := cmp.Diff(resp.Diagnostics, testCase.expectedDiagnostics); diff != "" {
				t.Errorf("diagnostics didn't match expectations: %s", diff)
			}

			if diff := cmp.Diff(resp.Functions, testCase.expectedFunctions); diff != "" {
				t.Errorf("functions didn't match expectations: %s", diff)
			}

			if diff := cmp.Diff(resp.Provider, testCase.expectedProviderSchema); diff != "" {
				t.Errorf("provider schema didn't match expectations: %s", diff)
			}

			if diff := cmp.Diff(resp.ProviderMeta, testCase.expectedProviderMetaSchema); diff != "" {
				t.Errorf("provider_meta schema didn't match expectations: %s", diff)
			}

			if diff := cmp.Diff(resp.ResourceSchemas, testCase.expectedResourceSchemas); diff != "" {
				t.Errorf("resource schemas didn't match expectations: %s", diff)
			}

			if diff := cmp.Diff(resp.ServerCapabilities, testCase.expectedServerCapabilities); diff != "" {
				t.Errorf("server capabilities didn't match expectations: %s", diff)
			}
		})
	}
}
