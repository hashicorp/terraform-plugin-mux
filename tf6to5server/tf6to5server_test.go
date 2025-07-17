// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package tf6to5server_test

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-go/tfprotov5"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-go/tftypes"

	"github.com/hashicorp/terraform-plugin-mux/internal/tf6testserver"
	"github.com/hashicorp/terraform-plugin-mux/internal/tfprotov6tov5"
	"github.com/hashicorp/terraform-plugin-mux/tf6to5server"
)

func TestDowngradeServer(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		v6Server      func() tfprotov6.ProviderServer
		expectedError error
	}{
		"compatible": {
			v6Server: (&tf6testserver.TestServer{
				GetProviderSchemaResponse: &tfprotov6.GetProviderSchemaResponse{
					ActionSchemas: map[string]*tfprotov6.ActionSchema{
						"test_action": {
							Type: tfprotov6.UnlinkedActionSchemaType{},
						},
					},
					DataSourceSchemas: map[string]*tfprotov6.Schema{
						"test_data_source": {},
					},
					EphemeralResourceSchemas: map[string]*tfprotov6.Schema{
						"test_ephemeral_resource": {},
					},
					Functions: map[string]*tfprotov6.Function{
						"test_function": {},
					},
					ListResourceSchemas: map[string]*tfprotov6.Schema{
						"test_list_resource": {},
					},
					Provider: &tfprotov6.Schema{
						Block: &tfprotov6.SchemaBlock{
							Attributes: []*tfprotov6.SchemaAttribute{
								{
									Name:            "test_attribute",
									Type:            tftypes.String,
									Required:        true,
									Description:     "test_attribute description",
									DescriptionKind: tfprotov6.StringKindPlain,
								},
							},
							BlockTypes: []*tfprotov6.SchemaNestedBlock{
								{
									TypeName: "test_block",
									Nesting:  tfprotov6.SchemaNestedBlockNestingModeList,
									Block: &tfprotov6.SchemaBlock{
										Description:     "test_block description",
										DescriptionKind: tfprotov6.StringKindPlain,
										Attributes: []*tfprotov6.SchemaAttribute{
											{
												Name:            "test_block_attribute",
												Type:            tftypes.Number,
												Required:        true,
												Description:     "test_block_attribute description",
												DescriptionKind: tfprotov6.StringKindPlain,
											},
										},
									},
								},
							},
						},
					},
					ProviderMeta: &tfprotov6.Schema{
						Block: &tfprotov6.SchemaBlock{
							Attributes: []*tfprotov6.SchemaAttribute{
								{
									Name:            "test_attribute",
									Type:            tftypes.String,
									Required:        true,
									Description:     "test_attribute description",
									DescriptionKind: tfprotov6.StringKindPlain,
								},
							},
							BlockTypes: []*tfprotov6.SchemaNestedBlock{
								{
									TypeName: "test_block",
									Nesting:  tfprotov6.SchemaNestedBlockNestingModeList,
									Block: &tfprotov6.SchemaBlock{
										Description:     "test_block description",
										DescriptionKind: tfprotov6.StringKindPlain,
										Attributes: []*tfprotov6.SchemaAttribute{
											{
												Name:            "test_block_attribute",
												Type:            tftypes.Number,
												Required:        true,
												Description:     "test_block_attribute description",
												DescriptionKind: tfprotov6.StringKindPlain,
											},
										},
									},
								},
							},
						},
					},
					ResourceSchemas: map[string]*tfprotov6.Schema{
						"test_resource": {},
					},
				},
			}).ProviderServer,
		},
		"SchemaAttribute-NestedType-data-source": {
			v6Server: (&tf6testserver.TestServer{
				GetProviderSchemaResponse: &tfprotov6.GetProviderSchemaResponse{
					DataSourceSchemas: map[string]*tfprotov6.Schema{
						"test_data_source": {
							Block: &tfprotov6.SchemaBlock{
								Attributes: []*tfprotov6.SchemaAttribute{
									{
										Name: "test_attribute",
										NestedType: &tfprotov6.SchemaObject{
											Attributes: []*tfprotov6.SchemaAttribute{},
											Nesting:    tfprotov6.SchemaObjectNestingModeSingle,
										},
										Required:        true,
										Description:     "test_attribute description",
										DescriptionKind: tfprotov6.StringKindPlain,
									},
								},
							},
						},
					},
				},
			}).ProviderServer,
			expectedError: fmt.Errorf("unable to convert data source \"test_data_source\" schema: unable to convert attribute \"test_attribute\" schema: %s", tfprotov6tov5.ErrSchemaAttributeNestedTypeNotImplemented),
		},
		"SchemaAttribute-NestedType-provider": {
			v6Server: (&tf6testserver.TestServer{
				GetProviderSchemaResponse: &tfprotov6.GetProviderSchemaResponse{
					Provider: &tfprotov6.Schema{
						Block: &tfprotov6.SchemaBlock{
							Attributes: []*tfprotov6.SchemaAttribute{
								{
									Name: "test_attribute",
									NestedType: &tfprotov6.SchemaObject{
										Attributes: []*tfprotov6.SchemaAttribute{},
										Nesting:    tfprotov6.SchemaObjectNestingModeSingle,
									},
									Required:        true,
									Description:     "test_attribute description",
									DescriptionKind: tfprotov6.StringKindPlain,
								},
							},
						},
					},
				},
			}).ProviderServer,
			expectedError: fmt.Errorf("unable to convert provider schema: unable to convert attribute \"test_attribute\" schema: %s", tfprotov6tov5.ErrSchemaAttributeNestedTypeNotImplemented),
		},
		"SchemaAttribute-NestedType-provider-meta": {
			v6Server: (&tf6testserver.TestServer{
				GetProviderSchemaResponse: &tfprotov6.GetProviderSchemaResponse{
					ProviderMeta: &tfprotov6.Schema{
						Block: &tfprotov6.SchemaBlock{
							Attributes: []*tfprotov6.SchemaAttribute{
								{
									Name: "test_attribute",
									NestedType: &tfprotov6.SchemaObject{
										Attributes: []*tfprotov6.SchemaAttribute{},
										Nesting:    tfprotov6.SchemaObjectNestingModeSingle,
									},
									Required:        true,
									Description:     "test_attribute description",
									DescriptionKind: tfprotov6.StringKindPlain,
								},
							},
						},
					},
				},
			}).ProviderServer,
			expectedError: fmt.Errorf("unable to convert provider meta schema: unable to convert attribute \"test_attribute\" schema: %s", tfprotov6tov5.ErrSchemaAttributeNestedTypeNotImplemented),
		},
		"SchemaAttribute-NestedType-resource": {
			v6Server: (&tf6testserver.TestServer{
				GetProviderSchemaResponse: &tfprotov6.GetProviderSchemaResponse{
					ResourceSchemas: map[string]*tfprotov6.Schema{
						"test_resource": {
							Block: &tfprotov6.SchemaBlock{
								Attributes: []*tfprotov6.SchemaAttribute{
									{
										Name: "test_attribute",
										NestedType: &tfprotov6.SchemaObject{
											Attributes: []*tfprotov6.SchemaAttribute{},
											Nesting:    tfprotov6.SchemaObjectNestingModeSingle,
										},
										Required:        true,
										Description:     "test_attribute description",
										DescriptionKind: tfprotov6.StringKindPlain,
									},
								},
							},
						},
					},
				},
			}).ProviderServer,
			expectedError: fmt.Errorf("unable to convert resource \"test_resource\" schema: unable to convert attribute \"test_attribute\" schema: %s", tfprotov6tov5.ErrSchemaAttributeNestedTypeNotImplemented),
		},
	}

	for name, testCase := range testCases {

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			_, err := tf6to5server.DowngradeServer(context.Background(), testCase.v6Server)

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

func TestV6ToV5ServerApplyResourceChange(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	v6server := &tf6testserver.TestServer{
		GetProviderSchemaResponse: &tfprotov6.GetProviderSchemaResponse{
			ResourceSchemas: map[string]*tfprotov6.Schema{
				"test_resource": {},
			},
		},
	}

	v5server, err := tf6to5server.DowngradeServer(context.Background(), v6server.ProviderServer)

	if err != nil {
		t.Fatalf("unexpected error downgrading server: %s", err)
	}

	_, err = v5server.ApplyResourceChange(ctx, &tfprotov5.ApplyResourceChangeRequest{
		TypeName: "test_resource",
	})

	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	if !v6server.ApplyResourceChangeCalled["test_resource"] {
		t.Errorf("expected test_resource ApplyResourceChange to be called")
	}
}

func TestV6ToV5ServerCallFunction(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	v6server := &tf6testserver.TestServer{
		GetProviderSchemaResponse: &tfprotov6.GetProviderSchemaResponse{
			Functions: map[string]*tfprotov6.Function{
				"test_function": {},
			},
		},
	}

	v5server, err := tf6to5server.DowngradeServer(context.Background(), v6server.ProviderServer)

	if err != nil {
		t.Fatalf("unexpected error downgrading server: %s", err)
	}

	_, err = v5server.CallFunction(ctx, &tfprotov5.CallFunctionRequest{
		Name: "test_function",
	})

	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	if !v6server.CallFunctionCalled["test_function"] {
		t.Errorf("expected test_function CallFunction to be called")
	}
}

func TestV6ToV5ServerCloseEphemeralResource(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	v6server := &tf6testserver.TestServer{
		GetProviderSchemaResponse: &tfprotov6.GetProviderSchemaResponse{
			ResourceSchemas: map[string]*tfprotov6.Schema{
				"test_ephemeral_resource": {},
			},
		},
	}

	v5server, err := tf6to5server.DowngradeServer(context.Background(), v6server.ProviderServer)

	if err != nil {
		t.Fatalf("unexpected error downgrading server: %s", err)
	}

	_, err = v5server.CloseEphemeralResource(ctx, &tfprotov5.CloseEphemeralResourceRequest{
		TypeName: "test_ephemeral_resource",
	})

	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	if !v6server.CloseEphemeralResourceCalled["test_ephemeral_resource"] {
		t.Errorf("expected test_ephemeral_resource CloseEphemeralResource to be called")
	}
}

func TestV6ToV5ServerConfigureProvider(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	v6server := &tf6testserver.TestServer{
		GetProviderSchemaResponse: &tfprotov6.GetProviderSchemaResponse{
			ResourceSchemas: map[string]*tfprotov6.Schema{
				"test_resource": {},
			},
		},
	}

	v5server, err := tf6to5server.DowngradeServer(context.Background(), v6server.ProviderServer)

	if err != nil {
		t.Fatalf("unexpected error downgrading server: %s", err)
	}

	_, err = v5server.ConfigureProvider(ctx, &tfprotov5.ConfigureProviderRequest{})

	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	if !v6server.ConfigureProviderCalled {
		t.Errorf("expected ConfigureProvider to be called")
	}
}

func TestV6ToV5ServerGetFunctions(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	v6server := &tf6testserver.TestServer{
		GetFunctionsResponse: &tfprotov6.GetFunctionsResponse{
			Functions: map[string]*tfprotov6.Function{
				"test_function": {},
			},
		},
	}

	v5server, err := tf6to5server.DowngradeServer(context.Background(), v6server.ProviderServer)

	if err != nil {
		t.Fatalf("unexpected error downgrading server: %s", err)
	}

	_, err = v5server.GetFunctions(ctx, &tfprotov5.GetFunctionsRequest{})

	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	if !v6server.GetFunctionsCalled {
		t.Errorf("expected GetFunctions to be called")
	}
}

func TestV6ToV5ServerGetMetadata(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	v6server := &tf6testserver.TestServer{
		GetMetadataResponse: &tfprotov6.GetMetadataResponse{
			Resources: []tfprotov6.ResourceMetadata{
				{
					TypeName: "test_resource",
				},
			},
		},
	}

	v5server, err := tf6to5server.DowngradeServer(context.Background(), v6server.ProviderServer)

	if err != nil {
		t.Fatalf("unexpected error downgrading server: %s", err)
	}

	_, err = v5server.GetMetadata(ctx, &tfprotov5.GetMetadataRequest{})

	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	if !v6server.GetMetadataCalled {
		t.Errorf("expected GetMetadata to be called")
	}
}

func TestV6ToV5ServerGetProviderSchema(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	v6server := &tf6testserver.TestServer{
		GetProviderSchemaResponse: &tfprotov6.GetProviderSchemaResponse{
			ResourceSchemas: map[string]*tfprotov6.Schema{
				"test_resource": {},
			},
		},
	}

	v5server, err := tf6to5server.DowngradeServer(context.Background(), v6server.ProviderServer)

	if err != nil {
		t.Fatalf("unexpected error downgrading server: %s", err)
	}

	_, err = v5server.GetProviderSchema(ctx, &tfprotov5.GetProviderSchemaRequest{})

	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	if !v6server.GetProviderSchemaCalled {
		t.Errorf("expected GetProviderSchema to be called")
	}
}

func TestV6ToV5ServerGetResourceIdentitySchemas(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	v6server := &tf6testserver.TestServer{
		GetResourceIdentitySchemasResponse: &tfprotov6.GetResourceIdentitySchemasResponse{},
	}

	v5server, err := tf6to5server.DowngradeServer(context.Background(), v6server.ProviderServer)

	if err != nil {
		t.Fatalf("unexpected error downgrading server: %s", err)
	}

	_, err = v5server.GetResourceIdentitySchemas(ctx, &tfprotov5.GetResourceIdentitySchemasRequest{})

	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	if !v6server.GetResourceIdentitySchemasCalled {
		t.Errorf("expected GetResourceIdentitySchemas to be called")
	}
}

func TestV6ToV5ServerImportResourceState(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	v6server := &tf6testserver.TestServer{
		GetProviderSchemaResponse: &tfprotov6.GetProviderSchemaResponse{
			ResourceSchemas: map[string]*tfprotov6.Schema{
				"test_resource": {},
			},
		},
	}

	v5server, err := tf6to5server.DowngradeServer(context.Background(), v6server.ProviderServer)

	if err != nil {
		t.Fatalf("unexpected error downgrading server: %s", err)
	}

	_, err = v5server.ImportResourceState(ctx, &tfprotov5.ImportResourceStateRequest{
		TypeName: "test_resource",
	})

	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	if !v6server.ImportResourceStateCalled["test_resource"] {
		t.Errorf("expected test_resource ImportResourceState to be called")
	}
}

func TestV6ToV5ServerMoveResourceState(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	v6server := &tf6testserver.TestServer{
		GetProviderSchemaResponse: &tfprotov6.GetProviderSchemaResponse{
			ResourceSchemas: map[string]*tfprotov6.Schema{
				"test_resource": {},
			},
		},
	}

	v5server, err := tf6to5server.DowngradeServer(context.Background(), v6server.ProviderServer)

	if err != nil {
		t.Fatalf("unexpected error downgrading server: %s", err)
	}

	_, err = v5server.MoveResourceState(ctx, &tfprotov5.MoveResourceStateRequest{
		TargetTypeName: "test_resource",
	})

	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	if !v6server.MoveResourceStateCalled["test_resource"] {
		t.Errorf("expected test_resource MoveResourceState to be called")
	}
}

func TestV6ToV5ServerOpenEphemeralResource(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	v6server := &tf6testserver.TestServer{
		GetProviderSchemaResponse: &tfprotov6.GetProviderSchemaResponse{
			ResourceSchemas: map[string]*tfprotov6.Schema{
				"test_ephemeral_resource": {},
			},
		},
	}

	v5server, err := tf6to5server.DowngradeServer(context.Background(), v6server.ProviderServer)

	if err != nil {
		t.Fatalf("unexpected error downgrading server: %s", err)
	}

	_, err = v5server.OpenEphemeralResource(ctx, &tfprotov5.OpenEphemeralResourceRequest{
		TypeName: "test_ephemeral_resource",
	})

	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	if !v6server.OpenEphemeralResourceCalled["test_ephemeral_resource"] {
		t.Errorf("expected test_ephemeral_resource OpenEphemeralResource to be called")
	}
}

func TestV6ToV5ServerPlanResourceChange(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	v6server := &tf6testserver.TestServer{
		GetProviderSchemaResponse: &tfprotov6.GetProviderSchemaResponse{
			ResourceSchemas: map[string]*tfprotov6.Schema{
				"test_resource": {},
			},
		},
	}

	v5server, err := tf6to5server.DowngradeServer(context.Background(), v6server.ProviderServer)

	if err != nil {
		t.Fatalf("unexpected error downgrading server: %s", err)
	}

	_, err = v5server.PlanResourceChange(ctx, &tfprotov5.PlanResourceChangeRequest{
		TypeName: "test_resource",
	})

	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	if !v6server.PlanResourceChangeCalled["test_resource"] {
		t.Errorf("expected test_resource PlanResourceChange to be called")
	}
}

func TestV6ToV5ServerPrepareProviderConfig(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	v6server := &tf6testserver.TestServer{
		GetProviderSchemaResponse: &tfprotov6.GetProviderSchemaResponse{
			ResourceSchemas: map[string]*tfprotov6.Schema{
				"test_resource": {},
			},
		},
	}

	v5server, err := tf6to5server.DowngradeServer(context.Background(), v6server.ProviderServer)

	if err != nil {
		t.Fatalf("unexpected error downgrading server: %s", err)
	}

	_, err = v5server.PrepareProviderConfig(ctx, &tfprotov5.PrepareProviderConfigRequest{})

	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	if !v6server.ValidateProviderConfigCalled {
		t.Errorf("expected ValidateProviderConfig to be called")
	}
}

func TestV6ToV5ServerReadDataSource(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	v6server := &tf6testserver.TestServer{
		GetProviderSchemaResponse: &tfprotov6.GetProviderSchemaResponse{
			DataSourceSchemas: map[string]*tfprotov6.Schema{
				"test_data_source": {},
			},
		},
	}

	v5server, err := tf6to5server.DowngradeServer(context.Background(), v6server.ProviderServer)

	if err != nil {
		t.Fatalf("unexpected error downgrading server: %s", err)
	}

	_, err = v5server.ReadDataSource(ctx, &tfprotov5.ReadDataSourceRequest{
		TypeName: "test_data_source",
	})

	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	if !v6server.ReadDataSourceCalled["test_data_source"] {
		t.Errorf("expected test_data_source ReadDataSource to be called")
	}
}

func TestV6ToV5ServerReadResource(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	v6server := &tf6testserver.TestServer{
		GetProviderSchemaResponse: &tfprotov6.GetProviderSchemaResponse{
			ResourceSchemas: map[string]*tfprotov6.Schema{
				"test_resource": {},
			},
		},
	}

	v5server, err := tf6to5server.DowngradeServer(context.Background(), v6server.ProviderServer)

	if err != nil {
		t.Fatalf("unexpected error downgrading server: %s", err)
	}

	_, err = v5server.ReadResource(ctx, &tfprotov5.ReadResourceRequest{
		TypeName: "test_resource",
	})

	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	if !v6server.ReadResourceCalled["test_resource"] {
		t.Errorf("expected test_resource ReadResource to be called")
	}
}

func TestV6ToV5ServerRenewEphemeralResource(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	v6server := &tf6testserver.TestServer{
		GetProviderSchemaResponse: &tfprotov6.GetProviderSchemaResponse{
			ResourceSchemas: map[string]*tfprotov6.Schema{
				"test_ephemeral_resource": {},
			},
		},
	}

	v5server, err := tf6to5server.DowngradeServer(context.Background(), v6server.ProviderServer)

	if err != nil {
		t.Fatalf("unexpected error downgrading server: %s", err)
	}

	_, err = v5server.RenewEphemeralResource(ctx, &tfprotov5.RenewEphemeralResourceRequest{
		TypeName: "test_ephemeral_resource",
	})

	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	if !v6server.RenewEphemeralResourceCalled["test_ephemeral_resource"] {
		t.Errorf("expected test_ephemeral_resource RenewEphemeralResource to be called")
	}
}

func TestV6ToV5ServerStopProvider(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	v6server := &tf6testserver.TestServer{
		GetProviderSchemaResponse: &tfprotov6.GetProviderSchemaResponse{
			ResourceSchemas: map[string]*tfprotov6.Schema{
				"test_resource": {},
			},
		},
	}

	v5server, err := tf6to5server.DowngradeServer(context.Background(), v6server.ProviderServer)

	if err != nil {
		t.Fatalf("unexpected error downgrading server: %s", err)
	}

	_, err = v5server.StopProvider(ctx, &tfprotov5.StopProviderRequest{})

	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	if !v6server.StopProviderCalled {
		t.Errorf("expected StopProvider to be called")
	}
}

func TestV6ToV5ServerUpgradeResourceState(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	v6server := &tf6testserver.TestServer{
		GetProviderSchemaResponse: &tfprotov6.GetProviderSchemaResponse{
			ResourceSchemas: map[string]*tfprotov6.Schema{
				"test_resource": {},
			},
		},
	}

	v5server, err := tf6to5server.DowngradeServer(context.Background(), v6server.ProviderServer)

	if err != nil {
		t.Fatalf("unexpected error downgrading server: %s", err)
	}

	_, err = v5server.UpgradeResourceState(ctx, &tfprotov5.UpgradeResourceStateRequest{
		TypeName: "test_resource",
	})

	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	if !v6server.UpgradeResourceStateCalled["test_resource"] {
		t.Errorf("expected test_resource UpgradeResourceState to be called")
	}
}

func TestV6ToV5ServerUpgradeResourceIdentity(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	v6server := &tf6testserver.TestServer{
		GetProviderSchemaResponse: &tfprotov6.GetProviderSchemaResponse{
			ResourceSchemas: map[string]*tfprotov6.Schema{
				"test_resource": {},
			},
		},
	}

	v5server, err := tf6to5server.DowngradeServer(context.Background(), v6server.ProviderServer)

	if err != nil {
		t.Fatalf("unexpected error downgrading server: %s", err)
	}

	_, err = v5server.UpgradeResourceIdentity(ctx, &tfprotov5.UpgradeResourceIdentityRequest{
		TypeName: "test_resource",
	})

	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	if !v6server.UpgradeResourceIdentityCalled["test_resource"] {
		t.Errorf("expected test_resource UpgradeResourceState to be called")
	}
}

func TestV6ToV5ServerValidateDataSourceConfig(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	v6server := &tf6testserver.TestServer{
		GetProviderSchemaResponse: &tfprotov6.GetProviderSchemaResponse{
			DataSourceSchemas: map[string]*tfprotov6.Schema{
				"test_data_source": {},
			},
		},
	}

	v5server, err := tf6to5server.DowngradeServer(context.Background(), v6server.ProviderServer)

	if err != nil {
		t.Fatalf("unexpected error downgrading server: %s", err)
	}

	_, err = v5server.ValidateDataSourceConfig(ctx, &tfprotov5.ValidateDataSourceConfigRequest{
		TypeName: "test_data_source",
	})

	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	if !v6server.ValidateDataResourceConfigCalled["test_data_source"] {
		t.Errorf("expected test_data_source ValidateDataResourceConfig to be called")
	}
}

func TestV6ToV5ServerValidateEphemeralResourceConfig(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	v6server := &tf6testserver.TestServer{
		GetProviderSchemaResponse: &tfprotov6.GetProviderSchemaResponse{
			ResourceSchemas: map[string]*tfprotov6.Schema{
				"test_ephemeral_resource": {},
			},
		},
	}

	v5server, err := tf6to5server.DowngradeServer(context.Background(), v6server.ProviderServer)

	if err != nil {
		t.Fatalf("unexpected error downgrading server: %s", err)
	}

	_, err = v5server.ValidateEphemeralResourceConfig(ctx, &tfprotov5.ValidateEphemeralResourceConfigRequest{
		TypeName: "test_ephemeral_resource",
	})

	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	if !v6server.ValidateEphemeralResourceConfigCalled["test_ephemeral_resource"] {
		t.Errorf("expected test_ephemeral_resource ValidateEphemeralResourceConfig to be called")
	}
}

func TestV6ToV5ServerValidateResourceTypeConfig(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	v6server := &tf6testserver.TestServer{
		GetProviderSchemaResponse: &tfprotov6.GetProviderSchemaResponse{
			ResourceSchemas: map[string]*tfprotov6.Schema{
				"test_resource": {},
			},
		},
	}

	v5server, err := tf6to5server.DowngradeServer(context.Background(), v6server.ProviderServer)

	if err != nil {
		t.Fatalf("unexpected error downgrading server: %s", err)
	}

	_, err = v5server.ValidateResourceTypeConfig(ctx, &tfprotov5.ValidateResourceTypeConfigRequest{
		TypeName: "test_resource",
	})

	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	if !v6server.ValidateResourceConfigCalled["test_resource"] {
		t.Errorf("expected test_resource ValidateResourceConfig to be called")
	}
}

func TestV6ToV5ServerValidateListResourceConfig(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	v6server := &tf6testserver.TestServer{
		GetProviderSchemaResponse: &tfprotov6.GetProviderSchemaResponse{
			ResourceSchemas: map[string]*tfprotov6.Schema{
				"test_list_resource": {},
			},
		},
	}

	v5server, err := tf6to5server.DowngradeServer(context.Background(), v6server.ProviderServer)

	if err != nil {
		t.Fatalf("unexpected error downgrading server: %s", err)
	}

	//nolint:staticcheck // Intentionally verifying interface implementation
	listResourceServer, ok := v5server.(tfprotov5.ProviderServerWithListResource)
	if !ok {
		t.Fatal("v6server should implement tfprotov5.ProviderServerWithListResource")
	}

	_, err = listResourceServer.ValidateListResourceConfig(ctx, &tfprotov5.ValidateListResourceConfigRequest{
		TypeName: "test_list_resource",
	})

	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	if !v6server.ValidateListResourceConfigCalled["test_list_resource"] {
		t.Errorf("expected test_list_resource ValidateListResourceConfig to be called")
	}
}

func TestV6ToV5ServerListResource(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	v6server := &tf6testserver.TestServer{
		GetProviderSchemaResponse: &tfprotov6.GetProviderSchemaResponse{
			ResourceSchemas: map[string]*tfprotov6.Schema{
				"test_list_resource": {},
			},
		},
	}

	v5server, err := tf6to5server.DowngradeServer(context.Background(), v6server.ProviderServer)

	if err != nil {
		t.Fatalf("unexpected error downgrading server: %s", err)
	}

	//nolint:staticcheck // Intentionally verifying interface implementation
	listResourceServer, ok := v5server.(tfprotov5.ProviderServerWithListResource)
	if !ok {
		t.Fatal("v6server should implement tfprotov5.ProviderServerWithListResource")
	}

	_, err = listResourceServer.ListResource(ctx, &tfprotov5.ListResourceRequest{
		TypeName: "test_list_resource",
	})

	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	if !v6server.ListResourceCalled["test_list_resource"] {
		t.Errorf("expected test_list_resource ListResource to be called")
	}
}

func TestV6ToV5ServerValidateActionConfig(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	v6server := &tf6testserver.TestServer{
		GetProviderSchemaResponse: &tfprotov6.GetProviderSchemaResponse{
			ResourceSchemas: map[string]*tfprotov6.Schema{
				"test_action": {},
			},
		},
	}

	v5server, err := tf6to5server.DowngradeServer(context.Background(), v6server.ProviderServer)

	if err != nil {
		t.Fatalf("unexpected error downgrading server: %s", err)
	}

	//nolint:staticcheck // Intentionally verifying interface implementation
	actionServer, ok := v5server.(tfprotov5.ProviderServerWithActions)
	if !ok {
		t.Fatal("v5server should implement tfprotov5.ProviderServerWithActions")
	}

	_, err = actionServer.ValidateActionConfig(ctx, &tfprotov5.ValidateActionConfigRequest{
		ActionType: "test_action",
	})

	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	if !v6server.ValidateActionConfigCalled["test_action"] {
		t.Errorf("expected test_action ValidateActionConfig to be called")
	}
}

func TestV6ToV5ServerPlanAction(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	v6server := &tf6testserver.TestServer{
		GetProviderSchemaResponse: &tfprotov6.GetProviderSchemaResponse{
			ActionSchemas: map[string]*tfprotov6.ActionSchema{
				"test_action": {
					Type: tfprotov6.UnlinkedActionSchemaType{},
				},
			},
		},
	}

	v5server, err := tf6to5server.DowngradeServer(context.Background(), v6server.ProviderServer)

	if err != nil {
		t.Fatalf("unexpected error downgrading server: %s", err)
	}

	//nolint:staticcheck // Intentionally verifying interface implementation
	actionServer, ok := v5server.(tfprotov5.ProviderServerWithActions)
	if !ok {
		t.Fatal("v5server should implement tfprotov5.ProviderServerWithActions")
	}

	_, err = actionServer.PlanAction(ctx, &tfprotov5.PlanActionRequest{
		ActionType: "test_action",
	})

	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	if !v6server.PlanActionCalled["test_action"] {
		t.Errorf("expected test_action PlanAction to be called")
	}
}

func TestV6ToV5ServerInvokeAction(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	v6server := &tf6testserver.TestServer{
		GetProviderSchemaResponse: &tfprotov6.GetProviderSchemaResponse{
			ActionSchemas: map[string]*tfprotov6.ActionSchema{
				"test_action": {
					Type: tfprotov6.UnlinkedActionSchemaType{},
				},
			},
		},
	}

	v5server, err := tf6to5server.DowngradeServer(context.Background(), v6server.ProviderServer)

	if err != nil {
		t.Fatalf("unexpected error upgrading server: %s", err)
	}

	//nolint:staticcheck // Intentionally verifying interface implementation
	actionServer, ok := v5server.(tfprotov5.ProviderServerWithActions)
	if !ok {
		t.Fatal("v5server should implement tfprotov5.ProviderServerWithActions")
	}

	_, err = actionServer.InvokeAction(ctx, &tfprotov5.InvokeActionRequest{
		ActionType: "test_action",
	})

	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	if !v6server.InvokeActionCalled["test_action"] {
		t.Errorf("expected test_action InvokeAction to be called")
	}
}
