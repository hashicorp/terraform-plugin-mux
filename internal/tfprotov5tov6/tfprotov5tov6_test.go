// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package tfprotov5tov6_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/terraform-plugin-go/tfprotov5"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/hashicorp/terraform-plugin-mux/internal/tfprotov5tov6"
)

var (
	testBytes []byte = []byte("test")

	testTfprotov5Diagnostics []*tfprotov5.Diagnostic = []*tfprotov5.Diagnostic{
		{
			Detail:   "test detail",
			Severity: tfprotov5.DiagnosticSeverityError,
			Summary:  "test summary",
		},
	}
	testTfprotov6Diagnostics []*tfprotov6.Diagnostic = []*tfprotov6.Diagnostic{
		{
			Detail:   "test detail",
			Severity: tfprotov6.DiagnosticSeverityError,
			Summary:  "test summary",
		},
	}

	testTfprotov5DynamicValue tfprotov5.DynamicValue
	testTfprotov6DynamicValue tfprotov6.DynamicValue

	testTfprotov5Schema *tfprotov5.Schema = &tfprotov5.Schema{
		Block: &tfprotov5.SchemaBlock{
			Attributes: []*tfprotov5.SchemaAttribute{
				{
					Name:     "test",
					Required: true,
				},
			},
			Version: 1,
		},
		Version: 1,
	}
	testTfprotov6Schema *tfprotov6.Schema = &tfprotov6.Schema{
		Block: &tfprotov6.SchemaBlock{
			Attributes: []*tfprotov6.SchemaAttribute{
				{
					Name:     "test",
					Required: true,
				},
			},
			Version: 1,
		},
		Version: 1,
	}
)

func init() {
	testTfprotov5DynamicValue, _ = tfprotov5.NewDynamicValue(tftypes.String, tftypes.NewValue(tftypes.String, "test"))
	testTfprotov6DynamicValue, _ = tfprotov6.NewDynamicValue(tftypes.String, tftypes.NewValue(tftypes.String, "test"))
}

func TestApplyResourceChangeRequest(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		in       *tfprotov5.ApplyResourceChangeRequest
		expected *tfprotov6.ApplyResourceChangeRequest
	}{
		"nil": {
			in:       nil,
			expected: nil,
		},
		"all-valid-fields": {
			in: &tfprotov5.ApplyResourceChangeRequest{
				Config:         &testTfprotov5DynamicValue,
				PlannedPrivate: testBytes,
				PlannedState:   &testTfprotov5DynamicValue,
				PriorState:     &testTfprotov5DynamicValue,
				ProviderMeta:   &testTfprotov5DynamicValue,
				TypeName:       "test",
			},
			expected: &tfprotov6.ApplyResourceChangeRequest{
				Config:         &testTfprotov6DynamicValue,
				PlannedPrivate: testBytes,
				PlannedState:   &testTfprotov6DynamicValue,
				PriorState:     &testTfprotov6DynamicValue,
				ProviderMeta:   &testTfprotov6DynamicValue,
				TypeName:       "test",
			},
		},
	}

	for name, testCase := range testCases {
		name, testCase := name, testCase

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got := tfprotov5tov6.ApplyResourceChangeRequest(testCase.in)

			if diff := cmp.Diff(got, testCase.expected); diff != "" {
				t.Errorf("unexpected difference: %s", diff)
			}
		})
	}
}

func TestApplyResourceChangeResponse(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		in       *tfprotov5.ApplyResourceChangeResponse
		expected *tfprotov6.ApplyResourceChangeResponse
	}{
		"nil": {
			in:       nil,
			expected: nil,
		},
		"all-valid-fields": {
			in: &tfprotov5.ApplyResourceChangeResponse{
				Diagnostics:                 testTfprotov5Diagnostics,
				NewState:                    &testTfprotov5DynamicValue,
				Private:                     testBytes,
				UnsafeToUseLegacyTypeSystem: true,
			},
			expected: &tfprotov6.ApplyResourceChangeResponse{
				Diagnostics:                 testTfprotov6Diagnostics,
				NewState:                    &testTfprotov6DynamicValue,
				Private:                     testBytes,
				UnsafeToUseLegacyTypeSystem: true,
			},
		},
	}

	for name, testCase := range testCases {
		name, testCase := name, testCase

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got := tfprotov5tov6.ApplyResourceChangeResponse(testCase.in)

			if diff := cmp.Diff(got, testCase.expected); diff != "" {
				t.Errorf("unexpected difference: %s", diff)
			}
		})
	}
}

