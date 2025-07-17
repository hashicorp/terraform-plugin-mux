// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package tf5to6server_test

import (
	"context"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-go/tfprotov5"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-go/tftypes"

	"github.com/hashicorp/terraform-plugin-mux/internal/tf5testserver"
	"github.com/hashicorp/terraform-plugin-mux/tf5to6server"
)

func TestUpgradeServer(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		v5Server      func() tfprotov5.ProviderServer
		expectedError error
	}{
		"compatible": {
			v5Server: (&tf5testserver.TestServer{
				GetProviderSchemaResponse: &tfprotov5.GetProviderSchemaResponse{
					ActionSchemas: map[string]*tfprotov5.ActionSchema{
						"test_action": {
							Type: tfprotov5.UnlinkedActionSchemaType{},
						},
					},
					DataSourceSchemas: map[string]*tfprotov5.Schema{
						"test_data_source": {},
					},
					EphemeralResourceSchemas: map[string]*tfprotov5.Schema{
						"test_ephemeral_resource": {},
					},
					Functions: map[string]*tfprotov5.Function{
						"test_function": {},
					},
					ListResourceSchemas: map[string]*tfprotov5.Schema{
						"test_list_resource": {},
					},
					Provider: &tfprotov5.Schema{
						Block: &tfprotov5.SchemaBlock{
							Attributes: []*tfprotov5.SchemaAttribute{
								{
									Name:            "test_attribute",
									Type:            tftypes.String,
									Required:        true,
									Description:     "test_attribute description",
									DescriptionKind: tfprotov5.StringKindPlain,
								},
							},
							BlockTypes: []*tfprotov5.SchemaNestedBlock{
								{
									TypeName: "test_block",
									Nesting:  tfprotov5.SchemaNestedBlockNestingModeList,
									Block: &tfprotov5.SchemaBlock{
										Description:     "test_block description",
										DescriptionKind: tfprotov5.StringKindPlain,
										Attributes: []*tfprotov5.SchemaAttribute{
											{
												Name:            "test_block_attribute",
												Type:            tftypes.Number,
												Required:        true,
												Description:     "test_block_attribute description",
												DescriptionKind: tfprotov5.StringKindPlain,
											},
										},
									},
								},
							},
						},
					},
					ProviderMeta: &tfprotov5.Schema{
						Block: &tfprotov5.SchemaBlock{
							Attributes: []*tfprotov5.SchemaAttribute{
								{
									Name:            "test_attribute",
									Type:            tftypes.String,
									Required:        true,
									Description:     "test_attribute description",
									DescriptionKind: tfprotov5.StringKindPlain,
								},
							},
							BlockTypes: []*tfprotov5.SchemaNestedBlock{
								{
									TypeName: "test_block",
									Nesting:  tfprotov5.SchemaNestedBlockNestingModeList,
									Block: &tfprotov5.SchemaBlock{
										Description:     "test_block description",
										DescriptionKind: tfprotov5.StringKindPlain,
										Attributes: []*tfprotov5.SchemaAttribute{
											{
												Name:            "test_block_attribute",
												Type:            tftypes.Number,
												Required:        true,
												Description:     "test_block_attribute description",
												DescriptionKind: tfprotov5.StringKindPlain,
											},
										},
									},
								},
							},
						},
					},
					ResourceSchemas: map[string]*tfprotov5.Schema{
						"test_resource": {},
					},
				},
			}).ProviderServer,
		},
	}

	for name, testCase := range testCases {

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			_, err := tf5to6server.UpgradeServer(context.Background(), testCase.v5Server)

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

func TestV5ToV6ServerApplyResourceChange(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	v5server := &tf5testserver.TestServer{
		GetProviderSchemaResponse: &tfprotov5.GetProviderSchemaResponse{
			ResourceSchemas: map[string]*tfprotov5.Schema{
				"test_resource": {},
			},
		},
	}

	v6server, err := tf5to6server.UpgradeServer(context.Background(), v5server.ProviderServer)

	if err != nil {
		t.Fatalf("unexpected error downgrading server: %s", err)
	}

	_, err = v6server.ApplyResourceChange(ctx, &tfprotov6.ApplyResourceChangeRequest{
		TypeName: "test_resource",
	})

	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	if !v5server.ApplyResourceChangeCalled["test_resource"] {
		t.Errorf("expected test_resource ApplyResourceChange to be called")
	}
}

func TestV5ToV6ServerCallFunction(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	v5server := &tf5testserver.TestServer{
		GetProviderSchemaResponse: &tfprotov5.GetProviderSchemaResponse{
			Functions: map[string]*tfprotov5.Function{
				"test_function": {},
			},
		},
	}

	v6server, err := tf5to6server.UpgradeServer(context.Background(), v5server.ProviderServer)

	if err != nil {
		t.Fatalf("unexpected error upgrading server: %s", err)
	}

	_, err = v6server.CallFunction(ctx, &tfprotov6.CallFunctionRequest{
		Name: "test_function",
	})

	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	if !v5server.CallFunctionCalled["test_function"] {
		t.Errorf("expected test_function CallFunction to be called")
	}
}

func TestV5ToV6ServerCloseEphemeralResource(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	v5server := &tf5testserver.TestServer{
		GetProviderSchemaResponse: &tfprotov5.GetProviderSchemaResponse{
			EphemeralResourceSchemas: map[string]*tfprotov5.Schema{
				"test_ephemeral_resource": {},
			},
		},
	}

	v6server, err := tf5to6server.UpgradeServer(context.Background(), v5server.ProviderServer)

	if err != nil {
		t.Fatalf("unexpected error downgrading server: %s", err)
	}

	_, err = v6server.CloseEphemeralResource(ctx, &tfprotov6.CloseEphemeralResourceRequest{
		TypeName: "test_ephemeral_resource",
	})

	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	if !v5server.CloseEphemeralResourceCalled["test_ephemeral_resource"] {
		t.Errorf("expected test_ephemeral_resource CloseEphemeralResource to be called")
	}
}

func TestV5ToV6ServerConfigureProvider(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	v5server := &tf5testserver.TestServer{
		GetProviderSchemaResponse: &tfprotov5.GetProviderSchemaResponse{
			ResourceSchemas: map[string]*tfprotov5.Schema{
				"test_resource": {},
			},
		},
	}

	v6server, err := tf5to6server.UpgradeServer(context.Background(), v5server.ProviderServer)

	if err != nil {
		t.Fatalf("unexpected error downgrading server: %s", err)
	}

	_, err = v6server.ConfigureProvider(ctx, &tfprotov6.ConfigureProviderRequest{})

	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	if !v5server.ConfigureProviderCalled {
		t.Errorf("expected ConfigureProvider to be called")
	}
}

func TestV5ToV6ServerGetFunctions(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	v5server := &tf5testserver.TestServer{
		GetFunctionsResponse: &tfprotov5.GetFunctionsResponse{
			Functions: map[string]*tfprotov5.Function{
				"test_function": {},
			},
		},
	}

	v6server, err := tf5to6server.UpgradeServer(context.Background(), v5server.ProviderServer)

	if err != nil {
		t.Fatalf("unexpected error upgrading server: %s", err)
	}

	_, err = v6server.GetFunctions(ctx, &tfprotov6.GetFunctionsRequest{})

	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	if !v5server.GetFunctionsCalled {
		t.Errorf("expected GetFunctions to be called")
	}
}

func TestV5ToV6ServerGetMetadata(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	v5server := &tf5testserver.TestServer{
		GetMetadataResponse: &tfprotov5.GetMetadataResponse{
			Resources: []tfprotov5.ResourceMetadata{
				{
					TypeName: "test_resource",
				},
			},
		},
	}

	v6server, err := tf5to6server.UpgradeServer(context.Background(), v5server.ProviderServer)

	if err != nil {
		t.Fatalf("unexpected error downgrading server: %s", err)
	}

	_, err = v6server.GetMetadata(ctx, &tfprotov6.GetMetadataRequest{})

	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	if !v5server.GetMetadataCalled {
		t.Errorf("expected GetMetadata to be called")
	}
}

func TestV5ToV6ServerGetProviderSchema(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	v5server := &tf5testserver.TestServer{
		GetProviderSchemaResponse: &tfprotov5.GetProviderSchemaResponse{
			ResourceSchemas: map[string]*tfprotov5.Schema{
				"test_resource": {},
			},
		},
	}

	v6server, err := tf5to6server.UpgradeServer(context.Background(), v5server.ProviderServer)

	if err != nil {
		t.Fatalf("unexpected error downgrading server: %s", err)
	}

	_, err = v6server.GetProviderSchema(ctx, &tfprotov6.GetProviderSchemaRequest{})

	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	if !v5server.GetProviderSchemaCalled {
		t.Errorf("expected GetProviderSchema to be called")
	}
}

func TestV5ToV6ServerGetResourceIdentitySchemas(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	v5server := &tf5testserver.TestServer{
		GetResourceIdentitySchemasResponse: &tfprotov5.GetResourceIdentitySchemasResponse{},
	}

	v6server, err := tf5to6server.UpgradeServer(context.Background(), v5server.ProviderServer)

	if err != nil {
		t.Fatalf("unexpected error downgrading server: %s", err)
	}

	_, err = v6server.GetResourceIdentitySchemas(ctx, &tfprotov6.GetResourceIdentitySchemasRequest{})

	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	if !v5server.GetResourceIdentitySchemasCalled {
		t.Errorf("expected GetResourceIdentitySchemas to be called")
	}
}

func TestV5ToV6ServerImportResourceState(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	v5server := &tf5testserver.TestServer{
		GetProviderSchemaResponse: &tfprotov5.GetProviderSchemaResponse{
			ResourceSchemas: map[string]*tfprotov5.Schema{
				"test_resource": {},
			},
		},
	}

	v6server, err := tf5to6server.UpgradeServer(context.Background(), v5server.ProviderServer)

	if err != nil {
		t.Fatalf("unexpected error downgrading server: %s", err)
	}

	_, err = v6server.ImportResourceState(ctx, &tfprotov6.ImportResourceStateRequest{
		TypeName: "test_resource",
	})

	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	if !v5server.ImportResourceStateCalled["test_resource"] {
		t.Errorf("expected test_resource ImportResourceState to be called")
	}
}

func TestV5ToV6ServerMoveResourceState(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	v5server := &tf5testserver.TestServer{
		GetProviderSchemaResponse: &tfprotov5.GetProviderSchemaResponse{
			ResourceSchemas: map[string]*tfprotov5.Schema{
				"test_resource": {},
			},
		},
	}

	v6server, err := tf5to6server.UpgradeServer(context.Background(), v5server.ProviderServer)

	if err != nil {
		t.Fatalf("unexpected error downgrading server: %s", err)
	}

	_, err = v6server.MoveResourceState(ctx, &tfprotov6.MoveResourceStateRequest{
		TargetTypeName: "test_resource",
	})

	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	if !v5server.MoveResourceStateCalled["test_resource"] {
		t.Errorf("expected test_resource MoveResourceState to be called")
	}
}

func TestV5ToV6ServerOpenEphemeralResource(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	v5server := &tf5testserver.TestServer{
		GetProviderSchemaResponse: &tfprotov5.GetProviderSchemaResponse{
			EphemeralResourceSchemas: map[string]*tfprotov5.Schema{
				"test_ephemeral_resource": {},
			},
		},
	}

	v6server, err := tf5to6server.UpgradeServer(context.Background(), v5server.ProviderServer)

	if err != nil {
		t.Fatalf("unexpected error downgrading server: %s", err)
	}

	_, err = v6server.OpenEphemeralResource(ctx, &tfprotov6.OpenEphemeralResourceRequest{
		TypeName: "test_ephemeral_resource",
	})

	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	if !v5server.OpenEphemeralResourceCalled["test_ephemeral_resource"] {
		t.Errorf("expected test_ephemeral_resource OpenEphemeralResource to be called")
	}
}

func TestV5ToV6ServerPlanResourceChange(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	v5server := &tf5testserver.TestServer{
		GetProviderSchemaResponse: &tfprotov5.GetProviderSchemaResponse{
			ResourceSchemas: map[string]*tfprotov5.Schema{
				"test_resource": {},
			},
		},
	}

	v6server, err := tf5to6server.UpgradeServer(context.Background(), v5server.ProviderServer)

	if err != nil {
		t.Fatalf("unexpected error downgrading server: %s", err)
	}

	_, err = v6server.PlanResourceChange(ctx, &tfprotov6.PlanResourceChangeRequest{
		TypeName: "test_resource",
	})

	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	if !v5server.PlanResourceChangeCalled["test_resource"] {
		t.Errorf("expected test_resource PlanResourceChange to be called")
	}
}

func TestV5ToV6ServerReadDataSource(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	v5server := &tf5testserver.TestServer{
		GetProviderSchemaResponse: &tfprotov5.GetProviderSchemaResponse{
			DataSourceSchemas: map[string]*tfprotov5.Schema{
				"test_data_source": {},
			},
		},
	}

	v6server, err := tf5to6server.UpgradeServer(context.Background(), v5server.ProviderServer)

	if err != nil {
		t.Fatalf("unexpected error downgrading server: %s", err)
	}

	_, err = v6server.ReadDataSource(ctx, &tfprotov6.ReadDataSourceRequest{
		TypeName: "test_data_source",
	})

	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	if !v5server.ReadDataSourceCalled["test_data_source"] {
		t.Errorf("expected test_data_source ReadDataSource to be called")
	}
}

func TestV5ToV6ServerReadResource(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	v5server := &tf5testserver.TestServer{
		GetProviderSchemaResponse: &tfprotov5.GetProviderSchemaResponse{
			ResourceSchemas: map[string]*tfprotov5.Schema{
				"test_resource": {},
			},
		},
	}

	v6server, err := tf5to6server.UpgradeServer(context.Background(), v5server.ProviderServer)

	if err != nil {
		t.Fatalf("unexpected error downgrading server: %s", err)
	}

	_, err = v6server.ReadResource(ctx, &tfprotov6.ReadResourceRequest{
		TypeName: "test_resource",
	})

	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	if !v5server.ReadResourceCalled["test_resource"] {
		t.Errorf("expected test_resource ReadResource to be called")
	}
}

func TestV5ToV6ServerRenewEphemeralResource(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	v5server := &tf5testserver.TestServer{
		GetProviderSchemaResponse: &tfprotov5.GetProviderSchemaResponse{
			EphemeralResourceSchemas: map[string]*tfprotov5.Schema{
				"test_ephemeral_resource": {},
			},
		},
	}

	v6server, err := tf5to6server.UpgradeServer(context.Background(), v5server.ProviderServer)

	if err != nil {
		t.Fatalf("unexpected error downgrading server: %s", err)
	}

	_, err = v6server.RenewEphemeralResource(ctx, &tfprotov6.RenewEphemeralResourceRequest{
		TypeName: "test_ephemeral_resource",
	})

	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	if !v5server.RenewEphemeralResourceCalled["test_ephemeral_resource"] {
		t.Errorf("expected test_ephemeral_resource RenewEphemeralResource to be called")
	}
}

func TestV5ToV6ServerStopProvider(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	v5server := &tf5testserver.TestServer{
		GetProviderSchemaResponse: &tfprotov5.GetProviderSchemaResponse{
			ResourceSchemas: map[string]*tfprotov5.Schema{
				"test_resource": {},
			},
		},
	}

	v6server, err := tf5to6server.UpgradeServer(context.Background(), v5server.ProviderServer)

	if err != nil {
		t.Fatalf("unexpected error downgrading server: %s", err)
	}

	_, err = v6server.StopProvider(ctx, &tfprotov6.StopProviderRequest{})

	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	if !v5server.StopProviderCalled {
		t.Errorf("expected StopProvider to be called")
	}
}

func TestV5ToV6ServerUpgradeResourceState(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	v5server := &tf5testserver.TestServer{
		GetProviderSchemaResponse: &tfprotov5.GetProviderSchemaResponse{
			ResourceSchemas: map[string]*tfprotov5.Schema{
				"test_resource": {},
			},
		},
	}

	v6server, err := tf5to6server.UpgradeServer(context.Background(), v5server.ProviderServer)

	if err != nil {
		t.Fatalf("unexpected error downgrading server: %s", err)
	}

	_, err = v6server.UpgradeResourceState(ctx, &tfprotov6.UpgradeResourceStateRequest{
		TypeName: "test_resource",
	})

	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	if !v5server.UpgradeResourceStateCalled["test_resource"] {
		t.Errorf("expected test_resource UpgradeResourceState to be called")
	}
}

func TestV5ToV6ServerUpgradeResourceIdentity(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	v5server := &tf5testserver.TestServer{
		GetProviderSchemaResponse: &tfprotov5.GetProviderSchemaResponse{
			ResourceSchemas: map[string]*tfprotov5.Schema{
				"test_resource": {},
			},
		},
	}

	v6server, err := tf5to6server.UpgradeServer(context.Background(), v5server.ProviderServer)

	if err != nil {
		t.Fatalf("unexpected error downgrading server: %s", err)
	}

	_, err = v6server.UpgradeResourceIdentity(ctx, &tfprotov6.UpgradeResourceIdentityRequest{
		TypeName: "test_resource",
	})

	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	if !v5server.UpgradeResourceIdentityCalled["test_resource"] {
		t.Errorf("expected test_resource UpgradeResourceState to be called")
	}
}

func TestV5ToV6ServerValidateDataResourceConfig(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	v5server := &tf5testserver.TestServer{
		GetProviderSchemaResponse: &tfprotov5.GetProviderSchemaResponse{
			DataSourceSchemas: map[string]*tfprotov5.Schema{
				"test_data_source": {},
			},
		},
	}

	v6server, err := tf5to6server.UpgradeServer(context.Background(), v5server.ProviderServer)

	if err != nil {
		t.Fatalf("unexpected error downgrading server: %s", err)
	}

	_, err = v6server.ValidateDataResourceConfig(ctx, &tfprotov6.ValidateDataResourceConfigRequest{
		TypeName: "test_data_source",
	})

	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	if !v5server.ValidateDataSourceConfigCalled["test_data_source"] {
		t.Errorf("expected test_data_source ValidateDataSourceConfig to be called")
	}
}

func TestV5ToV6ServerValidateEphemeralResourceConfig(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	v5server := &tf5testserver.TestServer{
		GetProviderSchemaResponse: &tfprotov5.GetProviderSchemaResponse{
			EphemeralResourceSchemas: map[string]*tfprotov5.Schema{
				"test_resource": {},
			},
		},
	}

	v6server, err := tf5to6server.UpgradeServer(context.Background(), v5server.ProviderServer)

	if err != nil {
		t.Fatalf("unexpected error downgrading server: %s", err)
	}

	_, err = v6server.ValidateEphemeralResourceConfig(ctx, &tfprotov6.ValidateEphemeralResourceConfigRequest{
		TypeName: "test_resource",
	})

	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	if !v5server.ValidateEphemeralResourceConfigCalled["test_resource"] {
		t.Errorf("expected test_resource ValidateEphemeralResourceConfig to be called")
	}
}

func TestV5ToV6ServerValidateProviderConfig(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	v5server := &tf5testserver.TestServer{
		GetProviderSchemaResponse: &tfprotov5.GetProviderSchemaResponse{
			ResourceSchemas: map[string]*tfprotov5.Schema{
				"test_resource": {},
			},
		},
	}

	v6server, err := tf5to6server.UpgradeServer(context.Background(), v5server.ProviderServer)

	if err != nil {
		t.Fatalf("unexpected error downgrading server: %s", err)
	}

	_, err = v6server.ValidateProviderConfig(ctx, &tfprotov6.ValidateProviderConfigRequest{})

	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	if !v5server.PrepareProviderConfigCalled {
		t.Errorf("expected PrepareProviderConfig to be called")
	}
}

func TestV5ToV6ServerValidateResourceConfig(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	v5server := &tf5testserver.TestServer{
		GetProviderSchemaResponse: &tfprotov5.GetProviderSchemaResponse{
			ResourceSchemas: map[string]*tfprotov5.Schema{
				"test_resource": {},
			},
		},
	}

	v6server, err := tf5to6server.UpgradeServer(context.Background(), v5server.ProviderServer)

	if err != nil {
		t.Fatalf("unexpected error downgrading server: %s", err)
	}

	_, err = v6server.ValidateResourceConfig(ctx, &tfprotov6.ValidateResourceConfigRequest{
		TypeName: "test_resource",
	})

	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	if !v5server.ValidateResourceTypeConfigCalled["test_resource"] {
		t.Errorf("expected test_resource ValidateResourceConfig to be called")
	}
}

func TestV5ToV6ServerValidateListResourceConfig(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	v5server := &tf5testserver.TestServer{
		GetProviderSchemaResponse: &tfprotov5.GetProviderSchemaResponse{
			ListResourceSchemas: map[string]*tfprotov5.Schema{
				"test_list_resource": {},
			},
		},
	}

	v6server, err := tf5to6server.UpgradeServer(context.Background(), v5server.ProviderServer)

	if err != nil {
		t.Fatalf("unexpected error upgrading server: %s", err)
	}

	//nolint:staticcheck // Intentionally verifying interface implementation
	listResourceServer, ok := v6server.(tfprotov6.ProviderServerWithListResource)
	if !ok {
		t.Fatal("v6server should implement tfprotov6.ProviderServerWithListResource")
	}

	_, err = listResourceServer.ValidateListResourceConfig(ctx, &tfprotov6.ValidateListResourceConfigRequest{
		TypeName: "test_list_resource",
	})

	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	if !v5server.ValidateListResourceConfigCalled["test_list_resource"] {
		t.Errorf("expected test_list_resource ValidateListResourceConfig to be called")
	}
}

func TestV5ToV6ServerListResource(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	v5server := &tf5testserver.TestServer{
		GetProviderSchemaResponse: &tfprotov5.GetProviderSchemaResponse{
			ListResourceSchemas: map[string]*tfprotov5.Schema{
				"test_list_resource": {},
			},
		},
	}

	v6server, err := tf5to6server.UpgradeServer(context.Background(), v5server.ProviderServer)

	if err != nil {
		t.Fatalf("unexpected error upgrading server: %s", err)
	}

	//nolint:staticcheck // Intentionally verifying interface implementation
	listResourceServer, ok := v6server.(tfprotov6.ProviderServerWithListResource)
	if !ok {
		t.Fatal("v6server should implement tfprotov6.ProviderServerWithListResource")
	}

	_, err = listResourceServer.ListResource(ctx, &tfprotov6.ListResourceRequest{
		TypeName: "test_list_resource",
	})

	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	if !v5server.ListResourceCalled["test_list_resource"] {
		t.Errorf("expected test_list_resource ListResourceConfig to be called")
	}
}

func TestV5ToV6ServerValidateActionConfig(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	v5server := &tf5testserver.TestServer{
		GetProviderSchemaResponse: &tfprotov5.GetProviderSchemaResponse{
			ActionSchemas: map[string]*tfprotov5.ActionSchema{
				"test_action": {},
			},
		},
	}

	v6server, err := tf5to6server.UpgradeServer(context.Background(), v5server.ProviderServer)

	if err != nil {
		t.Fatalf("unexpected error downgrading server: %s", err)
	}

	//nolint:staticcheck // Intentionally verifying interface implementation
	actionServer, ok := v6server.(tfprotov6.ProviderServerWithActions)
	if !ok {
		t.Fatal("v6server should implement tfprotov6.ProviderServerWithActions")
	}

	_, err = actionServer.ValidateActionConfig(ctx, &tfprotov6.ValidateActionConfigRequest{
		ActionType: "test_action",
	})

	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	if !v5server.ValidateActionConfigCalled["test_action"] {
		t.Errorf("expected test_action ValidateActionConfig to be called")
	}
}

func TestV5ToV6ServerPlanAction(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	v5server := &tf5testserver.TestServer{
		GetProviderSchemaResponse: &tfprotov5.GetProviderSchemaResponse{
			ActionSchemas: map[string]*tfprotov5.ActionSchema{
				"test_action": {
					Type: tfprotov5.UnlinkedActionSchemaType{},
				},
			},
		},
	}

	v6server, err := tf5to6server.UpgradeServer(context.Background(), v5server.ProviderServer)

	if err != nil {
		t.Fatalf("unexpected error downgrading server: %s", err)
	}

	//nolint:staticcheck // Intentionally verifying interface implementation
	actionServer, ok := v6server.(tfprotov6.ProviderServerWithActions)
	if !ok {
		t.Fatal("v6server should implement tfprotov6.ProviderServerWithActions")
	}

	_, err = actionServer.PlanAction(ctx, &tfprotov6.PlanActionRequest{
		ActionType: "test_action",
	})

	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	if !v5server.PlanActionCalled["test_action"] {
		t.Errorf("expected test_action PlanAction to be called")
	}
}

func TestV5ToV6ServerInvokeAction(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	v5server := &tf5testserver.TestServer{
		GetProviderSchemaResponse: &tfprotov5.GetProviderSchemaResponse{
			ActionSchemas: map[string]*tfprotov5.ActionSchema{
				"test_action": {
					Type: tfprotov5.UnlinkedActionSchemaType{},
				},
			},
		},
	}

	v6server, err := tf5to6server.UpgradeServer(context.Background(), v5server.ProviderServer)

	if err != nil {
		t.Fatalf("unexpected error upgrading server: %s", err)
	}

	//nolint:staticcheck // Intentionally verifying interface implementation
	actionServer, ok := v6server.(tfprotov6.ProviderServerWithActions)
	if !ok {
		t.Fatal("v6server should implement tfprotov6.ProviderServerWithActions")
	}

	_, err = actionServer.InvokeAction(ctx, &tfprotov6.InvokeActionRequest{
		ActionType: "test_action",
	})

	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	if !v5server.InvokeActionCalled["test_action"] {
		t.Errorf("expected test_action InvokeAction to be called")
	}
}
