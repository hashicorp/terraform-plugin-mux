package tfmux

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-mux/internal/tfplugin5"
)

var _ tfplugin5.ProviderServer = SchemaServer{}

type SchemaServerFactory struct {
	// determine which servers will respond to which requests
	resources   map[string]int
	dataSources map[string]int
	servers     []func() tfplugin5.ProviderServer

	// we respond to GetSchema requests using these schemas
	resourceSchemas    map[string]*tfplugin5.Schema
	dataSourceSchemas  map[string]*tfplugin5.Schema
	providerSchema     *tfplugin5.Schema
	providerMetaSchema *tfplugin5.Schema

	// any non-error diagnostics should get bubbled up, so we store them here
	diagnostics []*tfplugin5.Diagnostic

	// we just store these to surface better errors
	// track which server we got the provider schema and provider meta
	// schema from
	providerSchemaFrom     int
	providerMetaSchemaFrom int

	// this must be set manually, and will likely always be set to the
	// position of SDKv2 in the server list, as newer server
	// implementations don't need this RPC call.
	PrepareProviderConfigServer int
}

func NewSchemaServerFactory(ctx context.Context, servers ...func() tfplugin5.ProviderServer) (SchemaServerFactory, error) {
	var factory SchemaServerFactory

	// know when these are unset vs set to the element in pos 0
	factory.providerSchemaFrom = -1
	factory.providerMetaSchemaFrom = -1
	factory.PrepareProviderConfigServer = -1

	factory.servers = make([]func() tfplugin5.ProviderServer, len(servers))
	for pos, server := range servers {
		s := server()
		resp, err := s.GetSchema(ctx, &tfplugin5.GetProviderSchema_Request{})
		if err != nil {
			return factory, fmt.Errorf("error retrieving schema for %T: %w", s, err)
		}

		factory.servers[pos] = server

		for _, diag := range resp.Diagnostics {
			if diag == nil {
				continue
			}
			if diag.Severity != tfplugin5.Diagnostic_ERROR {
				factory.diagnostics = append(factory.diagnostics, diag)
				continue
			}
			return factory, fmt.Errorf("error retrieving schema for %T:\n\n\tAttribute: %s\n\tSummary: %s\n\tDetail: %s", s, diag.Attribute, diag.Summary, diag.Detail)
		}
		if resp.Provider != nil && factory.providerSchema != nil {
			return factory, fmt.Errorf("provider schema supported by multiple server implementations (%T, %T), remove support from one", factory.servers[factory.providerSchemaFrom], s)
		} else if resp.Provider != nil {
			factory.providerSchemaFrom = pos
			factory.providerSchema = resp.Provider
		}
		if resp.ProviderMeta != nil && factory.providerMetaSchema != nil {
			return factory, fmt.Errorf("provider_meta schema supported by multiple server implementations (%T, %T), remove support from one", factory.servers[factory.providerMetaSchemaFrom], s)
		} else if resp.ProviderMeta != nil {
			factory.providerMetaSchemaFrom = pos
			factory.providerMetaSchema = resp.ProviderMeta
		}
		for resource, schema := range resp.ResourceSchemas {
			if v, ok := factory.resources[resource]; ok {
				return factory, fmt.Errorf("resource %q supported by multiple server implementations (%T, %T); remove support from one", resource, factory.servers[v], s)
			}
			factory.resources[resource] = pos
			factory.resourceSchemas[resource] = schema
		}
		for data, schema := range resp.DataSourceSchemas {
			if v, ok := factory.dataSources[data]; ok {
				return factory, fmt.Errorf("data source %q supported by multiple server implementations (%T, %T); remove support from one", data, factory.servers[v], s)
			}
			factory.dataSources[data] = pos
			factory.dataSourceSchemas[data] = schema
		}
	}
	return factory, nil
}

func (s SchemaServerFactory) getSchemaHandler(_ context.Context, _ *tfplugin5.GetProviderSchema_Request) (*tfplugin5.GetProviderSchema_Response, error) {
	return &tfplugin5.GetProviderSchema_Response{
		Provider:          s.providerSchema,
		ResourceSchemas:   s.resourceSchemas,
		DataSourceSchemas: s.dataSourceSchemas,
		ProviderMeta:      s.providerMetaSchema,
	}, nil
}

func (s SchemaServerFactory) Server() SchemaServer {
	res := SchemaServer{
		getSchemaHandler: s.getSchemaHandler,
		servers:          make([]tfplugin5.ProviderServer, len(s.servers)),
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
	if len(res.servers) <= s.PrepareProviderConfigServer {
		panic(fmt.Sprintf("Tried to use PrepareProviderConfig from server in 0-indexed position %d, but only %d servers are available.", s.PrepareProviderConfigServer, len(res.servers)))
	}
	if s.PrepareProviderConfigServer >= 0 {
		res.prepareProviderConfigHandler = res.servers[s.PrepareProviderConfigServer].PrepareProviderConfig
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
	return s.getSchemaHandler(ctx, req)
}

func (s SchemaServer) PrepareProviderConfig(ctx context.Context, req *tfplugin5.PrepareProviderConfig_Request) (*tfplugin5.PrepareProviderConfig_Response, error) {
	if s.prepareProviderConfigHandler == nil {
		return &tfplugin5.PrepareProviderConfig_Response{
			PreparedConfig: req.Config,
		}, nil
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
