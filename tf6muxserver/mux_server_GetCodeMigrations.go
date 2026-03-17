// Copyright IBM Corp. 2020, 2026
// SPDX-License-Identifier: MPL-2.0

package tf6muxserver

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-go/tfprotov6"

	"github.com/hashicorp/terraform-plugin-mux/internal/logging"
)

func (s *muxServer) GetCodeMigrations(ctx context.Context, req *tfprotov6.GetCodeMigrationsRequest) (*tfprotov6.GetCodeMigrationsResponse, error) {
	rpc := "GetCodeMigrations"
	ctx = logging.InitContext(ctx)
	ctx = logging.RpcContext(ctx, rpc)

	resp := &tfprotov6.GetCodeMigrationsResponse{
		CodeMigrations: []tfprotov6.CodeMigration{},
		Diagnostics:    []*tfprotov6.Diagnostic{},
	}

	for _, server := range s.servers {
		// TODO: Remove and call server.GetCodeMigrations below directly once interface becomes required.
		codeMigrationServer, ok := server.(tfprotov6.ProviderServerWithCodeMigrations)
		if !ok {
			// Since it's a prototype, we'll ignore any servers that don't define code migrations rather than raising an error
			continue
		}

		ctx = logging.Tfprotov6ProviderServerContext(ctx, server)
		logging.MuxTrace(ctx, "calling downstream server")

		serverResp, err := codeMigrationServer.GetCodeMigrations(ctx, req)
		if err != nil {
			return resp, fmt.Errorf("error calling GetCodeMigrations for %T: %w", server, err)
		}

		// Maybe in the future we'd do some de-duping or validation errors here, but for now we'll just trust the downstream servers
		resp.CodeMigrations = append(resp.CodeMigrations, serverResp.CodeMigrations...)
		resp.Diagnostics = append(resp.Diagnostics, serverResp.Diagnostics...)
	}

	return resp, nil
}
