package tfmux

import (
	"context"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/terraform-plugin-go/tfprotov5"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

var _ tfprotov5.ProviderServer = &testServer{}

func testFactory(s *testServer) func() tfprotov5.ProviderServer {
	return func() tfprotov5.ProviderServer {
		return s
	}
}

type testServer struct {
	providerSchema     *tfprotov5.Schema
	providerMetaSchema *tfprotov5.Schema
	resourceSchemas    map[string]*tfprotov5.Schema
	dataSourceSchemas  map[string]*tfprotov5.Schema

	resourcesCalled   map[string]bool
	dataSourcesCalled map[string]bool
	configureCalled   bool
	stopCalled        bool
	stopError         string

	respondToPrepareProviderConfig bool
	errorOnPrepareProviderConfig   bool
	warnOnPrepareProviderConfig    bool
}

func (s *testServer) GetProviderSchema(_ context.Context, _ *tfprotov5.GetProviderSchemaRequest) (*tfprotov5.GetProviderSchemaResponse, error) {
	resources := s.resourceSchemas
	if resources == nil {
		resources = map[string]*tfprotov5.Schema{}
	}
	dataSources := s.dataSourceSchemas
	if dataSources == nil {
		dataSources = map[string]*tfprotov5.Schema{}
	}
	return &tfprotov5.GetProviderSchemaResponse{
		Provider:          s.providerSchema,
		ProviderMeta:      s.providerMetaSchema,
		ResourceSchemas:   resources,
		DataSourceSchemas: dataSources,
	}, nil
}

func (s *testServer) PrepareProviderConfig(_ context.Context, req *tfprotov5.PrepareProviderConfigRequest) (*tfprotov5.PrepareProviderConfigResponse, error) {
	if s.respondToPrepareProviderConfig {
		return &tfprotov5.PrepareProviderConfigResponse{
			PreparedConfig: req.Config,
		}, nil
	}
	if s.errorOnPrepareProviderConfig {
		return &tfprotov5.PrepareProviderConfigResponse{
			Diagnostics: []*tfprotov5.Diagnostic{
				{
					Severity: tfprotov5.DiagnosticSeverityError,
					Summary:  "hardcoded error for testing",
					Detail:   "testing that we only get errors from one thing",
				},
			},
		}, nil
	}
	if s.warnOnPrepareProviderConfig {
		return &tfprotov5.PrepareProviderConfigResponse{
			Diagnostics: []*tfprotov5.Diagnostic{
				{
					Severity: tfprotov5.DiagnosticSeverityWarning,
					Summary:  "hardcoded warning for testing",
					Detail:   "testing that we only get warnings from one thing",
				},
			},
		}, nil
	}
	return nil, nil
}

func (s *testServer) StopProvider(_ context.Context, _ *tfprotov5.StopProviderRequest) (*tfprotov5.StopProviderResponse, error) {
	s.stopCalled = true
	if s.stopError != "" {
		return &tfprotov5.StopProviderResponse{
			Error: s.stopError,
		}, nil
	}
	return &tfprotov5.StopProviderResponse{}, nil
}

func (s *testServer) ConfigureProvider(_ context.Context, _ *tfprotov5.ConfigureProviderRequest) (*tfprotov5.ConfigureProviderResponse, error) {
	s.configureCalled = true
	return &tfprotov5.ConfigureProviderResponse{}, nil
}

func (s *testServer) ValidateResourceTypeConfig(_ context.Context, req *tfprotov5.ValidateResourceTypeConfigRequest) (*tfprotov5.ValidateResourceTypeConfigResponse, error) {
	s.resourcesCalled[req.TypeName] = true
	return nil, nil
}

func (s *testServer) UpgradeResourceState(_ context.Context, req *tfprotov5.UpgradeResourceStateRequest) (*tfprotov5.UpgradeResourceStateResponse, error) {
	s.resourcesCalled[req.TypeName] = true
	return nil, nil
}

func (s *testServer) ReadResource(_ context.Context, req *tfprotov5.ReadResourceRequest) (*tfprotov5.ReadResourceResponse, error) {
	s.resourcesCalled[req.TypeName] = true
	return nil, nil
}

func (s *testServer) PlanResourceChange(_ context.Context, req *tfprotov5.PlanResourceChangeRequest) (*tfprotov5.PlanResourceChangeResponse, error) {
	s.resourcesCalled[req.TypeName] = true
	return nil, nil
}

func (s *testServer) ApplyResourceChange(_ context.Context, req *tfprotov5.ApplyResourceChangeRequest) (*tfprotov5.ApplyResourceChangeResponse, error) {
	s.resourcesCalled[req.TypeName] = true
	return nil, nil
}

func (s *testServer) ImportResourceState(_ context.Context, req *tfprotov5.ImportResourceStateRequest) (*tfprotov5.ImportResourceStateResponse, error) {
	s.resourcesCalled[req.TypeName] = true
	return nil, nil
}

func (s *testServer) ValidateDataSourceConfig(_ context.Context, req *tfprotov5.ValidateDataSourceConfigRequest) (*tfprotov5.ValidateDataSourceConfigResponse, error) {
	s.dataSourcesCalled[req.TypeName] = true
	return nil, nil
}

func (s *testServer) ReadDataSource(_ context.Context, req *tfprotov5.ReadDataSourceRequest) (*tfprotov5.ReadDataSourceResponse, error) {
	s.dataSourcesCalled[req.TypeName] = true
	return nil, nil
}

func TestSchemaServerGetProviderSchema_combined(t *testing.T) {
	server1 := testFactory(&testServer{
		providerSchema: &tfprotov5.Schema{
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
		providerMetaSchema: &tfprotov5.Schema{
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
		resourceSchemas: map[string]*tfprotov5.Schema{
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
		dataSourceSchemas: map[string]*tfprotov5.Schema{
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
	})
	server2 := testFactory(&testServer{
		providerSchema: &tfprotov5.Schema{
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
		providerMetaSchema: &tfprotov5.Schema{
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
		resourceSchemas: map[string]*tfprotov5.Schema{
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
		dataSourceSchemas: map[string]*tfprotov5.Schema{
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
	})

	expectedProviderSchema := &tfprotov5.Schema{
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
	}
	expectedProviderMetaSchema := &tfprotov5.Schema{
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
	}
	expectedResourceSchemas := map[string]*tfprotov5.Schema{
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
	}
	expectedDataSourceSchemas := map[string]*tfprotov5.Schema{
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
	}

	factory, err := NewSchemaServerFactory(context.Background(), server1, server2)
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	resp, err := factory.Server().GetProviderSchema(context.Background(), &tfprotov5.GetProviderSchemaRequest{})
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	if diff := cmp.Diff(resp.Provider, expectedProviderSchema); diff != "" {
		t.Errorf("provider schema didn't match expectations: %s", diff)
	}

	if diff := cmp.Diff(resp.ProviderMeta, expectedProviderMetaSchema); diff != "" {
		t.Errorf("provider_meta schema didn't match expectations: %s", diff)
	}

	if diff := cmp.Diff(resp.ResourceSchemas, expectedResourceSchemas); diff != "" {
		t.Errorf("resource schemas didn't match expectations: %s", diff)
	}

	if diff := cmp.Diff(resp.DataSourceSchemas, expectedDataSourceSchemas); diff != "" {
		t.Errorf("data source schemas didn't match expectations: %s", diff)
	}
}

func TestSchemaServerGetProviderSchema_errorDuplicateResource(t *testing.T) {
	server1 := testFactory(&testServer{
		resourceSchemas: map[string]*tfprotov5.Schema{
			"test_foo": {},
		},
	})
	server2 := testFactory(&testServer{
		resourceSchemas: map[string]*tfprotov5.Schema{
			"test_foo": {},
		},
	})

	_, err := NewSchemaServerFactory(context.Background(), server1, server2)
	if !strings.Contains(err.Error(), "resource \"test_foo\" supported by multiple server implementations") {
		t.Errorf("expected error about duplicated resources, got %q", err)
	}
}

func TestSchemaServerGetProviderSchema_errorDuplicateDataSource(t *testing.T) {
	server1 := testFactory(&testServer{
		dataSourceSchemas: map[string]*tfprotov5.Schema{
			"test_foo": {},
		},
	})
	server2 := testFactory(&testServer{
		dataSourceSchemas: map[string]*tfprotov5.Schema{
			"test_foo": {},
		},
	})

	_, err := NewSchemaServerFactory(context.Background(), server1, server2)
	if !strings.Contains(err.Error(), "data source \"test_foo\" supported by multiple server implementations") {
		t.Errorf("expected error about duplicated data sources, got %q", err)
	}
}

func TestSchemaServerGetProviderSchema_providerOutOfOrder(t *testing.T) {
	server1 := testFactory(&testServer{
		providerSchema: &tfprotov5.Schema{
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
	})
	server2 := testFactory(&testServer{
		providerSchema: &tfprotov5.Schema{
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
	})

	_, err := NewSchemaServerFactory(context.Background(), server1, server2)
	if err != nil {
		t.Error(err)
	}
}

func TestSchemaServerGetProviderSchema_providerMetaOutOfOrder(t *testing.T) {
	server1 := testFactory(&testServer{
		providerMetaSchema: &tfprotov5.Schema{
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
	})
	server2 := testFactory(&testServer{
		providerMetaSchema: &tfprotov5.Schema{
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
	})

	_, err := NewSchemaServerFactory(context.Background(), server1, server2)
	if err != nil {
		t.Error(err)
	}
}

func TestSchemaServerGetProviderSchema_errorProviderMismatch(t *testing.T) {
	server1 := testFactory(&testServer{
		providerSchema: &tfprotov5.Schema{
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
	})
	server2 := testFactory(&testServer{
		providerSchema: &tfprotov5.Schema{
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
	})

	_, err := NewSchemaServerFactory(context.Background(), server1, server2)
	if !strings.Contains(err.Error(), "got a different provider schema from two servers") {
		t.Errorf("expected error about mismatched provider schemas, got %q", err)
	}
}

func TestSchemaServerGetProviderSchema_errorProviderMetaMismatch(t *testing.T) {
	server1 := testFactory(&testServer{
		providerMetaSchema: &tfprotov5.Schema{
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
	})
	server2 := testFactory(&testServer{
		providerMetaSchema: &tfprotov5.Schema{
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
	})

	_, err := NewSchemaServerFactory(context.Background(), server1, server2)
	if !strings.Contains(err.Error(), "got a different provider_meta schema from two servers") {
		t.Errorf("expected error about mismatched provider_meta schemas, got %q", err)
	}
}

func TestSchemaServerPrepareProviderConfig_errorMultipleResponses(t *testing.T) {
	config, err := tfprotov5.NewDynamicValue(tftypes.Object{
		AttributeTypes: map[string]tftypes.Type{
			"hello": tftypes.String,
		},
	}, tftypes.NewValue(tftypes.Object{
		AttributeTypes: map[string]tftypes.Type{
			"hello": tftypes.String,
		},
	}, map[string]tftypes.Value{
		"hello": tftypes.NewValue(tftypes.String, "world"),
	}))
	if err != nil {
		t.Fatalf("error constructing config: %s", err)
	}

	server1 := testFactory(&testServer{
		respondToPrepareProviderConfig: true,
	})
	server2 := testFactory(&testServer{})
	server3 := testFactory(&testServer{})
	factory, err := NewSchemaServerFactory(context.Background(), server1, server2, server3)
	if err != nil {
		t.Fatalf("error setting up muxer: %s", err)
	}
	_, err = factory.Server().PrepareProviderConfig(context.Background(), &tfprotov5.PrepareProviderConfigRequest{
		Config: &config,
	})
	if err != nil {
		t.Errorf("unexpected error when only one server replied to PrepareProviderConfig: %s", err)
	}

	server1 = testFactory(&testServer{})
	server2 = testFactory(&testServer{})
	server3 = testFactory(&testServer{})
	factory, err = NewSchemaServerFactory(context.Background(), server1, server2, server3)
	if err != nil {
		t.Fatalf("error setting up muxer: %s", err)
	}
	_, err = factory.Server().PrepareProviderConfig(context.Background(), &tfprotov5.PrepareProviderConfigRequest{
		Config: &config,
	})
	if err != nil {
		t.Errorf("unexpected error when no servers replied to PrepareProviderConfig: %s", err)
	}

	for _, fs := range [][]func() tfprotov5.ProviderServer{
		{
			testFactory(&testServer{
				respondToPrepareProviderConfig: true,
			}),
			testFactory(&testServer{}),
			testFactory(&testServer{
				respondToPrepareProviderConfig: true,
			}),
		},
		{
			testFactory(&testServer{
				respondToPrepareProviderConfig: true,
			}),
			testFactory(&testServer{}),
			testFactory(&testServer{
				errorOnPrepareProviderConfig: true,
			}),
		},
		{
			testFactory(&testServer{
				respondToPrepareProviderConfig: true,
			}),
			testFactory(&testServer{}),
			testFactory(&testServer{
				warnOnPrepareProviderConfig: true,
			}),
		},
		{
			testFactory(&testServer{
				errorOnPrepareProviderConfig: true,
			}),
			testFactory(&testServer{}),
			testFactory(&testServer{
				warnOnPrepareProviderConfig: true,
			}),
		},
		{
			testFactory(&testServer{
				errorOnPrepareProviderConfig: true,
			}),
			testFactory(&testServer{}),
			testFactory(&testServer{
				errorOnPrepareProviderConfig: true,
			}),
		},
		{
			testFactory(&testServer{
				warnOnPrepareProviderConfig: true,
			}),
			testFactory(&testServer{}),
			testFactory(&testServer{
				warnOnPrepareProviderConfig: true,
			}),
		},
	} {
		factory, err = NewSchemaServerFactory(context.Background(), fs...)
		if err != nil {
			t.Fatalf("error setting up muxer: %s", err)
		}
		_, err = factory.Server().PrepareProviderConfig(context.Background(), &tfprotov5.PrepareProviderConfigRequest{
			Config: &config,
		})
		if !strings.Contains(err.Error(), "got a PrepareProviderConfig response from multiple servers") {
			t.Errorf("expected error about multiple servers returning PrepareProviderConfigResponses, got %q", err)
		}
	}
}

func TestSchemaServerConfigureProvider_configuredEveryone(t *testing.T) {
	server1 := testFactory(&testServer{})
	server2 := testFactory(&testServer{})
	server3 := testFactory(&testServer{})
	server4 := testFactory(&testServer{})
	server5 := testFactory(&testServer{})
	factory, err := NewSchemaServerFactory(context.Background(), server1, server2, server3, server4, server5)
	if err != nil {
		t.Fatalf("error setting up muxer: %s", err)
	}
	_, err = factory.Server().ConfigureProvider(context.Background(), &tfprotov5.ConfigureProviderRequest{})
	if err != nil {
		t.Fatalf("error calling ConfigureProvider: %s", err)
	}
	for num, f := range []func() tfprotov5.ProviderServer{
		server1, server2, server3, server4, server5,
	} {
		if !f().(*testServer).configureCalled {
			t.Errorf("configure not called on server%d", num+1)
		}
	}
}

func TestSchemaServerStopProvider_stoppedEveryone(t *testing.T) {
	server1 := testFactory(&testServer{})
	server2 := testFactory(&testServer{
		stopError: "error in server2",
	})
	server3 := testFactory(&testServer{})
	server4 := testFactory(&testServer{
		stopError: "error in server4",
	})
	server5 := testFactory(&testServer{})
	factory, err := NewSchemaServerFactory(context.Background(), server1, server2, server3, server4, server5)
	if err != nil {
		t.Fatalf("error setting up muxer: %s", err)
	}
	_, err = factory.Server().StopProvider(context.Background(), &tfprotov5.StopProviderRequest{})
	if err != nil {
		t.Fatalf("error calling StopProvider: %s", err)
	}
	for num, f := range []func() tfprotov5.ProviderServer{
		server1, server2, server3, server4, server5,
	} {
		if !f().(*testServer).stopCalled {
			t.Errorf("stop not called on server%d", num+1)
		}
	}
}

func TestSchemaServer_resourceRouting(t *testing.T) {
	server1 := testFactory(&testServer{
		resourcesCalled: map[string]bool{},
		resourceSchemas: map[string]*tfprotov5.Schema{
			"test_foo": {},
		},
	})
	server2 := testFactory(&testServer{
		resourcesCalled: map[string]bool{},
		resourceSchemas: map[string]*tfprotov5.Schema{
			"test_bar": {},
		},
	})

	factory, err := NewSchemaServerFactory(context.Background(), server1, server2)
	if err != nil {
		t.Fatalf("unexpected error setting up factory: %s", err)
	}

	_, err = factory.Server().ValidateResourceTypeConfig(context.Background(), &tfprotov5.ValidateResourceTypeConfigRequest{
		TypeName: "test_foo",
	})
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	if !server1().(*testServer).resourcesCalled["test_foo"] {
		t.Errorf("expected test_foo to be called on server1, was not")
	}
	if server2().(*testServer).resourcesCalled["test_foo"] {
		t.Errorf("expected test_foo not to be called on server2, was")
	}

	server1().(*testServer).resourcesCalled = map[string]bool{}
	server2().(*testServer).resourcesCalled = map[string]bool{}

	_, err = factory.Server().ValidateResourceTypeConfig(context.Background(), &tfprotov5.ValidateResourceTypeConfigRequest{
		TypeName: "test_bar",
	})
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	if !server2().(*testServer).resourcesCalled["test_bar"] {
		t.Errorf("expected test_bar to be called on server2, was not")
	}
	if server1().(*testServer).resourcesCalled["test_bar"] {
		t.Errorf("expected test_bar not to be called on server1, was")
	}

	server1().(*testServer).resourcesCalled = map[string]bool{}
	server2().(*testServer).resourcesCalled = map[string]bool{}

	_, err = factory.Server().UpgradeResourceState(context.Background(), &tfprotov5.UpgradeResourceStateRequest{
		TypeName: "test_foo",
	})
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	if !server1().(*testServer).resourcesCalled["test_foo"] {
		t.Errorf("expected test_foo to be called on server1, was not")
	}
	if server2().(*testServer).resourcesCalled["test_foo"] {
		t.Errorf("expected test_foo not to be called on server2, was")
	}

	server1().(*testServer).resourcesCalled = map[string]bool{}
	server2().(*testServer).resourcesCalled = map[string]bool{}

	_, err = factory.Server().UpgradeResourceState(context.Background(), &tfprotov5.UpgradeResourceStateRequest{
		TypeName: "test_bar",
	})
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	if !server2().(*testServer).resourcesCalled["test_bar"] {
		t.Errorf("expected test_bar to be called on server2, was not")
	}
	if server1().(*testServer).resourcesCalled["test_bar"] {
		t.Errorf("expected test_bar not to be called on server1, was")
	}

	server1().(*testServer).resourcesCalled = map[string]bool{}
	server2().(*testServer).resourcesCalled = map[string]bool{}

	_, err = factory.Server().ReadResource(context.Background(), &tfprotov5.ReadResourceRequest{
		TypeName: "test_foo",
	})
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	if !server1().(*testServer).resourcesCalled["test_foo"] {
		t.Errorf("expected test_foo to be called on server1, was not")
	}
	if server2().(*testServer).resourcesCalled["test_foo"] {
		t.Errorf("expected test_foo not to be called on server2, was")
	}

	server1().(*testServer).resourcesCalled = map[string]bool{}
	server2().(*testServer).resourcesCalled = map[string]bool{}

	_, err = factory.Server().ReadResource(context.Background(), &tfprotov5.ReadResourceRequest{
		TypeName: "test_bar",
	})
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	if !server2().(*testServer).resourcesCalled["test_bar"] {
		t.Errorf("expected test_bar to be called on server2, was not")
	}
	if server1().(*testServer).resourcesCalled["test_bar"] {
		t.Errorf("expected test_bar not to be called on server1, was")
	}

	server1().(*testServer).resourcesCalled = map[string]bool{}
	server2().(*testServer).resourcesCalled = map[string]bool{}

	_, err = factory.Server().PlanResourceChange(context.Background(), &tfprotov5.PlanResourceChangeRequest{
		TypeName: "test_foo",
	})
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	if !server1().(*testServer).resourcesCalled["test_foo"] {
		t.Errorf("expected test_foo to be called on server1, was not")
	}
	if server2().(*testServer).resourcesCalled["test_foo"] {
		t.Errorf("expected test_foo not to be called on server2, was")
	}

	server1().(*testServer).resourcesCalled = map[string]bool{}
	server2().(*testServer).resourcesCalled = map[string]bool{}

	_, err = factory.Server().PlanResourceChange(context.Background(), &tfprotov5.PlanResourceChangeRequest{
		TypeName: "test_bar",
	})
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	if !server2().(*testServer).resourcesCalled["test_bar"] {
		t.Errorf("expected test_bar to be called on server2, was not")
	}
	if server1().(*testServer).resourcesCalled["test_bar"] {
		t.Errorf("expected test_bar not to be called on server1, was")
	}

	server1().(*testServer).resourcesCalled = map[string]bool{}
	server2().(*testServer).resourcesCalled = map[string]bool{}

	_, err = factory.Server().ApplyResourceChange(context.Background(), &tfprotov5.ApplyResourceChangeRequest{
		TypeName: "test_foo",
	})
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	if !server1().(*testServer).resourcesCalled["test_foo"] {
		t.Errorf("expected test_foo to be called on server1, was not")
	}
	if server2().(*testServer).resourcesCalled["test_foo"] {
		t.Errorf("expected test_foo not to be called on server2, was")
	}

	server1().(*testServer).resourcesCalled = map[string]bool{}
	server2().(*testServer).resourcesCalled = map[string]bool{}

	_, err = factory.Server().ApplyResourceChange(context.Background(), &tfprotov5.ApplyResourceChangeRequest{
		TypeName: "test_bar",
	})
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	if !server2().(*testServer).resourcesCalled["test_bar"] {
		t.Errorf("expected test_bar to be called on server2, was not")
	}
	if server1().(*testServer).resourcesCalled["test_bar"] {
		t.Errorf("expected test_bar not to be called on server1, was")
	}

	server1().(*testServer).resourcesCalled = map[string]bool{}
	server2().(*testServer).resourcesCalled = map[string]bool{}

	_, err = factory.Server().ImportResourceState(context.Background(), &tfprotov5.ImportResourceStateRequest{
		TypeName: "test_foo",
	})
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	if !server1().(*testServer).resourcesCalled["test_foo"] {
		t.Errorf("expected test_foo to be called on server1, was not")
	}
	if server2().(*testServer).resourcesCalled["test_foo"] {
		t.Errorf("expected test_foo not to be called on server2, was")
	}

	server1().(*testServer).resourcesCalled = map[string]bool{}
	server2().(*testServer).resourcesCalled = map[string]bool{}

	_, err = factory.Server().ImportResourceState(context.Background(), &tfprotov5.ImportResourceStateRequest{
		TypeName: "test_bar",
	})
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	if !server2().(*testServer).resourcesCalled["test_bar"] {
		t.Errorf("expected test_bar to be called on server2, was not")
	}
	if server1().(*testServer).resourcesCalled["test_bar"] {
		t.Errorf("expected test_bar not to be called on server1, was")
	}
}

func TestSchemaServer_dataSourceRouting(t *testing.T) {
	server1 := testFactory(&testServer{
		dataSourcesCalled: map[string]bool{},
		dataSourceSchemas: map[string]*tfprotov5.Schema{
			"test_foo": {},
		},
	})
	server2 := testFactory(&testServer{
		dataSourcesCalled: map[string]bool{},
		dataSourceSchemas: map[string]*tfprotov5.Schema{
			"test_bar": {},
		},
	})

	factory, err := NewSchemaServerFactory(context.Background(), server1, server2)
	if err != nil {
		t.Fatalf("unexpected error setting up factory: %s", err)
	}

	_, err = factory.Server().ValidateDataSourceConfig(context.Background(), &tfprotov5.ValidateDataSourceConfigRequest{
		TypeName: "test_foo",
	})
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	if !server1().(*testServer).dataSourcesCalled["test_foo"] {
		t.Errorf("expected test_foo to be called on server1, was not")
	}
	if server2().(*testServer).dataSourcesCalled["test_foo"] {
		t.Errorf("expected test_foo not to be called on server2, was")
	}

	server1().(*testServer).dataSourcesCalled = map[string]bool{}
	server2().(*testServer).dataSourcesCalled = map[string]bool{}

	_, err = factory.Server().ReadDataSource(context.Background(), &tfprotov5.ReadDataSourceRequest{
		TypeName: "test_foo",
	})
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	if !server1().(*testServer).dataSourcesCalled["test_foo"] {
		t.Errorf("expected test_foo to be called on server1, was not")
	}
	if server2().(*testServer).dataSourcesCalled["test_foo"] {
		t.Errorf("expected test_foo not to be called on server2, was")
	}

	server1().(*testServer).dataSourcesCalled = map[string]bool{}
	server2().(*testServer).dataSourcesCalled = map[string]bool{}

	_, err = factory.Server().ValidateDataSourceConfig(context.Background(), &tfprotov5.ValidateDataSourceConfigRequest{
		TypeName: "test_bar",
	})
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	if !server2().(*testServer).dataSourcesCalled["test_bar"] {
		t.Errorf("expected test_bar to be called on server2, was not")
	}
	if server1().(*testServer).dataSourcesCalled["test_bar"] {
		t.Errorf("expected test_bar not to be called on server1, was")
	}

	server1().(*testServer).dataSourcesCalled = map[string]bool{}
	server2().(*testServer).dataSourcesCalled = map[string]bool{}

	_, err = factory.Server().ReadDataSource(context.Background(), &tfprotov5.ReadDataSourceRequest{
		TypeName: "test_bar",
	})
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	if !server2().(*testServer).dataSourcesCalled["test_bar"] {
		t.Errorf("expected test_bar to be called on server2, was not")
	}
	if server1().(*testServer).dataSourcesCalled["test_bar"] {
		t.Errorf("expected test_bar not to be called on server1, was")
	}
}
