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
		servers                    []func() tfprotov6.ProviderServer
		expectedDataSourceSchemas  map[string]*tfprotov6.Schema
		expectedDiagnostics        []*tfprotov6.Diagnostic
		expectedProviderSchema     *tfprotov6.Schema
		expectedProviderMetaSchema *tfprotov6.Schema
		expectedResourceSchemas    map[string]*tfprotov6.Schema
		expectedServerCapabilities *tfprotov6.ServerCapabilities
	}{
		"combined": {
			servers: []func() tfprotov6.ProviderServer{
				(&tf6testserver.TestServer{
					ProviderSchema: &tfprotov6.Schema{
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
					ProviderMetaSchema: &tfprotov6.Schema{
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
				}).ProviderServer,
				(&tf6testserver.TestServer{
					ProviderSchema: &tfprotov6.Schema{
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
					ProviderMetaSchema: &tfprotov6.Schema{
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
			expectedServerCapabilities: &tfprotov6.ServerCapabilities{
				PlanDestroy: true,
			},
		},
		"duplicate-data-source-type": {
			servers: []func() tfprotov6.ProviderServer{
				(&tf6testserver.TestServer{
					DataSourceSchemas: map[string]*tfprotov6.Schema{
						"test_foo": {},
					},
				}).ProviderServer,
				(&tf6testserver.TestServer{
					DataSourceSchemas: map[string]*tfprotov6.Schema{
						"test_foo": {},
					},
				}).ProviderServer,
			},
			expectedDataSourceSchemas: map[string]*tfprotov6.Schema{
				"test_foo": {},
			},
			expectedDiagnostics: []*tfprotov6.Diagnostic{
				{
					Summary: "Invalid Provider Server Combination",
					Detail: "The combined provider has multiple implementations of the same data source type across providers. " +
						"Data source types must be implemented by only one provider. " +
						"This is always an issue in the provider implementation and should be reported to the provider developers.\n\n" +
						"Duplicate data source type: test_foo",
				},
			},
			expectedResourceSchemas: map[string]*tfprotov6.Schema{},
			expectedServerCapabilities: &tfprotov6.ServerCapabilities{
				PlanDestroy: true,
			},
		},
		"duplicate-resource-type": {
			servers: []func() tfprotov6.ProviderServer{
				(&tf6testserver.TestServer{
					ResourceSchemas: map[string]*tfprotov6.Schema{
						"test_foo": {},
					},
				}).ProviderServer,
				(&tf6testserver.TestServer{
					ResourceSchemas: map[string]*tfprotov6.Schema{
						"test_foo": {},
					},
				}).ProviderServer,
			},
			expectedDataSourceSchemas: map[string]*tfprotov6.Schema{},
			expectedDiagnostics: []*tfprotov6.Diagnostic{
				{
					Summary: "Invalid Provider Server Combination",
					Detail: "The combined provider has multiple implementations of the same resource type across providers. " +
						"Resource types must be implemented by only one provider. " +
						"This is always an issue in the provider implementation and should be reported to the provider developers.\n\n" +
						"Duplicate resource type: test_foo",
				},
			},
			expectedResourceSchemas: map[string]*tfprotov6.Schema{
				"test_foo": {},
			},
			expectedServerCapabilities: &tfprotov6.ServerCapabilities{
				PlanDestroy: true,
			},
		},
		"provider-mismatch": {
			servers: []func() tfprotov6.ProviderServer{
				(&tf6testserver.TestServer{
					ProviderSchema: &tfprotov6.Schema{
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
				}).ProviderServer,
				(&tf6testserver.TestServer{
					ProviderSchema: &tfprotov6.Schema{
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
				}).ProviderServer,
			},
			expectedDataSourceSchemas: map[string]*tfprotov6.Schema{},
			expectedDiagnostics: []*tfprotov6.Diagnostic{
				{
					Summary: "Invalid Provider Server Combination",
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
				PlanDestroy: true,
			},
		},
		"provider-meta-mismatch": {
			servers: []func() tfprotov6.ProviderServer{
				(&tf6testserver.TestServer{
					ProviderMetaSchema: &tfprotov6.Schema{
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
				}).ProviderServer,
				(&tf6testserver.TestServer{
					ProviderMetaSchema: &tfprotov6.Schema{
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
				}).ProviderServer,
			},
			expectedDataSourceSchemas: map[string]*tfprotov6.Schema{},
			expectedDiagnostics: []*tfprotov6.Diagnostic{
				{
					Summary: "Invalid Provider Server Combination",
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
				PlanDestroy: true,
			},
		},
		"server-capabilities": {
			servers: []func() tfprotov6.ProviderServer{
				(&tf6testserver.TestServer{
					ResourceSchemas: map[string]*tfprotov6.Schema{
						"test_with_server_capabilities": {},
					},
					ServerCapabilities: &tfprotov6.ServerCapabilities{
						PlanDestroy: true,
					},
				}).ProviderServer,
				(&tf6testserver.TestServer{
					ResourceSchemas: map[string]*tfprotov6.Schema{
						"test_without_server_capabilities": {},
					},
				}).ProviderServer,
			},
			expectedDataSourceSchemas: map[string]*tfprotov6.Schema{},
			expectedResourceSchemas: map[string]*tfprotov6.Schema{
				"test_with_server_capabilities":    {},
				"test_without_server_capabilities": {},
			},
			expectedServerCapabilities: &tfprotov6.ServerCapabilities{
				PlanDestroy: true,
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

			resp, err := muxServer.ProviderServer().GetProviderSchema(context.Background(), &tfprotov6.GetProviderSchemaRequest{})

			if err != nil {
				t.Fatalf("unexpected error: %s", err)
			}

			if diff := cmp.Diff(resp.DataSourceSchemas, testCase.expectedDataSourceSchemas); diff != "" {
				t.Errorf("data source schemas didn't match expectations: %s", diff)
			}

			if diff := cmp.Diff(resp.Diagnostics, testCase.expectedDiagnostics); diff != "" {
				t.Errorf("diagnostics didn't match expectations: %s", diff)
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
		})
	}
}
