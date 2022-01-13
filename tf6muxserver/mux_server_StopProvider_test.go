package tf6muxserver_test

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-mux/internal/tf6testserver"
	"github.com/hashicorp/terraform-plugin-mux/tf6muxserver"
)

func TestMuxServerStopProvider(t *testing.T) {
	t.Parallel()

	servers := []func() tfprotov6.ProviderServer{
		(&tf6testserver.TestServer{}).ProviderServer,
		(&tf6testserver.TestServer{
			StopProviderError: "error in server2",
		}).ProviderServer,
		(&tf6testserver.TestServer{}).ProviderServer,
		(&tf6testserver.TestServer{
			StopProviderError: "error in server4",
		}).ProviderServer,
		(&tf6testserver.TestServer{}).ProviderServer,
	}

	muxServer, err := tf6muxserver.NewMuxServer(context.Background(), servers...)

	if err != nil {
		t.Fatalf("error setting up muxer: %s", err)
	}

	_, err = muxServer.ProviderServer().StopProvider(context.Background(), &tfprotov6.StopProviderRequest{})

	if err != nil {
		t.Fatalf("error calling StopProvider: %s", err)
	}

	for num, server := range servers {
		if !server().(*tf6testserver.TestServer).StopProviderCalled {
			t.Errorf("StopProvider not called on server%d", num+1)
		}
	}
}