func TestConfigureProviderRequest(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		in       *tfprotov5.ConfigureProviderRequest
		expected *tfprotov6.ConfigureProviderRequest
	}{
		"nil": {
			in:       nil,
			expected: nil,
		},
		"all-valid-fields": {
			in: &tfprotov5.ConfigureProviderRequest{
				Config:           &testTfprotov5DynamicValue,
				TerraformVersion: "1.2.3",
			},
			expected: &tfprotov6.ConfigureProviderRequest{
				Config:           &testTfprotov6DynamicValue,
				TerraformVersion: "1.2.3",
			},
		},
	}

	for name, testCase := range testCases {
		name, testCase := name, testCase

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got := tfprotov5tov6.ConfigureProviderRequest(testCase.in)

			if diff := cmp.Diff(got, testCase.expected); diff != "" {
				t.Errorf("unexpected difference: %s", diff)
			}
		})
	}
}

func TestConfigureProviderResponse(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		in       *tfprotov5.ConfigureProviderResponse
		expected *tfprotov6.ConfigureProviderResponse
	}{
		"nil": {
			in:       nil,
			expected: nil,
		},
		"all-valid-fields": {
			in: &tfprotov5.ConfigureProviderResponse{
				Diagnostics: testTfprotov5Diagnostics,
			},
			expected: &tfprotov6.ConfigureProviderResponse{
				Diagnostics: testTfprotov6Diagnostics,
			},
		},
	}

	for name, testCase := range testCases {
		name, testCase := name, testCase

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got := tfprotov5tov6.ConfigureProviderResponse(testCase.in)

			if diff := cmp.Diff(got, testCase.expected); diff != "" {
				t.Errorf("unexpected difference: %s", diff)
			}
		})
	}
}

func TestDiagnostics(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		in       []*tfprotov5.Diagnostic
		expected []*tfprotov6.Diagnostic
	}{
		"nil": {
			in:       nil,
			expected: nil,
		},
		"empty": {
			in:       []*tfprotov5.Diagnostic{},
			expected: []*tfprotov6.Diagnostic{},
		},
		"one": {
			in:       testTfprotov5Diagnostics,
			expected: testTfprotov6Diagnostics,
		},
		"multiple": {
			in: []*tfprotov5.Diagnostic{
				{
					Detail:   "test detail 1",
					Severity: tfprotov5.DiagnosticSeverityError,
					Summary:  "test summary 2",
				},
				{
					Detail:   "test detail 1",
					Severity: tfprotov5.DiagnosticSeverityError,
					Summary:  "test summary 2",
				},
			},
			expected: []*tfprotov6.Diagnostic{
				{
					Detail:   "test detail 1",
					Severity: tfprotov6.DiagnosticSeverityError,
					Summary:  "test summary 2",
				},
				{
					Detail:   "test detail 1",
					Severity: tfprotov6.DiagnosticSeverityError,
					Summary:  "test summary 2",
				},
			},
		},
	}

	for name, testCase := range testCases {
		name, testCase := name, testCase

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got := tfprotov5tov6.Diagnostics(testCase.in)

			if diff := cmp.Diff(got, testCase.expected); diff != "" {
				t.Errorf("unexpected difference: %s", diff)
			}
		})
	}
}

func TestDynamicValue(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		in       *tfprotov5.DynamicValue
		expected *tfprotov6.DynamicValue
	}{
		"nil": {
			in:       nil,
			expected: nil,
		},
		"all-valid-fields": {
			in: &tfprotov5.DynamicValue{
				JSON:    testBytes,
				MsgPack: testBytes,
			},
			expected: &tfprotov6.DynamicValue{
				JSON:    testBytes,
				MsgPack: testBytes,
			},
		},
	}

	for name, testCase := range testCases {
		name, testCase := name, testCase

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got := tfprotov5tov6.DynamicValue(testCase.in)

			if diff := cmp.Diff(got, testCase.expected); diff != "" {
				t.Errorf("unexpected difference: %s", diff)
			}
		})
	}
}

func TestGetProviderSchemaRequest(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		in       *tfprotov5.GetProviderSchemaRequest
		expected *tfprotov6.GetProviderSchemaRequest
	}{
		"nil": {
			in:       nil,
			expected: nil,
		},
		"all-valid-fields": {
			in:       &tfprotov5.GetProviderSchemaRequest{},
			expected: &tfprotov6.GetProviderSchemaRequest{},
		},
	}

	for name, testCase := range testCases {
		name, testCase := name, testCase

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got := tfprotov5tov6.GetProviderSchemaRequest(testCase.in)

			if diff := cmp.Diff(got, testCase.expected); diff != "" {
				t.Errorf("unexpected difference: %s", diff)
			}
		})
	}
}

