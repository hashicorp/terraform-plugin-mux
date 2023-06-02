// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package tfprotov6tov5_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/terraform-plugin-go/tfprotov5"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/hashicorp/terraform-plugin-mux/internal/tfprotov6tov5"
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
		in       *tfprotov6.ApplyResourceChangeRequest
		expected *tfprotov5.ApplyResourceChangeRequest
	}{
		"nil": {
			in:       nil,
			expected: nil,
		},
		"all-valid-fields": {
			in: &tfprotov6.ApplyResourceChangeRequest{
				Config:         &testTfprotov6DynamicValue,
				PlannedPrivate: testBytes,
				PlannedState:   &testTfprotov6DynamicValue,
				PriorState:     &testTfprotov6DynamicValue,
				ProviderMeta:   &testTfprotov6DynamicValue,
				TypeName:       "test",
			},
			expected: &tfprotov5.ApplyResourceChangeRequest{
				Config:         &testTfprotov5DynamicValue,
				PlannedPrivate: testBytes,
				PlannedState:   &testTfprotov5DynamicValue,
				PriorState:     &testTfprotov5DynamicValue,
				ProviderMeta:   &testTfprotov5DynamicValue,
				TypeName:       "test",
			},
		},
	}

	for name, testCase := range testCases {
		name, testCase := name, testCase

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got := tfprotov6tov5.ApplyResourceChangeRequest(testCase.in)

			if diff := cmp.Diff(got, testCase.expected); diff != "" {
				t.Errorf("unexpected difference: %s", diff)
			}
		})
	}
}

func TestApplyResourceChangeResponse(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		in       *tfprotov6.ApplyResourceChangeResponse
		expected *tfprotov5.ApplyResourceChangeResponse
	}{
		"nil": {
			in:       nil,
			expected: nil,
		},
		"all-valid-fields": {
			in: &tfprotov6.ApplyResourceChangeResponse{
				Diagnostics:                 testTfprotov6Diagnostics,
				NewState:                    &testTfprotov6DynamicValue,
				Private:                     testBytes,
				UnsafeToUseLegacyTypeSystem: true,
			},
			expected: &tfprotov5.ApplyResourceChangeResponse{
				Diagnostics:                 testTfprotov5Diagnostics,
				NewState:                    &testTfprotov5DynamicValue,
				Private:                     testBytes,
				UnsafeToUseLegacyTypeSystem: true,
			},
		},
	}

	for name, testCase := range testCases {
		name, testCase := name, testCase

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got := tfprotov6tov5.ApplyResourceChangeResponse(testCase.in)

			if diff := cmp.Diff(got, testCase.expected); diff != "" {
				t.Errorf("unexpected difference: %s", diff)
			}
		})
	}
}

func TestConfigureProviderRequest(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		in       *tfprotov6.ConfigureProviderRequest
		expected *tfprotov5.ConfigureProviderRequest
	}{
		"nil": {
			in:       nil,
			expected: nil,
		},
		"all-valid-fields": {
			in: &tfprotov6.ConfigureProviderRequest{
				Config:           &testTfprotov6DynamicValue,
				TerraformVersion: "1.2.3",
			},
			expected: &tfprotov5.ConfigureProviderRequest{
				Config:           &testTfprotov5DynamicValue,
				TerraformVersion: "1.2.3",
			},
		},
	}

	for name, testCase := range testCases {
		name, testCase := name, testCase

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got := tfprotov6tov5.ConfigureProviderRequest(testCase.in)

			if diff := cmp.Diff(got, testCase.expected); diff != "" {
				t.Errorf("unexpected difference: %s", diff)
			}
		})
	}
}

func TestConfigureProviderResponse(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		in       *tfprotov6.ConfigureProviderResponse
		expected *tfprotov5.ConfigureProviderResponse
	}{
		"nil": {
			in:       nil,
			expected: nil,
		},
		"all-valid-fields": {
			in: &tfprotov6.ConfigureProviderResponse{
				Diagnostics: testTfprotov6Diagnostics,
			},
			expected: &tfprotov5.ConfigureProviderResponse{
				Diagnostics: testTfprotov5Diagnostics,
			},
		},
	}

	for name, testCase := range testCases {
		name, testCase := name, testCase

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got := tfprotov6tov5.ConfigureProviderResponse(testCase.in)

			if diff := cmp.Diff(got, testCase.expected); diff != "" {
				t.Errorf("unexpected difference: %s", diff)
			}
		})
	}
}

