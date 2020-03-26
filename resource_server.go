package tfmux

import (
	"context"

	"github.com/hashicorp/terraform-plugin-mux/internal/tfplugin5"
)

// NewResourceServer returns a Server that will use `server` for all requests,
// except for resource-scoped requests for resources listed in `overrides`, for
// which `overrideServer` will be used, instead.
func NewResourceServer(server, overrideServer tfplugin5.ProviderServer, overrides []string) Server {

	validateResourceTypeConfig := map[string]func(ctx context.Context, req *tfplugin5.ValidateResourceTypeConfig_Request) (*tfplugin5.ValidateResourceTypeConfig_Response, error){}
	validateDataSourceConfig := map[string]func(ctx context.Context, req *tfplugin5.ValidateDataSourceConfig_Request) (*tfplugin5.ValidateDataSourceConfig_Response, error){}
	upgradeResourceState := map[string]func(ctx context.Context, req *tfplugin5.UpgradeResourceState_Request) (*tfplugin5.UpgradeResourceState_Response, error){}
	importResourceState := map[string]func(ctx context.Context, req *tfplugin5.ImportResourceState_Request) (*tfplugin5.ImportResourceState_Response, error){}
	readResource := map[string]func(ctx context.Context, req *tfplugin5.ReadResource_Request) (*tfplugin5.ReadResource_Response, error){}
	readDataSource := map[string]func(ctx context.Context, req *tfplugin5.ReadDataSource_Request) (*tfplugin5.ReadDataSource_Response, error){}
	planResourceChange := map[string]func(ctx context.Context, req *tfplugin5.PlanResourceChange_Request) (*tfplugin5.PlanResourceChange_Response, error){}
	applyResourceChange := map[string]func(ctx context.Context, req *tfplugin5.ApplyResourceChange_Request) (*tfplugin5.ApplyResourceChange_Response, error){}

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
		GetSchemaHandler: func(context.Context, *tfplugin5.GetProviderSchema_Request) func(context.Context, *tfplugin5.GetProviderSchema_Request) (*tfplugin5.GetProviderSchema_Response, error) {
			return server.GetSchema
		},
		PrepareProviderConfigHandler: func(context.Context, *tfplugin5.PrepareProviderConfig_Request) func(context.Context, *tfplugin5.PrepareProviderConfig_Request) (*tfplugin5.PrepareProviderConfig_Response, error) {
			return server.PrepareProviderConfig
		},
		ValidateResourceTypeConfigHandler: mux.ValidateResourceTypeConfig,
		ValidateDataSourceConfigHandler:   mux.ValidateDataSourceConfig,
		UpgradeResourceStateHandler:       mux.UpgradeResourceState,
		ImportResourceStateHandler:        mux.ImportResourceState,
		ConfigureHandler: func(context.Context, *tfplugin5.Configure_Request) func(context.Context, *tfplugin5.Configure_Request) (*tfplugin5.Configure_Response, error) {
			return server.Configure
		},
		ReadResourceHandler:        mux.ReadResource,
		ReadDataSourceHandler:      mux.ReadDataSource,
		PlanResourceChangeHandler:  mux.PlanResourceChange,
		ApplyResourceChangeHandler: mux.ApplyResourceChange,
		StopHandler: func(context.Context, *tfplugin5.Stop_Request) func(context.Context, *tfplugin5.Stop_Request) (*tfplugin5.Stop_Response, error) {
			return server.Stop
		},
	}
}
