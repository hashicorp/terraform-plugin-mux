// Copyright IBM Corp. 2020, 2026
// SPDX-License-Identifier: MPL-2.0

package tf6muxserver_test

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"

	"github.com/hashicorp/terraform-plugin-mux/internal/tf6testserver"
	"github.com/hashicorp/terraform-plugin-mux/tf6muxserver"
)

func TestMuxServerConfigureProvider(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	testServers := [5]*tf6testserver.TestServer{
		{},
		{
			ConfigureProviderResponse: &tfprotov6.ConfigureProviderResponse{
				Diagnostics: []*tfprotov6.Diagnostic{
					{
						Severity: tfprotov6.DiagnosticSeverityWarning,
						Summary:  "warning summary",
						Detail:   "warning detail",
					},
				},
			},
		},
		{},
		{
			ConfigureProviderResponse: &tfprotov6.ConfigureProviderResponse{
				Diagnostics: []*tfprotov6.Diagnostic{
					{
						Severity: tfprotov6.DiagnosticSeverityError,
						Summary:  "error summary",
						Detail:   "error detail",
					},
				},
			},
		},
		{
			ConfigureProviderResponse: &tfprotov6.ConfigureProviderResponse{
				Diagnostics: []*tfprotov6.Diagnostic{
					{
						Severity: tfprotov6.DiagnosticSeverityError,
						Summary:  "unexpected error summary",
						Detail:   "unexpected error detail",
					},
				},
			},
		},
	}

	servers := []func() tfprotov6.ProviderServer{
		testServers[0].ProviderServer,
		testServers[1].ProviderServer,
		testServers[2].ProviderServer,
		testServers[3].ProviderServer,
		testServers[4].ProviderServer,
	}

	muxServer, err := tf6muxserver.NewMuxServer(ctx, servers...)

	if err != nil {
		t.Fatalf("error setting up muxer: %s", err)
	}

	resp, err := muxServer.ProviderServer().ConfigureProvider(ctx, &tfprotov6.ConfigureProviderRequest{})

	if err != nil {
		t.Fatalf("error calling ConfigureProvider: %s", err)
	}

	expectedResp := &tfprotov6.ConfigureProviderResponse{
		Diagnostics: []*tfprotov6.Diagnostic{
			{
				Severity: tfprotov6.DiagnosticSeverityWarning,
				Summary:  "warning summary",
				Detail:   "warning detail",
			},
			{
				Severity: tfprotov6.DiagnosticSeverityError,
				Summary:  "error summary",
				Detail:   "error detail",
			},
		},
	}

	if diff := cmp.Diff(resp, expectedResp); diff != "" {
		t.Errorf("unexpected difference: %s", diff)
	}

	for num, testServer := range testServers {
		if num < 4 && !testServer.ConfigureProviderCalled {
			t.Errorf("configure not called on server%d", num+1)
		}
	}
}