func TestDiagnostics(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		in       []*tfprotov6.Diagnostic
		expected []*tfprotov5.Diagnostic
	}{
		"nil": {
			in:       nil,
			expected: nil,
		},
		"empty": {
			in:       []*tfprotov6.Diagnostic{},
			expected: []*tfprotov5.Diagnostic{},
		},
		"one": {
			in:       testTfprotov6Diagnostics,
			expected: testTfprotov5Diagnostics,
		},
		"multiple": {
			in: []*tfprotov6.Diagnostic{
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
			expected: []*tfprotov5.Diagnostic{
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
		},
	}

	for name, testCase := range testCases {
		name, testCase := name, testCase

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got := tfprotov6tov5.Diagnostics(testCase.in)

			if diff := cmp.Diff(got, testCase.expected); diff != "" {
				t.Errorf("unexpected difference: %s", diff)
			}
		})
	}
}

func TestDynamicValue(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		in       *tfprotov6.DynamicValue
		expected *tfprotov5.DynamicValue
	}{
		"nil": {
			in:       nil,
			expected: nil,
		},
		"all-valid-fields": {
			in: &tfprotov6.DynamicValue{
				JSON:    testBytes,
				MsgPack: testBytes,
			},
			expected: &tfprotov5.DynamicValue{
				JSON:    testBytes,
				MsgPack: testBytes,
			},
		},
	}

	for name, testCase := range testCases {
		name, testCase := name, testCase

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got := tfprotov6tov5.DynamicValue(testCase.in)

			if diff := cmp.Diff(got, testCase.expected); diff != "" {
				t.Errorf("unexpected difference: %s", diff)
			}
		})
	}
}

func TestGetProviderSchemaRequest(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		in       *tfprotov6.GetProviderSchemaRequest
		expected *tfprotov5.GetProviderSchemaRequest
	}{
		"nil": {
			in:       nil,
			expected: nil,
		},
		"all-valid-fields": {
			in:       &tfprotov6.GetProviderSchemaRequest{},
			expected: &tfprotov5.GetProviderSchemaRequest{},
		},
	}

	for name, testCase := range testCases {
		name, testCase := name, testCase

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got := tfprotov6tov5.GetProviderSchemaRequest(testCase.in)

			if diff := cmp.Diff(got, testCase.expected); diff != "" {
				t.Errorf("unexpected difference: %s", diff)
			}
		})
	}
}

func TestGetProviderSchemaResponse(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		in            *tfprotov6.GetProviderSchemaResponse
		expected      *tfprotov5.GetProviderSchemaResponse
		expectedError error
	}{
		"nil": {
			in:       nil,
			expected: nil,
		},
		"all-valid-fields": {
			in: &tfprotov6.GetProviderSchemaResponse{
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
			expected: &tfprotov5.GetProviderSchemaResponse{
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
		},
		"data-source-nested-attribute-error": {
			in: &tfprotov6.GetProviderSchemaResponse{
				DataSourceSchemas: map[string]*tfprotov6.Schema{
					"test_data_source": {
						Block: &tfprotov6.SchemaBlock{
							Attributes: []*tfprotov6.SchemaAttribute{
								{
									Name: "test_attribute",
									NestedType: &tfprotov6.SchemaObject{
										Attributes: []*tfprotov6.SchemaAttribute{
											{
												Name:     "test_nested_attribute",
												Required: true,
											},
										},
									},
									Required: true,
								},
							},
						},
					},
				},
			},
			expected:      nil,
			expectedError: fmt.Errorf("unable to convert data source \"test_data_source\" schema: unable to convert attribute \"test_attribute\" schema: %w", tfprotov6tov5.ErrSchemaAttributeNestedTypeNotImplemented),
		},
		"provider-nested-attribute-error": {
			in: &tfprotov6.GetProviderSchemaResponse{
				Provider: &tfprotov6.Schema{
					Block: &tfprotov6.SchemaBlock{
						Attributes: []*tfprotov6.SchemaAttribute{
							{
								Name: "test_attribute",
								NestedType: &tfprotov6.SchemaObject{
									Attributes: []*tfprotov6.SchemaAttribute{
										{
											Name:     "test_nested_attribute",
											Required: true,
										},
									},
								},
								Required: true,
							},
						},
					},
				},
			},
			expected:      nil,
			expectedError: fmt.Errorf("unable to convert provider schema: unable to convert attribute \"test_attribute\" schema: %w", tfprotov6tov5.ErrSchemaAttributeNestedTypeNotImplemented),
		},
		"provider-meta-nested-attribute-error": {
			in: &tfprotov6.GetProviderSchemaResponse{
				ProviderMeta: &tfprotov6.Schema{
					Block: &tfprotov6.SchemaBlock{
						Attributes: []*tfprotov6.SchemaAttribute{
							{
								Name: "test_attribute",
								NestedType: &tfprotov6.SchemaObject{
									Attributes: []*tfprotov6.SchemaAttribute{
										{
											Name:     "test_nested_attribute",
											Required: true,
										},
									},
								},
								Required: true,
							},
						},
					},
				},
			},
			expected:      nil,
			expectedError: fmt.Errorf("unable to convert provider meta schema: unable to convert attribute \"test_attribute\" schema: %w", tfprotov6tov5.ErrSchemaAttributeNestedTypeNotImplemented),
		},
		"resource-nested-attribute-error": {
			in: &tfprotov6.GetProviderSchemaResponse{
				ResourceSchemas: map[string]*tfprotov6.Schema{
					"test_resource": {
						Block: &tfprotov6.SchemaBlock{
							Attributes: []*tfprotov6.SchemaAttribute{
								{
									Name: "test_attribute",
									NestedType: &tfprotov6.SchemaObject{
										Attributes: []*tfprotov6.SchemaAttribute{
											{
												Name:     "test_nested_attribute",
												Required: true,
											},
										},
									},
									Required: true,
								},
							},
						},
					},
				},
			},
			expected:      nil,
			expectedError: fmt.Errorf("unable to convert resource \"test_resource\" schema: unable to convert attribute \"test_attribute\" schema: %w", tfprotov6tov5.ErrSchemaAttributeNestedTypeNotImplemented),
		},
	}

	for name, testCase := range testCases {
		name, testCase := name, testCase

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got, err := tfprotov6tov5.GetProviderSchemaResponse(testCase.in)

			if err != nil {
				if testCase.expectedError == nil {
					t.Fatalf("wanted no error, got unexpected error: %s", err)
				}

				if !strings.Contains(err.Error(), testCase.expectedError.Error()) {
					t.Fatalf("expected error %q, got: %s", testCase.expectedError, err)
				}
			} else if testCase.expectedError != nil {
				t.Fatalf("got no error, expected error: %s", testCase.expectedError)
			}

			if diff := cmp.Diff(got, testCase.expected); diff != "" {
				t.Errorf("unexpected difference: %s", diff)
			}
		})
	}
}

