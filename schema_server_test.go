package tfmux

import (
	"context"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-go/tfprotov5"
	"github.com/hashicorp/terraform-plugin-go/tfprotov5/tftypes"
)

var _ tfprotov5.ProviderServer = &testServer{}

func testFactory(s *testServer) func() tfprotov5.ProviderServer {
	return func() tfprotov5.ProviderServer {
		return s
	}
}

type testServer struct {
	providerSchema     *tfprotov5.Schema
	providerMetaSchema *tfprotov5.Schema
	resourceSchemas    map[string]*tfprotov5.Schema
	dataSourceSchemas  map[string]*tfprotov5.Schema

	resourcesCalled   map[string]bool
	dataSourcesCalled map[string]bool
	configureCalled   bool
	stopCalled        bool
	stopError         string

	respondToPrepareProviderConfig bool
	errorOnPrepareProviderConfig   bool
	warnOnPrepareProviderConfig    bool
}

func (s *testServer) GetProviderSchema(_ context.Context, _ *tfprotov5.GetProviderSchemaRequest) (*tfprotov5.GetProviderSchemaResponse, error) {
	resources := s.resourceSchemas
	if resources == nil {
		resources = map[string]*tfprotov5.Schema{}
	}
	dataSources := s.dataSourceSchemas
	if dataSources == nil {
		dataSources = map[string]*tfprotov5.Schema{}
	}
	return &tfprotov5.GetProviderSchemaResponse{
		Provider:          s.providerSchema,
		ProviderMeta:      s.providerMetaSchema,
		ResourceSchemas:   resources,
		DataSourceSchemas: dataSources,
	}, nil
}

func (s *testServer) PrepareProviderConfig(_ context.Context, req *tfprotov5.PrepareProviderConfigRequest) (*tfprotov5.PrepareProviderConfigResponse, error) {
	if s.respondToPrepareProviderConfig {
		return &tfprotov5.PrepareProviderConfigResponse{
			PreparedConfig: req.Config,
		}, nil
	}
	if s.errorOnPrepareProviderConfig {
		return &tfprotov5.PrepareProviderConfigResponse{
			Diagnostics: []*tfprotov5.Diagnostic{
				{
					Severity: tfprotov5.DiagnosticSeverityError,
					Summary:  "hardcoded error for testing",
					Detail:   "testing that we only get errors from one thing",
				},
			},
		}, nil
	}
	if s.warnOnPrepareProviderConfig {
		return &tfprotov5.PrepareProviderConfigResponse{
			Diagnostics: []*tfprotov5.Diagnostic{
				{
					Severity: tfprotov5.DiagnosticSeverityWarning,
					Summary:  "hardcoded warning for testing",
					Detail:   "testing that we only get warnings from one thing",
				},
			},
		}, nil
	}
	return nil, nil
}

func (s *testServer) StopProvider(_ context.Context, _ *tfprotov5.StopProviderRequest) (*tfprotov5.StopProviderResponse, error) {
	s.stopCalled = true
	if s.stopError != "" {
		return &tfprotov5.StopProviderResponse{
			Error: s.stopError,
		}, nil
	}
	return &tfprotov5.StopProviderResponse{}, nil
}

func (s *testServer) ConfigureProvider(_ context.Context, _ *tfprotov5.ConfigureProviderRequest) (*tfprotov5.ConfigureProviderResponse, error) {
	s.configureCalled = true
	return &tfprotov5.ConfigureProviderResponse{}, nil
}

func (s *testServer) ValidateResourceTypeConfig(_ context.Context, req *tfprotov5.ValidateResourceTypeConfigRequest) (*tfprotov5.ValidateResourceTypeConfigResponse, error) {
	s.resourcesCalled[req.TypeName] = true
	return nil, nil
}

func (s *testServer) UpgradeResourceState(_ context.Context, req *tfprotov5.UpgradeResourceStateRequest) (*tfprotov5.UpgradeResourceStateResponse, error) {
	s.resourcesCalled[req.TypeName] = true
	return nil, nil
}

func (s *testServer) ReadResource(_ context.Context, req *tfprotov5.ReadResourceRequest) (*tfprotov5.ReadResourceResponse, error) {
	s.resourcesCalled[req.TypeName] = true
	return nil, nil
}

