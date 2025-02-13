// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package tf5muxserver

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-go/tfprotov5"

	"github.com/hashicorp/terraform-plugin-mux/internal/logging"
)

// GetResourceIdentitySchemas merges the schemas returned by the
// tfprotov5.ResourceIdentitySchema associated with muxServer into a single schema.
// Everything must be returned from only one server.
// Schemas must be identical between all servers.
func (s *muxServer) GetResourceIdentitySchemas(ctx context.Context, req *tfprotov5.GetResourceIdentitySchemasRequest) (*tfprotov5.GetResourceIdentitySchemasResponse, error) {
	rpc := "GetResourceIdentitySchemas"
	ctx = logging.InitContext(ctx)
	ctx = logging.RpcContext(ctx, rpc)

	s.serverDiscoveryMutex.Lock()
	defer s.serverDiscoveryMutex.Unlock()

	resp := &tfprotov5.GetResourceIdentitySchemasResponse{
		IdentitySchemas: map[string]*tfprotov5.ResourceIdentitySchema{},
		Diagnostics:     []*tfprotov5.Diagnostic{},
	}

	for _, server := range s.servers {
		ctx := logging.Tfprotov5ProviderServerContext(ctx, server)
		logging.MuxTrace(ctx, "calling downstream server")

		// TODO: Remove and call server.GetResourceIdentitySchemas below directly once interface becomes required.
		//nolint:staticcheck // Intentionally verifying interface implementation
		resourceIdentityServer, err := server.(tfprotov5.ProviderServerWithResourceIdentity).GetResourceIdentitySchemas(ctx, req)

		if err != nil {
			return resp, fmt.Errorf("error calling GetResourceIdentitySchemas for %T: %w", server, err)
		}

		resp.Diagnostics = append(resp.Diagnostics, resourceIdentityServer.Diagnostics...)

		for resourceIdentityType, schema := range resourceIdentityServer.IdentitySchemas {
			if _, ok := resp.IdentitySchemas[resourceIdentityType]; ok {
				resp.Diagnostics = append(resp.Diagnostics, resourceIdentityDuplicateError(resourceIdentityType))

				continue
			}

			resp.IdentitySchemas[resourceIdentityType] = schema
		}
	}

	s.serverDiscoveryComplete = true

	return resp, nil
}