func TestGetProviderSchemaResponse(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		in       *tfprotov5.GetProviderSchemaResponse
		expected *tfprotov6.GetProviderSchemaResponse
	}{
		"nil": {
			in:       nil,
			expected: nil,
		},
		"all-valid-fields": {
			in: &tfprotov5.GetProviderSchemaResponse{
				DataSourceSchemas: map[string]*tfprotov5.Schema{
					"test_data_source": testTfprotov5Schema,
				},
				Diagnostics:  testTfprotov5Diagnostics,
				Provider:     testTfprotov5Schema,
				ProviderMeta: testTfprotov5Schema,
				ResourceSchemas: map[string]*tfprotov5.Schema{
					"test_resource": testTfprotov5Schema,
				},
			},
			expected: &tfprotov6.GetProviderSchemaResponse{
				DataSourceSchemas: map[string]*tfprotov6.Schema{
					"test_data_source": testTfprotov6Schema,
				},
				Diagnostics:  testTfprotov6Diagnostics,
				Provider:     testTfprotov6Schema,
				ProviderMeta: testTfprotov6Schema,
				ResourceSchemas: map[string]*tfprotov6.Schema{
					"test_resource": testTfprotov6Schema,
				},
			},
		},
	}

	for name, testCase := range testCases {
		name, testCase := name, testCase

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got := tfprotov5tov6.GetProviderSchemaResponse(testCase.in)

			if diff := cmp.Diff(got, testCase.expected); diff != "" {
				t.Errorf("unexpected difference: %s", diff)
			}
		})
	}
}

func TestImportResourceStateRequest(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		in       *tfprotov5.ImportResourceStateRequest
		expected *tfprotov6.ImportResourceStateRequest
	}{
		"nil": {
			in:       nil,
			expected: nil,
		},
		"all-valid-fields": {
			in: &tfprotov5.ImportResourceStateRequest{
				ID:       "test-id",
				TypeName: "test_resource",
			},
			expected: &tfprotov6.ImportResourceStateRequest{
				ID:       "test-id",
				TypeName: "test_resource",
			},
		},
	}

	for name, testCase := range testCases {
		name, testCase := name, testCase

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got := tfprotov5tov6.ImportResourceStateRequest(testCase.in)

			if diff := cmp.Diff(got, testCase.expected); diff != "" {
				t.Errorf("unexpected difference: %s", diff)
			}
		})
	}
}

func TestImportResourceStateResponse(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		in       *tfprotov5.ImportResourceStateResponse
		expected *tfprotov6.ImportResourceStateResponse
	}{
		"nil": {
			in:       nil,
			expected: nil,
		},
		"all-valid-fields": {
			in: &tfprotov5.ImportResourceStateResponse{
				Diagnostics: testTfprotov5Diagnostics,
				ImportedResources: []*tfprotov5.ImportedResource{
					{
						Private:  testBytes,
						State:    &testTfprotov5DynamicValue,
						TypeName: "test_resource1",
					},
				},
			},
			expected: &tfprotov6.ImportResourceStateResponse{
				Diagnostics: testTfprotov6Diagnostics,
				ImportedResources: []*tfprotov6.ImportedResource{
					{
						Private:  testBytes,
						State:    &testTfprotov6DynamicValue,
						TypeName: "test_resource1",
					},
				},
			},
		},
	}

	for name, testCase := range testCases {
		name, testCase := name, testCase

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got := tfprotov5tov6.ImportResourceStateResponse(testCase.in)

			if diff := cmp.Diff(got, testCase.expected); diff != "" {
				t.Errorf("unexpected difference: %s", diff)
			}
		})
	}
}

func TestImportedResources(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		in       []*tfprotov5.ImportedResource
		expected []*tfprotov6.ImportedResource
	}{
		"nil": {
			in:       nil,
			expected: nil,
		},
		"empty": {
			in:       []*tfprotov5.ImportedResource{},
			expected: []*tfprotov6.ImportedResource{},
		},
		"one": {
			in: []*tfprotov5.ImportedResource{
				{
					Private:  testBytes,
					State:    &testTfprotov5DynamicValue,
					TypeName: "test_resource1",
				},
			},
			expected: []*tfprotov6.ImportedResource{
				{
					Private:  testBytes,
					State:    &testTfprotov6DynamicValue,
					TypeName: "test_resource1",
				},
			},
		},

		"multiple": {
			in: []*tfprotov5.ImportedResource{
				{
					Private:  testBytes,
					State:    &testTfprotov5DynamicValue,
					TypeName: "test_resource1",
				},
				{
					Private:  testBytes,
					State:    &testTfprotov5DynamicValue,
					TypeName: "test_resource2",
				},
			},
			expected: []*tfprotov6.ImportedResource{
				{
					Private:  testBytes,
					State:    &testTfprotov6DynamicValue,
					TypeName: "test_resource1",
				},
				{
					Private:  testBytes,
					State:    &testTfprotov6DynamicValue,
					TypeName: "test_resource2",
				},
			},
		},
	}

	for name, testCase := range testCases {
		name, testCase := name, testCase

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got := tfprotov5tov6.ImportedResources(testCase.in)

			if diff := cmp.Diff(got, testCase.expected); diff != "" {
				t.Errorf("unexpected difference: %s", diff)
			}
		})
	}
}

