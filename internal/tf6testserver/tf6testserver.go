package tf6testserver

import (
	"context"

	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
)

var _ tfprotov6.ProviderServer = &TestServer{}

type TestServer struct {
	DataSourceSchemas  map[string]*tfprotov6.Schema
	ProviderMetaSchema *tfprotov6.Schema
	ProviderSchema     *tfprotov6.Schema
	ResourceSchemas    map[string]*tfprotov6.Schema
	ServerCapabilities *tfprotov6.ServerCapabilities

	ApplyResourceChangeCalled map[string]bool

	ConfigureProviderCalled bool

	GetProviderSchemaCalled bool

	ImportResourceStateCalled map[string]bool

	PlanResourceChangeCalled map[string]bool

	ReadDataSourceCalled map[string]bool

	ReadResourceCalled map[string]bool

	StopProviderCalled bool
	StopProviderError  string

	UpgradeResourceStateCalled map[string]bool

	ValidateDataResourceConfigCalled map[string]bool

	ValidateProviderConfigCalled   bool
	ValidateProviderConfigResponse *tfprotov6.ValidateProviderConfigResponse

	ValidateResourceConfigCalled map[string]bool
}

func (s *TestServer) ProviderServer() tfprotov6.ProviderServer {
	return s
}

func (s *TestServer) ApplyResourceChange(_ context.Context, req *tfprotov6.ApplyResourceChangeRequest) (*tfprotov6.ApplyResourceChangeResponse, error) {
	if s.ApplyResourceChangeCalled == nil {
		s.ApplyResourceChangeCalled = make(map[string]bool)
	}

	s.ApplyResourceChangeCalled[req.TypeName] = true
	return nil, nil
}

func (s *TestServer) ConfigureProvider(_ context.Context, _ *tfprotov6.ConfigureProviderRequest) (*tfprotov6.ConfigureProviderResponse, error) {
	s.ConfigureProviderCalled = true
	return &tfprotov6.ConfigureProviderResponse{}, nil
}

func (s *TestServer) GetProviderSchema(_ context.Context, _ *tfprotov6.GetProviderSchemaRequest) (*tfprotov6.GetProviderSchemaResponse, error) {
	s.GetProviderSchemaCalled = true

	if s.DataSourceSchemas == nil {
		s.DataSourceSchemas = make(map[string]*tfprotov6.Schema)
	}

	if s.ResourceSchemas == nil {
		s.ResourceSchemas = make(map[string]*tfprotov6.Schema)
	}

	return &tfprotov6.GetProviderSchemaResponse{
		Provider:           s.ProviderSchema,
		ProviderMeta:       s.ProviderMetaSchema,
		ResourceSchemas:    s.ResourceSchemas,
		DataSourceSchemas:  s.DataSourceSchemas,
		ServerCapabilities: s.ServerCapabilities,
	}, nil
}

func (s *TestServer) ImportResourceState(_ context.Context, req *tfprotov6.ImportResourceStateRequest) (*tfprotov6.ImportResourceStateResponse, error) {
	if s.ImportResourceStateCalled == nil {
		s.ImportResourceStateCalled = make(map[string]bool)
	}

	s.ImportResourceStateCalled[req.TypeName] = true
	return nil, nil
}

func (s *TestServer) PlanResourceChange(_ context.Context, req *tfprotov6.PlanResourceChangeRequest) (*tfprotov6.PlanResourceChangeResponse, error) {
	if s.PlanResourceChangeCalled == nil {
		s.PlanResourceChangeCalled = make(map[string]bool)
	}

	s.PlanResourceChangeCalled[req.TypeName] = true
	return nil, nil
}

func (s *TestServer) ReadDataSource(_ context.Context, req *tfprotov6.ReadDataSourceRequest) (*tfprotov6.ReadDataSourceResponse, error) {
	if s.ReadDataSourceCalled == nil {
		s.ReadDataSourceCalled = make(map[string]bool)
	}

	s.ReadDataSourceCalled[req.TypeName] = true
	return nil, nil
}

func (s *TestServer) ReadResource(_ context.Context, req *tfprotov6.ReadResourceRequest) (*tfprotov6.ReadResourceResponse, error) {
	if s.ReadResourceCalled == nil {
		s.ReadResourceCalled = make(map[string]bool)
	}

	s.ReadResourceCalled[req.TypeName] = true
	return nil, nil
}

func (s *TestServer) StopProvider(_ context.Context, _ *tfprotov6.StopProviderRequest) (*tfprotov6.StopProviderResponse, error) {
	s.StopProviderCalled = true

	if s.StopProviderError != "" {
		return &tfprotov6.StopProviderResponse{
			Error: s.StopProviderError,
		}, nil
	}

	return &tfprotov6.StopProviderResponse{}, nil
}

func (s *TestServer) UpgradeResourceState(_ context.Context, req *tfprotov6.UpgradeResourceStateRequest) (*tfprotov6.UpgradeResourceStateResponse, error) {
	if s.UpgradeResourceStateCalled == nil {
		s.UpgradeResourceStateCalled = make(map[string]bool)
	}

	s.UpgradeResourceStateCalled[req.TypeName] = true
	return nil, nil
}

func (s *TestServer) ValidateDataResourceConfig(_ context.Context, req *tfprotov6.ValidateDataResourceConfigRequest) (*tfprotov6.ValidateDataResourceConfigResponse, error) {
	if s.ValidateDataResourceConfigCalled == nil {
		s.ValidateDataResourceConfigCalled = make(map[string]bool)
	}

	s.ValidateDataResourceConfigCalled[req.TypeName] = true
	return nil, nil
}

func (s *TestServer) ValidateResourceConfig(_ context.Context, req *tfprotov6.ValidateResourceConfigRequest) (*tfprotov6.ValidateResourceConfigResponse, error) {
	if s.ValidateResourceConfigCalled == nil {
		s.ValidateResourceConfigCalled = make(map[string]bool)
	}

	s.ValidateResourceConfigCalled[req.TypeName] = true
	return nil, nil
}

func (s *TestServer) ValidateProviderConfig(_ context.Context, req *tfprotov6.ValidateProviderConfigRequest) (*tfprotov6.ValidateProviderConfigResponse, error) {
	s.ValidateProviderConfigCalled = true
	return s.ValidateProviderConfigResponse, nil
}
