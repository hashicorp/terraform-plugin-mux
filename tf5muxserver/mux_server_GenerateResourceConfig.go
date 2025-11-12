// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package tf5muxserver

import (
	"context"

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

	// TODO: Remove and call server.GenerateResourceConfig below directly once interface becomes required.
	generateResourceConfigServer, ok := server.(tfprotov5.GenerateResourceConfigServer)
	if !ok {
		resp := &tfprotov5.GenerateResourceConfigResponse{
			Diagnostics: []*tfprotov5.Diagnostic{
				{
					Severity: tfprotov5.DiagnosticSeverityError,
					Summary:  "GenerateResourceConfig Not Implemented",
					Detail: "A GenerateResourceConfig call was received by the provider, however the provider does not implement GenerateResourceConfig. " +
						"Either upgrade the provider to a version that implements GenerateResourceConfig or this is a bug in Terraform that should be reported to the Terraform maintainers.",
				},
			},
		}

		return resp, nil
	}

	ctx = logging.Tfprotov5ProviderServerContext(ctx, server)
	logging.MuxTrace(ctx, "calling downstream server")

	return generateResourceConfigServer.GenerateResourceConfig(ctx, req)
}
