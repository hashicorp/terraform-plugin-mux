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
				DataSourceSchemas: map[string]*tfprotov6.Schema{
					"test_data_source": {},
				},
				ProviderSchema: &tfprotov6.Schema{
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
				ProviderMetaSchema: &tfprotov6.Schema{
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
			}).ProviderServer,
		},
		"SchemaAttribute-NestedType-data-source": {
			v6Server: (&tf6testserver.TestServer{
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
			}).ProviderServer,
			expectedError: fmt.Errorf("unable to convert data source \"test_data_source\" schema: unable to convert attribute \"test_attribute\" schema: %s", tfprotov6tov5.ErrSchemaAttributeNestedTypeNotImplemented),
		},
		"SchemaAttribute-NestedType-provider": {
			v6Server: (&tf6testserver.TestServer{
				ProviderSchema: &tfprotov6.Schema{
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
			}).ProviderServer,
			expectedError: fmt.Errorf("unable to convert provider schema: unable to convert attribute \"test_attribute\" schema: %s", tfprotov6tov5.ErrSchemaAttributeNestedTypeNotImplemented),
		},
		"SchemaAttribute-NestedType-provider-meta": {
			v6Server: (&tf6testserver.TestServer{
				ProviderMetaSchema: &tfprotov6.Schema{
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
			}).ProviderServer,
			expectedError: fmt.Errorf("unable to convert provider meta schema: unable to convert attribute \"test_attribute\" schema: %s", tfprotov6tov5.ErrSchemaAttributeNestedTypeNotImplemented),
		},
		"SchemaAttribute-NestedType-resource": {
			v6Server: (&tf6testserver.TestServer{
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
			}).ProviderServer,
			expectedError: fmt.Errorf("unable to convert resource \"test_resource\" schema: unable to convert attribute \"test_attribute\" schema: %s", tfprotov6tov5.ErrSchemaAttributeNestedTypeNotImplemented),
		},
	}

	for name, testCase := range testCases {
		name, testCase := name, testCase

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
		ResourceSchemas: map[string]*tfprotov6.Schema{
			"test_resource": {},
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

func TestV6ToV5ServerConfigureProvider(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	v6server := &tf6testserver.TestServer{
		ResourceSchemas: map[string]*tfprotov6.Schema{
			"test_resource": {},
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

func TestV6ToV5ServerGetProviderSchema(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	v6server := &tf6testserver.TestServer{
		ResourceSchemas: map[string]*tfprotov6.Schema{
			"test_resource": {},
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

func TestV6ToV5ServerImportResourceState(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	v6server := &tf6testserver.TestServer{
		ResourceSchemas: map[string]*tfprotov6.Schema{
			"test_resource": {},
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

func TestV6ToV5ServerPlanResourceChange(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	v6server := &tf6testserver.TestServer{
		ResourceSchemas: map[string]*tfprotov6.Schema{
			"test_resource": {},
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
		ResourceSchemas: map[string]*tfprotov6.Schema{
			"test_resource": {},
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
		DataSourceSchemas: map[string]*tfprotov6.Schema{
			"test_data_source": {},
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
		ResourceSchemas: map[string]*tfprotov6.Schema{
			"test_resource": {},
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

func TestV6ToV5ServerStopProvider(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	v6server := &tf6testserver.TestServer{
		ResourceSchemas: map[string]*tfprotov6.Schema{
			"test_resource": {},
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
		ResourceSchemas: map[string]*tfprotov6.Schema{
			"test_resource": {},
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

func TestV6ToV5ServerValidateDataSourceConfig(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	v6server := &tf6testserver.TestServer{
		DataSourceSchemas: map[string]*tfprotov6.Schema{
			"test_data_source": {},
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

func TestV6ToV5ServerValidateResourceTypeConfig(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	v6server := &tf6testserver.TestServer{
		ResourceSchemas: map[string]*tfprotov6.Schema{
			"test_resource": {},
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
