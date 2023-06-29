// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package tf5muxserver

import (
	"context"

	"github.com/hashicorp/terraform-plugin-go/tfprotov5"
)

var _ tfprotov5.ProviderServer = &muxServer{}

// muxServer is a gRPC server implementation that stands in front of other
// gRPC servers, routing requests to them as if they were a single server. It
// should always be instantiated by calling NewMuxServer().
type muxServer struct {
	// Routing for data source types
	dataSources map[string]tfprotov5.ProviderServer

	// Provider schema is cached during GetProviderSchema for
	// ValidateProviderConfig equality checking.
	providerSchema *tfprotov5.Schema

	// Routing for resource types
	resources map[string]tfprotov5.ProviderServer

	// Resource capabilities are cached during GetProviderSchema
	resourceCapabilities map[string]*tfprotov5.ServerCapabilities

	// Underlying servers for requests that should be handled by all servers
	servers []tfprotov5.ProviderServer
}

// ProviderServer is a function compatible with tf6server.Serve.
func (s muxServer) ProviderServer() tfprotov5.ProviderServer {
	return &s
}

// NewMuxServer returns a muxed server that will route gRPC requests between
// tfprotov5.ProviderServers specified. The GetProviderSchema method of each
// is called to verify that the overall muxed server is compatible by ensuring:
//
//   - All provider schemas exactly match
//   - All provider meta schemas exactly match
//   - Only one provider implements each managed resource
//   - Only one provider implements each data source
func NewMuxServer(_ context.Context, servers ...func() tfprotov5.ProviderServer) (*muxServer, error) {
	result := muxServer{
		dataSources:          make(map[string]tfprotov5.ProviderServer),
		resources:            make(map[string]tfprotov5.ProviderServer),
		resourceCapabilities: make(map[string]*tfprotov5.ServerCapabilities),
	}

	for _, server := range servers {
		result.servers = append(result.servers, server())
	}

	return &result, nil
}
