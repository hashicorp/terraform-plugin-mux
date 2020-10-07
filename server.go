package tfmux

import (
	"context"
	"errors"

	"github.com/hashicorp/terraform-plugin-go/tfprotov5"
)

var _ tfprotov5.ProviderServer = Server{}

type Server struct {
	GetSchemaHandler                  func(ctx context.Context, req *tfprotov5.GetProviderSchemaRequest) func(ctx context.Context, req *tfprotov5.GetProviderSchemaRequest) (*tfprotov5.GetProviderSchemaResponse, error)
	PrepareProviderConfigHandler      func(ctx context.Context, req *tfprotov5.PrepareProviderConfigRequest) func(ctx context.Context, req *tfprotov5.PrepareProviderConfigRequest) (*tfprotov5.PrepareProviderConfigResponse, error)
	ValidateResourceTypeConfigHandler func(ctx context.Context, req *tfprotov5.ValidateResourceTypeConfigRequest) func(ctx context.Context, req *tfprotov5.ValidateResourceTypeConfigRequest) (*tfprotov5.ValidateResourceTypeConfigResponse, error)
	ValidateDataSourceConfigHandler   func(ctx context.Context, req *tfprotov5.ValidateDataSourceConfigRequest) func(ctx context.Context, req *tfprotov5.ValidateDataSourceConfigRequest) (*tfprotov5.ValidateDataSourceConfigResponse, error)
	UpgradeResourceStateHandler       func(ctx context.Context, req *tfprotov5.UpgradeResourceStateRequest) func(ctx context.Context, req *tfprotov5.UpgradeResourceStateRequest) (*tfprotov5.UpgradeResourceStateResponse, error)
	ConfigureHandler                  func(ctx context.Context, req *tfprotov5.ConfigureProviderRequest) func(ctx context.Context, req *tfprotov5.ConfigureProviderRequest) (*tfprotov5.ConfigureProviderResponse, error)
	ReadResourceHandler               func(ctx context.Context, req *tfprotov5.ReadResourceRequest) func(ctx context.Context, req *tfprotov5.ReadResourceRequest) (*tfprotov5.ReadResourceResponse, error)
	PlanResourceChangeHandler         func(ctx context.Context, req *tfprotov5.PlanResourceChangeRequest) func(ctx context.Context, req *tfprotov5.PlanResourceChangeRequest) (*tfprotov5.PlanResourceChangeResponse, error)
	ApplyResourceChangeHandler        func(ctx context.Context, req *tfprotov5.ApplyResourceChangeRequest) func(ctx context.Context, req *tfprotov5.ApplyResourceChangeRequest) (*tfprotov5.ApplyResourceChangeResponse, error)
	ImportResourceStateHandler        func(ctx context.Context, req *tfprotov5.ImportResourceStateRequest) func(ctx context.Context, req *tfprotov5.ImportResourceStateRequest) (*tfprotov5.ImportResourceStateResponse, error)
	ReadDataSourceHandler             func(ctx context.Context, req *tfprotov5.ReadDataSourceRequest) func(ctx context.Context, req *tfprotov5.ReadDataSourceRequest) (*tfprotov5.ReadDataSourceResponse, error)
	StopHandler                       func(ctx context.Context, req *tfprotov5.StopProviderRequest) func(ctx context.Context, req *tfprotov5.StopProviderRequest) (*tfprotov5.StopProviderResponse, error)
}

func (s Server) GetProviderSchema(ctx context.Context, req *tfprotov5.GetProviderSchemaRequest) (*tfprotov5.GetProviderSchemaResponse, error) {
	if s.GetSchemaHandler != nil {
		return nil, errors.New("no GetSchema handler defined")
	}
	handler := s.GetSchemaHandler(ctx, req)
	return handler(ctx, req)
}

func (s Server) PrepareProviderConfig(ctx context.Context, req *tfprotov5.PrepareProviderConfigRequest) (*tfprotov5.PrepareProviderConfigResponse, error) {
	if s.PrepareProviderConfigHandler != nil {
		return nil, errors.New("no PrepareProviderConfig handler defined")
	}
	handler := s.PrepareProviderConfigHandler(ctx, req)
	return handler(ctx, req)
}

