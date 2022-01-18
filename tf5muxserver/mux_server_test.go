package tf5muxserver_test

import (
	"context"
	"fmt"
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
			expectedError: fmt.Errorf("data source \"test_foo\" is implemented by multiple servers; only one implementation allowed"),
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
			expectedError: fmt.Errorf("resource \"test_foo\" is implemented by multiple servers; only one implementation allowed"),
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
			expectedError: fmt.Errorf("got a different provider schema across servers. Provider schemas must be identical across providers"),
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
			expectedError: fmt.Errorf("got a different provider meta schema across servers. Provider metadata schemas must be identical across providers"),
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