func (s *testServer) PlanResourceChange(_ context.Context, req *tfprotov5.PlanResourceChangeRequest) (*tfprotov5.PlanResourceChangeResponse, error) {
	s.resourcesCalled[req.TypeName] = true
	return nil, nil
}

func (s *testServer) ApplyResourceChange(_ context.Context, req *tfprotov5.ApplyResourceChangeRequest) (*tfprotov5.ApplyResourceChangeResponse, error) {
	s.resourcesCalled[req.TypeName] = true
	return nil, nil
}

func (s *testServer) ImportResourceState(_ context.Context, req *tfprotov5.ImportResourceStateRequest) (*tfprotov5.ImportResourceStateResponse, error) {
	s.resourcesCalled[req.TypeName] = true
	return nil, nil
}

func (s *testServer) ValidateDataSourceConfig(_ context.Context, req *tfprotov5.ValidateDataSourceConfigRequest) (*tfprotov5.ValidateDataSourceConfigResponse, error) {
	s.dataSourcesCalled[req.TypeName] = true
	return nil, nil
}

func (s *testServer) ReadDataSource(_ context.Context, req *tfprotov5.ReadDataSourceRequest) (*tfprotov5.ReadDataSourceResponse, error) {
	s.dataSourcesCalled[req.TypeName] = true
	return nil, nil
}

func TestSchemaServerGetProviderSchema_combined(t *testing.T) {
}

func TestSchemaServerGetProviderSchema_errorDuplicateResource(t *testing.T) {
}

func TestSchemaServerGetProviderSchema_errorDuplicateDataSource(t *testing.T) {
}

func TestSchemaServerGetProviderSchema_errorProviderMismatch(t *testing.T) {
}

func TestSchemaServerGetProviderSchema_errorProviderMetaMismatch(t *testing.T) {
}

func TestSchemaServerPrepareProviderConfig_errorMultipleResponses(t *testing.T) {
	config, err := tfprotov5.NewDynamicValue(tftypes.Object{
		AttributeTypes: map[string]tftypes.Type{
			"hello": tftypes.String,
		},
	}, tftypes.NewValue(tftypes.Object{
		AttributeTypes: map[string]tftypes.Type{
			"hello": tftypes.String,
		},
	}, map[string]tftypes.Value{
		"hello": tftypes.NewValue(tftypes.String, "world"),
	}))
	if err != nil {
		t.Fatalf("error constructing config: %s", err)
	}

	server1 := testFactory(&testServer{
		respondToPrepareProviderConfig: true,
	})
	server2 := testFactory(&testServer{})
	server3 := testFactory(&testServer{})
	factory, err := NewSchemaServerFactory(context.Background(), server1, server2, server3)
	if err != nil {
		t.Fatalf("error setting up muxer: %s", err)
	}
	_, err = factory.Server().PrepareProviderConfig(context.Background(), &tfprotov5.PrepareProviderConfigRequest{
		Config: &config,
	})
	if err != nil {
		t.Errorf("unexpected error when only one server replied to PrepareProviderConfig: %s", err)
	}

	server1 = testFactory(&testServer{})
	server2 = testFactory(&testServer{})
	server3 = testFactory(&testServer{})
	factory, err = NewSchemaServerFactory(context.Background(), server1, server2, server3)
	if err != nil {
		t.Fatalf("error setting up muxer: %s", err)
	}
	_, err = factory.Server().PrepareProviderConfig(context.Background(), &tfprotov5.PrepareProviderConfigRequest{
		Config: &config,
	})
	if err != nil {
		t.Errorf("unexpected error when no servers replied to PrepareProviderConfig: %s", err)
	}

	server1 = testFactory(&testServer{
		respondToPrepareProviderConfig: true,
	})
	server2 = testFactory(&testServer{})
	server3 = testFactory(&testServer{
		respondToPrepareProviderConfig: true,
	})
	factory, err = NewSchemaServerFactory(context.Background(), server1, server2, server3)
	if err != nil {
		t.Fatalf("error setting up muxer: %s", err)
	}
	_, err = factory.Server().PrepareProviderConfig(context.Background(), &tfprotov5.PrepareProviderConfigRequest{
		Config: &config,
	})
	if !strings.Contains(err.Error(), "got a PrepareProviderConfig response from multiple servers") {
		t.Errorf("expected error about multiple servers returning PrepareProviderConfigResponses, got %q", err)
	}

	server1 = testFactory(&testServer{
		respondToPrepareProviderConfig: true,
	})
	server2 = testFactory(&testServer{})
	server3 = testFactory(&testServer{
		errorOnPrepareProviderConfig: true,
	})
	factory, err = NewSchemaServerFactory(context.Background(), server1, server2, server3)
	if err != nil {
		t.Fatalf("error setting up muxer: %s", err)
	}
	_, err = factory.Server().PrepareProviderConfig(context.Background(), &tfprotov5.PrepareProviderConfigRequest{
		Config: &config,
	})
	if !strings.Contains(err.Error(), "got a PrepareProviderConfig response from multiple servers") {
		t.Errorf("expected error about multiple servers returning PrepareProviderConfigResponses, got %q", err)
	}

	server1 = testFactory(&testServer{
		respondToPrepareProviderConfig: true,
	})
	server2 = testFactory(&testServer{})
	server3 = testFactory(&testServer{
		warnOnPrepareProviderConfig: true,
	})
	factory, err = NewSchemaServerFactory(context.Background(), server1, server2, server3)
	if err != nil {
		t.Fatalf("error setting up muxer: %s", err)
	}
	_, err = factory.Server().PrepareProviderConfig(context.Background(), &tfprotov5.PrepareProviderConfigRequest{
		Config: &config,
	})
	if !strings.Contains(err.Error(), "got a PrepareProviderConfig response from multiple servers") {
		t.Errorf("expected error about multiple servers returning PrepareProviderConfigResponses, got %q", err)
	}

	server1 = testFactory(&testServer{
		errorOnPrepareProviderConfig: true,
	})
	server2 = testFactory(&testServer{})
	server3 = testFactory(&testServer{
		warnOnPrepareProviderConfig: true,
	})
	factory, err = NewSchemaServerFactory(context.Background(), server1, server2, server3)
	if err != nil {
		t.Fatalf("error setting up muxer: %s", err)
	}
	_, err = factory.Server().PrepareProviderConfig(context.Background(), &tfprotov5.PrepareProviderConfigRequest{
		Config: &config,
	})
	if !strings.Contains(err.Error(), "got a PrepareProviderConfig response from multiple servers") {
		t.Errorf("expected error about multiple servers returning PrepareProviderConfigResponses, got %q", err)
	}
}

