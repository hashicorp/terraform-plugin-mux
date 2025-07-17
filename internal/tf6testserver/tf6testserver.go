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

	CloseEphemeralResourceCalled map[string]bool

	ConfigureProviderCalled   bool
	ConfigureProviderResponse *tfprotov6.ConfigureProviderResponse

	GetFunctionsCalled   bool
	GetFunctionsResponse *tfprotov6.GetFunctionsResponse

	GetMetadataCalled   bool
	GetMetadataResponse *tfprotov6.GetMetadataResponse

	GetProviderSchemaCalled   bool
	GetProviderSchemaResponse *tfprotov6.GetProviderSchemaResponse

	GetResourceIdentitySchemasCalled   bool
	GetResourceIdentitySchemasResponse *tfprotov6.GetResourceIdentitySchemasResponse

	ImportResourceStateCalled map[string]bool

	MoveResourceStateCalled map[string]bool

	OpenEphemeralResourceCalled map[string]bool

	PlanResourceChangeCalled map[string]bool

	ReadDataSourceCalled map[string]bool

	ReadResourceCalled map[string]bool

	RenewEphemeralResourceCalled map[string]bool

	StopProviderCalled   bool
	StopProviderResponse *tfprotov6.StopProviderResponse

	UpgradeResourceIdentityCalled map[string]bool

	UpgradeResourceStateCalled map[string]bool

	ValidateDataResourceConfigCalled map[string]bool

	ValidateEphemeralResourceConfigCalled map[string]bool

	ValidateProviderConfigCalled   bool
	ValidateProviderConfigResponse *tfprotov6.ValidateProviderConfigResponse

	ValidateResourceConfigCalled map[string]bool

	ValidateListResourceConfigCalled map[string]bool

	ListResourceCalled map[string]bool

	ValidateActionConfigCalled map[string]bool

	PlanActionCalled map[string]bool

	InvokeActionCalled map[string]bool
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

