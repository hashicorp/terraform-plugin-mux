package tfmux

import (
	"context"
	"errors"

	"github.com/hashicorp/terraform-plugin-mux/internal/tfplugin5"
)

var _ tfplugin5.ProviderServer = Server{}

type Server struct{}

func (s Server) GetSchema(ctx context.Context, req *tfplugin5.GetProviderSchema_Request) (*tfplugin5.GetProviderSchema_Response, error) {
	return nil, errors.New("not implemented yet")
}

func (s Server) PrepareProviderConfig(ctx context.Context, req *tfplugin5.PrepareProviderConfig_Request) (*tfplugin5.PrepareProviderConfig_Response, error) {
	return nil, errors.New("not implemented yet")
}

func (s Server) ValidateResourceTypeConfig(ctx context.Context, req *tfplugin5.ValidateResourceTypeConfig_Request) (*tfplugin5.ValidateResourceTypeConfig_Response, error) {
	return nil, errors.New("not implemented yet")
}

func (s Server) ValidateDataSourceConfig(ctx context.Context, req *tfplugin5.ValidateDataSourceConfig_Request) (*tfplugin5.ValidateDataSourceConfig_Response, error) {
	return nil, errors.New("not implemented yet")
}

func (s Server) UpgradeResourceState(ctx context.Context, req *tfplugin5.UpgradeResourceState_Request) (*tfplugin5.UpgradeResourceState_Response, error) {
	return nil, errors.New("not implemented yet")
}

func (s Server) Configure(ctx context.Context, req *tfplugin5.Configure_Request) (*tfplugin5.Configure_Response, error) {
	return nil, errors.New("not implemented yet")
}

func (s Server) ReadResource(ctx context.Context, req *tfplugin5.ReadResource_Request) (*tfplugin5.ReadResource_Response, error) {
	return nil, errors.New("not implemented yet")
}

func (s Server) PlanResourceChange(ctx context.Context, req *tfplugin5.PlanResourceChange_Request) (*tfplugin5.PlanResourceChange_Response, error) {
	return nil, errors.New("not implemented yet")
}

func (s Server) ApplyResourceChange(ctx context.Context, req *tfplugin5.ApplyResourceChange_Request) (*tfplugin5.ApplyResourceChange_Response, error) {
	return nil, errors.New("not implemented yet")
}

func (s Server) ImportResourceState(ctx context.Context, req *tfplugin5.ImportResourceState_Request) (*tfplugin5.ImportResourceState_Response, error) {
	return nil, errors.New("not implemented yet")
}

func (s Server) ReadDataSource(ctx context.Context, req *tfplugin5.ReadDataSource_Request) (*tfplugin5.ReadDataSource_Response, error) {
	return nil, errors.New("not implemented yet")
}

func (s Server) Stop(ctx context.Context, req *tfplugin5.Stop_Request) (*tfplugin5.Stop_Response, error) {
	return nil, errors.New("not implemented yet")
}