func TestImportResourceStateRequest(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		in       *tfprotov6.ImportResourceStateRequest
		expected *tfprotov5.ImportResourceStateRequest
	}{
		"nil": {
			in:       nil,
			expected: nil,
		},
		"all-valid-fields": {
			in: &tfprotov6.ImportResourceStateRequest{
				ID:       "test-id",
				TypeName: "test_resource",
			},
			expected: &tfprotov5.ImportResourceStateRequest{
				ID:       "test-id",
				TypeName: "test_resource",
			},
		},
	}

	for name, testCase := range testCases {
		name, testCase := name, testCase

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got := tfprotov6tov5.ImportResourceStateRequest(testCase.in)

			if diff := cmp.Diff(got, testCase.expected); diff != "" {
				t.Errorf("unexpected difference: %s", diff)
			}
		})
	}
}

func TestImportResourceStateResponse(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		in       *tfprotov6.ImportResourceStateResponse
		expected *tfprotov5.ImportResourceStateResponse
	}{
		"nil": {
			in:       nil,
			expected: nil,
		},
		"all-valid-fields": {
			in: &tfprotov6.ImportResourceStateResponse{
				Diagnostics: testTfprotov6Diagnostics,
				ImportedResources: []*tfprotov6.ImportedResource{
					{
						Private:  testBytes,
						State:    &testTfprotov6DynamicValue,
						TypeName: "test_resource1",
					},
				},
			},
			expected: &tfprotov5.ImportResourceStateResponse{
				Diagnostics: testTfprotov5Diagnostics,
				ImportedResources: []*tfprotov5.ImportedResource{
					{
						Private:  testBytes,
						State:    &testTfprotov5DynamicValue,
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

			got := tfprotov6tov5.ImportResourceStateResponse(testCase.in)

			if diff := cmp.Diff(got, testCase.expected); diff != "" {
				t.Errorf("unexpected difference: %s", diff)
			}
		})
	}
}

func TestImportedResources(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		in       []*tfprotov6.ImportedResource
		expected []*tfprotov5.ImportedResource
	}{
		"nil": {
			in:       nil,
			expected: nil,
		},
		"empty": {
			in:       []*tfprotov6.ImportedResource{},
			expected: []*tfprotov5.ImportedResource{},
		},
		"one": {
			in: []*tfprotov6.ImportedResource{
				{
					Private:  testBytes,
					State:    &testTfprotov6DynamicValue,
					TypeName: "test_resource1",
				},
			},
			expected: []*tfprotov5.ImportedResource{
				{
					Private:  testBytes,
					State:    &testTfprotov5DynamicValue,
					TypeName: "test_resource1",
				},
			},
		},

		"multiple": {
			in: []*tfprotov6.ImportedResource{
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
			expected: []*tfprotov5.ImportedResource{
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
		},
	}

	for name, testCase := range testCases {
		name, testCase := name, testCase

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got := tfprotov6tov5.ImportedResources(testCase.in)

			if diff := cmp.Diff(got, testCase.expected); diff != "" {
				t.Errorf("unexpected difference: %s", diff)
			}
		})
	}
}

