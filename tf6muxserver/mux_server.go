package tf6muxserver

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-mux/internal/logging"
)

var _ tfprotov6.ProviderServer = muxServer{}

// muxServer is a gRPC server implementation that stands in front of other
// gRPC servers, routing requests to them as if they were a single server. It
// should always be instantiated by calling NewMuxServer().
type muxServer struct {
	// Routing for data source types
	dataSources map[string]tfprotov6.ProviderServer

	// Routing for resource types
	resources map[string]tfprotov6.ProviderServer

	// Underlying servers for requests that should be handled by all servers
	servers []tfprotov6.ProviderServer

	// Schemas are cached during server creation
	dataSourceSchemas  map[string]*tfprotov6.Schema
	providerMetaSchema *tfprotov6.Schema
	providerSchema     *tfprotov6.Schema
	resourceSchemas    map[string]*tfprotov6.Schema
}

// ProviderServer is a function compatible with tf6server.Serve.
func (s muxServer) ProviderServer() tfprotov6.ProviderServer {
	return s
}

// NewMuxServer returns a muxed server that will route gRPC requests between
// tfprotov6.ProviderServers specified. The GetProviderSchema method of each
// is called to verify that the overall muxed server is compatible by ensuring:
//
//  - All provider schemas exactly match
//  - All provider meta schemas exactly match
//  - Only one provider implements each managed resource
//  - Only one provider implements each data source
//
// The various schemas are cached and used to respond to the GetProviderSchema
// method of the muxed server.
func NewMuxServer(ctx context.Context, servers ...func() tfprotov6.ProviderServer) (muxServer, error) {
	ctx = logging.InitContext(ctx)
	result := muxServer{
		dataSources:       make(map[string]tfprotov6.ProviderServer),
		dataSourceSchemas: make(map[string]*tfprotov6.Schema),
		resources:         make(map[string]tfprotov6.ProviderServer),
		resourceSchemas:   make(map[string]*tfprotov6.Schema),
	}

	for _, serverFunc := range servers {
		server := serverFunc()

		ctx = logging.Tfprotov6ProviderServerContext(ctx, server)
		logging.MuxTrace(ctx, "calling downstream server")

		resp, err := server.GetProviderSchema(ctx, &tfprotov6.GetProviderSchemaRequest{})

		if err != nil {
			return result, fmt.Errorf("error retrieving schema for %T: %w", server, err)
		}

		for _, diag := range resp.Diagnostics {
			if diag == nil {
				continue
			}
			if diag.Severity != tfprotov6.DiagnosticSeverityError {
				continue
			}
			return result, fmt.Errorf("error retrieving schema for %T:\n\n\tAttribute: %s\n\tSummary: %s\n\tDetail: %s", server, diag.Attribute, diag.Summary, diag.Detail)
		}

		if resp.Provider != nil {
			if result.providerSchema != nil && !schemaEquals(resp.Provider, result.providerSchema) {
				return result, fmt.Errorf("got a different provider schema across servers. Provider schemas must be identical across providers. Diff: %s", schemaDiff(resp.Provider, result.providerSchema))
			}

			result.providerSchema = resp.Provider
		}

		if resp.ProviderMeta != nil {
			if result.providerMetaSchema != nil && !schemaEquals(resp.ProviderMeta, result.providerMetaSchema) {
				return result, fmt.Errorf("got a different provider meta schema across servers. Provider metadata schemas must be identical across providers. Diff: %s", schemaDiff(resp.ProviderMeta, result.providerMetaSchema))
			}

			result.providerMetaSchema = resp.ProviderMeta
		}

		for resourceType, schema := range resp.ResourceSchemas {
			if _, ok := result.resources[resourceType]; ok {
				return result, fmt.Errorf("resource %q is implemented by multiple servers; only one implementation allowed", resourceType)
			}

			result.resources[resourceType] = server
			result.resourceSchemas[resourceType] = schema
		}

		for dataSourceType, schema := range resp.DataSourceSchemas {
			if _, ok := result.dataSources[dataSourceType]; ok {
				return result, fmt.Errorf("data source %q is implemented by multiple servers; only one implementation allowed", dataSourceType)
			}

			result.dataSources[dataSourceType] = server
			result.dataSourceSchemas[dataSourceType] = schema
		}

		result.servers = append(result.servers, server)
	}

	return result, nil
}
