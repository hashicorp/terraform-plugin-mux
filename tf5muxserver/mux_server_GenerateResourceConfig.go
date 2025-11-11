package tf5muxserver

import (
	"context"
	//"fmt"

	//"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov5"
	"github.com/hashicorp/terraform-plugin-mux/internal/logging"
)

func (s *muxServer) GenerateResourceConfig(ctx context.Context, req *tfprotov5.GenerateResourceConfigRequest) (*tfprotov5.GenerateResourceConfigResponse, error) {
	rpc := "GenerateResourceConfig"
	ctx = logging.InitContext(ctx)
	ctx = logging.RpcContext(ctx, rpc)

	server, diags, err := s.getGenerateResourceConfigServer(ctx, req.TypeName)

	if err != nil {
		return nil, err
	}

	// If there is an error diagnostic, return it directly
	if diagnosticsHasError(diags) {
		return &tfprotov5.GenerateResourceConfigResponse{
			Diagnostics: diags,
		}, nil
	}

	return server.GenerateResourceConfig(ctx, req)
}