func TestPlanResourceChangeRequest(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		in       *tfprotov6.PlanResourceChangeRequest
		expected *tfprotov5.PlanResourceChangeRequest
	}{
		"nil": {
			in:       nil,
			expected: nil,
		},
		"all-valid-fields": {
			in: &tfprotov6.PlanResourceChangeRequest{
				Config:           &testTfprotov6DynamicValue,
				PriorPrivate:     testBytes,
				PriorState:       &testTfprotov6DynamicValue,
				ProposedNewState: &testTfprotov6DynamicValue,
				ProviderMeta:     &testTfprotov6DynamicValue,
				TypeName:         "test_resource",
			},
			expected: &tfprotov5.PlanResourceChangeRequest{
				Config:           &testTfprotov5DynamicValue,
				PriorPrivate:     testBytes,
				PriorState:       &testTfprotov5DynamicValue,
				ProposedNewState: &testTfprotov5DynamicValue,
				ProviderMeta:     &testTfprotov5DynamicValue,
				TypeName:         "test_resource",
			},
		},
	}

	for name, testCase := range testCases {
		name, testCase := name, testCase

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got := tfprotov6tov5.PlanResourceChangeRequest(testCase.in)

			if diff := cmp.Diff(got, testCase.expected); diff != "" {
				t.Errorf("unexpected difference: %s", diff)
			}
		})
	}
}

func TestPlanResourceChangeResponse(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		in       *tfprotov6.PlanResourceChangeResponse
		expected *tfprotov5.PlanResourceChangeResponse
	}{
		"nil": {
			in:       nil,
			expected: nil,
		},
		"all-valid-fields": {
			in: &tfprotov6.PlanResourceChangeResponse{
				Diagnostics:    testTfprotov6Diagnostics,
				PlannedPrivate: testBytes,
				PlannedState:   &testTfprotov6DynamicValue,
				RequiresReplace: []*tftypes.AttributePath{
					tftypes.NewAttributePath().WithAttributeName("test"),
				},
				UnsafeToUseLegacyTypeSystem: true,
			},
			expected: &tfprotov5.PlanResourceChangeResponse{
				Diagnostics:    testTfprotov5Diagnostics,
				PlannedPrivate: testBytes,
				PlannedState:   &testTfprotov5DynamicValue,
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

			got := tfprotov6tov5.PlanResourceChangeResponse(testCase.in)

			if diff := cmp.Diff(got, testCase.expected); diff != "" {
				t.Errorf("unexpected difference: %s", diff)
			}
		})
	}
}

func TestPrepareProviderConfigRequest(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		in       *tfprotov6.ValidateProviderConfigRequest
		expected *tfprotov5.PrepareProviderConfigRequest
	}{
		"nil": {
			in:       nil,
			expected: nil,
		},
		"all-valid-fields": {
			in: &tfprotov6.ValidateProviderConfigRequest{
				Config: &testTfprotov6DynamicValue,
			},
			expected: &tfprotov5.PrepareProviderConfigRequest{
				Config: &testTfprotov5DynamicValue,
			},
		},
	}

	for name, testCase := range testCases {
		name, testCase := name, testCase

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got := tfprotov6tov5.PrepareProviderConfigRequest(testCase.in)

			if diff := cmp.Diff(got, testCase.expected); diff != "" {
				t.Errorf("unexpected difference: %s", diff)
			}
		})
	}
}

func TestPrepareProviderConfigResponse(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		in       *tfprotov6.ValidateProviderConfigResponse
		expected *tfprotov5.PrepareProviderConfigResponse
	}{
		"nil": {
			in:       nil,
			expected: nil,
		},
		"all-valid-fields": {
			in: &tfprotov6.ValidateProviderConfigResponse{
				Diagnostics:    testTfprotov6Diagnostics,
				PreparedConfig: &testTfprotov6DynamicValue,
			},
			expected: &tfprotov5.PrepareProviderConfigResponse{
				Diagnostics:    testTfprotov5Diagnostics,
				PreparedConfig: &testTfprotov5DynamicValue,
			},
		},
	}

	for name, testCase := range testCases {
		name, testCase := name, testCase

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got := tfprotov6tov5.PrepareProviderConfigResponse(testCase.in)

			if diff := cmp.Diff(got, testCase.expected); diff != "" {
				t.Errorf("unexpected difference: %s", diff)
			}
		})
	}
}

