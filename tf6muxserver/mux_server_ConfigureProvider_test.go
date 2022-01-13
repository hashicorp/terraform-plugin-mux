package tf6muxserver_test

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-mux/internal/tf6testserver"
	"github.com/hashicorp/terraform-plugin-mux/tf6muxserver"
)

func TestMuxServerConfigureProvider(t *testing.T) {
	t.Parallel()

	servers := []func() tfprotov6.ProviderServer{
		(&tf6testserver.TestServer{}).ProviderServer,
		(&tf6testserver.TestServer{}).ProviderServer,
		(&tf6testserver.TestServer{}).ProviderServer,
		(&tf6testserver.TestServer{}).ProviderServer,
		(&tf6testserver.TestServer{}).ProviderServer,
	}

	muxServer, err := tf6muxserver.NewMuxServer(context.Background(), servers...)

	if err != nil {
		t.Fatalf("error setting up muxer: %s", err)
	}

	_, err = muxServer.ProviderServer().ConfigureProvider(context.Background(), &tfprotov6.ConfigureProviderRequest{})

	if err != nil {
		t.Fatalf("error calling ConfigureProvider: %s", err)
	}

	for num, server := range servers {
		if !server().(*tf6testserver.TestServer).ConfigureProviderCalled {
			t.Errorf("configure not called on server%d", num+1)
		}
	}
}
