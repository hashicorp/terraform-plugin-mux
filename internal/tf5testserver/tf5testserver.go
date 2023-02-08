package tf5testserver

import (
	"context"

	"github.com/hashicorp/terraform-plugin-go/tfprotov5"
)

var _ tfprotov5.ProviderServer = &TestServer{}

type TestServer struct {
	DataSourceSchemas  map[string]*tfprotov5.Schema
	ProviderMetaSchema *tfprotov5.Schema
	ProviderSchema     *tfprotov5.Schema
	ResourceSchemas    map[string]*tfprotov5.Schema
	ServerCapabilities *tfprotov5.ServerCapabilities

	ApplyResourceChangeCalled map[string]bool

	ConfigureProviderCalled bool

	GetProviderSchemaCalled bool

	ImportResourceStateCalled map[string]bool

	PlanResourceChangeCalled map[string]bool

	PrepareProviderConfigCalled   bool
	PrepareProviderConfigResponse *tfprotov5.PrepareProviderConfigResponse

	ReadDataSourceCalled map[string]bool

	ReadResourceCalled map[string]bool

	StopProviderCalled bool
	StopProviderError  string

	UpgradeResourceStateCalled map[string]bool

	ValidateDataSourceConfigCalled map[string]bool

	ValidateResourceTypeConfigCalled map[string]bool
}

func (s *TestServer) ProviderServer() tfprotov5.ProviderServer {
	return s
}

func (s *TestServer) ApplyResourceChange(_ context.Context, req *tfprotov5.ApplyResourceChangeRequest) (*tfprotov5.ApplyResourceChangeResponse, error) {
	if s.ApplyResourceChangeCalled == nil {
		s.ApplyResourceChangeCalled = make(map[string]bool)
	}

	s.ApplyResourceChangeCalled[req.TypeName] = true
	return nil, nil
}

func (s *TestServer) ConfigureProvider(_ context.Context, _ *tfprotov5.ConfigureProviderRequest) (*tfprotov5.ConfigureProviderResponse, error) {
	s.ConfigureProviderCalled = true
	return &tfprotov5.ConfigureProviderResponse{}, nil
}

func (s *TestServer) GetProviderSchema(_ context.Context, _ *tfprotov5.GetProviderSchemaRequest) (*tfprotov5.GetProviderSchemaResponse, error) {
	s.GetProviderSchemaCalled = true

	if s.DataSourceSchemas == nil {
		s.DataSourceSchemas = make(map[string]*tfprotov5.Schema)
	}

	if s.ResourceSchemas == nil {
		s.ResourceSchemas = make(map[string]*tfprotov5.Schema)
	}

	return &tfprotov5.GetProviderSchemaResponse{
		Provider:           s.ProviderSchema,
		ProviderMeta:       s.ProviderMetaSchema,
		ResourceSchemas:    s.ResourceSchemas,
		DataSourceSchemas:  s.DataSourceSchemas,
		ServerCapabilities: s.ServerCapabilities,
	}, nil
}

func (s *TestServer) ImportResourceState(_ context.Context, req *tfprotov5.ImportResourceStateRequest) (*tfprotov5.ImportResourceStateResponse, error) {
	if s.ImportResourceStateCalled == nil {
		s.ImportResourceStateCalled = make(map[string]bool)
	}

	s.ImportResourceStateCalled[req.TypeName] = true
	return nil, nil
}

func (s *TestServer) PlanResourceChange(_ context.Context, req *tfprotov5.PlanResourceChangeRequest) (*tfprotov5.PlanResourceChangeResponse, error) {
	if s.PlanResourceChangeCalled == nil {
		s.PlanResourceChangeCalled = make(map[string]bool)
	}

	s.PlanResourceChangeCalled[req.TypeName] = true
	return nil, nil
}

func (s *TestServer) ReadDataSource(_ context.Context, req *tfprotov5.ReadDataSourceRequest) (*tfprotov5.ReadDataSourceResponse, error) {
	if s.ReadDataSourceCalled == nil {
		s.ReadDataSourceCalled = make(map[string]bool)
	}

	s.ReadDataSourceCalled[req.TypeName] = true
	return nil, nil
}

func (s *TestServer) ReadResource(_ context.Context, req *tfprotov5.ReadResourceRequest) (*tfprotov5.ReadResourceResponse, error) {
	if s.ReadResourceCalled == nil {
		s.ReadResourceCalled = make(map[string]bool)
	}

	s.ReadResourceCalled[req.TypeName] = true
	return nil, nil
}

func (s *TestServer) StopProvider(_ context.Context, _ *tfprotov5.StopProviderRequest) (*tfprotov5.StopProviderResponse, error) {
	s.StopProviderCalled = true

	if s.StopProviderError != "" {
		return &tfprotov5.StopProviderResponse{
			Error: s.StopProviderError,
		}, nil
	}

	return &tfprotov5.StopProviderResponse{}, nil
}

func (s *TestServer) UpgradeResourceState(_ context.Context, req *tfprotov5.UpgradeResourceStateRequest) (*tfprotov5.UpgradeResourceStateResponse, error) {
	if s.UpgradeResourceStateCalled == nil {
		s.UpgradeResourceStateCalled = make(map[string]bool)
	}

	s.UpgradeResourceStateCalled[req.TypeName] = true
	return nil, nil
}

func (s *TestServer) ValidateDataSourceConfig(_ context.Context, req *tfprotov5.ValidateDataSourceConfigRequest) (*tfprotov5.ValidateDataSourceConfigResponse, error) {
	if s.ValidateDataSourceConfigCalled == nil {
		s.ValidateDataSourceConfigCalled = make(map[string]bool)
	}

	s.ValidateDataSourceConfigCalled[req.TypeName] = true
	return nil, nil
}

func (s *TestServer) ValidateResourceTypeConfig(_ context.Context, req *tfprotov5.ValidateResourceTypeConfigRequest) (*tfprotov5.ValidateResourceTypeConfigResponse, error) {
	if s.ValidateResourceTypeConfigCalled == nil {
		s.ValidateResourceTypeConfigCalled = make(map[string]bool)
	}

	s.ValidateResourceTypeConfigCalled[req.TypeName] = true
	return nil, nil
}

func (s *TestServer) PrepareProviderConfig(_ context.Context, req *tfprotov5.PrepareProviderConfigRequest) (*tfprotov5.PrepareProviderConfigResponse, error) {
	s.PrepareProviderConfigCalled = true
	return s.PrepareProviderConfigResponse, nil
}
