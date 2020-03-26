package tfmux

import (
	"context"
	"errors"

	"github.com/hashicorp/terraform-plugin-mux/internal/tfplugin5"
)

var _ tfplugin5.ProviderServer = Server{}

type Server struct {
	GetSchemaHandler                  func(ctx context.Context, req *tfplugin5.GetProviderSchema_Request) func(ctx context.Context, req *tfplugin5.GetProviderSchema_Request) (*tfplugin5.GetProviderSchema_Response, error)
	PrepareProviderConfigHandler      func(ctx context.Context, req *tfplugin5.PrepareProviderConfig_Request) func(ctx context.Context, req *tfplugin5.PrepareProviderConfig_Request) (*tfplugin5.PrepareProviderConfig_Response, error)
	ValidateResourceTypeConfigHandler func(ctx context.Context, req *tfplugin5.ValidateResourceTypeConfig_Request) func(ctx context.Context, req *tfplugin5.ValidateResourceTypeConfig_Request) (*tfplugin5.ValidateResourceTypeConfig_Response, error)
	ValidateDataSourceConfigHandler   func(ctx context.Context, req *tfplugin5.ValidateDataSourceConfig_Request) func(ctx context.Context, req *tfplugin5.ValidateDataSourceConfig_Request) (*tfplugin5.ValidateDataSourceConfig_Response, error)
	UpgradeResourceStateHandler       func(ctx context.Context, req *tfplugin5.UpgradeResourceState_Request) func(ctx context.Context, req *tfplugin5.UpgradeResourceState_Request) (*tfplugin5.UpgradeResourceState_Response, error)
	ConfigureHandler                  func(ctx context.Context, req *tfplugin5.Configure_Request) func(ctx context.Context, req *tfplugin5.Configure_Request) (*tfplugin5.Configure_Response, error)
	ReadResourceHandler               func(ctx context.Context, req *tfplugin5.ReadResource_Request) func(ctx context.Context, req *tfplugin5.ReadResource_Request) (*tfplugin5.ReadResource_Response, error)
	PlanResourceChangeHandler         func(ctx context.Context, req *tfplugin5.PlanResourceChange_Request) func(ctx context.Context, req *tfplugin5.PlanResourceChange_Request) (*tfplugin5.PlanResourceChange_Response, error)
	ApplyResourceChangeHandler        func(ctx context.Context, req *tfplugin5.ApplyResourceChange_Request) func(ctx context.Context, req *tfplugin5.ApplyResourceChange_Request) (*tfplugin5.ApplyResourceChange_Response, error)
	ImportResourceStateHandler        func(ctx context.Context, req *tfplugin5.ImportResourceState_Request) func(ctx context.Context, req *tfplugin5.ImportResourceState_Request) (*tfplugin5.ImportResourceState_Response, error)
	ReadDataSourceHandler             func(ctx context.Context, req *tfplugin5.ReadDataSource_Request) func(ctx context.Context, req *tfplugin5.ReadDataSource_Request) (*tfplugin5.ReadDataSource_Response, error)
	StopHandler                       func(ctx context.Context, req *tfplugin5.Stop_Request) func(ctx context.Context, req *tfplugin5.Stop_Request) (*tfplugin5.Stop_Response, error)
}

func (s Server) GetSchema(ctx context.Context, req *tfplugin5.GetProviderSchema_Request) (*tfplugin5.GetProviderSchema_Response, error) {
	if s.GetSchemaHandler != nil {
		return nil, errors.New("no GetSchema handler defined")
	}
	handler := s.GetSchemaHandler(ctx, req)
	return handler(ctx, req)
}

func (s Server) PrepareProviderConfig(ctx context.Context, req *tfplugin5.PrepareProviderConfig_Request) (*tfplugin5.PrepareProviderConfig_Response, error) {
	if s.PrepareProviderConfigHandler != nil {
		return nil, errors.New("no PrepareProviderConfig handler defined")
	}
	handler := s.PrepareProviderConfigHandler(ctx, req)
	return handler(ctx, req)
}