func TestRawState(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		in       *tfprotov6.RawState
		expected *tfprotov5.RawState
	}{
		"nil": {
			in:       nil,
			expected: nil,
		},
		"all-valid-fields": {
			in: &tfprotov6.RawState{
				Flatmap: map[string]string{"testkey": "testvalue"},
				JSON:    testBytes,
			},
			expected: &tfprotov5.RawState{
				Flatmap: map[string]string{"testkey": "testvalue"},
				JSON:    testBytes,
			},
		},
	}

	for name, testCase := range testCases {
		name, testCase := name, testCase

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got := tfprotov6tov5.RawState(testCase.in)

			if diff := cmp.Diff(got, testCase.expected); diff != "" {
				t.Errorf("unexpected difference: %s", diff)
			}
		})
	}
}

func TestReadDataSourceRequest(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		in       *tfprotov6.ReadDataSourceRequest
		expected *tfprotov5.ReadDataSourceRequest
	}{
		"nil": {
			in:       nil,
			expected: nil,
		},
		"all-valid-fields": {
			in: &tfprotov6.ReadDataSourceRequest{
				Config:       &testTfprotov6DynamicValue,
				ProviderMeta: &testTfprotov6DynamicValue,
				TypeName:     "test_data_source",
			},
			expected: &tfprotov5.ReadDataSourceRequest{
				Config:       &testTfprotov5DynamicValue,
				ProviderMeta: &testTfprotov5DynamicValue,
				TypeName:     "test_data_source",
			},
		},
	}

	for name, testCase := range testCases {
		name, testCase := name, testCase

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got := tfprotov6tov5.ReadDataSourceRequest(testCase.in)

			if diff := cmp.Diff(got, testCase.expected); diff != "" {
				t.Errorf("unexpected difference: %s", diff)
			}
		})
	}
}

func TestReadDataSourceResponse(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		in       *tfprotov6.ReadDataSourceResponse
		expected *tfprotov5.ReadDataSourceResponse
	}{
		"nil": {
			in:       nil,
			expected: nil,
		},
		"all-valid-fields": {
			in: &tfprotov6.ReadDataSourceResponse{
				Diagnostics: testTfprotov6Diagnostics,
				State:       &testTfprotov6DynamicValue,
			},
			expected: &tfprotov5.ReadDataSourceResponse{
				Diagnostics: testTfprotov5Diagnostics,
				State:       &testTfprotov5DynamicValue,
			},
		},
	}

	for name, testCase := range testCases {
		name, testCase := name, testCase

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got := tfprotov6tov5.ReadDataSourceResponse(testCase.in)

			if diff := cmp.Diff(got, testCase.expected); diff != "" {
				t.Errorf("unexpected difference: %s", diff)
			}
		})
	}
}

func TestReadResourceRequest(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		in       *tfprotov6.ReadResourceRequest
		expected *tfprotov5.ReadResourceRequest
	}{
		"nil": {
			in:       nil,
			expected: nil,
		},
		"all-valid-fields": {
			in: &tfprotov6.ReadResourceRequest{
				CurrentState: &testTfprotov6DynamicValue,
				Private:      testBytes,
				ProviderMeta: &testTfprotov6DynamicValue,
				TypeName:     "test_resource",
			},
			expected: &tfprotov5.ReadResourceRequest{
				CurrentState: &testTfprotov5DynamicValue,
				Private:      testBytes,
				ProviderMeta: &testTfprotov5DynamicValue,
				TypeName:     "test_resource",
			},
		},
	}

	for name, testCase := range testCases {
		name, testCase := name, testCase

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got := tfprotov6tov5.ReadResourceRequest(testCase.in)

			if diff := cmp.Diff(got, testCase.expected); diff != "" {
				t.Errorf("unexpected difference: %s", diff)
			}
		})
	}
}

func TestReadResourceResponse(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		in       *tfprotov6.ReadResourceResponse
		expected *tfprotov5.ReadResourceResponse
	}{
		"nil": {
			in:       nil,
			expected: nil,
		},
		"all-valid-fields": {
			in: &tfprotov6.ReadResourceResponse{
				Diagnostics: testTfprotov6Diagnostics,
				NewState:    &testTfprotov6DynamicValue,
				Private:     testBytes,
			},
			expected: &tfprotov5.ReadResourceResponse{
				Diagnostics: testTfprotov5Diagnostics,
				NewState:    &testTfprotov5DynamicValue,
				Private:     testBytes,
			},
		},
	}

	for name, testCase := range testCases {
		name, testCase := name, testCase

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got := tfprotov6tov5.ReadResourceResponse(testCase.in)

			if diff := cmp.Diff(got, testCase.expected); diff != "" {
				t.Errorf("unexpected difference: %s", diff)
			}
		})
	}
}