func (s *TestServer) CloseEphemeralResource(ctx context.Context, req *tfprotov6.CloseEphemeralResourceRequest) (*tfprotov6.CloseEphemeralResourceResponse, error) {
	if s.CloseEphemeralResourceCalled == nil {
		s.CloseEphemeralResourceCalled = make(map[string]bool)
	}

	s.CloseEphemeralResourceCalled[req.TypeName] = true
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

func (s *TestServer) GetResourceIdentitySchemas(_ context.Context, _ *tfprotov6.GetResourceIdentitySchemasRequest) (*tfprotov6.GetResourceIdentitySchemasResponse, error) {
	s.GetResourceIdentitySchemasCalled = true

	if s.GetResourceIdentitySchemasResponse != nil {
		return s.GetResourceIdentitySchemasResponse, nil
	}

	return &tfprotov6.GetResourceIdentitySchemasResponse{}, nil
}

func (s *TestServer) ImportResourceState(_ context.Context, req *tfprotov6.ImportResourceStateRequest) (*tfprotov6.ImportResourceStateResponse, error) {
	if s.ImportResourceStateCalled == nil {
		s.ImportResourceStateCalled = make(map[string]bool)
	}

	s.ImportResourceStateCalled[req.TypeName] = true
	return nil, nil
}

func (s *TestServer) MoveResourceState(_ context.Context, req *tfprotov6.MoveResourceStateRequest) (*tfprotov6.MoveResourceStateResponse, error) {
	if s.MoveResourceStateCalled == nil {
		s.MoveResourceStateCalled = make(map[string]bool)
	}

	s.MoveResourceStateCalled[req.TargetTypeName] = true
	return nil, nil
}

func (s *TestServer) OpenEphemeralResource(_ context.Context, req *tfprotov6.OpenEphemeralResourceRequest) (*tfprotov6.OpenEphemeralResourceResponse, error) {
	if s.OpenEphemeralResourceCalled == nil {
		s.OpenEphemeralResourceCalled = make(map[string]bool)
	}

	s.OpenEphemeralResourceCalled[req.TypeName] = true
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

func (s *TestServer) RenewEphemeralResource(_ context.Context, req *tfprotov6.RenewEphemeralResourceRequest) (*tfprotov6.RenewEphemeralResourceResponse, error) {
	if s.RenewEphemeralResourceCalled == nil {
		s.RenewEphemeralResourceCalled = make(map[string]bool)
	}

	s.RenewEphemeralResourceCalled[req.TypeName] = true
	return nil, nil
}

func (s *TestServer) StopProvider(_ context.Context, _ *tfprotov6.StopProviderRequest) (*tfprotov6.StopProviderResponse, error) {
	s.StopProviderCalled = true

	if s.StopProviderResponse != nil {
		return s.StopProviderResponse, nil
	}

	return &tfprotov6.StopProviderResponse{}, nil
}

func (s *TestServer) UpgradeResourceIdentity(_ context.Context, req *tfprotov6.UpgradeResourceIdentityRequest) (*tfprotov6.UpgradeResourceIdentityResponse, error) {
	if s.UpgradeResourceIdentityCalled == nil {
		s.UpgradeResourceIdentityCalled = make(map[string]bool)
	}

	s.UpgradeResourceIdentityCalled[req.TypeName] = true
	return nil, nil
}

func (s *TestServer) UpgradeResourceState(_ context.Context, req *tfprotov6.UpgradeResourceStateRequest) (*tfprotov6.UpgradeResourceStateResponse, error) {
	if s.UpgradeResourceStateCalled == nil {
		s.UpgradeResourceStateCalled = make(map[string]bool)
	}

	s.UpgradeResourceStateCalled[req.TypeName] = true
	return nil, nil
}

func (s *TestServer) ValidateEphemeralResourceConfig(_ context.Context, req *tfprotov6.ValidateEphemeralResourceConfigRequest) (*tfprotov6.ValidateEphemeralResourceConfigResponse, error) {
	if s.ValidateEphemeralResourceConfigCalled == nil {
		s.ValidateEphemeralResourceConfigCalled = make(map[string]bool)
	}

	s.ValidateEphemeralResourceConfigCalled[req.TypeName] = true
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

func (s *TestServer) ValidateListResourceConfig(_ context.Context, req *tfprotov6.ValidateListResourceConfigRequest) (*tfprotov6.ValidateListResourceConfigResponse, error) {
	if s.ValidateListResourceConfigCalled == nil {
		s.ValidateListResourceConfigCalled = make(map[string]bool)
	}

	s.ValidateListResourceConfigCalled[req.TypeName] = true
	return nil, nil
}

func (s *TestServer) ListResource(_ context.Context, req *tfprotov6.ListResourceRequest) (*tfprotov6.ListResourceServerStream, error) {
	if s.ListResourceCalled == nil {
		s.ListResourceCalled = make(map[string]bool)
	}

	s.ListResourceCalled[req.TypeName] = true
	return nil, nil
}

func (s *TestServer) ValidateActionConfig(_ context.Context, req *tfprotov6.ValidateActionConfigRequest) (*tfprotov6.ValidateActionConfigResponse, error) {
	if s.ValidateActionConfigCalled == nil {
		s.ValidateActionConfigCalled = make(map[string]bool)
	}

	s.ValidateActionConfigCalled[req.ActionType] = true
	return nil, nil
}

func (s *TestServer) PlanAction(ctx context.Context, req *tfprotov6.PlanActionRequest) (*tfprotov6.PlanActionResponse, error) {
	if s.PlanActionCalled == nil {
		s.PlanActionCalled = make(map[string]bool)
	}

	s.PlanActionCalled[req.ActionType] = true
	return nil, nil
}

func (s *TestServer) InvokeAction(ctx context.Context, req *tfprotov6.InvokeActionRequest) (*tfprotov6.InvokeActionServerStream, error) {
	if s.InvokeActionCalled == nil {
		s.InvokeActionCalled = make(map[string]bool)
	}

	s.InvokeActionCalled[req.ActionType] = true
	return nil, nil
}