func (s Server) ValidateResourceTypeConfig(ctx context.Context, req *tfplugin5.ValidateResourceTypeConfig_Request) (*tfplugin5.ValidateResourceTypeConfig_Response, error) {
	if s.ValidateResourceTypeConfigHandler != nil {
		return nil, errors.New("no ValidateResourceTypeConfig handler defined")
	}
	handler := s.ValidateResourceTypeConfigHandler(ctx, req)
	return handler(ctx, req)
}

func (s Server) ValidateDataSourceConfig(ctx context.Context, req *tfplugin5.ValidateDataSourceConfig_Request) (*tfplugin5.ValidateDataSourceConfig_Response, error) {
	if s.ValidateDataSourceConfigHandler != nil {
		return nil, errors.New("no ValidateDataSourceConfig handler defined")
	}
	handler := s.ValidateDataSourceConfigHandler(ctx, req)
	return handler(ctx, req)
}

func (s Server) UpgradeResourceState(ctx context.Context, req *tfplugin5.UpgradeResourceState_Request) (*tfplugin5.UpgradeResourceState_Response, error) {
	if s.UpgradeResourceStateHandler != nil {
		return nil, errors.New("no UpgradeResourceState handler defined")
	}
	handler := s.UpgradeResourceStateHandler(ctx, req)
	return handler(ctx, req)
}

func (s Server) Configure(ctx context.Context, req *tfplugin5.Configure_Request) (*tfplugin5.Configure_Response, error) {
	if s.ConfigureHandler != nil {
		return nil, errors.New("no Configure handler defined")
	}
	handler := s.ConfigureHandler(ctx, req)
	return handler(ctx, req)
}

func (s Server) ReadResource(ctx context.Context, req *tfplugin5.ReadResource_Request) (*tfplugin5.ReadResource_Response, error) {
	if s.ReadResourceHandler != nil {
		return nil, errors.New("no ReadResource handler defined")
	}
	handler := s.ReadResourceHandler(ctx, req)
	return handler(ctx, req)
}

func (s Server) PlanResourceChange(ctx context.Context, req *tfplugin5.PlanResourceChange_Request) (*tfplugin5.PlanResourceChange_Response, error) {
	if s.PlanResourceChangeHandler != nil {
		return nil, errors.New("no PlanResourceChange handler defined")
	}
	handler := s.PlanResourceChangeHandler(ctx, req)
	return handler(ctx, req)
}

func (s Server) ApplyResourceChange(ctx context.Context, req *tfplugin5.ApplyResourceChange_Request) (*tfplugin5.ApplyResourceChange_Response, error) {
	if s.ApplyResourceChangeHandler != nil {
		return nil, errors.New("no ApplyResourceChange handler defined")
	}
	handler := s.ApplyResourceChangeHandler(ctx, req)
	return handler(ctx, req)
}

func (s Server) ImportResourceState(ctx context.Context, req *tfplugin5.ImportResourceState_Request) (*tfplugin5.ImportResourceState_Response, error) {
	if s.ImportResourceStateHandler != nil {
		return nil, errors.New("no ImportResourceState handler defined")
	}
	handler := s.ImportResourceStateHandler(ctx, req)
	return handler(ctx, req)
}

func (s Server) ReadDataSource(ctx context.Context, req *tfplugin5.ReadDataSource_Request) (*tfplugin5.ReadDataSource_Response, error) {
	if s.ReadDataSourceHandler != nil {
		return nil, errors.New("no ReadDataSource handler defined")
	}
	handler := s.ReadDataSourceHandler(ctx, req)
	return handler(ctx, req)
}

func (s Server) Stop(ctx context.Context, req *tfplugin5.Stop_Request) (*tfplugin5.Stop_Response, error) {
	if s.StopHandler != nil {
		return nil, errors.New("no Stop handler defined")
	}
	handler := s.StopHandler(ctx, req)
	return handler(ctx, req)
}