func TestSchema(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		in            *tfprotov6.Schema
		expected      *tfprotov5.Schema
		expectedError error
	}{
		"nil": {
			in:       nil,
			expected: nil,
		},
		"all-valid-fields": {
			in:       testTfprotov6Schema,
			expected: testTfprotov5Schema,
		},
		"nested-attribute-error": {
			in: &tfprotov6.Schema{
				Block: &tfprotov6.SchemaBlock{
					Attributes: []*tfprotov6.SchemaAttribute{
						{
							Name: "test_attribute",
							NestedType: &tfprotov6.SchemaObject{
								Attributes: []*tfprotov6.SchemaAttribute{
									{
										Name:     "test_nested_attribute",
										Required: true,
									},
								},
							},
							Required: true,
						},
					},
				},
			},
			expected:      nil,
			expectedError: fmt.Errorf("unable to convert attribute \"test_attribute\" schema: %w", tfprotov6tov5.ErrSchemaAttributeNestedTypeNotImplemented),
		},
	}

	for name, testCase := range testCases {
		name, testCase := name, testCase

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got, err := tfprotov6tov5.Schema(testCase.in)

			if err != nil {
				if testCase.expectedError == nil {
					t.Fatalf("wanted no error, got unexpected error: %s", err)
				}

				if !strings.Contains(err.Error(), testCase.expectedError.Error()) {
					t.Fatalf("expected error %q, got: %s", testCase.expectedError, err)
				}
			} else if testCase.expectedError != nil {
				t.Fatalf("got no error, expected error: %s", testCase.expectedError)
			}

			if diff := cmp.Diff(got, testCase.expected); diff != "" {
				t.Errorf("unexpected difference: %s", diff)
			}
		})
	}
}

func TestSchemaAttribute(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		in            *tfprotov6.SchemaAttribute
		expected      *tfprotov5.SchemaAttribute
		expectedError error
	}{
		"nil": {
			in:       nil,
			expected: nil,
		},
		"all-valid-fields": {
			in: &tfprotov6.SchemaAttribute{
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
			expected: &tfprotov5.SchemaAttribute{
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
		},
		"NestedType-error": {
			in: &tfprotov6.SchemaAttribute{
				Name: "test",
				NestedType: &tfprotov6.SchemaObject{
					Attributes: []*tfprotov6.SchemaAttribute{
						{
							Name:     "test_nested_attribute",
							Required: true,
						},
					},
				},
				Required: true,
			},
			expected:      nil,
			expectedError: fmt.Errorf("unable to convert attribute \"test\" schema: %w", tfprotov6tov5.ErrSchemaAttributeNestedTypeNotImplemented),
		},
	}

	for name, testCase := range testCases {
		name, testCase := name, testCase

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got, err := tfprotov6tov5.SchemaAttribute(testCase.in)

			if err != nil {
				if testCase.expectedError == nil {
					t.Fatalf("wanted no error, got unexpected error: %s", err)
				}

				if !strings.Contains(err.Error(), testCase.expectedError.Error()) {
					t.Fatalf("expected error %q, got: %s", testCase.expectedError, err)
				}
			} else if testCase.expectedError != nil {
				t.Fatalf("got no error, expected error: %s", testCase.expectedError)
			}

			if diff := cmp.Diff(got, testCase.expected); diff != "" {
				t.Errorf("unexpected difference: %s", diff)
			}
		})
	}
}

func TestSchemaBlock(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		in            *tfprotov6.SchemaBlock
		expected      *tfprotov5.SchemaBlock
		expectedError error
	}{
		"nil": {
			in:       nil,
			expected: nil,
		},
		"all-valid-fields": {
			in: &tfprotov6.SchemaBlock{
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
			expected: &tfprotov5.SchemaBlock{
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
		},
		"nested-attribute-error": {
			in: &tfprotov6.SchemaBlock{
				Attributes: []*tfprotov6.SchemaAttribute{
					{
						Name: "test_attribute",
						NestedType: &tfprotov6.SchemaObject{
							Attributes: []*tfprotov6.SchemaAttribute{
								{
									Name:     "test_nested_attribute",
									Required: true,
								},
							},
						},
						Required: true,
					},
				},
			},
			expected:      nil,
			expectedError: fmt.Errorf("unable to convert attribute \"test_attribute\" schema: %w", tfprotov6tov5.ErrSchemaAttributeNestedTypeNotImplemented),
		},
	}

	for name, testCase := range testCases {
		name, testCase := name, testCase

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got, err := tfprotov6tov5.SchemaBlock(testCase.in)

			if err != nil {
				if testCase.expectedError == nil {
					t.Fatalf("wanted no error, got unexpected error: %s", err)
				}

				if !strings.Contains(err.Error(), testCase.expectedError.Error()) {
					t.Fatalf("expected error %q, got: %s", testCase.expectedError, err)
				}
			} else if testCase.expectedError != nil {
				t.Fatalf("got no error, expected error: %s", testCase.expectedError)
			}

			if diff := cmp.Diff(got, testCase.expected); diff != "" {
				t.Errorf("unexpected difference: %s", diff)
			}
		})
	}
}

