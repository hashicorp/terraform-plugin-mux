package tfmux

import (
	"context"

	"github.com/hashicorp/terraform-plugin-go/tfprotov5"
)

// NewResourceServer returns a Server that will use `server` for all requests,
// except for resource-scoped requests for resources listed in `overrides`, for
// which `overrideServer` will be used, instead.
func NewResourceServer(server, overrideServer tfprotov5.ProviderServer, overrides []string) Server {

	validateResourceTypeConfig := map[string]func(ctx context.Context, req *tfprotov5.ValidateResourceTypeConfigRequest) (*tfprotov5.ValidateResourceTypeConfigResponse, error){}
	validateDataSourceConfig := map[string]func(ctx context.Context, req *tfprotov5.ValidateDataSourceConfigRequest) (*tfprotov5.ValidateDataSourceConfigResponse, error){}
	upgradeResourceState := map[string]func(ctx context.Context, req *tfprotov5.UpgradeResourceStateRequest) (*tfprotov5.UpgradeResourceStateResponse, error){}
	importResourceState := map[string]func(ctx context.Context, req *tfprotov5.ImportResourceStateRequest) (*tfprotov5.ImportResourceStateResponse, error){}
	readResource := map[string]func(ctx context.Context, req *tfprotov5.ReadResourceRequest) (*tfprotov5.ReadResourceResponse, error){}
	readDataSource := map[string]func(ctx context.Context, req *tfprotov5.ReadDataSourceRequest) (*tfprotov5.ReadDataSourceResponse, error){}
	planResourceChange := map[string]func(ctx context.Context, req *tfprotov5.PlanResourceChangeRequest) (*tfprotov5.PlanResourceChangeResponse, error){}
	applyResourceChange := map[string]func(ctx context.Context, req *tfprotov5.ApplyResourceChangeRequest) (*tfprotov5.ApplyResourceChangeResponse, error){}

	for _, override := range overrides {
		validateResourceTypeConfig[override] = overrideServer.ValidateResourceTypeConfig
		validateDataSourceConfig[override] = overrideServer.ValidateDataSourceConfig
		upgradeResourceState[override] = overrideServer.UpgradeResourceState
		importResourceState[override] = overrideServer.ImportResourceState
		readResource[override] = overrideServer.ReadResource
		readDataSource[override] = overrideServer.ReadDataSource
		planResourceChange[override] = overrideServer.PlanResourceChange
		applyResourceChange[override] = overrideServer.ApplyResourceChange
	}
	mux := ResourceMuxer{
		ValidateResourceTypeConfigHandler:  server.ValidateResourceTypeConfig,
		OverrideValidateResourceTypeConfig: validateResourceTypeConfig,
		ValidateDataSourceConfigHandler:    server.ValidateDataSourceConfig,
		OverrideValidateDataSourceConfig:   validateDataSourceConfig,
		UpgradeResourceStateHandler:        server.UpgradeResourceState,
		OverrideUpgradeResourceState:       upgradeResourceState,
		ImportResourceStateHandler:         server.ImportResourceState,
		OverrideImportResourceState:        importResourceState,
		ReadResourceHandler:                server.ReadResource,
		OverrideReadResource:               readResource,
		ReadDataSourceHandler:              server.ReadDataSource,
		OverrideReadDataSource:             readDataSource,
		PlanResourceChangeHandler:          server.PlanResourceChange,
		OverridePlanResourceChange:         planResourceChange,
		ApplyResourceChangeHandler:         server.ApplyResourceChange,
		OverrideApplyResourceChange:        applyResourceChange,
	}

	return Server{
		GetSchemaHandler: func(context.Context, *tfprotov5.GetProviderSchemaRequest) func(context.Context, *tfprotov5.GetProviderSchemaRequest) (*tfprotov5.GetProviderSchemaResponse, error) {
			return server.GetProviderSchema
		},
		PrepareProviderConfigHandler: func(context.Context, *tfprotov5.PrepareProviderConfigRequest) func(context.Context, *tfprotov5.PrepareProviderConfigRequest) (*tfprotov5.PrepareProviderConfigResponse, error) {
			return server.PrepareProviderConfig
		},
		ValidateResourceTypeConfigHandler: mux.ValidateResourceTypeConfig,
		ValidateDataSourceConfigHandler:   mux.ValidateDataSourceConfig,
		UpgradeResourceStateHandler:       mux.UpgradeResourceState,
		ImportResourceStateHandler:        mux.ImportResourceState,
		ConfigureHandler: func(context.Context, *tfprotov5.ConfigureProviderRequest) func(context.Context, *tfprotov5.ConfigureProviderRequest) (*tfprotov5.ConfigureProviderResponse, error) {
			return server.ConfigureProvider
		},
		ReadResourceHandler:        mux.ReadResource,
		ReadDataSourceHandler:      mux.ReadDataSource,
		PlanResourceChangeHandler:  mux.PlanResourceChange,
		ApplyResourceChangeHandler: mux.ApplyResourceChange,
		StopHandler: func(context.Context, *tfprotov5.StopProviderRequest) func(context.Context, *tfprotov5.StopProviderRequest) (*tfprotov5.StopProviderResponse, error) {
			return server.StopProvider
		},
	}
}
