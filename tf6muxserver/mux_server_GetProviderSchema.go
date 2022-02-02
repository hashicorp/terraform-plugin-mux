package tf6muxserver

import (
	"context"

	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-mux/internal/logging"
)

// GetProviderSchema merges the schemas returned by the
// tfprotov6.ProviderServers associated with muxServer into a single schema.
// Resources and data sources must be returned from only one server. Provider
// and ProviderMeta schemas must be identical between all servers.
func (s muxServer) GetProviderSchema(ctx context.Context, req *tfprotov6.GetProviderSchemaRequest) (*tfprotov6.GetProviderSchemaResponse, error) {
	rpc := "GetProviderSchema"
	ctx = logging.InitContext(ctx)
	ctx = logging.RpcContext(ctx, rpc)
	logging.MuxTrace(ctx, "serving cached schema information")

	return &tfprotov6.GetProviderSchemaResponse{
		Provider:          s.providerSchema,
		ResourceSchemas:   s.resourceSchemas,
		DataSourceSchemas: s.dataSourceSchemas,
		ProviderMeta:      s.providerMetaSchema,
	}, nil
}
