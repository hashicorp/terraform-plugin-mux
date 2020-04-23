package tfmux

import (
	"context"
	"errors"
	"fmt"

	"github.com/hashicorp/terraform-plugin-mux/internal/tfplugin5"
)

var _ tfplugin5.ProviderServer = SchemaServer{}

type SchemaServerFactory struct {
	resources   map[string]int
	dataSources map[string]int
	servers     []func() tfplugin5.ProviderServer

	GetSchemaHandler             func(context.Context, *tfplugin5.GetProviderSchema_Request) (*tfplugin5.GetProviderSchema_Response, error)
	PrepareProviderConfigHandler func(context.Context, *tfplugin5.PrepareProviderConfig_Request) (*tfplugin5.PrepareProviderConfig_Response, error)
}

func NewSchemaServerFactory(ctx context.Context, servers ...func() tfplugin5.ProviderServer) (SchemaServerFactory, error) {
	var factory SchemaServerFactory
	factory.servers = make([]func() tfplugin5.ProviderServer, len(servers))
	for pos, server := range servers {
		s := server()
		resp, err := s.GetSchema(ctx, &tfplugin5.GetProviderSchema_Request{})
		if err != nil {
			return factory, fmt.Errorf("error retrieving schema for %T: %w", s, err)
		}

		factory.servers[pos] = server

		// TODO: handle Diagnostics in resp
		for resource := range resp.ResourceSchemas {
			if v, ok := factory.resources[resource]; ok {
				return factory, fmt.Errorf("resource %q supported by multiple server implementations (%T, %T); remove support from one", resource, factory.servers[v], s)
			}
			factory.resources[resource] = pos
		}
		for data := range resp.DataSourceSchemas {
			if v, ok := factory.dataSources[data]; ok {
				return factory, fmt.Errorf("data source %q supported by multiple server implementations (%T, %T); remove support from one", data, factory.servers[v], s)
			}
			factory.dataSources[data] = pos
		}
	}
	return factory, nil
}

func (s SchemaServerFactory) Server() SchemaServer {
	res := SchemaServer{
		getSchemaHandler:             s.GetSchemaHandler,
		prepareProviderConfigHandler: s.PrepareProviderConfigHandler,
		servers:                      make([]tfplugin5.ProviderServer, len(s.servers)),
	}
	for pos, server := range s.servers {
		res.servers[pos] = server()
	}
	for r, pos := range s.resources {
		res.resources[r] = res.servers[pos]
	}
	for ds, pos := range s.dataSources {
		res.dataSources[ds] = res.servers[pos]
	}
	return res
}

type SchemaServer struct {
	resources   map[string]tfplugin5.ProviderServer
	dataSources map[string]tfplugin5.ProviderServer
	servers     []tfplugin5.ProviderServer

	getSchemaHandler             func(context.Context, *tfplugin5.GetProviderSchema_Request) (*tfplugin5.GetProviderSchema_Response, error)
	prepareProviderConfigHandler func(context.Context, *tfplugin5.PrepareProviderConfig_Request) (*tfplugin5.PrepareProviderConfig_Response, error)
}

func (s SchemaServer) GetSchema(ctx context.Context, req *tfplugin5.GetProviderSchema_Request) (*tfplugin5.GetProviderSchema_Response, error) {
	if s.getSchemaHandler != nil {
		return nil, errors.New("no GetSchema handler defined")
	}
	return s.getSchemaHandler(ctx, req)
}

func (s SchemaServer) PrepareProviderConfig(ctx context.Context, req *tfplugin5.PrepareProviderConfig_Request) (*tfplugin5.PrepareProviderConfig_Response, error) {
	if s.prepareProviderConfigHandler != nil {
		return nil, errors.New("no PrepareProviderConfig handler defined")
	}
	return s.prepareProviderConfigHandler(ctx, req)
}

func (s SchemaServer) ValidateResourceTypeConfig(ctx context.Context, req *tfplugin5.ValidateResourceTypeConfig_Request) (*tfplugin5.ValidateResourceTypeConfig_Response, error) {
	h, ok := s.resources[req.TypeName]
	if !ok {
		return nil, fmt.Errorf("%q isn't supported by any servers", req.TypeName)
	}
	return h.ValidateResourceTypeConfig(ctx, req)
}