func TestSchemaServerConfigureProvider_configuredEveryone(t *testing.T) {
	server1 := testFactory(&testServer{})
	server2 := testFactory(&testServer{})
	server3 := testFactory(&testServer{})
	server4 := testFactory(&testServer{})
	server5 := testFactory(&testServer{})
	factory, err := NewSchemaServerFactory(context.Background(), server1, server2, server3, server4, server5)
	if err != nil {
		t.Fatalf("error setting up muxer: %s", err)
	}
	_, err = factory.Server().ConfigureProvider(context.Background(), &tfprotov5.ConfigureProviderRequest{})
	if err != nil {
		t.Fatalf("error calling ConfigureProvider: %s", err)
	}
	for num, f := range []func() tfprotov5.ProviderServer{
		server1, server2, server3, server4, server5,
	} {
		if !f().(*testServer).configureCalled {
			t.Errorf("configure not called on server%d", num+1)
		}
	}
}

func TestSchemaServerStopProvider_stoppedEveryone(t *testing.T) {
	server1 := testFactory(&testServer{})
	server2 := testFactory(&testServer{
		stopError: "error in server2",
	})
	server3 := testFactory(&testServer{})
	server4 := testFactory(&testServer{
		stopError: "error in server4",
	})
	server5 := testFactory(&testServer{})
	factory, err := NewSchemaServerFactory(context.Background(), server1, server2, server3, server4, server5)
	if err != nil {
		t.Fatalf("error setting up muxer: %s", err)
	}
	_, err = factory.Server().StopProvider(context.Background(), &tfprotov5.StopProviderRequest{})
	if err != nil {
		t.Fatalf("error calling StopProvider: %s", err)
	}
	for num, f := range []func() tfprotov5.ProviderServer{
		server1, server2, server3, server4, server5,
	} {
		if !f().(*testServer).stopCalled {
			t.Errorf("stop not called on server%d", num+1)
		}
	}
}

func TestSchemaServer_resourceRouting(t *testing.T) {
}

func TestSchemaServer_dataSourceRouting(t *testing.T) {
}