func TestSchemaNestedBlock(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		in            *tfprotov6.SchemaNestedBlock
		expected      *tfprotov5.SchemaNestedBlock
		expectedError error
	}{
		"nil": {
			in:       nil,
			expected: nil,
		},
		"all-valid-fields": {
			in: &tfprotov6.SchemaNestedBlock{
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
			expected: &tfprotov5.SchemaNestedBlock{
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
		},
		"nested-attribute-error": {
			in: &tfprotov6.SchemaNestedBlock{
				Block: &tfprotov6.SchemaBlock{
					Attributes: []*tfprotov6.SchemaAttribute{
						{
							Name: "test_attribute",
							NestedType: &tfprotov6.SchemaObject{
								Attributes: []*tfprotov6.SchemaAttribute{
									{
										Name:     "test_nested_attribute",
										Required: true,
									},
								},
							},
							Required: true,
						},
					},
				},
				Nesting:  tfprotov6.SchemaNestedBlockNestingModeList,
				TypeName: "test_block",
			},
			expected:      nil,
			expectedError: fmt.Errorf("unable to convert block \"test_block\" schema: unable to convert attribute \"test_attribute\" schema: %w", tfprotov6tov5.ErrSchemaAttributeNestedTypeNotImplemented),
		},
	}

	for name, testCase := range testCases {
		name, testCase := name, testCase

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got, err := tfprotov6tov5.SchemaNestedBlock(testCase.in)

			if err != nil {
				if testCase.expectedError == nil {
					t.Fatalf("wanted no error, got unexpected error: %s", err)
				}

				if !strings.Contains(err.Error(), testCase.expectedError.Error()) {
					t.Fatalf("expected error %q, got: %s", testCase.expectedError, err)
				}
			} else if testCase.expectedError != nil {
				t.Fatalf("got no error, expected error: %s", testCase.expectedError)
			}

			if diff := cmp.Diff(got, testCase.expected); diff != "" {
				t.Errorf("unexpected difference: %s", diff)
			}
		})
	}
}

func TestStopProviderRequest(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		in       *tfprotov6.StopProviderRequest
		expected *tfprotov5.StopProviderRequest
	}{
		"nil": {
			in:       nil,
			expected: nil,
		},
		"all-valid-fields": {
			in:       &tfprotov6.StopProviderRequest{},
			expected: &tfprotov5.StopProviderRequest{},
		},
	}

	for name, testCase := range testCases {
		name, testCase := name, testCase

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got := tfprotov6tov5.StopProviderRequest(testCase.in)

			if diff := cmp.Diff(got, testCase.expected); diff != "" {
				t.Errorf("unexpected difference: %s", diff)
			}
		})
	}
}

func TestStopProviderResponse(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		in       *tfprotov6.StopProviderResponse
		expected *tfprotov5.StopProviderResponse
	}{
		"nil": {
			in:       nil,
			expected: nil,
		},
		"all-valid-fields": {
			in: &tfprotov6.StopProviderResponse{
				Error: "test error",
			},
			expected: &tfprotov5.StopProviderResponse{
				Error: "test error",
			},
		},
	}

	for name, testCase := range testCases {
		name, testCase := name, testCase

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got := tfprotov6tov5.StopProviderResponse(testCase.in)

			if diff := cmp.Diff(got, testCase.expected); diff != "" {
				t.Errorf("unexpected difference: %s", diff)
			}
		})
	}
}

func TestStringKind(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		in       tfprotov6.StringKind
		expected tfprotov5.StringKind
	}{
		"markdown": {
			in:       tfprotov6.StringKindMarkdown,
			expected: tfprotov5.StringKindMarkdown,
		},
		"plain": {
			in:       tfprotov6.StringKindPlain,
			expected: tfprotov5.StringKindPlain,
		},
	}

	for name, testCase := range testCases {
		name, testCase := name, testCase

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got := tfprotov6tov5.StringKind(testCase.in)

			if diff := cmp.Diff(got, testCase.expected); diff != "" {
				t.Errorf("unexpected difference: %s", diff)
			}
		})
	}
}

func TestUpgradeResourceStateRequest(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		in       *tfprotov6.UpgradeResourceStateRequest
		expected *tfprotov5.UpgradeResourceStateRequest
	}{
		"nil": {
			in:       nil,
			expected: nil,
		},
		"all-valid-fields": {
			in: &tfprotov6.UpgradeResourceStateRequest{
				RawState: &tfprotov6.RawState{
					JSON: testBytes,
				},
				TypeName: "test_resource",
				Version:  1,
			},
			expected: &tfprotov5.UpgradeResourceStateRequest{
				RawState: &tfprotov5.RawState{
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

			got := tfprotov6tov5.UpgradeResourceStateRequest(testCase.in)

			if diff := cmp.Diff(got, testCase.expected); diff != "" {
				t.Errorf("unexpected difference: %s", diff)
			}
		})
	}
}

