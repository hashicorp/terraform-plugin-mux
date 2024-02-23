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

// Temporarily verify that v6tov5Server implements new RPCs correctly.
// Reference: https://github.com/hashicorp/terraform-plugin-mux/issues/210
// Reference: https://github.com/hashicorp/terraform-plugin-mux/issues/219
var (
	_ tfprotov5.FunctionServer = v6tov5Server{}
	//nolint:staticcheck // Intentional verification of interface implementation.
	_ tfprotov5.ResourceServerWithMoveResourceState = v6tov5Server{}
)

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

func (s v6tov5Server) CallFunction(ctx context.Context, req *tfprotov5.CallFunctionRequest) (*tfprotov5.CallFunctionResponse, error) {
	// Remove and call s.v6Server.CallFunction below directly.
	// Reference: https://github.com/hashicorp/terraform-plugin-mux/issues/210
	functionServer, ok := s.v6Server.(tfprotov6.FunctionServer)

	if !ok {
		v5Resp := &tfprotov5.CallFunctionResponse{
			Error: &tfprotov5.FunctionError{
				Text: "Provider Functions Not Implemented: A provider-defined function call was received by the provider, however the provider does not implement functions. " +
					"Either upgrade the provider to a version that implements provider-defined functions or this is a bug in Terraform that should be reported to the Terraform maintainers.",
			},
		}

		return v5Resp, nil
	}

	v6Req := tfprotov5tov6.CallFunctionRequest(req)
	// Reference: https://github.com/hashicorp/terraform-plugin-mux/issues/210
	// v6Resp, err := s.v6Server.CallFunction(ctx, v6Req)
	v6Resp, err := functionServer.CallFunction(ctx, v6Req)

	if err != nil {
		return nil, err
	}

	return tfprotov6tov5.CallFunctionResponse(v6Resp), nil
}

func (s v6tov5Server) ConfigureProvider(ctx context.Context, req *tfprotov5.ConfigureProviderRequest) (*tfprotov5.ConfigureProviderResponse, error) {
	v6Req := tfprotov5tov6.ConfigureProviderRequest(req)
	v6Resp, err := s.v6Server.ConfigureProvider(ctx, v6Req)

	if err != nil {
		return nil, err
	}

	return tfprotov6tov5.ConfigureProviderResponse(v6Resp), nil
}

func (s v6tov5Server) GetFunctions(ctx context.Context, req *tfprotov5.GetFunctionsRequest) (*tfprotov5.GetFunctionsResponse, error) {
	// Remove and call s.v6Server.GetFunctions below directly.
	// Reference: https://github.com/hashicorp/terraform-plugin-mux/issues/210
	functionServer, ok := s.v6Server.(tfprotov6.FunctionServer)

	if !ok {
		v5Resp := &tfprotov5.GetFunctionsResponse{
			Functions: map[string]*tfprotov5.Function{},
		}

		return v5Resp, nil
	}

	v6Req := tfprotov5tov6.GetFunctionsRequest(req)
	// Reference: https://github.com/hashicorp/terraform-plugin-mux/issues/210
	// v6Resp, err := s.v6Server.GetFunctions(ctx, v6Req)
	v6Resp, err := functionServer.GetFunctions(ctx, v6Req)

	if err != nil {
		return nil, err
	}

	return tfprotov6tov5.GetFunctionsResponse(v6Resp), nil
}

func (s v6tov5Server) GetMetadata(ctx context.Context, req *tfprotov5.GetMetadataRequest) (*tfprotov5.GetMetadataResponse, error) {
	v6Req := tfprotov5tov6.GetMetadataRequest(req)
	v6Resp, err := s.v6Server.GetMetadata(ctx, v6Req)

	if err != nil {
		return nil, err
	}

	return tfprotov6tov5.GetMetadataResponse(v6Resp), nil
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

func (s v6tov5Server) MoveResourceState(ctx context.Context, req *tfprotov5.MoveResourceStateRequest) (*tfprotov5.MoveResourceStateResponse, error) {
	// Remove and call s.v6Server.MoveResourceState below directly.
	// Reference: https://github.com/hashicorp/terraform-plugin-mux/issues/219
	//nolint:staticcheck // Intentional verification of interface implementation.
	resourceServer, ok := s.v6Server.(tfprotov6.ResourceServerWithMoveResourceState)

	if !ok {
		v5Resp := &tfprotov5.MoveResourceStateResponse{
			Diagnostics: []*tfprotov5.Diagnostic{
				{
					Severity: tfprotov5.DiagnosticSeverityError,
					Summary:  "MoveResourceState Not Implemented",
					Detail: "A MoveResourceState call was received by the provider, however the provider does not implement the RPC. " +
						"Either upgrade the provider to a version that implements MoveResourceState or this is a bug in Terraform that should be reported to the Terraform maintainers.",
				},
			},
		}

		return v5Resp, nil
	}

	v6Req := tfprotov5tov6.MoveResourceStateRequest(req)
	// Reference: https://github.com/hashicorp/terraform-plugin-mux/issues/219
	// v6Resp, err := s.v6Server.MoveResourceState(ctx, v6Req)
	v6Resp, err := resourceServer.MoveResourceState(ctx, v6Req)

	if err != nil {
		return nil, err
	}

	return tfprotov6tov5.MoveResourceStateResponse(v6Resp), nil
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