func TestPlanResourceChangeRequest(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		in       *tfprotov5.PlanResourceChangeRequest
		expected *tfprotov6.PlanResourceChangeRequest
	}{
		"nil": {
			in:       nil,
			expected: nil,
		},
		"all-valid-fields": {
			in: &tfprotov5.PlanResourceChangeRequest{
				Config:           &testTfprotov5DynamicValue,
				PriorPrivate:     testBytes,
				PriorState:       &testTfprotov5DynamicValue,
				ProposedNewState: &testTfprotov5DynamicValue,
				ProviderMeta:     &testTfprotov5DynamicValue,
				TypeName:         "test_resource",
			},
			expected: &tfprotov6.PlanResourceChangeRequest{
				Config:           &testTfprotov6DynamicValue,
				PriorPrivate:     testBytes,
				PriorState:       &testTfprotov6DynamicValue,
				ProposedNewState: &testTfprotov6DynamicValue,
				ProviderMeta:     &testTfprotov6DynamicValue,
				TypeName:         "test_resource",
			},
		},
	}

	for name, testCase := range testCases {
		name, testCase := name, testCase

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got := tfprotov5tov6.PlanResourceChangeRequest(testCase.in)

			if diff := cmp.Diff(got, testCase.expected); diff != "" {
				t.Errorf("unexpected difference: %s", diff)
			}
		})
	}
}

func TestPlanResourceChangeResponse(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		in       *tfprotov5.PlanResourceChangeResponse
		expected *tfprotov6.PlanResourceChangeResponse
	}{
		"nil": {
			in:       nil,
			expected: nil,
		},
		"all-valid-fields": {
			in: &tfprotov5.PlanResourceChangeResponse{
				Diagnostics:    testTfprotov5Diagnostics,
				PlannedPrivate: testBytes,
				PlannedState:   &testTfprotov5DynamicValue,
				RequiresReplace: []*tftypes.AttributePath{
					tftypes.NewAttributePath().WithAttributeName("test"),
				},
				UnsafeToUseLegacyTypeSystem: true,
			},
			expected: &tfprotov6.PlanResourceChangeResponse{
				Diagnostics:    testTfprotov6Diagnostics,
				PlannedPrivate: testBytes,
				PlannedState:   &testTfprotov6DynamicValue,
				RequiresReplace: []*tftypes.AttributePath{
					tftypes.NewAttributePath().WithAttributeName("test"),
				},
				UnsafeToUseLegacyTypeSystem: true,
			},
		},
	}

	for name, testCase := range testCases {
		name, testCase := name, testCase

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got := tfprotov5tov6.PlanResourceChangeResponse(testCase.in)

			if diff := cmp.Diff(got, testCase.expected); diff != "" {
				t.Errorf("unexpected difference: %s", diff)
			}
		})
	}
}

func TestRawState(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		in       *tfprotov5.RawState
		expected *tfprotov6.RawState
	}{
		"nil": {
			in:       nil,
			expected: nil,
		},
		"all-valid-fields": {
			in: &tfprotov5.RawState{
				Flatmap: map[string]string{"testkey": "testvalue"},
				JSON:    testBytes,
			},
			expected: &tfprotov6.RawState{
				Flatmap: map[string]string{"testkey": "testvalue"},
				JSON:    testBytes,
			},
		},
	}

	for name, testCase := range testCases {
		name, testCase := name, testCase

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got := tfprotov5tov6.RawState(testCase.in)

			if diff := cmp.Diff(got, testCase.expected); diff != "" {
				t.Errorf("unexpected difference: %s", diff)
			}
		})
	}
}

