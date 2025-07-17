// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package tf6to5server

import (
	"context"
	"slices"

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

func (s v6tov5Server) CallFunction(ctx context.Context, req *tfprotov5.CallFunctionRequest) (*tfprotov5.CallFunctionResponse, error) {
	v6Req := tfprotov5tov6.CallFunctionRequest(req)

	v6Resp, err := s.v6Server.CallFunction(ctx, v6Req)
	if err != nil {
		return nil, err
	}

	return tfprotov6tov5.CallFunctionResponse(v6Resp), nil
}

func (s v6tov5Server) CloseEphemeralResource(ctx context.Context, req *tfprotov5.CloseEphemeralResourceRequest) (*tfprotov5.CloseEphemeralResourceResponse, error) {
	v6Req := tfprotov5tov6.CloseEphemeralResourceRequest(req)

	v6Resp, err := s.v6Server.CloseEphemeralResource(ctx, v6Req)
	if err != nil {
		return nil, err
	}

	return tfprotov6tov5.CloseEphemeralResourceResponse(v6Resp), nil
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
	v6Req := tfprotov5tov6.GetFunctionsRequest(req)

	v6Resp, err := s.v6Server.GetFunctions(ctx, v6Req)
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

func (s v6tov5Server) GetResourceIdentitySchemas(ctx context.Context, req *tfprotov5.GetResourceIdentitySchemasRequest) (*tfprotov5.GetResourceIdentitySchemasResponse, error) {
	v6Req := tfprotov5tov6.GetResourceIdentitySchemasRequest(req)
	v6Resp, err := s.v6Server.GetResourceIdentitySchemas(ctx, v6Req)

	if err != nil {
		return nil, err
	}

	return tfprotov6tov5.GetResourceIdentitySchemasResponse(v6Resp), nil
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
	v6Req := tfprotov5tov6.MoveResourceStateRequest(req)

	v6Resp, err := s.v6Server.MoveResourceState(ctx, v6Req)
	if err != nil {
		return nil, err
	}

	return tfprotov6tov5.MoveResourceStateResponse(v6Resp), nil
}

func (s v6tov5Server) OpenEphemeralResource(ctx context.Context, req *tfprotov5.OpenEphemeralResourceRequest) (*tfprotov5.OpenEphemeralResourceResponse, error) {
	v6Req := tfprotov5tov6.OpenEphemeralResourceRequest(req)

	v6Resp, err := s.v6Server.OpenEphemeralResource(ctx, v6Req)
	if err != nil {
		return nil, err
	}

	return tfprotov6tov5.OpenEphemeralResourceResponse(v6Resp), nil
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

func (s v6tov5Server) RenewEphemeralResource(ctx context.Context, req *tfprotov5.RenewEphemeralResourceRequest) (*tfprotov5.RenewEphemeralResourceResponse, error) {
	v6Req := tfprotov5tov6.RenewEphemeralResourceRequest(req)

	v6Resp, err := s.v6Server.RenewEphemeralResource(ctx, v6Req)
	if err != nil {
		return nil, err
	}

	return tfprotov6tov5.RenewEphemeralResourceResponse(v6Resp), nil
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

func (s v6tov5Server) UpgradeResourceIdentity(ctx context.Context, req *tfprotov5.UpgradeResourceIdentityRequest) (*tfprotov5.UpgradeResourceIdentityResponse, error) {
	v6Req := tfprotov5tov6.UpgradeResourceIdentityRequest(req)
	v6Resp, err := s.v6Server.UpgradeResourceIdentity(ctx, v6Req)

	if err != nil {
		return nil, err
	}

	return tfprotov6tov5.UpgradeResourceIdentityResponse(v6Resp), nil
}

func (s v6tov5Server) ValidateDataSourceConfig(ctx context.Context, req *tfprotov5.ValidateDataSourceConfigRequest) (*tfprotov5.ValidateDataSourceConfigResponse, error) {
	v6Req := tfprotov5tov6.ValidateDataResourceConfigRequest(req)
	v6Resp, err := s.v6Server.ValidateDataResourceConfig(ctx, v6Req)

	if err != nil {
		return nil, err
	}

	return tfprotov6tov5.ValidateDataSourceConfigResponse(v6Resp), nil
}

func (s v6tov5Server) ValidateEphemeralResourceConfig(ctx context.Context, req *tfprotov5.ValidateEphemeralResourceConfigRequest) (*tfprotov5.ValidateEphemeralResourceConfigResponse, error) {
	v6Req := tfprotov5tov6.ValidateEphemeralResourceConfigRequest(req)

	v6Resp, err := s.v6Server.ValidateEphemeralResourceConfig(ctx, v6Req)
	if err != nil {
		return nil, err
	}

	return tfprotov6tov5.ValidateEphemeralResourceConfigResponse(v6Resp), nil
}

func (s v6tov5Server) ValidateResourceTypeConfig(ctx context.Context, req *tfprotov5.ValidateResourceTypeConfigRequest) (*tfprotov5.ValidateResourceTypeConfigResponse, error) {
	v6Req := tfprotov5tov6.ValidateResourceConfigRequest(req)
	v6Resp, err := s.v6Server.ValidateResourceConfig(ctx, v6Req)

	if err != nil {
		return nil, err
	}

	return tfprotov6tov5.ValidateResourceTypeConfigResponse(v6Resp), nil
}

func (s v6tov5Server) ValidateListResourceConfig(ctx context.Context, req *tfprotov5.ValidateListResourceConfigRequest) (*tfprotov5.ValidateListResourceConfigResponse, error) {
	// TODO: Remove and call s.v6Server.ValidateListResourceConfig below directly once interface becomes required
	listResourceServer, ok := s.v6Server.(tfprotov6.ListResourceServer)
	if !ok {
		v5Resp := &tfprotov5.ValidateListResourceConfigResponse{
			Diagnostics: []*tfprotov5.Diagnostic{
				{
					Severity: tfprotov5.DiagnosticSeverityError,
					Summary:  "ValidateListResourceConfig Not Implemented",
					Detail: "A ValidateListResourceConfig call was received by the provider, however the provider does not implement the RPC. " +
						"Either upgrade the provider to a version that implements ValidateListResourceConfig or this is a bug in Terraform that should be reported to the Terraform maintainers.",
				},
			},
		}

		return v5Resp, nil
	}

	v6Req := tfprotov5tov6.ValidateListResourceConfigRequest(req)

	v6Resp, err := listResourceServer.ValidateListResourceConfig(ctx, v6Req)
	if err != nil {
		return nil, err
	}

	return tfprotov6tov5.ValidateListResourceConfigResponse(v6Resp), nil
}

func (s v6tov5Server) ListResource(ctx context.Context, req *tfprotov5.ListResourceRequest) (*tfprotov5.ListResourceServerStream, error) {
	// TODO: Remove and call s.v6Server.ListResource below directly once interface becomes required
	listResourceServer, ok := s.v6Server.(tfprotov6.ListResourceServer)
	if !ok {
		v5Resp := &tfprotov5.ListResourceServerStream{
			Results: slices.Values([]tfprotov5.ListResourceResult{
				{
					Diagnostics: []*tfprotov5.Diagnostic{
						{
							Severity: tfprotov5.DiagnosticSeverityError,
							Summary:  "ListResource Not Implemented",
							Detail: "A ListResource call was received by the provider, however the provider does not implement the RPC. " +
								"Either upgrade the provider to a version that implements ListResource or this is a bug in Terraform that should be reported to the Terraform maintainers.",
						},
					},
				},
			}),
		}

		return v5Resp, nil
	}

	v6Req := tfprotov5tov6.ListResourceRequest(req)

	v6Resp, err := listResourceServer.ListResource(ctx, v6Req)
	if err != nil {
		return nil, err
	}

	return tfprotov6tov5.ListResourceServerStream(v6Resp), nil
}

func (s v6tov5Server) ValidateActionConfig(ctx context.Context, req *tfprotov5.ValidateActionConfigRequest) (*tfprotov5.ValidateActionConfigResponse, error) {
	// TODO: Remove and call s.v6Server.ValidateActionConfig below directly once interface becomes required
	actionServer, ok := s.v6Server.(tfprotov6.ActionServer)
	if !ok {
		v5Resp := &tfprotov5.ValidateActionConfigResponse{
			Diagnostics: []*tfprotov5.Diagnostic{
				{
					Severity: tfprotov5.DiagnosticSeverityError,
					Summary:  "ValidateActionConfig Not Implemented",
					Detail: "A ValidateActionConfig call was received by the provider, however the provider does not implement the RPC. " +
						"Either upgrade the provider to a version that implements ValidateActionConfig or this is a bug in Terraform that should be reported to the Terraform maintainers.",
				},
			},
		}

		return v5Resp, nil
	}

	v6Req := tfprotov5tov6.ValidateActionConfigRequest(req)

	// v6Resp, err := s.v6Server.ValidateActionConfig(ctx, v6Req)
	v6Resp, err := actionServer.ValidateActionConfig(ctx, v6Req)
	if err != nil {
		return nil, err
	}

	return tfprotov6tov5.ValidateActionConfigResponse(v6Resp), nil
}

func (s v6tov5Server) PlanAction(ctx context.Context, req *tfprotov5.PlanActionRequest) (*tfprotov5.PlanActionResponse, error) {
	// TODO: Remove and call s.v6Server.PlanAction below directly once interface becomes required
	actionServer, ok := s.v6Server.(tfprotov6.ActionServer)
	if !ok {
		v5Resp := &tfprotov5.PlanActionResponse{
			Diagnostics: []*tfprotov5.Diagnostic{
				{
					Severity: tfprotov5.DiagnosticSeverityError,
					Summary:  "PlanAction Not Implemented",
					Detail: "A PlanAction call was received by the provider, however the provider does not implement the RPC. " +
						"Either upgrade the provider to a version that implements PlanAction or this is a bug in Terraform that should be reported to the Terraform maintainers.",
				},
			},
		}

		return v5Resp, nil
	}

	v6Req := tfprotov5tov6.PlanActionRequest(req)

	// v6Resp, err := s.v6Server.PlanAction(ctx, v6Req)
	v6Resp, err := actionServer.PlanAction(ctx, v6Req)
	if err != nil {
		return nil, err
	}

	return tfprotov6tov5.PlanActionResponse(v6Resp), nil
}

func (s v6tov5Server) InvokeAction(ctx context.Context, req *tfprotov5.InvokeActionRequest) (*tfprotov5.InvokeActionServerStream, error) {
	// TODO: Remove and call s.v6Server.InvokeAction below directly once interface becomes required
	actionServer, ok := s.v6Server.(tfprotov6.ActionServer)
	if !ok {
		v5Resp := &tfprotov5.InvokeActionServerStream{
			Events: slices.Values([]tfprotov5.InvokeActionEvent{
				{
					Type: tfprotov5.CompletedInvokeActionEventType{
						Diagnostics: []*tfprotov5.Diagnostic{
							{
								Severity: tfprotov5.DiagnosticSeverityError,
								Summary:  "InvokeAction Not Implemented",
								Detail: "An InvokeAction call was received by the provider, however the provider does not implement the RPC. " +
									"Either upgrade the provider to a version that implements InvokeAction or this is a bug in Terraform that should be reported to the Terraform maintainers.",
							},
						},
					},
				},
			}),
		}

		return v5Resp, nil
	}

	v6Req := tfprotov5tov6.InvokeActionRequest(req)

	// v6Resp, err := s.v6Server.InvokeAction(ctx, v6Req)
	v6Resp, err := actionServer.InvokeAction(ctx, v6Req)
	if err != nil {
		return nil, err
	}

	return tfprotov6tov5.InvokeActionServerStream(v6Resp), nil
}