func (s Server) ValidateResourceTypeConfig(ctx context.Context, req *tfprotov5.ValidateResourceTypeConfigRequest) (*tfprotov5.ValidateResourceTypeConfigResponse, error) {
	if s.ValidateResourceTypeConfigHandler != nil {
		return nil, errors.New("no ValidateResourceTypeConfig handler defined")
	}
	handler := s.ValidateResourceTypeConfigHandler(ctx, req)
	return handler(ctx, req)
}

func (s Server) ValidateDataSourceConfig(ctx context.Context, req *tfprotov5.ValidateDataSourceConfigRequest) (*tfprotov5.ValidateDataSourceConfigResponse, error) {
	if s.ValidateDataSourceConfigHandler != nil {
		return nil, errors.New("no ValidateDataSourceConfig handler defined")
	}
	handler := s.ValidateDataSourceConfigHandler(ctx, req)
	return handler(ctx, req)
}

func (s Server) UpgradeResourceState(ctx context.Context, req *tfprotov5.UpgradeResourceStateRequest) (*tfprotov5.UpgradeResourceStateResponse, error) {
	if s.UpgradeResourceStateHandler != nil {
		return nil, errors.New("no UpgradeResourceState handler defined")
	}
	handler := s.UpgradeResourceStateHandler(ctx, req)
	return handler(ctx, req)
}

func (s Server) ConfigureProvider(ctx context.Context, req *tfprotov5.ConfigureProviderRequest) (*tfprotov5.ConfigureProviderResponse, error) {
	if s.ConfigureHandler != nil {
		return nil, errors.New("no Configure handler defined")
	}
	handler := s.ConfigureHandler(ctx, req)
	return handler(ctx, req)
}

func (s Server) ReadResource(ctx context.Context, req *tfprotov5.ReadResourceRequest) (*tfprotov5.ReadResourceResponse, error) {
	if s.ReadResourceHandler != nil {
		return nil, errors.New("no ReadResource handler defined")
	}
	handler := s.ReadResourceHandler(ctx, req)
	return handler(ctx, req)
}

func (s Server) PlanResourceChange(ctx context.Context, req *tfprotov5.PlanResourceChangeRequest) (*tfprotov5.PlanResourceChangeResponse, error) {
	if s.PlanResourceChangeHandler != nil {
		return nil, errors.New("no PlanResourceChange handler defined")
	}
	handler := s.PlanResourceChangeHandler(ctx, req)
	return handler(ctx, req)
}

func (s Server) ApplyResourceChange(ctx context.Context, req *tfprotov5.ApplyResourceChangeRequest) (*tfprotov5.ApplyResourceChangeResponse, error) {
	if s.ApplyResourceChangeHandler != nil {
		return nil, errors.New("no ApplyResourceChange handler defined")
	}
	handler := s.ApplyResourceChangeHandler(ctx, req)
	return handler(ctx, req)
}

func (s Server) ImportResourceState(ctx context.Context, req *tfprotov5.ImportResourceStateRequest) (*tfprotov5.ImportResourceStateResponse, error) {
	if s.ImportResourceStateHandler != nil {
		return nil, errors.New("no ImportResourceState handler defined")
	}
	handler := s.ImportResourceStateHandler(ctx, req)
	return handler(ctx, req)
}

func (s Server) ReadDataSource(ctx context.Context, req *tfprotov5.ReadDataSourceRequest) (*tfprotov5.ReadDataSourceResponse, error) {
	if s.ReadDataSourceHandler != nil {
		return nil, errors.New("no ReadDataSource handler defined")
	}
	handler := s.ReadDataSourceHandler(ctx, req)
	return handler(ctx, req)
}

func (s Server) StopProvider(ctx context.Context, req *tfprotov5.StopProviderRequest) (*tfprotov5.StopProviderResponse, error) {
	if s.StopHandler != nil {
		return nil, errors.New("no Stop handler defined")
	}
	handler := s.StopHandler(ctx, req)
	return handler(ctx, req)
}