func TestReadDataSourceRequest(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		in       *tfprotov5.ReadDataSourceRequest
		expected *tfprotov6.ReadDataSourceRequest
	}{
		"nil": {
			in:       nil,
			expected: nil,
		},
		"all-valid-fields": {
			in: &tfprotov5.ReadDataSourceRequest{
				Config:       &testTfprotov5DynamicValue,
				ProviderMeta: &testTfprotov5DynamicValue,
				TypeName:     "test_data_source",
			},
			expected: &tfprotov6.ReadDataSourceRequest{
				Config:       &testTfprotov6DynamicValue,
				ProviderMeta: &testTfprotov6DynamicValue,
				TypeName:     "test_data_source",
			},
		},
	}

	for name, testCase := range testCases {
		name, testCase := name, testCase

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got := tfprotov5tov6.ReadDataSourceRequest(testCase.in)

			if diff := cmp.Diff(got, testCase.expected); diff != "" {
				t.Errorf("unexpected difference: %s", diff)
			}
		})
	}
}

func TestReadDataSourceResponse(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		in       *tfprotov5.ReadDataSourceResponse
		expected *tfprotov6.ReadDataSourceResponse
	}{
		"nil": {
			in:       nil,
			expected: nil,
		},
		"all-valid-fields": {
			in: &tfprotov5.ReadDataSourceResponse{
				Diagnostics: testTfprotov5Diagnostics,
				State:       &testTfprotov5DynamicValue,
			},
			expected: &tfprotov6.ReadDataSourceResponse{
				Diagnostics: testTfprotov6Diagnostics,
				State:       &testTfprotov6DynamicValue,
			},
		},
	}

	for name, testCase := range testCases {
		name, testCase := name, testCase

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got := tfprotov5tov6.ReadDataSourceResponse(testCase.in)

			if diff := cmp.Diff(got, testCase.expected); diff != "" {
				t.Errorf("unexpected difference: %s", diff)
			}
		})
	}
}

func TestReadResourceRequest(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		in       *tfprotov5.ReadResourceRequest
		expected *tfprotov6.ReadResourceRequest
	}{
		"nil": {
			in:       nil,
			expected: nil,
		},
		"all-valid-fields": {
			in: &tfprotov5.ReadResourceRequest{
				CurrentState: &testTfprotov5DynamicValue,
				Private:      testBytes,
				ProviderMeta: &testTfprotov5DynamicValue,
				TypeName:     "test_resource",
			},
			expected: &tfprotov6.ReadResourceRequest{
				CurrentState: &testTfprotov6DynamicValue,
				Private:      testBytes,
				ProviderMeta: &testTfprotov6DynamicValue,
				TypeName:     "test_resource",
			},
		},
	}

	for name, testCase := range testCases {
		name, testCase := name, testCase

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got := tfprotov5tov6.ReadResourceRequest(testCase.in)

			if diff := cmp.Diff(got, testCase.expected); diff != "" {
				t.Errorf("unexpected difference: %s", diff)
			}
		})
	}
}

func TestReadResourceResponse(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		in       *tfprotov5.ReadResourceResponse
		expected *tfprotov6.ReadResourceResponse
	}{
		"nil": {
			in:       nil,
			expected: nil,
		},
		"all-valid-fields": {
			in: &tfprotov5.ReadResourceResponse{
				Diagnostics: testTfprotov5Diagnostics,
				NewState:    &testTfprotov5DynamicValue,
				Private:     testBytes,
			},
			expected: &tfprotov6.ReadResourceResponse{
				Diagnostics: testTfprotov6Diagnostics,
				NewState:    &testTfprotov6DynamicValue,
				Private:     testBytes,
			},
		},
	}

	for name, testCase := range testCases {
		name, testCase := name, testCase

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got := tfprotov5tov6.ReadResourceResponse(testCase.in)

			if diff := cmp.Diff(got, testCase.expected); diff != "" {
				t.Errorf("unexpected difference: %s", diff)
			}
		})
	}
}

func TestSchema(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		in       *tfprotov5.Schema
		expected *tfprotov6.Schema
	}{
		"nil": {
			in:       nil,
			expected: nil,
		},
		"all-valid-fields": {
			in:       testTfprotov5Schema,
			expected: testTfprotov6Schema,
		},
	}

	for name, testCase := range testCases {
		name, testCase := name, testCase

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got := tfprotov5tov6.Schema(testCase.in)

			if diff := cmp.Diff(got, testCase.expected); diff != "" {
				t.Errorf("unexpected difference: %s", diff)
			}
		})
	}
}

