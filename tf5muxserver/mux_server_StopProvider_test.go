package tf5muxserver_test

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-go/tfprotov5"
	"github.com/hashicorp/terraform-plugin-mux/internal/tf5testserver"
	"github.com/hashicorp/terraform-plugin-mux/tf5muxserver"
)

func TestMuxServerStopProvider(t *testing.T) {
	t.Parallel()

	servers := []func() tfprotov5.ProviderServer{
		(&tf5testserver.TestServer{}).ProviderServer,
		(&tf5testserver.TestServer{
			StopProviderError: "error in server2",
		}).ProviderServer,
		(&tf5testserver.TestServer{}).ProviderServer,
		(&tf5testserver.TestServer{
			StopProviderError: "error in server4",
		}).ProviderServer,
		(&tf5testserver.TestServer{}).ProviderServer,
	}

	muxServer, err := tf5muxserver.NewMuxServer(context.Background(), servers...)

	if err != nil {
		t.Fatalf("error setting up muxer: %s", err)
	}

	_, err = muxServer.ProviderServer().StopProvider(context.Background(), &tfprotov5.StopProviderRequest{})

	if err != nil {
		t.Fatalf("error calling StopProvider: %s", err)
	}

	for num, server := range servers {
		if !server().(*tf5testserver.TestServer).StopProviderCalled {
			t.Errorf("StopProvider not called on server%d", num+1)
		}
	}
}
