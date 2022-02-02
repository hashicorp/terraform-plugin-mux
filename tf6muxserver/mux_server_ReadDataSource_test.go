package tf6muxserver_test

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-mux/internal/tf6testserver"
	"github.com/hashicorp/terraform-plugin-mux/tf6muxserver"
)

func TestMuxServerReadDataSource(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	servers := []func() tfprotov6.ProviderServer{
		(&tf6testserver.TestServer{
			DataSourceSchemas: map[string]*tfprotov6.Schema{
				"test_data_source_server1": {},
			},
		}).ProviderServer,
		(&tf6testserver.TestServer{
			DataSourceSchemas: map[string]*tfprotov6.Schema{
				"test_data_source_server2": {},
			},
		}).ProviderServer,
	}

	muxServer, err := tf6muxserver.NewMuxServer(ctx, servers...)

	if err != nil {
		t.Fatalf("unexpected error setting up factory: %s", err)
	}

	_, err = muxServer.ProviderServer().ReadDataSource(ctx, &tfprotov6.ReadDataSourceRequest{
		TypeName: "test_data_source_server1",
	})

	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	if !servers[0]().(*tf6testserver.TestServer).ReadDataSourceCalled["test_data_source_server1"] {
		t.Errorf("expected test_data_source_server1 ReadDataSource to be called on server1")
	}

	if servers[1]().(*tf6testserver.TestServer).ReadDataSourceCalled["test_data_source_server1"] {
		t.Errorf("unexpected test_data_source_server1 ReadDataSource called on server2")
	}

	_, err = muxServer.ProviderServer().ReadDataSource(ctx, &tfprotov6.ReadDataSourceRequest{
		TypeName: "test_data_source_server2",
	})

	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	if servers[0]().(*tf6testserver.TestServer).ReadDataSourceCalled["test_data_source_server2"] {
		t.Errorf("unexpected test_data_source_server2 ReadDataSource called on server1")
	}

	if !servers[1]().(*tf6testserver.TestServer).ReadDataSourceCalled["test_data_source_server2"] {
		t.Errorf("expected test_data_source_server2 ReadDataSource to be called on server2")
	}
}
