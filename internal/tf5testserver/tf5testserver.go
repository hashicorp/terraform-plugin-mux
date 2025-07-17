// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package tf5testserver

import (
	"context"

	"github.com/hashicorp/terraform-plugin-go/tfprotov5"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var _ tfprotov5.ProviderServer = &TestServer{}

type TestServer struct {
	ApplyResourceChangeCalled map[string]bool

	CallFunctionCalled map[string]bool

	CloseEphemeralResourceCalled map[string]bool

	ConfigureProviderCalled   bool
	ConfigureProviderResponse *tfprotov5.ConfigureProviderResponse

	GetFunctionsCalled   bool
	GetFunctionsResponse *tfprotov5.GetFunctionsResponse

	GetMetadataCalled   bool
	GetMetadataResponse *tfprotov5.GetMetadataResponse

	GetProviderSchemaCalled   bool
	GetProviderSchemaResponse *tfprotov5.GetProviderSchemaResponse

	GetResourceIdentitySchemasCalled   bool
	GetResourceIdentitySchemasResponse *tfprotov5.GetResourceIdentitySchemasResponse

	ImportResourceStateCalled map[string]bool

	MoveResourceStateCalled map[string]bool

	OpenEphemeralResourceCalled map[string]bool

	PlanResourceChangeCalled map[string]bool

	PrepareProviderConfigCalled   bool
	PrepareProviderConfigResponse *tfprotov5.PrepareProviderConfigResponse

	ReadDataSourceCalled map[string]bool

	ReadResourceCalled map[string]bool

	RenewEphemeralResourceCalled map[string]bool

	StopProviderCalled   bool
	StopProviderResponse *tfprotov5.StopProviderResponse

	UpgradeResourceIdentityCalled map[string]bool

	UpgradeResourceStateCalled map[string]bool

	ValidateEphemeralResourceConfigCalled map[string]bool

	ValidateDataSourceConfigCalled map[string]bool

	ValidateResourceTypeConfigCalled map[string]bool

	ValidateListResourceConfigCalled map[string]bool

	ListResourceCalled map[string]bool

	ValidateActionConfigCalled map[string]bool

	PlanActionCalled map[string]bool

	InvokeActionCalled map[string]bool
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

func (s *TestServer) CallFunction(_ context.Context, req *tfprotov5.CallFunctionRequest) (*tfprotov5.CallFunctionResponse, error) {
	if s.CallFunctionCalled == nil {
		s.CallFunctionCalled = make(map[string]bool)
	}

	s.CallFunctionCalled[req.Name] = true
	return nil, nil
}

func (s *TestServer) CloseEphemeralResource(ctx context.Context, req *tfprotov5.CloseEphemeralResourceRequest) (*tfprotov5.CloseEphemeralResourceResponse, error) {
	if s.CloseEphemeralResourceCalled == nil {
		s.CloseEphemeralResourceCalled = make(map[string]bool)
	}

	s.CloseEphemeralResourceCalled[req.TypeName] = true
	return nil, nil
}

func (s *TestServer) ConfigureProvider(_ context.Context, _ *tfprotov5.ConfigureProviderRequest) (*tfprotov5.ConfigureProviderResponse, error) {
	s.ConfigureProviderCalled = true

	if s.ConfigureProviderResponse != nil {
		return s.ConfigureProviderResponse, nil
	}

	return &tfprotov5.ConfigureProviderResponse{}, nil
}

func (s *TestServer) GetFunctions(_ context.Context, _ *tfprotov5.GetFunctionsRequest) (*tfprotov5.GetFunctionsResponse, error) {
	s.GetFunctionsCalled = true

	if s.GetFunctionsResponse != nil {
		return s.GetFunctionsResponse, nil
	}

	return &tfprotov5.GetFunctionsResponse{}, nil
}

func (s *TestServer) GetMetadata(_ context.Context, _ *tfprotov5.GetMetadataRequest) (*tfprotov5.GetMetadataResponse, error) {
	s.GetMetadataCalled = true

	if s.GetMetadataResponse != nil {
		return s.GetMetadataResponse, nil
	}

	if s.GetProviderSchemaResponse != nil {
		return nil, status.Error(codes.Unimplemented, "only GetProviderSchemaResponse set, simulating GetMetadata as unimplemented")
	}

	return &tfprotov5.GetMetadataResponse{}, nil
}

func (s *TestServer) GetProviderSchema(_ context.Context, _ *tfprotov5.GetProviderSchemaRequest) (*tfprotov5.GetProviderSchemaResponse, error) {
	s.GetProviderSchemaCalled = true

	if s.GetProviderSchemaResponse != nil {
		return s.GetProviderSchemaResponse, nil
	}

	return &tfprotov5.GetProviderSchemaResponse{}, nil
}

func (s *TestServer) GetResourceIdentitySchemas(_ context.Context, _ *tfprotov5.GetResourceIdentitySchemasRequest) (*tfprotov5.GetResourceIdentitySchemasResponse, error) {
	s.GetResourceIdentitySchemasCalled = true

	if s.GetResourceIdentitySchemasResponse != nil {
		return s.GetResourceIdentitySchemasResponse, nil
	}

	return &tfprotov5.GetResourceIdentitySchemasResponse{}, nil
}

func (s *TestServer) ImportResourceState(_ context.Context, req *tfprotov5.ImportResourceStateRequest) (*tfprotov5.ImportResourceStateResponse, error) {
	if s.ImportResourceStateCalled == nil {
		s.ImportResourceStateCalled = make(map[string]bool)
	}

	s.ImportResourceStateCalled[req.TypeName] = true
	return nil, nil
}

func (s *TestServer) MoveResourceState(_ context.Context, req *tfprotov5.MoveResourceStateRequest) (*tfprotov5.MoveResourceStateResponse, error) {
	if s.MoveResourceStateCalled == nil {
		s.MoveResourceStateCalled = make(map[string]bool)
	}

	s.MoveResourceStateCalled[req.TargetTypeName] = true
	return nil, nil
}

func (s *TestServer) OpenEphemeralResource(_ context.Context, req *tfprotov5.OpenEphemeralResourceRequest) (*tfprotov5.OpenEphemeralResourceResponse, error) {
	if s.OpenEphemeralResourceCalled == nil {
		s.OpenEphemeralResourceCalled = make(map[string]bool)
	}

	s.OpenEphemeralResourceCalled[req.TypeName] = true
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

func (s *TestServer) RenewEphemeralResource(_ context.Context, req *tfprotov5.RenewEphemeralResourceRequest) (*tfprotov5.RenewEphemeralResourceResponse, error) {
	if s.RenewEphemeralResourceCalled == nil {
		s.RenewEphemeralResourceCalled = make(map[string]bool)
	}

	s.RenewEphemeralResourceCalled[req.TypeName] = true
	return nil, nil
}

func (s *TestServer) StopProvider(_ context.Context, _ *tfprotov5.StopProviderRequest) (*tfprotov5.StopProviderResponse, error) {
	s.StopProviderCalled = true

	if s.StopProviderResponse != nil {
		return s.StopProviderResponse, nil
	}

	return &tfprotov5.StopProviderResponse{}, nil
}

func (s *TestServer) UpgradeResourceIdentity(_ context.Context, req *tfprotov5.UpgradeResourceIdentityRequest) (*tfprotov5.UpgradeResourceIdentityResponse, error) {
	if s.UpgradeResourceIdentityCalled == nil {
		s.UpgradeResourceIdentityCalled = make(map[string]bool)
	}

	s.UpgradeResourceIdentityCalled[req.TypeName] = true
	return nil, nil
}

func (s *TestServer) UpgradeResourceState(_ context.Context, req *tfprotov5.UpgradeResourceStateRequest) (*tfprotov5.UpgradeResourceStateResponse, error) {
	if s.UpgradeResourceStateCalled == nil {
		s.UpgradeResourceStateCalled = make(map[string]bool)
	}

	s.UpgradeResourceStateCalled[req.TypeName] = true
	return nil, nil
}

func (s *TestServer) ValidateEphemeralResourceConfig(_ context.Context, req *tfprotov5.ValidateEphemeralResourceConfigRequest) (*tfprotov5.ValidateEphemeralResourceConfigResponse, error) {
	if s.ValidateEphemeralResourceConfigCalled == nil {
		s.ValidateEphemeralResourceConfigCalled = make(map[string]bool)
	}

	s.ValidateEphemeralResourceConfigCalled[req.TypeName] = true
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

func (s *TestServer) ValidateListResourceConfig(_ context.Context, req *tfprotov5.ValidateListResourceConfigRequest) (*tfprotov5.ValidateListResourceConfigResponse, error) {
	if s.ValidateListResourceConfigCalled == nil {
		s.ValidateListResourceConfigCalled = make(map[string]bool)
	}

	s.ValidateListResourceConfigCalled[req.TypeName] = true
	return nil, nil
}

func (s *TestServer) PrepareProviderConfig(_ context.Context, req *tfprotov5.PrepareProviderConfigRequest) (*tfprotov5.PrepareProviderConfigResponse, error) {
	s.PrepareProviderConfigCalled = true
	return s.PrepareProviderConfigResponse, nil
}

func (s *TestServer) ListResource(_ context.Context, req *tfprotov5.ListResourceRequest) (*tfprotov5.ListResourceServerStream, error) {
	if s.ListResourceCalled == nil {
		s.ListResourceCalled = make(map[string]bool)
	}

	s.ListResourceCalled[req.TypeName] = true
	return nil, nil
}

func (s *TestServer) ValidateActionConfig(_ context.Context, req *tfprotov5.ValidateActionConfigRequest) (*tfprotov5.ValidateActionConfigResponse, error) {
	if s.ValidateActionConfigCalled == nil {
		s.ValidateActionConfigCalled = make(map[string]bool)
	}

	s.ValidateActionConfigCalled[req.ActionType] = true
	return nil, nil
}

func (s *TestServer) PlanAction(ctx context.Context, req *tfprotov5.PlanActionRequest) (*tfprotov5.PlanActionResponse, error) {
	if s.PlanActionCalled == nil {
		s.PlanActionCalled = make(map[string]bool)
	}

	s.PlanActionCalled[req.ActionType] = true
	return nil, nil
}

func (s *TestServer) InvokeAction(ctx context.Context, req *tfprotov5.InvokeActionRequest) (*tfprotov5.InvokeActionServerStream, error) {
	if s.InvokeActionCalled == nil {
		s.InvokeActionCalled = make(map[string]bool)
	}

	s.InvokeActionCalled[req.ActionType] = true
	return nil, nil
}
