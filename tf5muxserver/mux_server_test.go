package tf5muxserver_test

import (
	"context"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-go/tfprotov5"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/hashicorp/terraform-plugin-mux/internal/tf5testserver"
	"github.com/hashicorp/terraform-plugin-mux/tf5muxserver"
)

func TestNewMuxServer(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		servers       []func() tfprotov5.ProviderServer
		expectedError error
	}{
		"duplicate-data-source": {
			servers: []func() tfprotov5.ProviderServer{
				(&tf5testserver.TestServer{
					DataSourceSchemas: map[string]*tfprotov5.Schema{
						"test_foo": {},
					},
				}).ProviderServer,
				(&tf5testserver.TestServer{
					DataSourceSchemas: map[string]*tfprotov5.Schema{
						"test_foo": {},
					},
				}).ProviderServer,
			},
			expectedError: nil, // deferred to GetProviderSchema
		},
		"duplicate-resource": {
			servers: []func() tfprotov5.ProviderServer{
				(&tf5testserver.TestServer{
					ResourceSchemas: map[string]*tfprotov5.Schema{
						"test_foo": {},
					},
				}).ProviderServer,
				(&tf5testserver.TestServer{
					ResourceSchemas: map[string]*tfprotov5.Schema{
						"test_foo": {},
					},
				}).ProviderServer,
			},
			expectedError: nil, // deferred to GetProviderSchema
		},
		"provider-mismatch": {
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
												Required:        true,
												Description:     "whether the feature is enabled",
												DescriptionKind: tfprotov5.StringKindPlain,
											},
										},
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
				}).ProviderServer,
			},
			expectedError: nil, // deferred to GetProviderSchema
		},
		"provider-ordering": {
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
								{
									Name:            "secret",
									Type:            tftypes.String,
									Required:        true,
									Description:     "the secret to authenticate with",
									DescriptionKind: tfprotov5.StringKindPlain,
								},
							},
							BlockTypes: []*tfprotov5.SchemaNestedBlock{
								{
									TypeName: "other_feature",
									Nesting:  tfprotov5.SchemaNestedBlockNestingModeList,
									Block: &tfprotov5.SchemaBlock{
										Version:         1,
										Description:     "features to enable on the provider",
										DescriptionKind: tfprotov5.StringKindPlain,
										Attributes: []*tfprotov5.SchemaAttribute{
											{
												Name:            "enabled",
												Type:            tftypes.Bool,
												Required:        true,
												Description:     "whether the feature is enabled",
												DescriptionKind: tfprotov5.StringKindPlain,
											},
											{
												Name:            "feature_id",
												Type:            tftypes.Number,
												Required:        true,
												Description:     "The ID of the feature",
												DescriptionKind: tfprotov5.StringKindPlain,
											},
										},
									},
								},
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
												Required:        true,
												Description:     "whether the feature is enabled",
												DescriptionKind: tfprotov5.StringKindPlain,
											},
										},
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
									Name:            "secret",
									Type:            tftypes.String,
									Required:        true,
									Description:     "the secret to authenticate with",
									DescriptionKind: tfprotov5.StringKindPlain,
								},
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
												Name:            "enabled",
												Type:            tftypes.Bool,
												Required:        true,
												Description:     "whether the feature is enabled",
												DescriptionKind: tfprotov5.StringKindPlain,
											},
											{
												Name:            "feature_id",
												Type:            tftypes.Number,
												Required:        true,
												Description:     "The ID of the feature",
												DescriptionKind: tfprotov5.StringKindPlain,
											},
										},
									},
								},
								{
									TypeName: "other_feature",
									Nesting:  tfprotov5.SchemaNestedBlockNestingModeList,
									Block: &tfprotov5.SchemaBlock{
										Version:         1,
										Description:     "features to enable on the provider",
										DescriptionKind: tfprotov5.StringKindPlain,
										Attributes: []*tfprotov5.SchemaAttribute{
											{
												Name:            "enabled",
												Type:            tftypes.Bool,
												Required:        true,
												Description:     "whether the feature is enabled",
												DescriptionKind: tfprotov5.StringKindPlain,
											},
											{
												Name:            "feature_id",
												Type:            tftypes.Number,
												Required:        true,
												Description:     "The ID of the feature",
												DescriptionKind: tfprotov5.StringKindPlain,
											},
										},
									},
								},
							},
						},
					},
				}).ProviderServer,
			},
		},
		"provider-meta-mismatch": {
			servers: []func() tfprotov5.ProviderServer{
				(&tf5testserver.TestServer{
					ProviderMetaSchema: &tfprotov5.Schema{
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
												Required:        true,
												Description:     "whether the feature is enabled",
												DescriptionKind: tfprotov5.StringKindPlain,
											},
										},
									},
								},
							},
						},
					},
				}).ProviderServer,
				(&tf5testserver.TestServer{
					ProviderMetaSchema: &tfprotov5.Schema{
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
				}).ProviderServer,
			},
			expectedError: nil, // deferred to GetProviderSchema
		},
		"provider-meta-ordering": {
			servers: []func() tfprotov5.ProviderServer{
				(&tf5testserver.TestServer{
					ProviderMetaSchema: &tfprotov5.Schema{
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
								{
									Name:            "secret",
									Type:            tftypes.String,
									Required:        true,
									Description:     "the secret to authenticate with",
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
												Required:        true,
												Description:     "whether the feature is enabled",
												DescriptionKind: tfprotov5.StringKindPlain,
											},
										},
									},
								},
								{
									TypeName: "other_feature",
									Nesting:  tfprotov5.SchemaNestedBlockNestingModeList,
									Block: &tfprotov5.SchemaBlock{
										Version:         1,
										Description:     "features to enable on the provider",
										DescriptionKind: tfprotov5.StringKindPlain,
										Attributes: []*tfprotov5.SchemaAttribute{
											{
												Name:            "enabled",
												Type:            tftypes.Bool,
												Required:        true,
												Description:     "whether the feature is enabled",
												DescriptionKind: tfprotov5.StringKindPlain,
											},
											{
												Name:            "feature_id",
												Type:            tftypes.Number,
												Required:        true,
												Description:     "The ID of the feature",
												DescriptionKind: tfprotov5.StringKindPlain,
											},
										},
									},
								},
							},
						},
					},
				}).ProviderServer,
				(&tf5testserver.TestServer{
					ProviderMetaSchema: &tfprotov5.Schema{
						Version: 1,
						Block: &tfprotov5.SchemaBlock{
							Version: 1,
							Attributes: []*tfprotov5.SchemaAttribute{
								{
									Name:            "secret",
									Type:            tftypes.String,
									Required:        true,
									Description:     "the secret to authenticate with",
									DescriptionKind: tfprotov5.StringKindPlain,
								},
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
									TypeName: "other_feature",
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
												Required:        true,
												Description:     "whether the feature is enabled",
												DescriptionKind: tfprotov5.StringKindPlain,
											},
										},
									},
								},
								{
									TypeName: "feature",
									Nesting:  tfprotov5.SchemaNestedBlockNestingModeList,
									Block: &tfprotov5.SchemaBlock{
										Version:         1,
										Description:     "features to enable on the provider",
										DescriptionKind: tfprotov5.StringKindPlain,
										Attributes: []*tfprotov5.SchemaAttribute{
											{
												Name:            "enabled",
												Type:            tftypes.Bool,
												Required:        true,
												Description:     "whether the feature is enabled",
												DescriptionKind: tfprotov5.StringKindPlain,
											},
											{
												Name:            "feature_id",
												Type:            tftypes.Number,
												Required:        true,
												Description:     "The ID of the feature",
												DescriptionKind: tfprotov5.StringKindPlain,
											},
										},
									},
								},
							},
						},
					},
				}).ProviderServer,
			},
		},
		"nested block MinItems and MaxItems diff are ignored": {
			servers: []func() tfprotov5.ProviderServer{
				(&tf5testserver.TestServer{
					ProviderSchema: &tfprotov5.Schema{
						Version: 1,
						Block: &tfprotov5.SchemaBlock{
							Version: 1,
							BlockTypes: []*tfprotov5.SchemaNestedBlock{
								{
									TypeName: "feature",
									Nesting:  tfprotov5.SchemaNestedBlockNestingModeList,
									Block: &tfprotov5.SchemaBlock{
										Version:         1,
										Description:     "features to enable on the provider",
										DescriptionKind: tfprotov5.StringKindPlain,
										Attributes:      []*tfprotov5.SchemaAttribute{},
									},
									MinItems: 1,
									MaxItems: 2,
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
							BlockTypes: []*tfprotov5.SchemaNestedBlock{
								{
									TypeName: "feature",
									Nesting:  tfprotov5.SchemaNestedBlockNestingModeList,
									Block: &tfprotov5.SchemaBlock{
										Version:         1,
										Description:     "features to enable on the provider",
										DescriptionKind: tfprotov5.StringKindPlain,
										Attributes:      []*tfprotov5.SchemaAttribute{},
									},
									MinItems: 0,
									MaxItems: 0,
								},
							},
						},
					},
				}).ProviderServer,
			},
		},
		"server-capabilities": {
			servers: []func() tfprotov5.ProviderServer{
				(&tf5testserver.TestServer{
					ResourceSchemas: map[string]*tfprotov5.Schema{
						"test_with_server_capabilities": {},
					},
					ServerCapabilities: &tfprotov5.ServerCapabilities{
						PlanDestroy: true,
					},
				}).ProviderServer,
				(&tf5testserver.TestServer{
					ResourceSchemas: map[string]*tfprotov5.Schema{
						"test_without_server_capabilities": {},
					},
				}).ProviderServer,
			},
			expectedError: nil,
		},
	}

	for name, testCase := range testCases {
		name, testCase := name, testCase

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			_, err := tf5muxserver.NewMuxServer(context.Background(), testCase.servers...)

			if err != nil {
				if testCase.expectedError == nil {
					t.Fatalf("unexpected error: %s", err)
				}

				if !strings.Contains(err.Error(), testCase.expectedError.Error()) {
					t.Fatalf("expected error %q, got: %s", testCase.expectedError, err)
				}
			}

			if err == nil && testCase.expectedError != nil {
				t.Fatalf("expected error: %s", testCase.expectedError)
			}
		})
	}
}
