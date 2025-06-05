// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package tf6muxserver

import (
	"context"

	"github.com/hashicorp/terraform-plugin-go/tfprotov6"

	"github.com/hashicorp/terraform-plugin-mux/internal/logging"
)

func (s *muxServer) ValidateListResourceConfig(ctx context.Context, req *tfprotov6.ValidateListResourceConfigRequest) (*tfprotov6.ValidateListResourceConfigResponse, error) {
	rpc := "ValidateListResourceTypeConfig"
	ctx = logging.InitContext(ctx)
	ctx = logging.RpcContext(ctx, rpc)

	server, diags, err := s.getListResourceServer(ctx, req.TypeName)

	if err != nil {
		return nil, err
	}

	if diagnosticsHasError(diags) {
		return &tfprotov6.ValidateListResourceConfigResponse{
			Diagnostics: diags,
		}, nil
	}

	ctx = logging.Tfprotov6ProviderServerContext(ctx, server)
	logging.MuxTrace(ctx, "calling downstream server")

	return server.ValidateListResourceConfig(ctx, req)
}
