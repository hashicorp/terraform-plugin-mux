package tfmux

import (
	"context"

	"github.com/hashicorp/terraform-plugin-go/tfprotov5"
)

type ResourceMuxer struct {
	ValidateResourceTypeConfigHandler  func(ctx context.Context, req *tfprotov5.ValidateResourceTypeConfigRequest) (*tfprotov5.ValidateResourceTypeConfigResponse, error)
	OverrideValidateResourceTypeConfig map[string]func(ctx context.Context, req *tfprotov5.ValidateResourceTypeConfigRequest) (*tfprotov5.ValidateResourceTypeConfigResponse, error)

	ValidateDataSourceConfigHandler  func(ctx context.Context, req *tfprotov5.ValidateDataSourceConfigRequest) (*tfprotov5.ValidateDataSourceConfigResponse, error)
	OverrideValidateDataSourceConfig map[string]func(ctx context.Context, req *tfprotov5.ValidateDataSourceConfigRequest) (*tfprotov5.ValidateDataSourceConfigResponse, error)

	UpgradeResourceStateHandler  func(ctx context.Context, req *tfprotov5.UpgradeResourceStateRequest) (*tfprotov5.UpgradeResourceStateResponse, error)
	OverrideUpgradeResourceState map[string]func(ctx context.Context, req *tfprotov5.UpgradeResourceStateRequest) (*tfprotov5.UpgradeResourceStateResponse, error)

	ImportResourceStateHandler  func(ctx context.Context, req *tfprotov5.ImportResourceStateRequest) (*tfprotov5.ImportResourceStateResponse, error)
	OverrideImportResourceState map[string]func(ctx context.Context, req *tfprotov5.ImportResourceStateRequest) (*tfprotov5.ImportResourceStateResponse, error)

	ReadResourceHandler  func(ctx context.Context, req *tfprotov5.ReadResourceRequest) (*tfprotov5.ReadResourceResponse, error)
	OverrideReadResource map[string]func(ctx context.Context, req *tfprotov5.ReadResourceRequest) (*tfprotov5.ReadResourceResponse, error)

	ReadDataSourceHandler  func(ctx context.Context, req *tfprotov5.ReadDataSourceRequest) (*tfprotov5.ReadDataSourceResponse, error)
	OverrideReadDataSource map[string]func(ctx context.Context, req *tfprotov5.ReadDataSourceRequest) (*tfprotov5.ReadDataSourceResponse, error)

	PlanResourceChangeHandler  func(ctx context.Context, req *tfprotov5.PlanResourceChangeRequest) (*tfprotov5.PlanResourceChangeResponse, error)
	OverridePlanResourceChange map[string]func(ctx context.Context, req *tfprotov5.PlanResourceChangeRequest) (*tfprotov5.PlanResourceChangeResponse, error)

	ApplyResourceChangeHandler  func(ctx context.Context, req *tfprotov5.ApplyResourceChangeRequest) (*tfprotov5.ApplyResourceChangeResponse, error)
	OverrideApplyResourceChange map[string]func(ctx context.Context, req *tfprotov5.ApplyResourceChangeRequest) (*tfprotov5.ApplyResourceChangeResponse, error)
}

func (r ResourceMuxer) ValidateResourceTypeConfig(ctx context.Context, req *tfprotov5.ValidateResourceTypeConfigRequest) func(ctx context.Context, req *tfprotov5.ValidateResourceTypeConfigRequest) (*tfprotov5.ValidateResourceTypeConfigResponse, error) {
	if h, ok := r.OverrideValidateResourceTypeConfig[req.TypeName]; ok {
		return h
	}
	return r.ValidateResourceTypeConfigHandler
}

func (r ResourceMuxer) ValidateDataSourceConfig(ctx context.Context, req *tfprotov5.ValidateDataSourceConfigRequest) func(ctx context.Context, req *tfprotov5.ValidateDataSourceConfigRequest) (*tfprotov5.ValidateDataSourceConfigResponse, error) {
	if h, ok := r.OverrideValidateDataSourceConfig[req.TypeName]; ok {
		return h
	}
	return r.ValidateDataSourceConfigHandler
}

func (r ResourceMuxer) UpgradeResourceState(ctx context.Context, req *tfprotov5.UpgradeResourceStateRequest) func(ctx context.Context, req *tfprotov5.UpgradeResourceStateRequest) (*tfprotov5.UpgradeResourceStateResponse, error) {
	if h, ok := r.OverrideUpgradeResourceState[req.TypeName]; ok {
		return h
	}
	return r.UpgradeResourceStateHandler
}

func (r ResourceMuxer) ImportResourceState(ctx context.Context, req *tfprotov5.ImportResourceStateRequest) func(ctx context.Context, req *tfprotov5.ImportResourceStateRequest) (*tfprotov5.ImportResourceStateResponse, error) {
	if h, ok := r.OverrideImportResourceState[req.TypeName]; ok {
		return h
	}
	return r.ImportResourceStateHandler
}

func (r ResourceMuxer) ReadResource(ctx context.Context, req *tfprotov5.ReadResourceRequest) func(ctx context.Context, req *tfprotov5.ReadResourceRequest) (*tfprotov5.ReadResourceResponse, error) {
	if h, ok := r.OverrideReadResource[req.TypeName]; ok {
		return h
	}
	return r.ReadResourceHandler
}

func (r ResourceMuxer) ReadDataSource(ctx context.Context, req *tfprotov5.ReadDataSourceRequest) func(ctx context.Context, req *tfprotov5.ReadDataSourceRequest) (*tfprotov5.ReadDataSourceResponse, error) {
	if h, ok := r.OverrideReadDataSource[req.TypeName]; ok {
		return h
	}
	return r.ReadDataSourceHandler
}

func (r ResourceMuxer) PlanResourceChange(ctx context.Context, req *tfprotov5.PlanResourceChangeRequest) func(ctx context.Context, req *tfprotov5.PlanResourceChangeRequest) (*tfprotov5.PlanResourceChangeResponse, error) {
	if h, ok := r.OverridePlanResourceChange[req.TypeName]; ok {
		return h
	}
	return r.PlanResourceChangeHandler
}

func (r ResourceMuxer) ApplyResourceChange(ctx context.Context, req *tfprotov5.ApplyResourceChangeRequest) func(ctx context.Context, req *tfprotov5.ApplyResourceChangeRequest) (*tfprotov5.ApplyResourceChangeResponse, error) {
	if h, ok := r.OverrideApplyResourceChange[req.TypeName]; ok {
		return h
	}
	return r.ApplyResourceChangeHandler
}