func TestSchemaAttribute(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		in       *tfprotov5.SchemaAttribute
		expected *tfprotov6.SchemaAttribute
	}{
		"nil": {
			in:       nil,
			expected: nil,
		},
		"all-valid-fields": {
			in: &tfprotov5.SchemaAttribute{
				Computed:        true,
				Deprecated:      true,
				Description:     "test description",
				DescriptionKind: tfprotov5.StringKindPlain,
				Name:            "test",
				Optional:        true,
				Required:        true,
				Sensitive:       true,
				Type:            tftypes.String,
			},
			expected: &tfprotov6.SchemaAttribute{
				Computed:        true,
				Deprecated:      true,
				Description:     "test description",
				DescriptionKind: tfprotov6.StringKindPlain,
				Name:            "test",
				Optional:        true,
				Required:        true,
				Sensitive:       true,
				Type:            tftypes.String,
			},
		},
	}

	for name, testCase := range testCases {
		name, testCase := name, testCase

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got := tfprotov5tov6.SchemaAttribute(testCase.in)

			if diff := cmp.Diff(got, testCase.expected); diff != "" {
				t.Errorf("unexpected difference: %s", diff)
			}
		})
	}
}

func TestSchemaBlock(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		in       *tfprotov5.SchemaBlock
		expected *tfprotov6.SchemaBlock
	}{
		"nil": {
			in:       nil,
			expected: nil,
		},
		"all-valid-fields": {
			in: &tfprotov5.SchemaBlock{
				Attributes: []*tfprotov5.SchemaAttribute{
					{
						Name:     "test_attribute",
						Required: true,
					},
				},
				BlockTypes: []*tfprotov5.SchemaNestedBlock{
					{
						Block: &tfprotov5.SchemaBlock{
							Attributes: []*tfprotov5.SchemaAttribute{
								{
									Name:     "test_attribute",
									Required: true,
								},
							},
						},
						TypeName: "test_block",
					},
				},
				Deprecated:      true,
				Description:     "test description",
				DescriptionKind: tfprotov5.StringKindPlain,
				Version:         1,
			},
			expected: &tfprotov6.SchemaBlock{
				Attributes: []*tfprotov6.SchemaAttribute{
					{
						Name:     "test_attribute",
						Required: true,
					},
				},
				BlockTypes: []*tfprotov6.SchemaNestedBlock{
					{
						Block: &tfprotov6.SchemaBlock{
							Attributes: []*tfprotov6.SchemaAttribute{
								{
									Name:     "test_attribute",
									Required: true,
								},
							},
						},
						TypeName: "test_block",
					},
				},
				Deprecated:      true,
				Description:     "test description",
				DescriptionKind: tfprotov6.StringKindPlain,
				Version:         1,
			},
		},
	}

	for name, testCase := range testCases {
		name, testCase := name, testCase

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got := tfprotov5tov6.SchemaBlock(testCase.in)

			if diff := cmp.Diff(got, testCase.expected); diff != "" {
				t.Errorf("unexpected difference: %s", diff)
			}
		})
	}
}

func TestSchemaNestedBlock(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		in       *tfprotov5.SchemaNestedBlock
		expected *tfprotov6.SchemaNestedBlock
	}{
		"nil": {
			in:       nil,
			expected: nil,
		},
		"all-valid-fields": {
			in: &tfprotov5.SchemaNestedBlock{
				Block: &tfprotov5.SchemaBlock{
					Attributes: []*tfprotov5.SchemaAttribute{
						{
							Name:     "test_attribute",
							Required: true,
						},
					},
				},
				MaxItems: 1,
				MinItems: 1,
				Nesting:  tfprotov5.SchemaNestedBlockNestingModeList,
				TypeName: "test_block",
			},
			expected: &tfprotov6.SchemaNestedBlock{
				Block: &tfprotov6.SchemaBlock{
					Attributes: []*tfprotov6.SchemaAttribute{
						{
							Name:     "test_attribute",
							Required: true,
						},
					},
				},
				MaxItems: 1,
				MinItems: 1,
				Nesting:  tfprotov6.SchemaNestedBlockNestingModeList,
				TypeName: "test_block",
			},
		},
	}

	for name, testCase := range testCases {
		name, testCase := name, testCase

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got := tfprotov5tov6.SchemaNestedBlock(testCase.in)

			if diff := cmp.Diff(got, testCase.expected); diff != "" {
				t.Errorf("unexpected difference: %s", diff)
			}
		})
	}
}

func TestStopProviderRequest(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		in       *tfprotov5.StopProviderRequest
		expected *tfprotov6.StopProviderRequest
	}{
		"nil": {
			in:       nil,
			expected: nil,
		},
		"all-valid-fields": {
			in:       &tfprotov5.StopProviderRequest{},
			expected: &tfprotov6.StopProviderRequest{},
		},
	}

	for name, testCase := range testCases {
		name, testCase := name, testCase

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got := tfprotov5tov6.StopProviderRequest(testCase.in)

			if diff := cmp.Diff(got, testCase.expected); diff != "" {
				t.Errorf("unexpected difference: %s", diff)
			}
		})
	}
}

