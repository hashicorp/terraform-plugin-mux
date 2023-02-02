package tf6muxserver_test

import (
	"context"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/hashicorp/terraform-plugin-mux/internal/tf6testserver"
	"github.com/hashicorp/terraform-plugin-mux/tf6muxserver"
)

func TestNewMuxServer(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		servers       []func() tfprotov6.ProviderServer
		expectedError error
	}{
		"duplicate-data-source": {
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
			expectedError: nil, // deferred to GetProviderSchema
		},
		"duplicate-resource": {
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
			expectedError: nil, // deferred to GetProviderSchema
		},
		"provider-mismatch": {
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
												Required:        true,
												Description:     "whether the feature is enabled",
												DescriptionKind: tfprotov6.StringKindPlain,
											},
										},
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
				}).ProviderServer,
			},
			expectedError: nil, // deferred to GetProviderSchema
		},
		"provider-ordering": {
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
								{
									Name:            "secret",
									Type:            tftypes.String,
									Required:        true,
									Description:     "the secret to authenticate with",
									DescriptionKind: tfprotov6.StringKindPlain,
								},
							},
							BlockTypes: []*tfprotov6.SchemaNestedBlock{
								{
									TypeName: "other_feature",
									Nesting:  tfprotov6.SchemaNestedBlockNestingModeList,
									Block: &tfprotov6.SchemaBlock{
										Version:         1,
										Description:     "features to enable on the provider",
										DescriptionKind: tfprotov6.StringKindPlain,
										Attributes: []*tfprotov6.SchemaAttribute{
											{
												Name:            "enabled",
												Type:            tftypes.Bool,
												Required:        true,
												Description:     "whether the feature is enabled",
												DescriptionKind: tfprotov6.StringKindPlain,
											},
											{
												Name:            "feature_id",
												Type:            tftypes.Number,
												Required:        true,
												Description:     "The ID of the feature",
												DescriptionKind: tfprotov6.StringKindPlain,
											},
										},
									},
								},
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
												Required:        true,
												Description:     "whether the feature is enabled",
												DescriptionKind: tfprotov6.StringKindPlain,
											},
										},
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
									Name:            "secret",
									Type:            tftypes.String,
									Required:        true,
									Description:     "the secret to authenticate with",
									DescriptionKind: tfprotov6.StringKindPlain,
								},
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
												Name:            "enabled",
												Type:            tftypes.Bool,
												Required:        true,
												Description:     "whether the feature is enabled",
												DescriptionKind: tfprotov6.StringKindPlain,
											},
											{
												Name:            "feature_id",
												Type:            tftypes.Number,
												Required:        true,
												Description:     "The ID of the feature",
												DescriptionKind: tfprotov6.StringKindPlain,
											},
										},
									},
								},
								{
									TypeName: "other_feature",
									Nesting:  tfprotov6.SchemaNestedBlockNestingModeList,
									Block: &tfprotov6.SchemaBlock{
										Version:         1,
										Description:     "features to enable on the provider",
										DescriptionKind: tfprotov6.StringKindPlain,
										Attributes: []*tfprotov6.SchemaAttribute{
											{
												Name:            "enabled",
												Type:            tftypes.Bool,
												Required:        true,
												Description:     "whether the feature is enabled",
												DescriptionKind: tfprotov6.StringKindPlain,
											},
											{
												Name:            "feature_id",
												Type:            tftypes.Number,
												Required:        true,
												Description:     "The ID of the feature",
												DescriptionKind: tfprotov6.StringKindPlain,
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
			servers: []func() tfprotov6.ProviderServer{
				(&tf6testserver.TestServer{
					ProviderMetaSchema: &tfprotov6.Schema{
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
												Required:        true,
												Description:     "whether the feature is enabled",
												DescriptionKind: tfprotov6.StringKindPlain,
											},
										},
									},
								},
							},
						},
					},
				}).ProviderServer,
				(&tf6testserver.TestServer{
					ProviderMetaSchema: &tfprotov6.Schema{
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
				}).ProviderServer,
			},
			expectedError: nil, // deferred to GetProviderSchema
		},
		"provider-meta-ordering": {
			servers: []func() tfprotov6.ProviderServer{
				(&tf6testserver.TestServer{
					ProviderMetaSchema: &tfprotov6.Schema{
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
								{
									Name:            "secret",
									Type:            tftypes.String,
									Required:        true,
									Description:     "the secret to authenticate with",
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
												Required:        true,
												Description:     "whether the feature is enabled",
												DescriptionKind: tfprotov6.StringKindPlain,
											},
										},
									},
								},
								{
									TypeName: "other_feature",
									Nesting:  tfprotov6.SchemaNestedBlockNestingModeList,
									Block: &tfprotov6.SchemaBlock{
										Version:         1,
										Description:     "features to enable on the provider",
										DescriptionKind: tfprotov6.StringKindPlain,
										Attributes: []*tfprotov6.SchemaAttribute{
											{
												Name:            "enabled",
												Type:            tftypes.Bool,
												Required:        true,
												Description:     "whether the feature is enabled",
												DescriptionKind: tfprotov6.StringKindPlain,
											},
											{
												Name:            "feature_id",
												Type:            tftypes.Number,
												Required:        true,
												Description:     "The ID of the feature",
												DescriptionKind: tfprotov6.StringKindPlain,
											},
										},
									},
								},
							},
						},
					},
				}).ProviderServer,
				(&tf6testserver.TestServer{
					ProviderMetaSchema: &tfprotov6.Schema{
						Version: 1,
						Block: &tfprotov6.SchemaBlock{
							Version: 1,
							Attributes: []*tfprotov6.SchemaAttribute{
								{
									Name:            "secret",
									Type:            tftypes.String,
									Required:        true,
									Description:     "the secret to authenticate with",
									DescriptionKind: tfprotov6.StringKindPlain,
								},
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
									TypeName: "other_feature",
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
												Required:        true,
												Description:     "whether the feature is enabled",
												DescriptionKind: tfprotov6.StringKindPlain,
											},
										},
									},
								},
								{
									TypeName: "feature",
									Nesting:  tfprotov6.SchemaNestedBlockNestingModeList,
									Block: &tfprotov6.SchemaBlock{
										Version:         1,
										Description:     "features to enable on the provider",
										DescriptionKind: tfprotov6.StringKindPlain,
										Attributes: []*tfprotov6.SchemaAttribute{
											{
												Name:            "enabled",
												Type:            tftypes.Bool,
												Required:        true,
												Description:     "whether the feature is enabled",
												DescriptionKind: tfprotov6.StringKindPlain,
											},
											{
												Name:            "feature_id",
												Type:            tftypes.Number,
												Required:        true,
												Description:     "The ID of the feature",
												DescriptionKind: tfprotov6.StringKindPlain,
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
			servers: []func() tfprotov6.ProviderServer{
				(&tf6testserver.TestServer{
					ProviderSchema: &tfprotov6.Schema{
						Version: 1,
						Block: &tfprotov6.SchemaBlock{
							Version: 1,
							BlockTypes: []*tfprotov6.SchemaNestedBlock{
								{
									TypeName: "feature",
									Nesting:  tfprotov6.SchemaNestedBlockNestingModeList,
									Block: &tfprotov6.SchemaBlock{
										Version:         1,
										Description:     "features to enable on the provider",
										DescriptionKind: tfprotov6.StringKindPlain,
										Attributes:      []*tfprotov6.SchemaAttribute{},
									},
									MinItems: 1,
									MaxItems: 2,
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
							BlockTypes: []*tfprotov6.SchemaNestedBlock{
								{
									TypeName: "feature",
									Nesting:  tfprotov6.SchemaNestedBlockNestingModeList,
									Block: &tfprotov6.SchemaBlock{
										Version:         1,
										Description:     "features to enable on the provider",
										DescriptionKind: tfprotov6.StringKindPlain,
										Attributes:      []*tfprotov6.SchemaAttribute{},
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
			expectedError: nil,
		},
	}

	for name, testCase := range testCases {
		name, testCase := name, testCase

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			_, err := tf6muxserver.NewMuxServer(context.Background(), testCase.servers...)

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
