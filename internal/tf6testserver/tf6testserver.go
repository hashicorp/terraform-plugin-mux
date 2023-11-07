// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package tf6testserver

import (
	"context"

	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var _ tfprotov6.ProviderServer = &TestServer{}

type TestServer struct {
	ApplyResourceChangeCalled map[string]bool

	CallFunctionCalled map[string]bool

	ConfigureProviderCalled   bool
	ConfigureProviderResponse *tfprotov6.ConfigureProviderResponse

	GetFunctionsCalled   bool
	GetFunctionsResponse *tfprotov6.GetFunctionsResponse

	GetMetadataCalled   bool
	GetMetadataResponse *tfprotov6.GetMetadataResponse

	GetProviderSchemaCalled   bool
	GetProviderSchemaResponse *tfprotov6.GetProviderSchemaResponse

	ImportResourceStateCalled map[string]bool

	PlanResourceChangeCalled map[string]bool

	ReadDataSourceCalled map[string]bool

	ReadResourceCalled map[string]bool

	StopProviderCalled   bool
	StopProviderResponse *tfprotov6.StopProviderResponse

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

func (s *TestServer) CallFunction(_ context.Context, req *tfprotov6.CallFunctionRequest) (*tfprotov6.CallFunctionResponse, error) {
	if s.CallFunctionCalled == nil {
		s.CallFunctionCalled = make(map[string]bool)
	}

	s.CallFunctionCalled[req.Name] = true
	return nil, nil
}

func (s *TestServer) ConfigureProvider(_ context.Context, _ *tfprotov6.ConfigureProviderRequest) (*tfprotov6.ConfigureProviderResponse, error) {
	s.ConfigureProviderCalled = true

	if s.ConfigureProviderResponse != nil {
		return s.ConfigureProviderResponse, nil
	}

	return &tfprotov6.ConfigureProviderResponse{}, nil
}

func (s *TestServer) GetFunctions(_ context.Context, _ *tfprotov6.GetFunctionsRequest) (*tfprotov6.GetFunctionsResponse, error) {
	s.GetFunctionsCalled = true

	if s.GetFunctionsResponse != nil {
		return s.GetFunctionsResponse, nil
	}

	return &tfprotov6.GetFunctionsResponse{}, nil
}

func (s *TestServer) GetMetadata(_ context.Context, _ *tfprotov6.GetMetadataRequest) (*tfprotov6.GetMetadataResponse, error) {
	s.GetMetadataCalled = true

	if s.GetMetadataResponse != nil {
		return s.GetMetadataResponse, nil
	}

	if s.GetProviderSchemaResponse != nil {
		return nil, status.Error(codes.Unimplemented, "only GetProviderSchemaResponse set, simulating GetMetadata as unimplemented")
	}

	return &tfprotov6.GetMetadataResponse{}, nil
}

func (s *TestServer) GetProviderSchema(_ context.Context, _ *tfprotov6.GetProviderSchemaRequest) (*tfprotov6.GetProviderSchemaResponse, error) {
	s.GetProviderSchemaCalled = true

	if s.GetProviderSchemaResponse != nil {
		return s.GetProviderSchemaResponse, nil
	}

	return &tfprotov6.GetProviderSchemaResponse{}, nil
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

	if s.StopProviderResponse != nil {
		return s.StopProviderResponse, nil
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