func TestStopProviderResponse(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		in       *tfprotov5.StopProviderResponse
		expected *tfprotov6.StopProviderResponse
	}{
		"nil": {
			in:       nil,
			expected: nil,
		},
		"all-valid-fields": {
			in: &tfprotov5.StopProviderResponse{
				Error: "test error",
			},
			expected: &tfprotov6.StopProviderResponse{
				Error: "test error",
			},
		},
	}

	for name, testCase := range testCases {
		name, testCase := name, testCase

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got := tfprotov5tov6.StopProviderResponse(testCase.in)

			if diff := cmp.Diff(got, testCase.expected); diff != "" {
				t.Errorf("unexpected difference: %s", diff)
			}
		})
	}
}

func TestStringKind(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		in       tfprotov5.StringKind
		expected tfprotov6.StringKind
	}{
		"markdown": {
			in:       tfprotov5.StringKindMarkdown,
			expected: tfprotov6.StringKindMarkdown,
		},
		"plain": {
			in:       tfprotov5.StringKindPlain,
			expected: tfprotov6.StringKindPlain,
		},
	}

	for name, testCase := range testCases {
		name, testCase := name, testCase

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got := tfprotov5tov6.StringKind(testCase.in)

			if diff := cmp.Diff(got, testCase.expected); diff != "" {
				t.Errorf("unexpected difference: %s", diff)
			}
		})
	}
}

func TestUpgradeResourceStateRequest(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		in       *tfprotov5.UpgradeResourceStateRequest
		expected *tfprotov6.UpgradeResourceStateRequest
	}{
		"nil": {
			in:       nil,
			expected: nil,
		},
		"all-valid-fields": {
			in: &tfprotov5.UpgradeResourceStateRequest{
				RawState: &tfprotov5.RawState{
					JSON: testBytes,
				},
				TypeName: "test_resource",
				Version:  1,
			},
			expected: &tfprotov6.UpgradeResourceStateRequest{
				RawState: &tfprotov6.RawState{
					JSON: testBytes,
				},
				TypeName: "test_resource",
				Version:  1,
			},
		},
	}

	for name, testCase := range testCases {
		name, testCase := name, testCase

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got := tfprotov5tov6.UpgradeResourceStateRequest(testCase.in)

			if diff := cmp.Diff(got, testCase.expected); diff != "" {
				t.Errorf("unexpected difference: %s", diff)
			}
		})
	}
}

func TestUpgradeResourceStateResponse(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		in       *tfprotov5.UpgradeResourceStateResponse
		expected *tfprotov6.UpgradeResourceStateResponse
	}{
		"nil": {
			in:       nil,
			expected: nil,
		},
		"all-valid-fields": {
			in: &tfprotov5.UpgradeResourceStateResponse{
				Diagnostics:   testTfprotov5Diagnostics,
				UpgradedState: &testTfprotov5DynamicValue,
			},
			expected: &tfprotov6.UpgradeResourceStateResponse{
				Diagnostics:   testTfprotov6Diagnostics,
				UpgradedState: &testTfprotov6DynamicValue,
			},
		},
	}

	for name, testCase := range testCases {
		name, testCase := name, testCase

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got := tfprotov5tov6.UpgradeResourceStateResponse(testCase.in)

			if diff := cmp.Diff(got, testCase.expected); diff != "" {
				t.Errorf("unexpected difference: %s", diff)
			}
		})
	}
}

func TestValidateDataResourceConfigRequest(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		in       *tfprotov5.ValidateDataSourceConfigRequest
		expected *tfprotov6.ValidateDataResourceConfigRequest
	}{
		"nil": {
			in:       nil,
			expected: nil,
		},
		"all-valid-fields": {
			in: &tfprotov5.ValidateDataSourceConfigRequest{
				Config:   &testTfprotov5DynamicValue,
				TypeName: "test_data_source",
			},
			expected: &tfprotov6.ValidateDataResourceConfigRequest{
				Config:   &testTfprotov6DynamicValue,
				TypeName: "test_data_source",
			},
		},
	}

	for name, testCase := range testCases {
		name, testCase := name, testCase

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got := tfprotov5tov6.ValidateDataResourceConfigRequest(testCase.in)

			if diff := cmp.Diff(got, testCase.expected); diff != "" {
				t.Errorf("unexpected difference: %s", diff)
			}
		})
	}
}

