package tfmux

import (
	"github.com/hashicorp/terraform-plugin-mux/internal/tfplugin5"
)

func ExampleNewResourceServer_v2v3() {
	// the ProviderServer from SDKv2
	var sdkv2 tfplugin5.ProviderServer

	// the ProviderServer from SDKv3
	var sdkv3 tfplugin5.ProviderServer

	// all RPC requests will be handled by SDKv2, except
	// for requests related to test_instance, test_resource,
	// and test_foo, which will be handled by SDKv3.
	_ = NewResourceServer(sdkv2, sdkv3, []string{
		"test_instance",
		"test_resource",
		"test_foo",
	})
	// use the result when instantiating the terraform-plugin-sdk.Plugin
	/*
		plugin.Serve(&plugin.ServeOpts{
			GRPCProviderFunc: plugin.GRPCProviderFunc(func() tfplugin5.ProviderServer {
				return server
			}),
		})
	*/
}