func (s SchemaServer) ValidateDataSourceConfig(ctx context.Context, req *tfplugin5.ValidateDataSourceConfig_Request) (*tfplugin5.ValidateDataSourceConfig_Response, error) {
	h, ok := s.dataSources[req.TypeName]
	if !ok {
		return nil, fmt.Errorf("%q isn't supported by any servers", req.TypeName)
	}
	return h.ValidateDataSourceConfig(ctx, req)
}

func (s SchemaServer) UpgradeResourceState(ctx context.Context, req *tfplugin5.UpgradeResourceState_Request) (*tfplugin5.UpgradeResourceState_Response, error) {
	h, ok := s.resources[req.TypeName]
	if !ok {
		return nil, fmt.Errorf("%q isn't supported by any servers", req.TypeName)
	}
	return h.UpgradeResourceState(ctx, req)
}

func (s SchemaServer) Configure(ctx context.Context, req *tfplugin5.Configure_Request) (*tfplugin5.Configure_Response, error) {
	var diags []*tfplugin5.Diagnostic
	for _, server := range s.servers {
		resp, err := server.Configure(ctx, req)
		if err != nil {
			return resp, fmt.Errorf("error configuring %T: %w", server, err)
		}
		for _, diag := range resp.Diagnostics {
			if diag == nil {
				continue
			}
			diags = append(diags, diag)
			if diag.Severity != tfplugin5.Diagnostic_ERROR {
				continue
			}
			resp.Diagnostics = diags
			return resp, err
		}
	}
	return &tfplugin5.Configure_Response{Diagnostics: diags}, nil
}

func (s SchemaServer) ReadResource(ctx context.Context, req *tfplugin5.ReadResource_Request) (*tfplugin5.ReadResource_Response, error) {
	h, ok := s.resources[req.TypeName]
	if !ok {
		return nil, fmt.Errorf("%q isn't supported by any servers", req.TypeName)
	}
	return h.ReadResource(ctx, req)
}

func (s SchemaServer) PlanResourceChange(ctx context.Context, req *tfplugin5.PlanResourceChange_Request) (*tfplugin5.PlanResourceChange_Response, error) {
	h, ok := s.resources[req.TypeName]
	if !ok {
		return nil, fmt.Errorf("%q isn't supported by any servers", req.TypeName)
	}
	return h.PlanResourceChange(ctx, req)
}

func (s SchemaServer) ApplyResourceChange(ctx context.Context, req *tfplugin5.ApplyResourceChange_Request) (*tfplugin5.ApplyResourceChange_Response, error) {
	h, ok := s.resources[req.TypeName]
	if !ok {
		return nil, fmt.Errorf("%q isn't supported by any servers", req.TypeName)
	}
	return h.ApplyResourceChange(ctx, req)
}

func (s SchemaServer) ImportResourceState(ctx context.Context, req *tfplugin5.ImportResourceState_Request) (*tfplugin5.ImportResourceState_Response, error) {
	h, ok := s.resources[req.TypeName]
	if !ok {
		return nil, fmt.Errorf("%q isn't supported by any servers", req.TypeName)
	}
	return h.ImportResourceState(ctx, req)
}

func (s SchemaServer) ReadDataSource(ctx context.Context, req *tfplugin5.ReadDataSource_Request) (*tfplugin5.ReadDataSource_Response, error) {
	h, ok := s.dataSources[req.TypeName]
	if !ok {
		return nil, fmt.Errorf("%q isn't supported by any servers", req.TypeName)
	}
	return h.ReadDataSource(ctx, req)
}

func (s SchemaServer) Stop(ctx context.Context, req *tfplugin5.Stop_Request) (*tfplugin5.Stop_Response, error) {
	for _, server := range s.servers {
		resp, err := server.Stop(ctx, req)
		if err != nil {
			return resp, fmt.Errorf("error stopping %T: %w", server, err)
		}
		if resp.Error != "" {
			return resp, err
		}
	}
	return &tfplugin5.Stop_Response{}, nil
}
