// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package tf6to5server

import (
	"context"

	"github.com/hashicorp/terraform-plugin-go/tfprotov5"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-mux/internal/tfprotov5tov6"
	"github.com/hashicorp/terraform-plugin-mux/internal/tfprotov6tov5"
)

// DowngradeServer wraps a protocol version 6 ProviderServer in a protocol
// version 5 server. Protocol version 5 is not backwards compatible with
// protocol version 6, so additional validation is performed:
//
//   - GetProviderSchema is called to ensure SchemaAttribute.NestedType
//     (nested attributes) are not implemented.
//
// Protocol version 5 servers require Terraform CLI 0.12 or later.
func DowngradeServer(ctx context.Context, v6server func() tfprotov6.ProviderServer) (tfprotov5.ProviderServer, error) {
	v6Resp, err := v6server().GetProviderSchema(ctx, &tfprotov6.GetProviderSchemaRequest{})

	if err != nil {
		return nil, err
	}

	_, err = tfprotov6tov5.GetProviderSchemaResponse(v6Resp)

	if err != nil {
		return nil, err
	}

	return v6tov5Server{
		v6Server: v6server(),
	}, nil
}

var _ tfprotov5.ProviderServer = v6tov5Server{}

type v6tov5Server struct {
	v6Server tfprotov6.ProviderServer
}

func (s v6tov5Server) ApplyResourceChange(ctx context.Context, req *tfprotov5.ApplyResourceChangeRequest) (*tfprotov5.ApplyResourceChangeResponse, error) {
	v6Req := tfprotov5tov6.ApplyResourceChangeRequest(req)
	v6Resp, err := s.v6Server.ApplyResourceChange(ctx, v6Req)

	if err != nil {
		return nil, err
	}

	return tfprotov6tov5.ApplyResourceChangeResponse(v6Resp), nil
}

func (s v6tov5Server) ConfigureProvider(ctx context.Context, req *tfprotov5.ConfigureProviderRequest) (*tfprotov5.ConfigureProviderResponse, error) {
	v6Req := tfprotov5tov6.ConfigureProviderRequest(req)
	v6Resp, err := s.v6Server.ConfigureProvider(ctx, v6Req)

	if err != nil {
		return nil, err
	}

	return tfprotov6tov5.ConfigureProviderResponse(v6Resp), nil
}

func (s v6tov5Server) GetProviderSchema(ctx context.Context, req *tfprotov5.GetProviderSchemaRequest) (*tfprotov5.GetProviderSchemaResponse, error) {
	v6Req := tfprotov5tov6.GetProviderSchemaRequest(req)
	v6Resp, err := s.v6Server.GetProviderSchema(ctx, v6Req)

	if err != nil {
		return nil, err
	}

	return tfprotov6tov5.GetProviderSchemaResponse(v6Resp)
}

func (s v6tov5Server) ImportResourceState(ctx context.Context, req *tfprotov5.ImportResourceStateRequest) (*tfprotov5.ImportResourceStateResponse, error) {
	v6Req := tfprotov5tov6.ImportResourceStateRequest(req)
	v6Resp, err := s.v6Server.ImportResourceState(ctx, v6Req)

	if err != nil {
		return nil, err
	}

	return tfprotov6tov5.ImportResourceStateResponse(v6Resp), nil
}

func (s v6tov5Server) PlanResourceChange(ctx context.Context, req *tfprotov5.PlanResourceChangeRequest) (*tfprotov5.PlanResourceChangeResponse, error) {
	v6Req := tfprotov5tov6.PlanResourceChangeRequest(req)
	v6Resp, err := s.v6Server.PlanResourceChange(ctx, v6Req)

	if err != nil {
		return nil, err
	}

	return tfprotov6tov5.PlanResourceChangeResponse(v6Resp), nil
}

func (s v6tov5Server) PrepareProviderConfig(ctx context.Context, req *tfprotov5.PrepareProviderConfigRequest) (*tfprotov5.PrepareProviderConfigResponse, error) {
	v6Req := tfprotov5tov6.ValidateProviderConfigRequest(req)
	v6Resp, err := s.v6Server.ValidateProviderConfig(ctx, v6Req)

	if err != nil {
		return nil, err
	}

	return tfprotov6tov5.PrepareProviderConfigResponse(v6Resp), nil
}

// ProviderServer is a function compatible with tf5server.Serve.
func (s v6tov5Server) ProviderServer() tfprotov5.ProviderServer {
	return s
}

func (s v6tov5Server) ReadDataSource(ctx context.Context, req *tfprotov5.ReadDataSourceRequest) (*tfprotov5.ReadDataSourceResponse, error) {
	v6Req := tfprotov5tov6.ReadDataSourceRequest(req)
	v6Resp, err := s.v6Server.ReadDataSource(ctx, v6Req)

	if err != nil {
		return nil, err
	}

	return tfprotov6tov5.ReadDataSourceResponse(v6Resp), nil
}

func (s v6tov5Server) ReadResource(ctx context.Context, req *tfprotov5.ReadResourceRequest) (*tfprotov5.ReadResourceResponse, error) {
	v6Req := tfprotov5tov6.ReadResourceRequest(req)
	v6Resp, err := s.v6Server.ReadResource(ctx, v6Req)

	if err != nil {
		return nil, err
	}

	return tfprotov6tov5.ReadResourceResponse(v6Resp), nil
}

func (s v6tov5Server) StopProvider(ctx context.Context, req *tfprotov5.StopProviderRequest) (*tfprotov5.StopProviderResponse, error) {
	v6Req := tfprotov5tov6.StopProviderRequest(req)
	v6Resp, err := s.v6Server.StopProvider(ctx, v6Req)

	if err != nil {
		return nil, err
	}

	return tfprotov6tov5.StopProviderResponse(v6Resp), nil
}

func (s v6tov5Server) UpgradeResourceState(ctx context.Context, req *tfprotov5.UpgradeResourceStateRequest) (*tfprotov5.UpgradeResourceStateResponse, error) {
	v6Req := tfprotov5tov6.UpgradeResourceStateRequest(req)
	v6Resp, err := s.v6Server.UpgradeResourceState(ctx, v6Req)

	if err != nil {
		return nil, err
	}

	return tfprotov6tov5.UpgradeResourceStateResponse(v6Resp), nil
}

func (s v6tov5Server) ValidateDataSourceConfig(ctx context.Context, req *tfprotov5.ValidateDataSourceConfigRequest) (*tfprotov5.ValidateDataSourceConfigResponse, error) {
	v6Req := tfprotov5tov6.ValidateDataResourceConfigRequest(req)
	v6Resp, err := s.v6Server.ValidateDataResourceConfig(ctx, v6Req)

	if err != nil {
		return nil, err
	}

	return tfprotov6tov5.ValidateDataSourceConfigResponse(v6Resp), nil
}

func (s v6tov5Server) ValidateResourceTypeConfig(ctx context.Context, req *tfprotov5.ValidateResourceTypeConfigRequest) (*tfprotov5.ValidateResourceTypeConfigResponse, error) {
	v6Req := tfprotov5tov6.ValidateResourceConfigRequest(req)
	v6Resp, err := s.v6Server.ValidateResourceConfig(ctx, v6Req)

	if err != nil {
		return nil, err
	}

	return tfprotov6tov5.ValidateResourceTypeConfigResponse(v6Resp), nil
}