func TestValidateDataResourceConfigResponse(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		in       *tfprotov5.ValidateDataSourceConfigResponse
		expected *tfprotov6.ValidateDataResourceConfigResponse
	}{
		"nil": {
			in:       nil,
			expected: nil,
		},
		"all-valid-fields": {
			in: &tfprotov5.ValidateDataSourceConfigResponse{
				Diagnostics: testTfprotov5Diagnostics,
			},
			expected: &tfprotov6.ValidateDataResourceConfigResponse{
				Diagnostics: testTfprotov6Diagnostics,
			},
		},
	}

	for name, testCase := range testCases {
		name, testCase := name, testCase

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got := tfprotov5tov6.ValidateDataResourceConfigResponse(testCase.in)

			if diff := cmp.Diff(got, testCase.expected); diff != "" {
				t.Errorf("unexpected difference: %s", diff)
			}
		})
	}
}

func TestValidateProviderConfigRequest(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		in       *tfprotov5.PrepareProviderConfigRequest
		expected *tfprotov6.ValidateProviderConfigRequest
	}{
		"nil": {
			in:       nil,
			expected: nil,
		},
		"all-valid-fields": {
			in: &tfprotov5.PrepareProviderConfigRequest{
				Config: &testTfprotov5DynamicValue,
			},
			expected: &tfprotov6.ValidateProviderConfigRequest{
				Config: &testTfprotov6DynamicValue,
			},
		},
	}

	for name, testCase := range testCases {
		name, testCase := name, testCase

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got := tfprotov5tov6.ValidateProviderConfigRequest(testCase.in)

			if diff := cmp.Diff(got, testCase.expected); diff != "" {
				t.Errorf("unexpected difference: %s", diff)
			}
		})
	}
}

func TestValidateProviderConfigResponse(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		in       *tfprotov5.PrepareProviderConfigResponse
		expected *tfprotov6.ValidateProviderConfigResponse
	}{
		"nil": {
			in:       nil,
			expected: nil,
		},
		"all-valid-fields": {
			in: &tfprotov5.PrepareProviderConfigResponse{
				Diagnostics:    testTfprotov5Diagnostics,
				PreparedConfig: &testTfprotov5DynamicValue,
			},
			expected: &tfprotov6.ValidateProviderConfigResponse{
				Diagnostics:    testTfprotov6Diagnostics,
				PreparedConfig: &testTfprotov6DynamicValue,
			},
		},
	}

	for name, testCase := range testCases {
		name, testCase := name, testCase

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got := tfprotov5tov6.ValidateProviderConfigResponse(testCase.in)

			if diff := cmp.Diff(got, testCase.expected); diff != "" {
				t.Errorf("unexpected difference: %s", diff)
			}
		})
	}
}

func TestValidateResourceConfigRequest(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		in       *tfprotov5.ValidateResourceTypeConfigRequest
		expected *tfprotov6.ValidateResourceConfigRequest
	}{
		"nil": {
			in:       nil,
			expected: nil,
		},
		"all-valid-fields": {
			in: &tfprotov5.ValidateResourceTypeConfigRequest{
				Config:   &testTfprotov5DynamicValue,
				TypeName: "test_resource",
			},
			expected: &tfprotov6.ValidateResourceConfigRequest{
				Config:   &testTfprotov6DynamicValue,
				TypeName: "test_resource",
			},
		},
	}

	for name, testCase := range testCases {
		name, testCase := name, testCase

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got := tfprotov5tov6.ValidateResourceConfigRequest(testCase.in)

			if diff := cmp.Diff(got, testCase.expected); diff != "" {
				t.Errorf("unexpected difference: %s", diff)
			}
		})
	}
}

func TestValidateResourceConfigResponse(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		in       *tfprotov5.ValidateResourceTypeConfigResponse
		expected *tfprotov6.ValidateResourceConfigResponse
	}{
		"nil": {
			in:       nil,
			expected: nil,
		},
		"all-valid-fields": {
			in: &tfprotov5.ValidateResourceTypeConfigResponse{
				Diagnostics: testTfprotov5Diagnostics,
			},
			expected: &tfprotov6.ValidateResourceConfigResponse{
				Diagnostics: testTfprotov6Diagnostics,
			},
		},
	}

	for name, testCase := range testCases {
		name, testCase := name, testCase

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got := tfprotov5tov6.ValidateResourceConfigResponse(testCase.in)

			if diff := cmp.Diff(got, testCase.expected); diff != "" {
				t.Errorf("unexpected difference: %s", diff)
			}
		})
	}
}
