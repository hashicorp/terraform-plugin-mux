package tfmux

import (
	"context"

	"github.com/hashicorp/terraform-plugin-mux/internal/tfplugin5"
)

type ResourceMuxer struct {
	ValidateResourceTypeConfigHandler  func(ctx context.Context, req *tfplugin5.ValidateResourceTypeConfig_Request) (*tfplugin5.ValidateResourceTypeConfig_Response, error)
	OverrideValidateResourceTypeConfig map[string]func(ctx context.Context, req *tfplugin5.ValidateResourceTypeConfig_Request) (*tfplugin5.ValidateResourceTypeConfig_Response, error)

	ValidateDataSourceConfigHandler  func(ctx context.Context, req *tfplugin5.ValidateDataSourceConfig_Request) (*tfplugin5.ValidateDataSourceConfig_Response, error)
	OverrideValidateDataSourceConfig map[string]func(ctx context.Context, req *tfplugin5.ValidateDataSourceConfig_Request) (*tfplugin5.ValidateDataSourceConfig_Response, error)

	UpgradeResourceStateHandler  func(ctx context.Context, req *tfplugin5.UpgradeResourceState_Request) (*tfplugin5.UpgradeResourceState_Response, error)
	OverrideUpgradeResourceState map[string]func(ctx context.Context, req *tfplugin5.UpgradeResourceState_Request) (*tfplugin5.UpgradeResourceState_Response, error)

	ImportResourceStateHandler  func(ctx context.Context, req *tfplugin5.ImportResourceState_Request) (*tfplugin5.ImportResourceState_Response, error)
	OverrideImportResourceState map[string]func(ctx context.Context, req *tfplugin5.ImportResourceState_Request) (*tfplugin5.ImportResourceState_Response, error)

	ReadResourceHandler  func(ctx context.Context, req *tfplugin5.ReadResource_Request) (*tfplugin5.ReadResource_Response, error)
	OverrideReadResource map[string]func(ctx context.Context, req *tfplugin5.ReadResource_Request) (*tfplugin5.ReadResource_Response, error)

	ReadDataSourceHandler  func(ctx context.Context, req *tfplugin5.ReadDataSource_Request) (*tfplugin5.ReadDataSource_Response, error)
	OverrideReadDataSource map[string]func(ctx context.Context, req *tfplugin5.ReadDataSource_Request) (*tfplugin5.ReadDataSource_Response, error)

	PlanResourceChangeHandler  func(ctx context.Context, req *tfplugin5.PlanResourceChange_Request) (*tfplugin5.PlanResourceChange_Response, error)
	OverridePlanResourceChange map[string]func(ctx context.Context, req *tfplugin5.PlanResourceChange_Request) (*tfplugin5.PlanResourceChange_Response, error)

	ApplyResourceChangeHandler  func(ctx context.Context, req *tfplugin5.ApplyResourceChange_Request) (*tfplugin5.ApplyResourceChange_Response, error)
	OverrideApplyResourceChange map[string]func(ctx context.Context, req *tfplugin5.ApplyResourceChange_Request) (*tfplugin5.ApplyResourceChange_Response, error)
}

func (r ResourceMuxer) ValidateResourceTypeConfig(ctx context.Context, req *tfplugin5.ValidateResourceTypeConfig_Request) func(ctx context.Context, req *tfplugin5.ValidateResourceTypeConfig_Request) (*tfplugin5.ValidateResourceTypeConfig_Response, error) {
	if h, ok := r.OverrideValidateResourceTypeConfig[req.TypeName]; ok {
		return h
	}
	return r.ValidateResourceTypeConfigHandler
}

func (r ResourceMuxer) ValidateDataSourceConfig(ctx context.Context, req *tfplugin5.ValidateDataSourceConfig_Request) func(ctx context.Context, req *tfplugin5.ValidateDataSourceConfig_Request) (*tfplugin5.ValidateDataSourceConfig_Response, error) {
	if h, ok := r.OverrideValidateDataSourceConfig[req.TypeName]; ok {
		return h
	}
	return r.ValidateDataSourceConfigHandler
}

func (r ResourceMuxer) UpgradeResourceState(ctx context.Context, req *tfplugin5.UpgradeResourceState_Request) func(ctx context.Context, req *tfplugin5.UpgradeResourceState_Request) (*tfplugin5.UpgradeResourceState_Response, error) {
	if h, ok := r.OverrideUpgradeResourceState[req.TypeName]; ok {
		return h
	}
	return r.UpgradeResourceStateHandler
}

func (r ResourceMuxer) ImportResourceState(ctx context.Context, req *tfplugin5.ImportResourceState_Request) func(ctx context.Context, req *tfplugin5.ImportResourceState_Request) (*tfplugin5.ImportResourceState_Response, error) {
	if h, ok := r.OverrideImportResourceState[req.TypeName]; ok {
		return h
	}
	return r.ImportResourceStateHandler
}

func (r ResourceMuxer) ReadResource(ctx context.Context, req *tfplugin5.ReadResource_Request) func(ctx context.Context, req *tfplugin5.ReadResource_Request) (*tfplugin5.ReadResource_Response, error) {
	if h, ok := r.OverrideReadResource[req.TypeName]; ok {
		return h
	}
	return r.ReadResourceHandler
}

func (r ResourceMuxer) ReadDataSource(ctx context.Context, req *tfplugin5.ReadDataSource_Request) func(ctx context.Context, req *tfplugin5.ReadDataSource_Request) (*tfplugin5.ReadDataSource_Response, error) {
	if h, ok := r.OverrideReadDataSource[req.TypeName]; ok {
		return h
	}
	return r.ReadDataSourceHandler
}

func (r ResourceMuxer) PlanResourceChange(ctx context.Context, req *tfplugin5.PlanResourceChange_Request) func(ctx context.Context, req *tfplugin5.PlanResourceChange_Request) (*tfplugin5.PlanResourceChange_Response, error) {
	if h, ok := r.OverridePlanResourceChange[req.TypeName]; ok {
		return h
	}
	return r.PlanResourceChangeHandler
}

func (r ResourceMuxer) ApplyResourceChange(ctx context.Context, req *tfplugin5.ApplyResourceChange_Request) func(ctx context.Context, req *tfplugin5.ApplyResourceChange_Request) (*tfplugin5.ApplyResourceChange_Response, error) {
	if h, ok := r.OverrideApplyResourceChange[req.TypeName]; ok {
		return h
	}
	return r.ApplyResourceChangeHandler
}
