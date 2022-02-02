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
		expectedProviderSchema     *tfprotov5.Schema
		expectedProviderMetaSchema *tfprotov5.Schema
		expectedResourceSchemas    map[string]*tfprotov5.Schema
	}{
		"combined": {
			servers: []func() tfprotov5.ProviderServer{
				(&tf5testserver.TestServer{
					ProviderSchema: &tfprotov5.Schema{
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
					ProviderMetaSchema: &tfprotov5.Schema{
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
				}).ProviderServer,
				(&tf5testserver.TestServer{
					ProviderSchema: &tfprotov5.Schema{
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
					ProviderMetaSchema: &tfprotov5.Schema{
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