func TestUpgradeResourceStateResponse(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		in       *tfprotov6.UpgradeResourceStateResponse
		expected *tfprotov5.UpgradeResourceStateResponse
	}{
		"nil": {
			in:       nil,
			expected: nil,
		},
		"all-valid-fields": {
			in: &tfprotov6.UpgradeResourceStateResponse{
				Diagnostics:   testTfprotov6Diagnostics,
				UpgradedState: &testTfprotov6DynamicValue,
			},
			expected: &tfprotov5.UpgradeResourceStateResponse{
				Diagnostics:   testTfprotov5Diagnostics,
				UpgradedState: &testTfprotov5DynamicValue,
			},
		},
	}

	for name, testCase := range testCases {
		name, testCase := name, testCase

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got := tfprotov6tov5.UpgradeResourceStateResponse(testCase.in)

			if diff := cmp.Diff(got, testCase.expected); diff != "" {
				t.Errorf("unexpected difference: %s", diff)
			}
		})
	}
}

func TestValidateDataSourceConfigRequest(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		in       *tfprotov6.ValidateDataResourceConfigRequest
		expected *tfprotov5.ValidateDataSourceConfigRequest
	}{
		"nil": {
			in:       nil,
			expected: nil,
		},
		"all-valid-fields": {
			in: &tfprotov6.ValidateDataResourceConfigRequest{
				Config:   &testTfprotov6DynamicValue,
				TypeName: "test_data_source",
			},
			expected: &tfprotov5.ValidateDataSourceConfigRequest{
				Config:   &testTfprotov5DynamicValue,
				TypeName: "test_data_source",
			},
		},
	}

	for name, testCase := range testCases {
		name, testCase := name, testCase

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got := tfprotov6tov5.ValidateDataSourceConfigRequest(testCase.in)

			if diff := cmp.Diff(got, testCase.expected); diff != "" {
				t.Errorf("unexpected difference: %s", diff)
			}
		})
	}
}

func TestValidateDataSourceConfigResponse(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		in       *tfprotov6.ValidateDataResourceConfigResponse
		expected *tfprotov5.ValidateDataSourceConfigResponse
	}{
		"nil": {
			in:       nil,
			expected: nil,
		},
		"all-valid-fields": {
			in: &tfprotov6.ValidateDataResourceConfigResponse{
				Diagnostics: testTfprotov6Diagnostics,
			},
			expected: &tfprotov5.ValidateDataSourceConfigResponse{
				Diagnostics: testTfprotov5Diagnostics,
			},
		},
	}

	for name, testCase := range testCases {
		name, testCase := name, testCase

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got := tfprotov6tov5.ValidateDataSourceConfigResponse(testCase.in)

			if diff := cmp.Diff(got, testCase.expected); diff != "" {
				t.Errorf("unexpected difference: %s", diff)
			}
		})
	}
}

func TestValidateResourceTypeConfigRequest(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		in       *tfprotov6.ValidateResourceConfigRequest
		expected *tfprotov5.ValidateResourceTypeConfigRequest
	}{
		"nil": {
			in:       nil,
			expected: nil,
		},
		"all-valid-fields": {
			in: &tfprotov6.ValidateResourceConfigRequest{
				Config:   &testTfprotov6DynamicValue,
				TypeName: "test_resource",
			},
			expected: &tfprotov5.ValidateResourceTypeConfigRequest{
				Config:   &testTfprotov5DynamicValue,
				TypeName: "test_resource",
			},
		},
	}

	for name, testCase := range testCases {
		name, testCase := name, testCase

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got := tfprotov6tov5.ValidateResourceTypeConfigRequest(testCase.in)

			if diff := cmp.Diff(got, testCase.expected); diff != "" {
				t.Errorf("unexpected difference: %s", diff)
			}
		})
	}
}

func TestValidateResourceTypeConfigResponse(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		in       *tfprotov6.ValidateResourceConfigResponse
		expected *tfprotov5.ValidateResourceTypeConfigResponse
	}{
		"nil": {
			in:       nil,
			expected: nil,
		},
		"all-valid-fields": {
			in: &tfprotov6.ValidateResourceConfigResponse{
				Diagnostics: testTfprotov6Diagnostics,
			},
			expected: &tfprotov5.ValidateResourceTypeConfigResponse{
				Diagnostics: testTfprotov5Diagnostics,
			},
		},
	}

	for name, testCase := range testCases {
		name, testCase := name, testCase

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got := tfprotov6tov5.ValidateResourceTypeConfigResponse(testCase.in)

			if diff := cmp.Diff(got, testCase.expected); diff != "" {
				t.Errorf("unexpected difference: %s", diff)
			}
		})
	}
}
