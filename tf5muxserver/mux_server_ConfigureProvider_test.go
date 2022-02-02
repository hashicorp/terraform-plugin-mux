package tf5muxserver_test

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-go/tfprotov5"
	"github.com/hashicorp/terraform-plugin-mux/internal/tf5testserver"
	"github.com/hashicorp/terraform-plugin-mux/tf5muxserver"
)

func TestMuxServerConfigureProvider(t *testing.T) {
	t.Parallel()

	servers := []func() tfprotov5.ProviderServer{
		(&tf5testserver.TestServer{}).ProviderServer,
		(&tf5testserver.TestServer{}).ProviderServer,
		(&tf5testserver.TestServer{}).ProviderServer,
		(&tf5testserver.TestServer{}).ProviderServer,
		(&tf5testserver.TestServer{}).ProviderServer,
	}

	muxServer, err := tf5muxserver.NewMuxServer(context.Background(), servers...)

	if err != nil {
		t.Fatalf("error setting up muxer: %s", err)
	}

	_, err = muxServer.ProviderServer().ConfigureProvider(context.Background(), &tfprotov5.ConfigureProviderRequest{})

	if err != nil {
		t.Fatalf("error calling ConfigureProvider: %s", err)
	}

	for num, server := range servers {
		if !server().(*tf5testserver.TestServer).ConfigureProviderCalled {
			t.Errorf("configure not called on server%d", num+1)
		}
	}
}
