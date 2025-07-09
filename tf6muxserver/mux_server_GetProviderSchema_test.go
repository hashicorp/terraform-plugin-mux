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

func TestMuxServerGetProviderSchema(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		servers                           []func() tfprotov6.ProviderServer
		expectedActionSchemas             map[string]*tfprotov6.ActionSchema
		expectedDataSourceSchemas         map[string]*tfprotov6.Schema
		expectedDiagnostics               []*tfprotov6.Diagnostic
		expectedEphemeralResourcesSchemas map[string]*tfprotov6.Schema
		expectedListResourcesSchemas      map[string]*tfprotov6.Schema
		expectedFunctions                 map[string]*tfprotov6.Function
		expectedProviderSchema            *tfprotov6.Schema
		expectedProviderMetaSchema        *tfprotov6.Schema
		expectedResourceSchemas           map[string]*tfprotov6.Schema
		expectedServerCapabilities        *tfprotov6.ServerCapabilities
	}{
		"combined": {
			servers: []func() tfprotov6.ProviderServer{
				(&tf6testserver.TestServer{
					GetProviderSchemaResponse: &tfprotov6.GetProviderSchemaResponse{
						Provider: &tfprotov6.Schema{
							Version: 1,
							Block: &tfprotov6.SchemaBlock{
								Version: 1,
								Attributes: []*tfprotov6.SchemaAttribute{
									{
										Name:            "account_id",
										Type:            tftypes.String,
										Required:        true,
										Description:     "the account ID to make requests for",
										DescriptionKind: tfprotov6.StringKindPlain,
									},
								},
								BlockTypes: []*tfprotov6.SchemaNestedBlock{
									{
										TypeName: "feature",
										Nesting:  tfprotov6.SchemaNestedBlockNestingModeList,
										Block: &tfprotov6.SchemaBlock{
											Version:         1,
											Description:     "features to enable on the provider",
											DescriptionKind: tfprotov6.StringKindPlain,
											Attributes: []*tfprotov6.SchemaAttribute{
												{
													Name:            "feature_id",
													Type:            tftypes.Number,
													Required:        true,
													Description:     "The ID of the feature",
													DescriptionKind: tfprotov6.StringKindPlain,
												},
												{
													Name:            "enabled",
													Type:            tftypes.Bool,
													Optional:        true,
													Description:     "whether the feature is enabled",
													DescriptionKind: tfprotov6.StringKindPlain,
												},
											},
										},
									},
								},
							},
						},
						ProviderMeta: &tfprotov6.Schema{
							Version: 4,
							Block: &tfprotov6.SchemaBlock{
								Version: 4,
								Attributes: []*tfprotov6.SchemaAttribute{
									{
										Name:            "module_id",
										Type:            tftypes.String,
										Optional:        true,
										Description:     "a unique identifier for the module",
										DescriptionKind: tfprotov6.StringKindPlain,
									},
								},
							},
						},
						ResourceSchemas: map[string]*tfprotov6.Schema{
							"test_foo": {
								Version: 1,
								Block: &tfprotov6.SchemaBlock{
									Version: 1,
									Attributes: []*tfprotov6.SchemaAttribute{
										{
											Name:            "airspeed_velocity",
											Type:            tftypes.Number,
											Required:        true,
											Description:     "the airspeed velocity of a swallow",
											DescriptionKind: tfprotov6.StringKindPlain,
										},
									},
								},
							},
							"test_bar": {
								Version: 1,
								Block: &tfprotov6.SchemaBlock{
									Version: 1,
									Attributes: []*tfprotov6.SchemaAttribute{
										{
											Name:            "name",
											Type:            tftypes.String,
											Optional:        true,
											Description:     "your name",
											DescriptionKind: tfprotov6.StringKindPlain,
										},
										{
											Name:            "color",
											Type:            tftypes.String,
											Optional:        true,
											Description:     "your favorite color",
											DescriptionKind: tfprotov6.StringKindPlain,
										},
									},
								},
							},
						},
						ActionSchemas: map[string]*tfprotov6.ActionSchema{
							"test_foo": {
								Schema: &tfprotov6.Schema{
									Version: 1,
									Block: &tfprotov6.SchemaBlock{
										Version: 1,
										Attributes: []*tfprotov6.SchemaAttribute{
											{
												Name:            "current_time",
												Type:            tftypes.String,
												Computed:        true,
												Description:     "the current time in RFC 3339 format",
												DescriptionKind: tfprotov6.StringKindPlain,
											},
										},
									},
								},
								Type: tfprotov6.UnlinkedActionSchemaType{},
							},
						},
						DataSourceSchemas: map[string]*tfprotov6.Schema{
							"test_foo": {
								Version: 1,
								Block: &tfprotov6.SchemaBlock{
									Version: 1,
									Attributes: []*tfprotov6.SchemaAttribute{
										{
											Name:            "current_time",
											Type:            tftypes.String,
											Computed:        true,
											Description:     "the current time in RFC 3339 format",
											DescriptionKind: tfprotov6.StringKindPlain,
										},
									},
								},
							},
						},
						Functions: map[string]*tfprotov6.Function{
							"test_function1": {
								Return: &tfprotov6.FunctionReturn{
									Type: tftypes.String,
								},
							},
						},
						EphemeralResourceSchemas: map[string]*tfprotov6.Schema{
							"test_ephemeral_foo": {
								Version: 1,
								Block: &tfprotov6.SchemaBlock{
									Version: 1,
									Attributes: []*tfprotov6.SchemaAttribute{
										{
											Name:            "secret_number",
											Type:            tftypes.Number,
											Required:        true,
											Description:     "input the secret number",
											DescriptionKind: tfprotov6.StringKindPlain,
										},
									},
								},
							},
							"test_ephemeral_bar": {
								Version: 1,
								Block: &tfprotov6.SchemaBlock{
									Version: 1,
									Attributes: []*tfprotov6.SchemaAttribute{
										{
											Name:            "username",
											Type:            tftypes.String,
											Optional:        true,
											Description:     "your username",
											DescriptionKind: tfprotov6.StringKindPlain,
										},
										{
											Name:            "password",
											Type:            tftypes.String,
											Optional:        true,
											Description:     "your password",
											DescriptionKind: tfprotov6.StringKindPlain,
										},
									},
								},
							},
						},
						ListResourceSchemas: map[string]*tfprotov6.Schema{
							"test_list_foo": {
								Version: 1,
								Block: &tfprotov6.SchemaBlock{
									Version: 1,
									Attributes: []*tfprotov6.SchemaAttribute{
										{
											Name:            "query_name",
											Type:            tftypes.String,
											Required:        true,
											Description:     "input the query name",
											DescriptionKind: tfprotov6.StringKindPlain,
										},
									},
								},
							},
							"test_list_bar": {
								Version: 1,
								Block: &tfprotov6.SchemaBlock{
									Version: 1,
									Attributes: []*tfprotov6.SchemaAttribute{
										{
											Name:            "filter",
											Type:            tftypes.String,
											Optional:        true,
											Description:     "search filter",
											DescriptionKind: tfprotov6.StringKindPlain,
										},
										{
											Name:            "prefix",
											Type:            tftypes.String,
											Optional:        true,
											Description:     "name prefix",
											DescriptionKind: tfprotov6.StringKindPlain,
										},
									},
								},
							},
						},
					},
				}).ProviderServer,
				(&tf6testserver.TestServer{
					GetProviderSchemaResponse: &tfprotov6.GetProviderSchemaResponse{
						Provider: &tfprotov6.Schema{
							Version: 1,
							Block: &tfprotov6.SchemaBlock{
								Version: 1,
								Attributes: []*tfprotov6.SchemaAttribute{
									{
										Name:            "account_id",
										Type:            tftypes.String,
										Required:        true,
										Description:     "the account ID to make requests for",
										DescriptionKind: tfprotov6.StringKindPlain,
									},
								},
								BlockTypes: []*tfprotov6.SchemaNestedBlock{
									{
										TypeName: "feature",
										Nesting:  tfprotov6.SchemaNestedBlockNestingModeList,
										Block: &tfprotov6.SchemaBlock{
											Version:         1,
											Description:     "features to enable on the provider",
											DescriptionKind: tfprotov6.StringKindPlain,
											Attributes: []*tfprotov6.SchemaAttribute{
												{
													Name:            "feature_id",
													Type:            tftypes.Number,
													Required:        true,
													Description:     "The ID of the feature",
													DescriptionKind: tfprotov6.StringKindPlain,
												},
												{
													Name:            "enabled",
													Type:            tftypes.Bool,
													Optional:        true,
													Description:     "whether the feature is enabled",
													DescriptionKind: tfprotov6.StringKindPlain,
												},
											},
										},
									},
								},
							},
						},
						ProviderMeta: &tfprotov6.Schema{
							Version: 4,
							Block: &tfprotov6.SchemaBlock{
								Version: 4,
								Attributes: []*tfprotov6.SchemaAttribute{
									{
										Name:            "module_id",
										Type:            tftypes.String,
										Optional:        true,
										Description:     "a unique identifier for the module",
										DescriptionKind: tfprotov6.StringKindPlain,
									},
								},
							},
						},
						ResourceSchemas: map[string]*tfprotov6.Schema{
							"test_quux": {
								Version: 1,
								Block: &tfprotov6.SchemaBlock{
									Version: 1,
									Attributes: []*tfprotov6.SchemaAttribute{
										{
											Name:            "a",
											Type:            tftypes.String,
											Required:        true,
											Description:     "the account ID to make requests for",
											DescriptionKind: tfprotov6.StringKindPlain,
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
						ActionSchemas: map[string]*tfprotov6.ActionSchema{
							"test_bar": {
								Schema: &tfprotov6.Schema{
									Version: 1,
									Block: &tfprotov6.SchemaBlock{
										Version: 1,
										Attributes: []*tfprotov6.SchemaAttribute{
											{
												Name:            "a",
												Type:            tftypes.Number,
												Computed:        true,
												Description:     "some field that's set by the provider",
												DescriptionKind: tfprotov6.StringKindMarkdown,
											},
										},
									},
								},
								Type: tfprotov6.UnlinkedActionSchemaType{},
							},
							"test_quux": {
								Schema: &tfprotov6.Schema{
									Version: 1,
									Block: &tfprotov6.SchemaBlock{
										Version: 1,
										Attributes: []*tfprotov6.SchemaAttribute{
											{
												Name:            "abc",
												Type:            tftypes.Number,
												Computed:        true,
												Description:     "some other field that's set by the provider",
												DescriptionKind: tfprotov6.StringKindMarkdown,
											},
										},
									},
								},
								Type: tfprotov6.UnlinkedActionSchemaType{},
							},
						},
						DataSourceSchemas: map[string]*tfprotov6.Schema{
							"test_bar": {
								Version: 1,
								Block: &tfprotov6.SchemaBlock{
									Version: 1,
									Attributes: []*tfprotov6.SchemaAttribute{
										{
											Name:            "a",
											Type:            tftypes.Number,
											Computed:        true,
											Description:     "some field that's set by the provider",
											DescriptionKind: tfprotov6.StringKindMarkdown,
										},
									},
								},
							},
							"test_quux": {
								Version: 1,
								Block: &tfprotov6.SchemaBlock{
									Version: 1,
									Attributes: []*tfprotov6.SchemaAttribute{
										{
											Name:            "abc",
											Type:            tftypes.Number,
											Computed:        true,
											Description:     "some other field that's set by the provider",
											DescriptionKind: tfprotov6.StringKindMarkdown,
										},
									},
								},
							},
						},
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
						EphemeralResourceSchemas: map[string]*tfprotov6.Schema{
							"test_ephemeral_foobar": {
								Version: 1,
								Block: &tfprotov6.SchemaBlock{
									Version: 1,
									Attributes: []*tfprotov6.SchemaAttribute{
										{
											Name:            "secret_number",
											Type:            tftypes.Number,
											Computed:        true,
											Description:     "A generated secret number",
											DescriptionKind: tfprotov6.StringKindPlain,
										},
									},
								},
							},
						},
						ListResourceSchemas: map[string]*tfprotov6.Schema{
							"test_list_foobar": {
								Version: 1,
								Block: &tfprotov6.SchemaBlock{
									Version: 1,
									Attributes: []*tfprotov6.SchemaAttribute{
										{
											Name:            "query_name",
											Type:            tftypes.String,
											Computed:        true,
											Description:     "A generated query name",
											DescriptionKind: tfprotov6.StringKindPlain,
										},
									},
								},
							},
						},
					},
				}).ProviderServer,
			},
			expectedProviderSchema: &tfprotov6.Schema{
				Version: 1,
				Block: &tfprotov6.SchemaBlock{
					Version: 1,
					Attributes: []*tfprotov6.SchemaAttribute{
						{
							Name:            "account_id",
							Type:            tftypes.String,
							Required:        true,
							Description:     "the account ID to make requests for",
							DescriptionKind: tfprotov6.StringKindPlain,
						},
					},
					BlockTypes: []*tfprotov6.SchemaNestedBlock{
						{
							TypeName: "feature",
							Nesting:  tfprotov6.SchemaNestedBlockNestingModeList,
							Block: &tfprotov6.SchemaBlock{
								Version:         1,
								Description:     "features to enable on the provider",
								DescriptionKind: tfprotov6.StringKindPlain,
								Attributes: []*tfprotov6.SchemaAttribute{
									{
										Name:            "feature_id",
										Type:            tftypes.Number,
										Required:        true,
										Description:     "The ID of the feature",
										DescriptionKind: tfprotov6.StringKindPlain,
									},
									{
										Name:            "enabled",
										Type:            tftypes.Bool,
										Optional:        true,
										Description:     "whether the feature is enabled",
										DescriptionKind: tfprotov6.StringKindPlain,
									},
								},
							},
						},
					},
				},
			},
			expectedProviderMetaSchema: &tfprotov6.Schema{
				Version: 4,
				Block: &tfprotov6.SchemaBlock{
					Version: 4,
					Attributes: []*tfprotov6.SchemaAttribute{
						{
							Name:            "module_id",
							Type:            tftypes.String,
							Optional:        true,
							Description:     "a unique identifier for the module",
							DescriptionKind: tfprotov6.StringKindPlain,
						},
					},
				},
			},
			expectedResourceSchemas: map[string]*tfprotov6.Schema{
				"test_foo": {
					Version: 1,
					Block: &tfprotov6.SchemaBlock{
						Version: 1,
						Attributes: []*tfprotov6.SchemaAttribute{
							{
								Name:            "airspeed_velocity",
								Type:            tftypes.Number,
								Required:        true,
								Description:     "the airspeed velocity of a swallow",
								DescriptionKind: tfprotov6.StringKindPlain,
							},
						},
					},
				},
				"test_bar": {
					Version: 1,
					Block: &tfprotov6.SchemaBlock{
						Version: 1,
						Attributes: []*tfprotov6.SchemaAttribute{
							{
								Name:            "name",
								Type:            tftypes.String,
								Optional:        true,
								Description:     "your name",
								DescriptionKind: tfprotov6.StringKindPlain,
							},
							{
								Name:            "color",
								Type:            tftypes.String,
								Optional:        true,
								Description:     "your favorite color",
								DescriptionKind: tfprotov6.StringKindPlain,
							},
						},
					},
				},
				"test_quux": {
					Version: 1,
					Block: &tfprotov6.SchemaBlock{
						Version: 1,
						Attributes: []*tfprotov6.SchemaAttribute{
							{
								Name:            "a",
								Type:            tftypes.String,
								Required:        true,
								Description:     "the account ID to make requests for",
								DescriptionKind: tfprotov6.StringKindPlain,
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
			expectedActionSchemas: map[string]*tfprotov6.ActionSchema{
				"test_foo": {
					Schema: &tfprotov6.Schema{
						Version: 1,
						Block: &tfprotov6.SchemaBlock{
							Version: 1,
							Attributes: []*tfprotov6.SchemaAttribute{
								{
									Name:            "current_time",
									Type:            tftypes.String,
									Computed:        true,
									Description:     "the current time in RFC 3339 format",
									DescriptionKind: tfprotov6.StringKindPlain,
								},
							},
						},
					},
					Type: tfprotov6.UnlinkedActionSchemaType{},
				},
				"test_bar": {
					Schema: &tfprotov6.Schema{
						Version: 1,
						Block: &tfprotov6.SchemaBlock{
							Version: 1,
							Attributes: []*tfprotov6.SchemaAttribute{
								{
									Name:            "a",
									Type:            tftypes.Number,
									Computed:        true,
									Description:     "some field that's set by the provider",
									DescriptionKind: tfprotov6.StringKindMarkdown,
								},
							},
						},
					},
					Type: tfprotov6.UnlinkedActionSchemaType{},
				},
				"test_quux": {
					Schema: &tfprotov6.Schema{
						Version: 1,
						Block: &tfprotov6.SchemaBlock{
							Version: 1,
							Attributes: []*tfprotov6.SchemaAttribute{
								{
									Name:            "abc",
									Type:            tftypes.Number,
									Computed:        true,
									Description:     "some other field that's set by the provider",
									DescriptionKind: tfprotov6.StringKindMarkdown,
								},
							},
						},
					},
					Type: tfprotov6.UnlinkedActionSchemaType{},
				},
			},
			expectedDataSourceSchemas: map[string]*tfprotov6.Schema{
				"test_foo": {
					Version: 1,
					Block: &tfprotov6.SchemaBlock{
						Version: 1,
						Attributes: []*tfprotov6.SchemaAttribute{
							{
								Name:            "current_time",
								Type:            tftypes.String,
								Computed:        true,
								Description:     "the current time in RFC 3339 format",
								DescriptionKind: tfprotov6.StringKindPlain,
							},
						},
					},
				},
				"test_bar": {
					Version: 1,
					Block: &tfprotov6.SchemaBlock{
						Version: 1,
						Attributes: []*tfprotov6.SchemaAttribute{
							{
								Name:            "a",
								Type:            tftypes.Number,
								Computed:        true,
								Description:     "some field that's set by the provider",
								DescriptionKind: tfprotov6.StringKindMarkdown,
							},
						},
					},
				},
				"test_quux": {
					Version: 1,
					Block: &tfprotov6.SchemaBlock{
						Version: 1,
						Attributes: []*tfprotov6.SchemaAttribute{
							{
								Name:            "abc",
								Type:            tftypes.Number,
								Computed:        true,
								Description:     "some other field that's set by the provider",
								DescriptionKind: tfprotov6.StringKindMarkdown,
							},
						},
					},
				},
			},
			expectedFunctions: map[string]*tfprotov6.Function{
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
			expectedEphemeralResourcesSchemas: map[string]*tfprotov6.Schema{
				"test_ephemeral_foo": {
					Version: 1,
					Block: &tfprotov6.SchemaBlock{
						Version: 1,
						Attributes: []*tfprotov6.SchemaAttribute{
							{
								Name:            "secret_number",
								Type:            tftypes.Number,
								Required:        true,
								Description:     "input the secret number",
								DescriptionKind: tfprotov6.StringKindPlain,
							},
						},
					},
				},
				"test_ephemeral_bar": {
					Version: 1,
					Block: &tfprotov6.SchemaBlock{
						Version: 1,
						Attributes: []*tfprotov6.SchemaAttribute{
							{
								Name:            "username",
								Type:            tftypes.String,
								Optional:        true,
								Description:     "your username",
								DescriptionKind: tfprotov6.StringKindPlain,
							},
							{
								Name:            "password",
								Type:            tftypes.String,
								Optional:        true,
								Description:     "your password",
								DescriptionKind: tfprotov6.StringKindPlain,
							},
						},
					},
				},
				"test_ephemeral_foobar": {
					Version: 1,
					Block: &tfprotov6.SchemaBlock{
						Version: 1,
						Attributes: []*tfprotov6.SchemaAttribute{
							{
								Name:            "secret_number",
								Type:            tftypes.Number,
								Computed:        true,
								Description:     "A generated secret number",
								DescriptionKind: tfprotov6.StringKindPlain,
							},
						},
					},
				},
			},
			expectedListResourcesSchemas: map[string]*tfprotov6.Schema{
				"test_list_foo": {
					Version: 1,
					Block: &tfprotov6.SchemaBlock{
						Version: 1,
						Attributes: []*tfprotov6.SchemaAttribute{
							{
								Name:            "query_name",
								Type:            tftypes.String,
								Required:        true,
								Description:     "input the query name",
								DescriptionKind: tfprotov6.StringKindPlain,
							},
						},
					},
				},
				"test_list_bar": {
					Version: 1,
					Block: &tfprotov6.SchemaBlock{
						Version: 1,
						Attributes: []*tfprotov6.SchemaAttribute{
							{
								Name:            "filter",
								Type:            tftypes.String,
								Optional:        true,
								Description:     "search filter",
								DescriptionKind: tfprotov6.StringKindPlain,
							},
							{
								Name:            "prefix",
								Type:            tftypes.String,
								Optional:        true,
								Description:     "name prefix",
								DescriptionKind: tfprotov6.StringKindPlain,
							},
						},
					},
				},
				"test_list_foobar": {
					Version: 1,
					Block: &tfprotov6.SchemaBlock{
						Version: 1,
						Attributes: []*tfprotov6.SchemaAttribute{
							{
								Name:            "query_name",
								Type:            tftypes.String,
								Computed:        true,
								Description:     "A generated query name",
								DescriptionKind: tfprotov6.StringKindPlain,
							},
						},
					},
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
					GetProviderSchemaResponse: &tfprotov6.GetProviderSchemaResponse{
						ActionSchemas: map[string]*tfprotov6.ActionSchema{
							"test_foo": {},
						},
					},
				}).ProviderServer,
				(&tf6testserver.TestServer{
					GetProviderSchemaResponse: &tfprotov6.GetProviderSchemaResponse{
						ActionSchemas: map[string]*tfprotov6.ActionSchema{
							"test_foo": {},
						},
					},
				}).ProviderServer,
			},
			expectedActionSchemas: map[string]*tfprotov6.ActionSchema{
				"test_foo": {},
			},
			expectedDataSourceSchemas: map[string]*tfprotov6.Schema{},
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
			expectedEphemeralResourcesSchemas: map[string]*tfprotov6.Schema{},
			expectedListResourcesSchemas:      map[string]*tfprotov6.Schema{},
			expectedFunctions:                 map[string]*tfprotov6.Function{},
			expectedResourceSchemas:           map[string]*tfprotov6.Schema{},
			expectedServerCapabilities: &tfprotov6.ServerCapabilities{
				GetProviderSchemaOptional: true,
				MoveResourceState:         true,
				PlanDestroy:               true,
			},
		},
		"duplicate-data-source-type": {
			servers: []func() tfprotov6.ProviderServer{
				(&tf6testserver.TestServer{
					GetProviderSchemaResponse: &tfprotov6.GetProviderSchemaResponse{
						DataSourceSchemas: map[string]*tfprotov6.Schema{
							"test_foo": {},
						},
					},
				}).ProviderServer,
				(&tf6testserver.TestServer{
					GetProviderSchemaResponse: &tfprotov6.GetProviderSchemaResponse{
						DataSourceSchemas: map[string]*tfprotov6.Schema{
							"test_foo": {},
						},
					},
				}).ProviderServer,
			},
			expectedActionSchemas: map[string]*tfprotov6.ActionSchema{},
			expectedDataSourceSchemas: map[string]*tfprotov6.Schema{
				"test_foo": {},
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
			expectedEphemeralResourcesSchemas: map[string]*tfprotov6.Schema{},
			expectedListResourcesSchemas:      map[string]*tfprotov6.Schema{},
			expectedFunctions:                 map[string]*tfprotov6.Function{},
			expectedResourceSchemas:           map[string]*tfprotov6.Schema{},
			expectedServerCapabilities: &tfprotov6.ServerCapabilities{
				GetProviderSchemaOptional: true,
				MoveResourceState:         true,
				PlanDestroy:               true,
			},
		},
		"duplicate-ephemeral-resource-type": {
			servers: []func() tfprotov6.ProviderServer{
				(&tf6testserver.TestServer{
					GetProviderSchemaResponse: &tfprotov6.GetProviderSchemaResponse{
						EphemeralResourceSchemas: map[string]*tfprotov6.Schema{
							"test_foo": {},
						},
					},
				}).ProviderServer,
				(&tf6testserver.TestServer{
					GetProviderSchemaResponse: &tfprotov6.GetProviderSchemaResponse{
						EphemeralResourceSchemas: map[string]*tfprotov6.Schema{
							"test_foo": {},
						},
					},
				}).ProviderServer,
			},
			expectedActionSchemas:     map[string]*tfprotov6.ActionSchema{},
			expectedDataSourceSchemas: map[string]*tfprotov6.Schema{},
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
			expectedEphemeralResourcesSchemas: map[string]*tfprotov6.Schema{
				"test_foo": {},
			},
			expectedListResourcesSchemas: map[string]*tfprotov6.Schema{},
			expectedFunctions:            map[string]*tfprotov6.Function{},
			expectedResourceSchemas:      map[string]*tfprotov6.Schema{},
			expectedServerCapabilities: &tfprotov6.ServerCapabilities{
				GetProviderSchemaOptional: true,
				MoveResourceState:         true,
				PlanDestroy:               true,
			},
		},
		"duplicate-list-resource-type": {
			servers: []func() tfprotov6.ProviderServer{
				(&tf6testserver.TestServer{
					GetProviderSchemaResponse: &tfprotov6.GetProviderSchemaResponse{
						ListResourceSchemas: map[string]*tfprotov6.Schema{
							"test_foo": {},
						},
					},
				}).ProviderServer,
				(&tf6testserver.TestServer{
					GetProviderSchemaResponse: &tfprotov6.GetProviderSchemaResponse{
						ListResourceSchemas: map[string]*tfprotov6.Schema{
							"test_foo": {},
						},
					},
				}).ProviderServer,
			},
			expectedActionSchemas:     map[string]*tfprotov6.ActionSchema{},
			expectedDataSourceSchemas: map[string]*tfprotov6.Schema{},
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
			expectedListResourcesSchemas: map[string]*tfprotov6.Schema{
				"test_foo": {},
			},
			expectedEphemeralResourcesSchemas: map[string]*tfprotov6.Schema{},
			expectedFunctions:                 map[string]*tfprotov6.Function{},
			expectedResourceSchemas:           map[string]*tfprotov6.Schema{},
			expectedServerCapabilities: &tfprotov6.ServerCapabilities{
				GetProviderSchemaOptional: true,
				MoveResourceState:         true,
				PlanDestroy:               true,
			},
		},
		"duplicate-function": {
			servers: []func() tfprotov6.ProviderServer{
				(&tf6testserver.TestServer{
					GetProviderSchemaResponse: &tfprotov6.GetProviderSchemaResponse{
						Functions: map[string]*tfprotov6.Function{
							"test_function": {},
						},
					},
				}).ProviderServer,
				(&tf6testserver.TestServer{
					GetProviderSchemaResponse: &tfprotov6.GetProviderSchemaResponse{
						Functions: map[string]*tfprotov6.Function{
							"test_function": {},
						},
					},
				}).ProviderServer,
			},
			expectedActionSchemas:     map[string]*tfprotov6.ActionSchema{},
			expectedDataSourceSchemas: map[string]*tfprotov6.Schema{},
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
			expectedEphemeralResourcesSchemas: map[string]*tfprotov6.Schema{},
			expectedListResourcesSchemas:      map[string]*tfprotov6.Schema{},
			expectedFunctions: map[string]*tfprotov6.Function{
				"test_function": {},
			},
			expectedResourceSchemas: map[string]*tfprotov6.Schema{},
			expectedServerCapabilities: &tfprotov6.ServerCapabilities{
				GetProviderSchemaOptional: true,
				MoveResourceState:         true,
				PlanDestroy:               true,
			},
		},
		"duplicate-resource-type": {
			servers: []func() tfprotov6.ProviderServer{
				(&tf6testserver.TestServer{
					GetProviderSchemaResponse: &tfprotov6.GetProviderSchemaResponse{
						ResourceSchemas: map[string]*tfprotov6.Schema{
							"test_foo": {},
						},
					},
				}).ProviderServer,
				(&tf6testserver.TestServer{
					GetProviderSchemaResponse: &tfprotov6.GetProviderSchemaResponse{
						ResourceSchemas: map[string]*tfprotov6.Schema{
							"test_foo": {},
						},
					},
				}).ProviderServer,
			},
			expectedActionSchemas:     map[string]*tfprotov6.ActionSchema{},
			expectedDataSourceSchemas: map[string]*tfprotov6.Schema{},
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
			expectedEphemeralResourcesSchemas: map[string]*tfprotov6.Schema{},
			expectedListResourcesSchemas:      map[string]*tfprotov6.Schema{},
			expectedFunctions:                 map[string]*tfprotov6.Function{},
			expectedResourceSchemas: map[string]*tfprotov6.Schema{
				"test_foo": {},
			},
			expectedServerCapabilities: &tfprotov6.ServerCapabilities{
				GetProviderSchemaOptional: true,
				MoveResourceState:         true,
				PlanDestroy:               true,
			},
		},
		"provider-mismatch": {
			servers: []func() tfprotov6.ProviderServer{
				(&tf6testserver.TestServer{
					GetProviderSchemaResponse: &tfprotov6.GetProviderSchemaResponse{
						Provider: &tfprotov6.Schema{
							Block: &tfprotov6.SchemaBlock{
								Attributes: []*tfprotov6.SchemaAttribute{
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
				(&tf6testserver.TestServer{
					GetProviderSchemaResponse: &tfprotov6.GetProviderSchemaResponse{
						Provider: &tfprotov6.Schema{
							Block: &tfprotov6.SchemaBlock{
								Attributes: []*tfprotov6.SchemaAttribute{
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
			expectedActionSchemas:     map[string]*tfprotov6.ActionSchema{},
			expectedDataSourceSchemas: map[string]*tfprotov6.Schema{},
			expectedDiagnostics: []*tfprotov6.Diagnostic{
				{
					Severity: tfprotov6.DiagnosticSeverityError,
					Summary:  "Invalid Provider Server Combination",
					Detail: "The combined provider has differing provider schema implementations across providers. " +
						"Provider schemas must be identical across providers. " +
						"This is always an issue in the provider implementation and should be reported to the provider developers.\n\n" +
						"Provider schema difference: " + cmp.Diff(
						&tfprotov6.Schema{
							Block: &tfprotov6.SchemaBlock{
								Attributes: []*tfprotov6.SchemaAttribute{
									{
										Name:     "testattribute2",
										Type:     tftypes.String,
										Required: true,
									},
								},
							},
						},
						&tfprotov6.Schema{
							Block: &tfprotov6.SchemaBlock{
								Attributes: []*tfprotov6.SchemaAttribute{
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
			expectedEphemeralResourcesSchemas: map[string]*tfprotov6.Schema{},
			expectedListResourcesSchemas:      map[string]*tfprotov6.Schema{},
			expectedFunctions:                 map[string]*tfprotov6.Function{},
			expectedProviderSchema: &tfprotov6.Schema{
				Block: &tfprotov6.SchemaBlock{
					Attributes: []*tfprotov6.SchemaAttribute{
						{
							Name:     "testattribute1",
							Type:     tftypes.String,
							Required: true,
						},
					},
				},
			},
			expectedResourceSchemas: map[string]*tfprotov6.Schema{},
			expectedServerCapabilities: &tfprotov6.ServerCapabilities{
				GetProviderSchemaOptional: true,
				MoveResourceState:         true,
				PlanDestroy:               true,
			},
		},
		"provider-meta-mismatch": {
			servers: []func() tfprotov6.ProviderServer{
				(&tf6testserver.TestServer{
					GetProviderSchemaResponse: &tfprotov6.GetProviderSchemaResponse{
						ProviderMeta: &tfprotov6.Schema{
							Block: &tfprotov6.SchemaBlock{
								Attributes: []*tfprotov6.SchemaAttribute{
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
				(&tf6testserver.TestServer{
					GetProviderSchemaResponse: &tfprotov6.GetProviderSchemaResponse{
						ProviderMeta: &tfprotov6.Schema{
							Block: &tfprotov6.SchemaBlock{
								Attributes: []*tfprotov6.SchemaAttribute{
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
			expectedActionSchemas:     map[string]*tfprotov6.ActionSchema{},
			expectedDataSourceSchemas: map[string]*tfprotov6.Schema{},
			expectedDiagnostics: []*tfprotov6.Diagnostic{
				{
					Severity: tfprotov6.DiagnosticSeverityError,
					Summary:  "Invalid Provider Server Combination",
					Detail: "The combined provider has differing provider meta schema implementations across providers. " +
						"Provider meta schemas must be identical across providers. " +
						"This is always an issue in the provider implementation and should be reported to the provider developers.\n\n" +
						"Provider meta schema difference: " + cmp.Diff(
						&tfprotov6.Schema{
							Block: &tfprotov6.SchemaBlock{
								Attributes: []*tfprotov6.SchemaAttribute{
									{
										Name:     "testattribute2",
										Type:     tftypes.String,
										Required: true,
									},
								},
							},
						},
						&tfprotov6.Schema{
							Block: &tfprotov6.SchemaBlock{
								Attributes: []*tfprotov6.SchemaAttribute{
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
			expectedEphemeralResourcesSchemas: map[string]*tfprotov6.Schema{},
			expectedListResourcesSchemas:      map[string]*tfprotov6.Schema{},
			expectedFunctions:                 map[string]*tfprotov6.Function{},
			expectedProviderMetaSchema: &tfprotov6.Schema{
				Block: &tfprotov6.SchemaBlock{
					Attributes: []*tfprotov6.SchemaAttribute{
						{
							Name:     "testattribute1",
							Type:     tftypes.String,
							Required: true,
						},
					},
				},
			},
			expectedResourceSchemas: map[string]*tfprotov6.Schema{},
			expectedServerCapabilities: &tfprotov6.ServerCapabilities{
				GetProviderSchemaOptional: true,
				MoveResourceState:         true,
				PlanDestroy:               true,
			},
		},
		"server-capabilities": {
			servers: []func() tfprotov6.ProviderServer{
				(&tf6testserver.TestServer{
					GetProviderSchemaResponse: &tfprotov6.GetProviderSchemaResponse{
						ResourceSchemas: map[string]*tfprotov6.Schema{
							"test_with_server_capabilities": {},
						},
						ServerCapabilities: &tfprotov6.ServerCapabilities{
							GetProviderSchemaOptional: true,
							PlanDestroy:               true,
						},
					},
				}).ProviderServer,
				(&tf6testserver.TestServer{
					GetProviderSchemaResponse: &tfprotov6.GetProviderSchemaResponse{
						ResourceSchemas: map[string]*tfprotov6.Schema{
							"test_without_server_capabilities": {},
						},
					},
				}).ProviderServer,
			},
			expectedActionSchemas:             map[string]*tfprotov6.ActionSchema{},
			expectedDataSourceSchemas:         map[string]*tfprotov6.Schema{},
			expectedEphemeralResourcesSchemas: map[string]*tfprotov6.Schema{},
			expectedListResourcesSchemas:      map[string]*tfprotov6.Schema{},
			expectedFunctions:                 map[string]*tfprotov6.Function{},
			expectedResourceSchemas: map[string]*tfprotov6.Schema{
				"test_with_server_capabilities":    {},
				"test_without_server_capabilities": {},
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
					GetProviderSchemaResponse: &tfprotov6.GetProviderSchemaResponse{
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
			expectedActionSchemas:     map[string]*tfprotov6.ActionSchema{},
			expectedDataSourceSchemas: map[string]*tfprotov6.Schema{},
			expectedDiagnostics: []*tfprotov6.Diagnostic{
				{
					Severity: tfprotov6.DiagnosticSeverityError,
					Summary:  "test error summary",
					Detail:   "test error details",
				},
			},
			expectedEphemeralResourcesSchemas: map[string]*tfprotov6.Schema{},
			expectedListResourcesSchemas:      map[string]*tfprotov6.Schema{},
			expectedFunctions:                 map[string]*tfprotov6.Function{},
			expectedResourceSchemas:           map[string]*tfprotov6.Schema{},
			expectedServerCapabilities: &tfprotov6.ServerCapabilities{
				GetProviderSchemaOptional: true,
				MoveResourceState:         true,
				PlanDestroy:               true,
			},
		},
		"error-multiple": {
			servers: []func() tfprotov6.ProviderServer{
				(&tf6testserver.TestServer{
					GetProviderSchemaResponse: &tfprotov6.GetProviderSchemaResponse{
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
					GetProviderSchemaResponse: &tfprotov6.GetProviderSchemaResponse{
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
			expectedActionSchemas:     map[string]*tfprotov6.ActionSchema{},
			expectedDataSourceSchemas: map[string]*tfprotov6.Schema{},
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
			expectedEphemeralResourcesSchemas: map[string]*tfprotov6.Schema{},
			expectedListResourcesSchemas:      map[string]*tfprotov6.Schema{},
			expectedFunctions:                 map[string]*tfprotov6.Function{},
			expectedResourceSchemas:           map[string]*tfprotov6.Schema{},
			expectedServerCapabilities: &tfprotov6.ServerCapabilities{
				GetProviderSchemaOptional: true,
				MoveResourceState:         true,
				PlanDestroy:               true,
			},
		},
		"warning-once": {
			servers: []func() tfprotov6.ProviderServer{
				(&tf6testserver.TestServer{
					GetProviderSchemaResponse: &tfprotov6.GetProviderSchemaResponse{
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
			expectedActionSchemas:     map[string]*tfprotov6.ActionSchema{},
			expectedDataSourceSchemas: map[string]*tfprotov6.Schema{},
			expectedDiagnostics: []*tfprotov6.Diagnostic{
				{
					Severity: tfprotov6.DiagnosticSeverityWarning,
					Summary:  "test warning summary",
					Detail:   "test warning details",
				},
			},
			expectedEphemeralResourcesSchemas: map[string]*tfprotov6.Schema{},
			expectedListResourcesSchemas:      map[string]*tfprotov6.Schema{},
			expectedFunctions:                 map[string]*tfprotov6.Function{},
			expectedResourceSchemas:           map[string]*tfprotov6.Schema{},
			expectedServerCapabilities: &tfprotov6.ServerCapabilities{
				GetProviderSchemaOptional: true,
				MoveResourceState:         true,
				PlanDestroy:               true,
			},
		},
		"warning-multiple": {
			servers: []func() tfprotov6.ProviderServer{
				(&tf6testserver.TestServer{
					GetProviderSchemaResponse: &tfprotov6.GetProviderSchemaResponse{
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
					GetProviderSchemaResponse: &tfprotov6.GetProviderSchemaResponse{
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
			expectedActionSchemas:     map[string]*tfprotov6.ActionSchema{},
			expectedDataSourceSchemas: map[string]*tfprotov6.Schema{},
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
			expectedEphemeralResourcesSchemas: map[string]*tfprotov6.Schema{},
			expectedListResourcesSchemas:      map[string]*tfprotov6.Schema{},
			expectedFunctions:                 map[string]*tfprotov6.Function{},
			expectedResourceSchemas:           map[string]*tfprotov6.Schema{},
			expectedServerCapabilities: &tfprotov6.ServerCapabilities{
				GetProviderSchemaOptional: true,
				MoveResourceState:         true,
				PlanDestroy:               true,
			},
		},
		"warning-then-error": {
			servers: []func() tfprotov6.ProviderServer{
				(&tf6testserver.TestServer{
					GetProviderSchemaResponse: &tfprotov6.GetProviderSchemaResponse{
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
					GetProviderSchemaResponse: &tfprotov6.GetProviderSchemaResponse{
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
			expectedActionSchemas:     map[string]*tfprotov6.ActionSchema{},
			expectedDataSourceSchemas: map[string]*tfprotov6.Schema{},
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
			expectedEphemeralResourcesSchemas: map[string]*tfprotov6.Schema{},
			expectedListResourcesSchemas:      map[string]*tfprotov6.Schema{},
			expectedFunctions:                 map[string]*tfprotov6.Function{},
			expectedResourceSchemas:           map[string]*tfprotov6.Schema{},
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

			resp, err := muxServer.ProviderServer().GetProviderSchema(context.Background(), &tfprotov6.GetProviderSchemaRequest{})

			if err != nil {
				t.Fatalf("unexpected error: %s", err)
			}

			if diff := cmp.Diff(resp.ActionSchemas, testCase.expectedActionSchemas); diff != "" {
				t.Errorf("action schemas didn't match expectations: %s", diff)
			}

			if diff := cmp.Diff(resp.DataSourceSchemas, testCase.expectedDataSourceSchemas); diff != "" {
				t.Errorf("data source schemas didn't match expectations: %s", diff)
			}

			if diff := cmp.Diff(resp.Diagnostics, testCase.expectedDiagnostics); diff != "" {
				t.Errorf("diagnostics didn't match expectations: %s", diff)
			}

			if diff := cmp.Diff(resp.EphemeralResourceSchemas, testCase.expectedEphemeralResourcesSchemas); diff != "" {
				t.Errorf("ephemeral resources schemas didn't match expectations: %s", diff)
			}

			if diff := cmp.Diff(resp.ListResourceSchemas, testCase.expectedListResourcesSchemas); diff != "" {
				t.Errorf("list resources schemas didn't match expectations: %s", diff)
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
