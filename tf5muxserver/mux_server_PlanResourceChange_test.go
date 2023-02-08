package tf5muxserver_test

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-go/tfprotov5"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/hashicorp/terraform-plugin-mux/internal/tf5dynamicvalue"
	"github.com/hashicorp/terraform-plugin-mux/internal/tf5testserver"
	"github.com/hashicorp/terraform-plugin-mux/tf5muxserver"
)

func TestMuxServerPlanResourceChange_Routing(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	testServer1 := &tf5testserver.TestServer{
		ResourceSchemas: map[string]*tfprotov5.Schema{
			"test_resource_server1": {
				Block: &tfprotov5.SchemaBlock{
					Attributes: []*tfprotov5.SchemaAttribute{
						{
							Name: "test_string_attribute",
							Type: tftypes.String,
						},
					},
				},
			},
		},
	}
	testServer2 := &tf5testserver.TestServer{
		ResourceSchemas: map[string]*tfprotov5.Schema{
			"test_resource_server2": {
				Block: &tfprotov5.SchemaBlock{
					Attributes: []*tfprotov5.SchemaAttribute{
						{
							Name: "test_string_attribute",
							Type: tftypes.String,
						},
					},
				},
			},
		},
	}

	testProposedNewState := tf5dynamicvalue.Must(
		tftypes.Object{
			AttributeTypes: map[string]tftypes.Type{
				"test_string_attribute": tftypes.String,
			},
		},
		tftypes.NewValue(
			tftypes.Object{
				AttributeTypes: map[string]tftypes.Type{
					"test_string_attribute": tftypes.String,
				},
			},
			// intentionally set for create/update plan
			map[string]tftypes.Value{
				"test_string_attribute": tftypes.NewValue(tftypes.String, "test-value"),
			},
		),
	)

	servers := []func() tfprotov5.ProviderServer{testServer1.ProviderServer, testServer2.ProviderServer}
	muxServer, err := tf5muxserver.NewMuxServer(ctx, servers...)

	if err != nil {
		t.Fatalf("unexpected error setting up factory: %s", err)
	}

	_, err = muxServer.ProviderServer().PlanResourceChange(ctx, &tfprotov5.PlanResourceChangeRequest{
		ProposedNewState: testProposedNewState,
		TypeName:         "test_resource_server1",
	})

	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	if !testServer1.PlanResourceChangeCalled["test_resource_server1"] {
		t.Errorf("expected test_resource_server1 PlanResourceChange to be called on server1")
	}

	if testServer2.PlanResourceChangeCalled["test_resource_server1"] {
		t.Errorf("unexpected test_resource_server1 PlanResourceChange called on server2")
	}

	_, err = muxServer.ProviderServer().PlanResourceChange(ctx, &tfprotov5.PlanResourceChangeRequest{
		ProposedNewState: testProposedNewState,
		TypeName:         "test_resource_server2",
	})

	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	if testServer1.PlanResourceChangeCalled["test_resource_server2"] {
		t.Errorf("unexpected test_resource_server2 PlanResourceChange called on server1")
	}

	if !testServer2.PlanResourceChangeCalled["test_resource_server2"] {
		t.Errorf("expected test_resource_server2 PlanResourceChange to be called on server2")
	}
}

func TestMuxServerPlanResourceChange_ServerCapabilities_PlanDestroy(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	testServer1 := &tf5testserver.TestServer{
		ResourceSchemas: map[string]*tfprotov5.Schema{
			"test_resource_server1": {
				Block: &tfprotov5.SchemaBlock{
					Attributes: []*tfprotov5.SchemaAttribute{
						{
							Name: "test_string_attribute",
							Type: tftypes.String,
						},
					},
				},
			},
		},
		ServerCapabilities: &tfprotov5.ServerCapabilities{
			PlanDestroy: true,
		},
	}
	testServer2 := &tf5testserver.TestServer{
		ResourceSchemas: map[string]*tfprotov5.Schema{
			"test_resource_server2": {
				Block: &tfprotov5.SchemaBlock{
					Attributes: []*tfprotov5.SchemaAttribute{
						{
							Name: "test_string_attribute",
							Type: tftypes.String,
						},
					},
				},
			},
		},
		// Intentionally no ServerCapabilities on this server
	}

	testProposedNewState := tf5dynamicvalue.Must(
		tftypes.Object{
			AttributeTypes: map[string]tftypes.Type{
				"test_string_attribute": tftypes.String,
			},
		},
		tftypes.NewValue(
			tftypes.Object{
				AttributeTypes: map[string]tftypes.Type{
					"test_string_attribute": tftypes.String,
				},
			},
			nil, // intentionally null for destroy plan
		),
	)

	servers := []func() tfprotov5.ProviderServer{testServer1.ProviderServer, testServer2.ProviderServer}
	muxServer, err := tf5muxserver.NewMuxServer(ctx, servers...)

	if err != nil {
		t.Fatalf("unexpected error setting up factory: %s", err)
	}

	_, err = muxServer.ProviderServer().PlanResourceChange(ctx, &tfprotov5.PlanResourceChangeRequest{
		ProposedNewState: testProposedNewState,
		TypeName:         "test_resource_server1",
	})

	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	if !testServer1.PlanResourceChangeCalled["test_resource_server1"] {
		t.Errorf("expected test_resource_server1 PlanResourceChange to be called on server1")
	}

	if testServer2.PlanResourceChangeCalled["test_resource_server1"] {
		t.Errorf("unexpected test_resource_server1 PlanResourceChange called on server2")
	}

	_, err = muxServer.ProviderServer().PlanResourceChange(ctx, &tfprotov5.PlanResourceChangeRequest{
		ProposedNewState: testProposedNewState,
		TypeName:         "test_resource_server2",
	})

	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	if testServer1.PlanResourceChangeCalled["test_resource_server2"] {
		t.Errorf("unexpected test_resource_server2 PlanResourceChange called on server1")
	}

	// Server does not enable ServerCapabilities.PlanDestroy
	if testServer2.PlanResourceChangeCalled["test_resource_server2"] {
		t.Errorf("unexpected test_resource_server2 PlanResourceChange called on server2")
	}
}
