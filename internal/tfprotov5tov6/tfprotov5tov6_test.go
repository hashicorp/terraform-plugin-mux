// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package tfprotov5tov6_test

import (
	"slices"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/terraform-plugin-go/tfprotov5"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-go/tftypes"

	"github.com/hashicorp/terraform-plugin-mux/internal/tfprotov5tov6"
)

var (
	testBytes []byte = []byte("test")

	testTfprotov5ActionMetadata tfprotov5.ActionMetadata = tfprotov5.ActionMetadata{
		TypeName: "test_action",
	}

	testTfprotov6ActionMetadata tfprotov6.ActionMetadata = tfprotov6.ActionMetadata{
		TypeName: "test_action",
	}

	testTfprotov5DataSourceMetadata tfprotov5.DataSourceMetadata = tfprotov5.DataSourceMetadata{
		TypeName: "test_data_source",
	}

	testTfprotov6DataSourceMetadata tfprotov6.DataSourceMetadata = tfprotov6.DataSourceMetadata{
		TypeName: "test_data_source",
	}

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

	testTfprotov5ResourceIdentityData tfprotov5.ResourceIdentityData
	testTfprotov6ResourceIdentityData tfprotov6.ResourceIdentityData

	testTfprotov5EphemeralResourceMetadata tfprotov5.EphemeralResourceMetadata = tfprotov5.EphemeralResourceMetadata{
		TypeName: "test_ephemeral_resource",
	}

	testTfprotov6EphemeralResourceMetadata tfprotov6.EphemeralResourceMetadata = tfprotov6.EphemeralResourceMetadata{
		TypeName: "test_ephemeral_resource",
	}

	testTfprotov5ListResourceMetadata tfprotov5.ListResourceMetadata = tfprotov5.ListResourceMetadata{
		TypeName: "test_list_resource",
	}

	testTfprotov6ListResourceMetadata tfprotov6.ListResourceMetadata = tfprotov6.ListResourceMetadata{
		TypeName: "test_list_resource",
	}

	testTfprotov5Function *tfprotov5.Function = &tfprotov5.Function{
		Parameters: []*tfprotov5.FunctionParameter{},
		Return: &tfprotov5.FunctionReturn{
			Type: tftypes.String,
		},
	}

	testTfprotov5FunctionError *tfprotov5.FunctionError = &tfprotov5.FunctionError{
		Text:             "test error",
		FunctionArgument: pointer(int64(0)),
	}

	testTfprotov6FunctionError *tfprotov6.FunctionError = &tfprotov6.FunctionError{
		Text:             "test error",
		FunctionArgument: pointer(int64(0)),
	}

	testTfprotov5FunctionMetadata tfprotov5.FunctionMetadata = tfprotov5.FunctionMetadata{
		Name: "test_function",
	}

	testTfprotov6Function *tfprotov6.Function = &tfprotov6.Function{
		Parameters: []*tfprotov6.FunctionParameter{},
		Return: &tfprotov6.FunctionReturn{
			Type: tftypes.String,
		},
	}

	testTfprotov6FunctionMetadata tfprotov6.FunctionMetadata = tfprotov6.FunctionMetadata{
		Name: "test_function",
	}

	testTfprotov5ResourceMetadata tfprotov5.ResourceMetadata = tfprotov5.ResourceMetadata{
		TypeName: "test_resource",
	}

	testTfprotov6ResourceMetadata tfprotov6.ResourceMetadata = tfprotov6.ResourceMetadata{
		TypeName: "test_resource",
	}

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

	testTfprotov5ResourceIdentitySchema *tfprotov5.ResourceIdentitySchema = &tfprotov5.ResourceIdentitySchema{
		Version: 1,
		IdentityAttributes: []*tfprotov5.ResourceIdentitySchemaAttribute{
			{
				Name:              "req",
				Type:              tftypes.String,
				RequiredForImport: true,
				Description:       "this one's required",
			},
			{
				Name:              "opt",
				Type:              tftypes.String,
				OptionalForImport: true,
				Description:       "this one's optional",
			},
		},
	}
	testTfprotov6ResourceIdentitySchema *tfprotov6.ResourceIdentitySchema = &tfprotov6.ResourceIdentitySchema{
		Version: 1,
		IdentityAttributes: []*tfprotov6.ResourceIdentitySchemaAttribute{
			{
				Name:              "req",
				Type:              tftypes.String,
				RequiredForImport: true,
				Description:       "this one's required",
			},
			{
				Name:              "opt",
				Type:              tftypes.String,
				OptionalForImport: true,
				Description:       "this one's optional",
			},
		},
	}

	testTime time.Time = time.Date(2000, 1, 2, 3, 4, 5, 6, time.UTC)

	testStreamv5 *tfprotov5.ListResourceServerStream = &tfprotov5.ListResourceServerStream{
		Results: slices.Values([]tfprotov5.ListResourceResult{
			{
				DisplayName: "test",
				Resource:    &testTfprotov5DynamicValue,
				Identity:    &testTfprotov5ResourceIdentityData,
				Diagnostics: testTfprotov5Diagnostics,
			},
		}),
	}

	testStreamv6 *tfprotov6.ListResourceServerStream = &tfprotov6.ListResourceServerStream{
		Results: slices.Values([]tfprotov6.ListResourceResult{
			{
				DisplayName: "test",
				Resource:    &testTfprotov6DynamicValue,
				Identity:    &testTfprotov6ResourceIdentityData,
				Diagnostics: testTfprotov6Diagnostics,
			},
		}),
	}
)

func init() {
	testTfprotov5DynamicValue, _ = tfprotov5.NewDynamicValue(tftypes.String, tftypes.NewValue(tftypes.String, "test"))
	testTfprotov6DynamicValue, _ = tfprotov6.NewDynamicValue(tftypes.String, tftypes.NewValue(tftypes.String, "test"))
	testTfprotov5ResourceIdentityData = tfprotov5.ResourceIdentityData{
		IdentityData: &testTfprotov5DynamicValue,
	}
	testTfprotov6ResourceIdentityData = tfprotov6.ResourceIdentityData{
		IdentityData: &testTfprotov6DynamicValue,
	}
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
				Config:          &testTfprotov5DynamicValue,
				PlannedPrivate:  testBytes,
				PlannedState:    &testTfprotov5DynamicValue,
				PriorState:      &testTfprotov5DynamicValue,
				ProviderMeta:    &testTfprotov5DynamicValue,
				TypeName:        "test",
				PlannedIdentity: &testTfprotov5ResourceIdentityData,
			},
			expected: &tfprotov6.ApplyResourceChangeRequest{
				Config:          &testTfprotov6DynamicValue,
				PlannedPrivate:  testBytes,
				PlannedState:    &testTfprotov6DynamicValue,
				PriorState:      &testTfprotov6DynamicValue,
				ProviderMeta:    &testTfprotov6DynamicValue,
				TypeName:        "test",
				PlannedIdentity: &testTfprotov6ResourceIdentityData,
			},
		},
	}

	for name, testCase := range testCases {

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
				NewIdentity:                 &testTfprotov5ResourceIdentityData,
			},
			expected: &tfprotov6.ApplyResourceChangeResponse{
				Diagnostics:                 testTfprotov6Diagnostics,
				NewState:                    &testTfprotov6DynamicValue,
				Private:                     testBytes,
				UnsafeToUseLegacyTypeSystem: true,
				NewIdentity:                 &testTfprotov6ResourceIdentityData,
			},
		},
	}

	for name, testCase := range testCases {

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got := tfprotov5tov6.ApplyResourceChangeResponse(testCase.in)

			if diff := cmp.Diff(got, testCase.expected); diff != "" {
				t.Errorf("unexpected difference: %s", diff)
			}
		})
	}
}

func TestCallFunctionRequest(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		in       *tfprotov5.CallFunctionRequest
		expected *tfprotov6.CallFunctionRequest
	}{
		"nil": {
			in:       nil,
			expected: nil,
		},
		"all-valid-fields": {
			in: &tfprotov5.CallFunctionRequest{
				Arguments: []*tfprotov5.DynamicValue{
					&testTfprotov5DynamicValue,
				},
				Name: "test_function",
			},
			expected: &tfprotov6.CallFunctionRequest{
				Arguments: []*tfprotov6.DynamicValue{
					&testTfprotov6DynamicValue,
				},
				Name: "test_function",
			},
		},
	}

	for name, testCase := range testCases {

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got := tfprotov5tov6.CallFunctionRequest(testCase.in)

			if diff := cmp.Diff(got, testCase.expected); diff != "" {
				t.Errorf("unexpected difference: %s", diff)
			}
		})
	}
}

func TestCallFunctionResponse(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		in       *tfprotov5.CallFunctionResponse
		expected *tfprotov6.CallFunctionResponse
	}{
		"nil": {
			in:       nil,
			expected: nil,
		},
		"all-valid-fields": {
			in: &tfprotov5.CallFunctionResponse{
				Error:  testTfprotov5FunctionError,
				Result: &testTfprotov5DynamicValue,
			},
			expected: &tfprotov6.CallFunctionResponse{
				Error:  testTfprotov6FunctionError,
				Result: &testTfprotov6DynamicValue,
			},
		},
	}

	for name, testCase := range testCases {

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got := tfprotov5tov6.CallFunctionResponse(testCase.in)

			if diff := cmp.Diff(got, testCase.expected); diff != "" {
				t.Errorf("unexpected difference: %s", diff)
			}
		})
	}
}

func TestCloseEphemeralResourceRequest(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		in       *tfprotov5.CloseEphemeralResourceRequest
		expected *tfprotov6.CloseEphemeralResourceRequest
	}{
		"nil": {
			in:       nil,
			expected: nil,
		},
		"all-valid-fields": {
			in: &tfprotov5.CloseEphemeralResourceRequest{
				Private:  testBytes,
				TypeName: "test_ephemeral_resource",
			},
			expected: &tfprotov6.CloseEphemeralResourceRequest{
				Private:  testBytes,
				TypeName: "test_ephemeral_resource",
			},
		},
	}

	for name, testCase := range testCases {

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got := tfprotov5tov6.CloseEphemeralResourceRequest(testCase.in)

			if diff := cmp.Diff(got, testCase.expected); diff != "" {
				t.Errorf("unexpected difference: %s", diff)
			}
		})
	}
}

func TestCloseEphemeralResourceResponse(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		in       *tfprotov5.CloseEphemeralResourceResponse
		expected *tfprotov6.CloseEphemeralResourceResponse
	}{
		"nil": {
			in:       nil,
			expected: nil,
		},
		"all-valid-fields": {
			in: &tfprotov5.CloseEphemeralResourceResponse{
				Diagnostics: testTfprotov5Diagnostics,
			},
			expected: &tfprotov6.CloseEphemeralResourceResponse{
				Diagnostics: testTfprotov6Diagnostics,
			},
		},
	}

	for name, testCase := range testCases {

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got := tfprotov5tov6.CloseEphemeralResourceResponse(testCase.in)

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
		"client-capabilities-deferral-allowed": {
			in: &tfprotov5.ConfigureProviderRequest{
				Config:           &testTfprotov5DynamicValue,
				TerraformVersion: "1.2.3",
				ClientCapabilities: &tfprotov5.ConfigureProviderClientCapabilities{
					DeferralAllowed: true,
				},
			},
			expected: &tfprotov6.ConfigureProviderRequest{
				Config:           &testTfprotov6DynamicValue,
				TerraformVersion: "1.2.3",
				ClientCapabilities: &tfprotov6.ConfigureProviderClientCapabilities{
					DeferralAllowed: true,
				},
			},
		},
	}

	for name, testCase := range testCases {

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

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got := tfprotov5tov6.DynamicValue(testCase.in)

			if diff := cmp.Diff(got, testCase.expected); diff != "" {
				t.Errorf("unexpected difference: %s", diff)
			}
		})
	}
}

func TestResourceIdentityData(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		in       *tfprotov5.ResourceIdentityData
		expected *tfprotov6.ResourceIdentityData
	}{
		"nil": {
			in:       nil,
			expected: nil,
		},
		"all-valid-fields": {
			in: &tfprotov5.ResourceIdentityData{
				IdentityData: &testTfprotov5DynamicValue,
			},
			expected: &tfprotov6.ResourceIdentityData{
				IdentityData: &testTfprotov6DynamicValue,
			},
		},
	}

	for name, testCase := range testCases {

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got := tfprotov5tov6.ResourceIdentityData(testCase.in)

			if diff := cmp.Diff(got, testCase.expected); diff != "" {
				t.Errorf("unexpected difference: %s", diff)
			}
		})
	}
}

func TestFunction(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		in       *tfprotov5.Function
		expected *tfprotov6.Function
	}{
		"nil": {
			in:       nil,
			expected: nil,
		},
		"all-valid-fields": {
			in: &tfprotov5.Function{
				DeprecationMessage: "test deprecation message",
				Description:        "test description",
				DescriptionKind:    tfprotov5.StringKindPlain,
				Parameters: []*tfprotov5.FunctionParameter{
					{
						Type: tftypes.String,
					},
				},
				Return: &tfprotov5.FunctionReturn{
					Type: tftypes.String,
				},
				Summary: "test summary",
				VariadicParameter: &tfprotov5.FunctionParameter{
					Type: tftypes.String,
				},
			},
			expected: &tfprotov6.Function{
				DeprecationMessage: "test deprecation message",
				Description:        "test description",
				DescriptionKind:    tfprotov6.StringKindPlain,
				Parameters: []*tfprotov6.FunctionParameter{
					{
						Type: tftypes.String,
					},
				},
				Return: &tfprotov6.FunctionReturn{
					Type: tftypes.String,
				},
				Summary: "test summary",
				VariadicParameter: &tfprotov6.FunctionParameter{
					Type: tftypes.String,
				},
			},
		},
	}

	for name, testCase := range testCases {

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got := tfprotov5tov6.Function(testCase.in)

			if diff := cmp.Diff(got, testCase.expected); diff != "" {
				t.Errorf("unexpected difference: %s", diff)
			}
		})
	}
}

func TestFunctionMetadata(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		in       tfprotov5.FunctionMetadata
		expected tfprotov6.FunctionMetadata
	}{
		"all-valid-fields": {
			in: tfprotov5.FunctionMetadata{
				Name: "test_function",
			},
			expected: tfprotov6.FunctionMetadata{
				Name: "test_function",
			},
		},
	}

	for name, testCase := range testCases {

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got := tfprotov5tov6.FunctionMetadata(testCase.in)

			if diff := cmp.Diff(got, testCase.expected); diff != "" {
				t.Errorf("unexpected difference: %s", diff)
			}
		})
	}
}

func TestFunctionParameter(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		in       *tfprotov5.FunctionParameter
		expected *tfprotov6.FunctionParameter
	}{
		"nil": {
			in:       nil,
			expected: nil,
		},
		"all-valid-fields": {
			in: &tfprotov5.FunctionParameter{
				Description:     "test description",
				DescriptionKind: tfprotov5.StringKindPlain,
				Type:            tftypes.String,
			},
			expected: &tfprotov6.FunctionParameter{
				Description:     "test description",
				DescriptionKind: tfprotov6.StringKindPlain,
				Type:            tftypes.String,
			},
		},
	}

	for name, testCase := range testCases {

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got := tfprotov5tov6.FunctionParameter(testCase.in)

			if diff := cmp.Diff(got, testCase.expected); diff != "" {
				t.Errorf("unexpected difference: %s", diff)
			}
		})
	}
}

func TestFunctionReturn(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		in       *tfprotov5.FunctionReturn
		expected *tfprotov6.FunctionReturn
	}{
		"nil": {
			in:       nil,
			expected: nil,
		},
		"all-valid-fields": {
			in: &tfprotov5.FunctionReturn{
				Type: tftypes.String,
			},
			expected: &tfprotov6.FunctionReturn{
				Type: tftypes.String,
			},
		},
	}

	for name, testCase := range testCases {

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got := tfprotov5tov6.FunctionReturn(testCase.in)

			if diff := cmp.Diff(got, testCase.expected); diff != "" {
				t.Errorf("unexpected difference: %s", diff)
			}
		})
	}
}

func TestGetFunctionsRequest(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		in       *tfprotov5.GetFunctionsRequest
		expected *tfprotov6.GetFunctionsRequest
	}{
		"nil": {
			in:       nil,
			expected: nil,
		},
		"all-valid-fields": {
			in:       &tfprotov5.GetFunctionsRequest{},
			expected: &tfprotov6.GetFunctionsRequest{},
		},
	}

	for name, testCase := range testCases {

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got := tfprotov5tov6.GetFunctionsRequest(testCase.in)

			if diff := cmp.Diff(got, testCase.expected); diff != "" {
				t.Errorf("unexpected difference: %s", diff)
			}
		})
	}
}

func TestGetFunctionsResponse(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		in       *tfprotov5.GetFunctionsResponse
		expected *tfprotov6.GetFunctionsResponse
	}{
		"nil": {
			in:       nil,
			expected: nil,
		},
		"all-valid-fields": {
			in: &tfprotov5.GetFunctionsResponse{
				Diagnostics: testTfprotov5Diagnostics,
				Functions: map[string]*tfprotov5.Function{
					"test_function": testTfprotov5Function,
				},
			},
			expected: &tfprotov6.GetFunctionsResponse{
				Diagnostics: testTfprotov6Diagnostics,
				Functions: map[string]*tfprotov6.Function{
					"test_function": testTfprotov6Function,
				},
			},
		},
	}

	for name, testCase := range testCases {

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got := tfprotov5tov6.GetFunctionsResponse(testCase.in)

			if diff := cmp.Diff(got, testCase.expected); diff != "" {
				t.Errorf("unexpected difference: %s", diff)
			}
		})
	}
}

func TestGetMetadataRequest(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		in       *tfprotov5.GetMetadataRequest
		expected *tfprotov6.GetMetadataRequest
	}{
		"nil": {
			in:       nil,
			expected: nil,
		},
		"all-valid-fields": {
			in:       &tfprotov5.GetMetadataRequest{},
			expected: &tfprotov6.GetMetadataRequest{},
		},
	}

	for name, testCase := range testCases {

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got := tfprotov5tov6.GetMetadataRequest(testCase.in)

			if diff := cmp.Diff(got, testCase.expected); diff != "" {
				t.Errorf("unexpected difference: %s", diff)
			}
		})
	}
}

func TestGetMetadataResponse(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		in       *tfprotov5.GetMetadataResponse
		expected *tfprotov6.GetMetadataResponse
	}{
		"nil": {
			in:       nil,
			expected: nil,
		},
		"all-valid-fields": {
			in: &tfprotov5.GetMetadataResponse{
				Actions: []tfprotov5.ActionMetadata{
					testTfprotov5ActionMetadata,
				},
				DataSources: []tfprotov5.DataSourceMetadata{
					testTfprotov5DataSourceMetadata,
				},
				Diagnostics: testTfprotov5Diagnostics,
				EphemeralResources: []tfprotov5.EphemeralResourceMetadata{
					testTfprotov5EphemeralResourceMetadata,
				},
				ListResources: []tfprotov5.ListResourceMetadata{
					testTfprotov5ListResourceMetadata,
				},
				Functions: []tfprotov5.FunctionMetadata{
					testTfprotov5FunctionMetadata,
				},
				Resources: []tfprotov5.ResourceMetadata{
					testTfprotov5ResourceMetadata,
				},
			},
			expected: &tfprotov6.GetMetadataResponse{
				Actions: []tfprotov6.ActionMetadata{
					testTfprotov6ActionMetadata,
				},
				DataSources: []tfprotov6.DataSourceMetadata{
					testTfprotov6DataSourceMetadata,
				},
				Diagnostics: testTfprotov6Diagnostics,
				EphemeralResources: []tfprotov6.EphemeralResourceMetadata{
					testTfprotov6EphemeralResourceMetadata,
				},
				ListResources: []tfprotov6.ListResourceMetadata{
					testTfprotov6ListResourceMetadata,
				},
				Functions: []tfprotov6.FunctionMetadata{
					testTfprotov6FunctionMetadata,
				},
				Resources: []tfprotov6.ResourceMetadata{
					testTfprotov6ResourceMetadata,
				},
			},
		},
	}

	for name, testCase := range testCases {

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got := tfprotov5tov6.GetMetadataResponse(testCase.in)

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
				ActionSchemas: map[string]*tfprotov5.ActionSchema{
					"test_action": {
						Schema: testTfprotov5Schema,
						Type:   tfprotov5.UnlinkedActionSchemaType{},
					},
				},
				DataSourceSchemas: map[string]*tfprotov5.Schema{
					"test_data_source": testTfprotov5Schema,
				},
				Diagnostics: testTfprotov5Diagnostics,
				EphemeralResourceSchemas: map[string]*tfprotov5.Schema{
					"test_ephemeral_resource": testTfprotov5Schema,
				},
				ListResourceSchemas: map[string]*tfprotov5.Schema{
					"test_list_resource": testTfprotov5Schema,
				},
				Functions: map[string]*tfprotov5.Function{
					"test_function": testTfprotov5Function,
				},
				Provider:     testTfprotov5Schema,
				ProviderMeta: testTfprotov5Schema,
				ResourceSchemas: map[string]*tfprotov5.Schema{
					"test_resource": testTfprotov5Schema,
				},
			},
			expected: &tfprotov6.GetProviderSchemaResponse{
				ActionSchemas: map[string]*tfprotov6.ActionSchema{
					"test_action": {
						Schema: testTfprotov6Schema,
						Type:   tfprotov6.UnlinkedActionSchemaType{},
					},
				},
				DataSourceSchemas: map[string]*tfprotov6.Schema{
					"test_data_source": testTfprotov6Schema,
				},
				Diagnostics: testTfprotov6Diagnostics,
				EphemeralResourceSchemas: map[string]*tfprotov6.Schema{
					"test_ephemeral_resource": testTfprotov6Schema,
				},
				ListResourceSchemas: map[string]*tfprotov6.Schema{
					"test_list_resource": testTfprotov6Schema,
				},
				Functions: map[string]*tfprotov6.Function{
					"test_function": testTfprotov6Function,
				},
				Provider:     testTfprotov6Schema,
				ProviderMeta: testTfprotov6Schema,
				ResourceSchemas: map[string]*tfprotov6.Schema{
					"test_resource": testTfprotov6Schema,
				},
			},
		},
	}

	for name, testCase := range testCases {

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got := tfprotov5tov6.GetProviderSchemaResponse(testCase.in)

			if diff := cmp.Diff(got, testCase.expected); diff != "" {
				t.Errorf("unexpected difference: %s", diff)
			}
		})
	}
}

func TestGetResourceIdentitySchemasRequest(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		in       *tfprotov5.GetResourceIdentitySchemasRequest
		expected *tfprotov6.GetResourceIdentitySchemasRequest
	}{
		"nil": {
			in:       nil,
			expected: nil,
		},
		"all-valid-fields": {
			in:       &tfprotov5.GetResourceIdentitySchemasRequest{},
			expected: &tfprotov6.GetResourceIdentitySchemasRequest{},
		},
	}

	for name, testCase := range testCases {

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got := tfprotov5tov6.GetResourceIdentitySchemasRequest(testCase.in)

			if diff := cmp.Diff(got, testCase.expected); diff != "" {
				t.Errorf("unexpected difference: %s", diff)
			}
		})
	}
}

func TestGetResourceIdentitySchemasResponse(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		in       *tfprotov5.GetResourceIdentitySchemasResponse
		expected *tfprotov6.GetResourceIdentitySchemasResponse
	}{
		"nil": {
			in:       nil,
			expected: nil,
		},
		"all-valid-fields": {
			in: &tfprotov5.GetResourceIdentitySchemasResponse{
				Diagnostics: testTfprotov5Diagnostics,
				IdentitySchemas: map[string]*tfprotov5.ResourceIdentitySchema{
					"test_resource": testTfprotov5ResourceIdentitySchema,
				},
			},
			expected: &tfprotov6.GetResourceIdentitySchemasResponse{
				Diagnostics: testTfprotov6Diagnostics,
				IdentitySchemas: map[string]*tfprotov6.ResourceIdentitySchema{
					"test_resource": testTfprotov6ResourceIdentitySchema,
				},
			},
		},
	}

	for name, testCase := range testCases {

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got := tfprotov5tov6.GetResourceIdentitySchemasResponse(testCase.in)

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
				Identity: &testTfprotov5ResourceIdentityData,
			},
			expected: &tfprotov6.ImportResourceStateRequest{
				ID:       "test-id",
				TypeName: "test_resource",
				Identity: &testTfprotov6ResourceIdentityData,
			},
		},
		"client-capabilities-deferral-allowed": {
			in: &tfprotov5.ImportResourceStateRequest{
				ID:       "test-id",
				TypeName: "test_resource",
				ClientCapabilities: &tfprotov5.ImportResourceStateClientCapabilities{
					DeferralAllowed: true,
				},
			},
			expected: &tfprotov6.ImportResourceStateRequest{
				ID:       "test-id",
				TypeName: "test_resource",
				ClientCapabilities: &tfprotov6.ImportResourceStateClientCapabilities{
					DeferralAllowed: true,
				},
			},
		},
	}

	for name, testCase := range testCases {

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
						Identity: &testTfprotov5ResourceIdentityData,
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
						Identity: &testTfprotov6ResourceIdentityData,
					},
				},
			},
		},
		"deferred-reason": {
			in: &tfprotov5.ImportResourceStateResponse{
				Diagnostics: testTfprotov5Diagnostics,
				ImportedResources: []*tfprotov5.ImportedResource{
					{
						Private:  testBytes,
						State:    &testTfprotov5DynamicValue,
						TypeName: "test_resource1",
					},
				},
				Deferred: &tfprotov5.Deferred{
					Reason: tfprotov5.DeferredReasonResourceConfigUnknown,
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
				Deferred: &tfprotov6.Deferred{
					Reason: tfprotov6.DeferredReasonResourceConfigUnknown,
				},
			},
		},
	}

	for name, testCase := range testCases {

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
					Identity: &testTfprotov5ResourceIdentityData,
				},
			},
			expected: []*tfprotov6.ImportedResource{
				{
					Private:  testBytes,
					State:    &testTfprotov6DynamicValue,
					TypeName: "test_resource1",
					Identity: &testTfprotov6ResourceIdentityData,
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

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got := tfprotov5tov6.ImportedResources(testCase.in)

			if diff := cmp.Diff(got, testCase.expected); diff != "" {
				t.Errorf("unexpected difference: %s", diff)
			}
		})
	}
}

func TestMoveResourceStateRequest(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		in       *tfprotov5.MoveResourceStateRequest
		expected *tfprotov6.MoveResourceStateRequest
	}{
		"nil": {
			in:       nil,
			expected: nil,
		},
		"all-valid-fields": {
			in: &tfprotov5.MoveResourceStateRequest{
				SourcePrivate:         testBytes,
				SourceProviderAddress: "example.com/namespace/test",
				SourceSchemaVersion:   1,
				SourceState: &tfprotov5.RawState{
					JSON: testBytes,
				},
				SourceTypeName: "test_source",
				TargetTypeName: "test_target",
				SourceIdentity: &tfprotov5.RawState{
					JSON: testBytes,
				},
				SourceIdentitySchemaVersion: 1,
			},
			expected: &tfprotov6.MoveResourceStateRequest{
				SourcePrivate:         testBytes,
				SourceProviderAddress: "example.com/namespace/test",
				SourceSchemaVersion:   1,
				SourceState: &tfprotov6.RawState{
					JSON: testBytes,
				},
				SourceTypeName: "test_source",
				TargetTypeName: "test_target",
				SourceIdentity: &tfprotov6.RawState{
					JSON: testBytes,
				},
				SourceIdentitySchemaVersion: 1,
			},
		},
	}

	for name, testCase := range testCases {

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got := tfprotov5tov6.MoveResourceStateRequest(testCase.in)

			if diff := cmp.Diff(got, testCase.expected); diff != "" {
				t.Errorf("unexpected difference: %s", diff)
			}
		})
	}
}

func TestMoveResourceStateResponse(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		in       *tfprotov5.MoveResourceStateResponse
		expected *tfprotov6.MoveResourceStateResponse
	}{
		"nil": {
			in:       nil,
			expected: nil,
		},
		"all-valid-fields": {
			in: &tfprotov5.MoveResourceStateResponse{
				Diagnostics:    testTfprotov5Diagnostics,
				TargetPrivate:  testBytes,
				TargetState:    &testTfprotov5DynamicValue,
				TargetIdentity: &testTfprotov5ResourceIdentityData,
			},
			expected: &tfprotov6.MoveResourceStateResponse{
				Diagnostics:    testTfprotov6Diagnostics,
				TargetState:    &testTfprotov6DynamicValue,
				TargetPrivate:  testBytes,
				TargetIdentity: &testTfprotov6ResourceIdentityData,
			},
		},
	}

	for name, testCase := range testCases {

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got := tfprotov5tov6.MoveResourceStateResponse(testCase.in)

			if diff := cmp.Diff(got, testCase.expected); diff != "" {
				t.Errorf("unexpected difference: %s", diff)
			}
		})
	}
}

func TestOpenEphemeralResourceRequest(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		in       *tfprotov5.OpenEphemeralResourceRequest
		expected *tfprotov6.OpenEphemeralResourceRequest
	}{
		"nil": {
			in:       nil,
			expected: nil,
		},
		"all-valid-fields": {
			in: &tfprotov5.OpenEphemeralResourceRequest{
				Config:   &testTfprotov5DynamicValue,
				TypeName: "test_ephemeral_resource",
			},
			expected: &tfprotov6.OpenEphemeralResourceRequest{
				Config:   &testTfprotov6DynamicValue,
				TypeName: "test_ephemeral_resource",
			},
		},
		"client-capabilities-deferral-allowed": {
			in: &tfprotov5.OpenEphemeralResourceRequest{
				Config:   &testTfprotov5DynamicValue,
				TypeName: "test_ephemeral_resource",
				ClientCapabilities: &tfprotov5.OpenEphemeralResourceClientCapabilities{
					DeferralAllowed: true,
				},
			},
			expected: &tfprotov6.OpenEphemeralResourceRequest{
				Config:   &testTfprotov6DynamicValue,
				TypeName: "test_ephemeral_resource",
				ClientCapabilities: &tfprotov6.OpenEphemeralResourceClientCapabilities{
					DeferralAllowed: true,
				},
			},
		},
	}

	for name, testCase := range testCases {

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got := tfprotov5tov6.OpenEphemeralResourceRequest(testCase.in)

			if diff := cmp.Diff(got, testCase.expected); diff != "" {
				t.Errorf("unexpected difference: %s", diff)
			}
		})
	}
}

func TestOpenEphemeralResourceResponse(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		in       *tfprotov5.OpenEphemeralResourceResponse
		expected *tfprotov6.OpenEphemeralResourceResponse
	}{
		"nil": {
			in:       nil,
			expected: nil,
		},
		"all-valid-fields": {
			in: &tfprotov5.OpenEphemeralResourceResponse{
				Diagnostics: testTfprotov5Diagnostics,
				Private:     testBytes,
				RenewAt:     testTime,
				Result:      &testTfprotov5DynamicValue,
			},
			expected: &tfprotov6.OpenEphemeralResourceResponse{
				Diagnostics: testTfprotov6Diagnostics,
				Private:     testBytes,
				RenewAt:     testTime,
				Result:      &testTfprotov6DynamicValue,
			},
		},
		"deferred-reason": {
			in: &tfprotov5.OpenEphemeralResourceResponse{
				Diagnostics: testTfprotov5Diagnostics,
				Result:      &testTfprotov5DynamicValue,
				Deferred: &tfprotov5.Deferred{
					Reason: tfprotov5.DeferredReasonResourceConfigUnknown,
				},
			},
			expected: &tfprotov6.OpenEphemeralResourceResponse{
				Diagnostics: testTfprotov6Diagnostics,
				Result:      &testTfprotov6DynamicValue,
				Deferred: &tfprotov6.Deferred{
					Reason: tfprotov6.DeferredReasonResourceConfigUnknown,
				},
			},
		},
	}

	for name, testCase := range testCases {

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got := tfprotov5tov6.OpenEphemeralResourceResponse(testCase.in)

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
				PriorIdentity:    &testTfprotov5ResourceIdentityData,
			},
			expected: &tfprotov6.PlanResourceChangeRequest{
				Config:           &testTfprotov6DynamicValue,
				PriorPrivate:     testBytes,
				PriorState:       &testTfprotov6DynamicValue,
				ProposedNewState: &testTfprotov6DynamicValue,
				ProviderMeta:     &testTfprotov6DynamicValue,
				TypeName:         "test_resource",
				PriorIdentity:    &testTfprotov6ResourceIdentityData,
			},
		},
		"client-capabilities-deferral-allowed": {
			in: &tfprotov5.PlanResourceChangeRequest{
				Config:           &testTfprotov5DynamicValue,
				PriorPrivate:     testBytes,
				PriorState:       &testTfprotov5DynamicValue,
				ProposedNewState: &testTfprotov5DynamicValue,
				ProviderMeta:     &testTfprotov5DynamicValue,
				TypeName:         "test_resource",
				PriorIdentity:    &testTfprotov5ResourceIdentityData,
				ClientCapabilities: &tfprotov5.PlanResourceChangeClientCapabilities{
					DeferralAllowed: true,
				},
			},
			expected: &tfprotov6.PlanResourceChangeRequest{
				Config:           &testTfprotov6DynamicValue,
				PriorPrivate:     testBytes,
				PriorState:       &testTfprotov6DynamicValue,
				ProposedNewState: &testTfprotov6DynamicValue,
				ProviderMeta:     &testTfprotov6DynamicValue,
				TypeName:         "test_resource",
				PriorIdentity:    &testTfprotov6ResourceIdentityData,
				ClientCapabilities: &tfprotov6.PlanResourceChangeClientCapabilities{
					DeferralAllowed: true,
				},
			},
		},
	}

	for name, testCase := range testCases {

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
				PlannedIdentity:             &testTfprotov5ResourceIdentityData,
			},
			expected: &tfprotov6.PlanResourceChangeResponse{
				Diagnostics:    testTfprotov6Diagnostics,
				PlannedPrivate: testBytes,
				PlannedState:   &testTfprotov6DynamicValue,
				RequiresReplace: []*tftypes.AttributePath{
					tftypes.NewAttributePath().WithAttributeName("test"),
				},
				UnsafeToUseLegacyTypeSystem: true,
				PlannedIdentity:             &testTfprotov6ResourceIdentityData,
			},
		},
		"deferred-reason": {
			in: &tfprotov5.PlanResourceChangeResponse{
				Diagnostics:    testTfprotov5Diagnostics,
				PlannedPrivate: testBytes,
				PlannedState:   &testTfprotov5DynamicValue,
				RequiresReplace: []*tftypes.AttributePath{
					tftypes.NewAttributePath().WithAttributeName("test"),
				},
				UnsafeToUseLegacyTypeSystem: true,
				PlannedIdentity:             &testTfprotov5ResourceIdentityData,
				Deferred: &tfprotov5.Deferred{
					Reason: tfprotov5.DeferredReasonResourceConfigUnknown,
				},
			},
			expected: &tfprotov6.PlanResourceChangeResponse{
				Diagnostics:    testTfprotov6Diagnostics,
				PlannedPrivate: testBytes,
				PlannedState:   &testTfprotov6DynamicValue,
				RequiresReplace: []*tftypes.AttributePath{
					tftypes.NewAttributePath().WithAttributeName("test"),
				},
				UnsafeToUseLegacyTypeSystem: true,
				PlannedIdentity:             &testTfprotov6ResourceIdentityData,
				Deferred: &tfprotov6.Deferred{
					Reason: tfprotov6.DeferredReasonResourceConfigUnknown,
				},
			},
		},
	}

	for name, testCase := range testCases {

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
		"client-capabilities-deferral-allowed": {
			in: &tfprotov5.ReadDataSourceRequest{
				Config:       &testTfprotov5DynamicValue,
				ProviderMeta: &testTfprotov5DynamicValue,
				TypeName:     "test_data_source",
				ClientCapabilities: &tfprotov5.ReadDataSourceClientCapabilities{
					DeferralAllowed: true,
				},
			},
			expected: &tfprotov6.ReadDataSourceRequest{
				Config:       &testTfprotov6DynamicValue,
				ProviderMeta: &testTfprotov6DynamicValue,
				TypeName:     "test_data_source",
				ClientCapabilities: &tfprotov6.ReadDataSourceClientCapabilities{
					DeferralAllowed: true,
				},
			},
		},
	}

	for name, testCase := range testCases {

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
		"deferred-reason": {
			in: &tfprotov5.ReadDataSourceResponse{
				Diagnostics: testTfprotov5Diagnostics,
				State:       &testTfprotov5DynamicValue,
				Deferred: &tfprotov5.Deferred{
					Reason: tfprotov5.DeferredReasonResourceConfigUnknown,
				},
			},
			expected: &tfprotov6.ReadDataSourceResponse{
				Diagnostics: testTfprotov6Diagnostics,
				State:       &testTfprotov6DynamicValue,
				Deferred: &tfprotov6.Deferred{
					Reason: tfprotov6.DeferredReasonResourceConfigUnknown,
				},
			},
		},
	}

	for name, testCase := range testCases {

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
				CurrentState:    &testTfprotov5DynamicValue,
				Private:         testBytes,
				ProviderMeta:    &testTfprotov5DynamicValue,
				TypeName:        "test_resource",
				CurrentIdentity: &testTfprotov5ResourceIdentityData,
			},
			expected: &tfprotov6.ReadResourceRequest{
				CurrentState:    &testTfprotov6DynamicValue,
				Private:         testBytes,
				ProviderMeta:    &testTfprotov6DynamicValue,
				TypeName:        "test_resource",
				CurrentIdentity: &testTfprotov6ResourceIdentityData,
			},
		},
		"client-capabilities-deferral-allowed": {
			in: &tfprotov5.ReadResourceRequest{
				CurrentState:    &testTfprotov5DynamicValue,
				Private:         testBytes,
				ProviderMeta:    &testTfprotov5DynamicValue,
				TypeName:        "test_resource",
				CurrentIdentity: &testTfprotov5ResourceIdentityData,
				ClientCapabilities: &tfprotov5.ReadResourceClientCapabilities{
					DeferralAllowed: true,
				},
			},
			expected: &tfprotov6.ReadResourceRequest{
				CurrentState:    &testTfprotov6DynamicValue,
				Private:         testBytes,
				ProviderMeta:    &testTfprotov6DynamicValue,
				TypeName:        "test_resource",
				CurrentIdentity: &testTfprotov6ResourceIdentityData,
				ClientCapabilities: &tfprotov6.ReadResourceClientCapabilities{
					DeferralAllowed: true,
				},
			},
		},
	}

	for name, testCase := range testCases {

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
				NewIdentity: &testTfprotov5ResourceIdentityData,
			},
			expected: &tfprotov6.ReadResourceResponse{
				Diagnostics: testTfprotov6Diagnostics,
				NewState:    &testTfprotov6DynamicValue,
				Private:     testBytes,
				NewIdentity: &testTfprotov6ResourceIdentityData,
			},
		},
		"deferred-reason": {
			in: &tfprotov5.ReadResourceResponse{
				Diagnostics: testTfprotov5Diagnostics,
				NewState:    &testTfprotov5DynamicValue,
				Private:     testBytes,
				NewIdentity: &testTfprotov5ResourceIdentityData,
				Deferred: &tfprotov5.Deferred{
					Reason: tfprotov5.DeferredReasonResourceConfigUnknown,
				},
			},
			expected: &tfprotov6.ReadResourceResponse{
				Diagnostics: testTfprotov6Diagnostics,
				NewState:    &testTfprotov6DynamicValue,
				Private:     testBytes,
				NewIdentity: &testTfprotov6ResourceIdentityData,
				Deferred: &tfprotov6.Deferred{
					Reason: tfprotov6.DeferredReasonResourceConfigUnknown,
				},
			},
		},
	}

	for name, testCase := range testCases {

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got := tfprotov5tov6.ReadResourceResponse(testCase.in)

			if diff := cmp.Diff(got, testCase.expected); diff != "" {
				t.Errorf("unexpected difference: %s", diff)
			}
		})
	}
}

func TestRenewEphemeralResourceRequest(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		in       *tfprotov5.RenewEphemeralResourceRequest
		expected *tfprotov6.RenewEphemeralResourceRequest
	}{
		"nil": {
			in:       nil,
			expected: nil,
		},
		"all-valid-fields": {
			in: &tfprotov5.RenewEphemeralResourceRequest{
				Private:  testBytes,
				TypeName: "test_ephemeral_resource",
			},
			expected: &tfprotov6.RenewEphemeralResourceRequest{
				Private:  testBytes,
				TypeName: "test_ephemeral_resource",
			},
		},
	}

	for name, testCase := range testCases {

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got := tfprotov5tov6.RenewEphemeralResourceRequest(testCase.in)

			if diff := cmp.Diff(got, testCase.expected); diff != "" {
				t.Errorf("unexpected difference: %s", diff)
			}
		})
	}
}

func TestRenewEphemeralResourceResponse(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		in       *tfprotov5.RenewEphemeralResourceResponse
		expected *tfprotov6.RenewEphemeralResourceResponse
	}{
		"nil": {
			in:       nil,
			expected: nil,
		},
		"all-valid-fields": {
			in: &tfprotov5.RenewEphemeralResourceResponse{
				Diagnostics: testTfprotov5Diagnostics,
				Private:     testBytes,
				RenewAt:     testTime,
			},
			expected: &tfprotov6.RenewEphemeralResourceResponse{
				Diagnostics: testTfprotov6Diagnostics,
				Private:     testBytes,
				RenewAt:     testTime,
			},
		},
	}

	for name, testCase := range testCases {

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got := tfprotov5tov6.RenewEphemeralResourceResponse(testCase.in)

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
				WriteOnly:       true,
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
				WriteOnly:       true,
			},
		},
	}

	for name, testCase := range testCases {

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

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got := tfprotov5tov6.SchemaNestedBlock(testCase.in)

			if diff := cmp.Diff(got, testCase.expected); diff != "" {
				t.Errorf("unexpected difference: %s", diff)
			}
		})
	}
}

func TestResourceIdentitySchema(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		in       *tfprotov5.ResourceIdentitySchema
		expected *tfprotov6.ResourceIdentitySchema
	}{
		"nil": {
			in:       nil,
			expected: nil,
		},
		"all-valid-fields": {
			in:       testTfprotov5ResourceIdentitySchema,
			expected: testTfprotov6ResourceIdentitySchema,
		},
	}

	for name, testCase := range testCases {

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got := tfprotov5tov6.ResourceIdentitySchema(testCase.in)

			if diff := cmp.Diff(got, testCase.expected); diff != "" {
				t.Errorf("unexpected difference: %s", diff)
			}
		})
	}
}

func TestResourceIdentitySchemaAttribute(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		in       *tfprotov5.ResourceIdentitySchemaAttribute
		expected *tfprotov6.ResourceIdentitySchemaAttribute
	}{
		"nil": {
			in:       nil,
			expected: nil,
		},
		"all-valid-fields": {
			in: &tfprotov5.ResourceIdentitySchemaAttribute{
				Name:              "test",
				Description:       "test description",
				Type:              tftypes.String,
				RequiredForImport: true,
				OptionalForImport: true,
			},
			expected: &tfprotov6.ResourceIdentitySchemaAttribute{
				Name:              "test",
				Description:       "test description",
				Type:              tftypes.String,
				RequiredForImport: true,
				OptionalForImport: true,
			},
		},
	}

	for name, testCase := range testCases {

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got := tfprotov5tov6.ResourceIdentitySchemaAttribute(testCase.in)

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

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got := tfprotov5tov6.UpgradeResourceStateResponse(testCase.in)

			if diff := cmp.Diff(got, testCase.expected); diff != "" {
				t.Errorf("unexpected difference: %s", diff)
			}
		})
	}
}

func TestUpgradeResourceIdentityRequest(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		in       *tfprotov5.UpgradeResourceIdentityRequest
		expected *tfprotov6.UpgradeResourceIdentityRequest
	}{
		"nil": {
			in:       nil,
			expected: nil,
		},
		"all-valid-fields": {
			in: &tfprotov5.UpgradeResourceIdentityRequest{
				RawIdentity: &tfprotov5.RawState{
					JSON: testBytes,
				},
				TypeName: "test_resource",
				Version:  1,
			},
			expected: &tfprotov6.UpgradeResourceIdentityRequest{
				RawIdentity: &tfprotov6.RawState{
					JSON: testBytes,
				},
				TypeName: "test_resource",
				Version:  1,
			},
		},
	}

	for name, testCase := range testCases {

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got := tfprotov5tov6.UpgradeResourceIdentityRequest(testCase.in)

			if diff := cmp.Diff(got, testCase.expected); diff != "" {
				t.Errorf("unexpected difference: %s", diff)
			}
		})
	}
}

func TestUpgradeResourceIdentityResponse(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		in       *tfprotov5.UpgradeResourceIdentityResponse
		expected *tfprotov6.UpgradeResourceIdentityResponse
	}{
		"nil": {
			in:       nil,
			expected: nil,
		},
		"all-valid-fields": {
			in: &tfprotov5.UpgradeResourceIdentityResponse{
				Diagnostics:      testTfprotov5Diagnostics,
				UpgradedIdentity: &testTfprotov5ResourceIdentityData,
			},
			expected: &tfprotov6.UpgradeResourceIdentityResponse{
				Diagnostics:      testTfprotov6Diagnostics,
				UpgradedIdentity: &testTfprotov6ResourceIdentityData,
			},
		},
	}

	for name, testCase := range testCases {

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got := tfprotov5tov6.UpgradeResourceIdentityResponse(testCase.in)

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

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got := tfprotov5tov6.ValidateDataResourceConfigResponse(testCase.in)

			if diff := cmp.Diff(got, testCase.expected); diff != "" {
				t.Errorf("unexpected difference: %s", diff)
			}
		})
	}
}

func TestValidateEphemeralResourceConfigRequest(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		in       *tfprotov5.ValidateEphemeralResourceConfigRequest
		expected *tfprotov6.ValidateEphemeralResourceConfigRequest
	}{
		"nil": {
			in:       nil,
			expected: nil,
		},
		"all-valid-fields": {
			in: &tfprotov5.ValidateEphemeralResourceConfigRequest{
				Config:   &testTfprotov5DynamicValue,
				TypeName: "test_ephemeral_resource",
			},
			expected: &tfprotov6.ValidateEphemeralResourceConfigRequest{
				Config:   &testTfprotov6DynamicValue,
				TypeName: "test_ephemeral_resource",
			},
		},
	}

	for name, testCase := range testCases {

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got := tfprotov5tov6.ValidateEphemeralResourceConfigRequest(testCase.in)

			if diff := cmp.Diff(got, testCase.expected); diff != "" {
				t.Errorf("unexpected difference: %s", diff)
			}
		})
	}
}

func TestValidateEphemeralResourceConfigResponse(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		in       *tfprotov5.ValidateEphemeralResourceConfigResponse
		expected *tfprotov6.ValidateEphemeralResourceConfigResponse
	}{
		"nil": {
			in:       nil,
			expected: nil,
		},
		"all-valid-fields": {
			in: &tfprotov5.ValidateEphemeralResourceConfigResponse{
				Diagnostics: testTfprotov5Diagnostics,
			},
			expected: &tfprotov6.ValidateEphemeralResourceConfigResponse{
				Diagnostics: testTfprotov6Diagnostics,
			},
		},
	}

	for name, testCase := range testCases {

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got := tfprotov5tov6.ValidateEphemeralResourceConfigResponse(testCase.in)

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
		"client-capabilities-write-only-attributes-allowed": {
			in: &tfprotov5.ValidateResourceTypeConfigRequest{
				ClientCapabilities: &tfprotov5.ValidateResourceTypeConfigClientCapabilities{
					WriteOnlyAttributesAllowed: true,
				},
				Config:   &testTfprotov5DynamicValue,
				TypeName: "test_resource",
			},
			expected: &tfprotov6.ValidateResourceConfigRequest{
				ClientCapabilities: &tfprotov6.ValidateResourceConfigClientCapabilities{
					WriteOnlyAttributesAllowed: true,
				},
				Config:   &testTfprotov6DynamicValue,
				TypeName: "test_resource",
			},
		},
	}

	for name, testCase := range testCases {

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

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got := tfprotov5tov6.ValidateResourceConfigResponse(testCase.in)

			if diff := cmp.Diff(got, testCase.expected); diff != "" {
				t.Errorf("unexpected difference: %s", diff)
			}
		})
	}
}

func TestValidateListResourceConfigRequest(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		in       *tfprotov5.ValidateListResourceConfigRequest
		expected *tfprotov6.ValidateListResourceConfigRequest
	}{
		"nil": {
			in:       nil,
			expected: nil,
		},
		"all-valid-fields": {
			in: &tfprotov5.ValidateListResourceConfigRequest{
				Config:   &testTfprotov5DynamicValue,
				TypeName: "test_list_resource",
			},
			expected: &tfprotov6.ValidateListResourceConfigRequest{
				Config:   &testTfprotov6DynamicValue,
				TypeName: "test_list_resource",
			},
		},
	}

	for name, testCase := range testCases {

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got := tfprotov5tov6.ValidateListResourceConfigRequest(testCase.in)

			if diff := cmp.Diff(got, testCase.expected); diff != "" {
				t.Errorf("unexpected difference: %s", diff)
			}
		})
	}
}

func TestValidateListResourceConfigResponse(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		in       *tfprotov5.ValidateListResourceConfigResponse
		expected *tfprotov6.ValidateListResourceConfigResponse
	}{
		"nil": {
			in:       nil,
			expected: nil,
		},
		"all-valid-fields": {
			in: &tfprotov5.ValidateListResourceConfigResponse{
				Diagnostics: testTfprotov5Diagnostics,
			},
			expected: &tfprotov6.ValidateListResourceConfigResponse{
				Diagnostics: testTfprotov6Diagnostics,
			},
		},
	}

	for name, testCase := range testCases {

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got := tfprotov5tov6.ValidateListResourceConfigResponse(testCase.in)

			if diff := cmp.Diff(got, testCase.expected); diff != "" {
				t.Errorf("unexpected difference: %s", diff)
			}
		})
	}
}

func TestListResourceRequest(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		in       *tfprotov5.ListResourceRequest
		expected *tfprotov6.ListResourceRequest
	}{
		"nil": {
			in:       nil,
			expected: nil,
		},
		"all-valid-fields": {
			in: &tfprotov5.ListResourceRequest{
				Config:   &testTfprotov5DynamicValue,
				TypeName: "test_list_resource",
			},
			expected: &tfprotov6.ListResourceRequest{
				Config:   &testTfprotov6DynamicValue,
				TypeName: "test_list_resource",
			},
		},
	}

	for name, testCase := range testCases {

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got := tfprotov5tov6.ListResourceRequest(testCase.in)

			if diff := cmp.Diff(got, testCase.expected); diff != "" {
				t.Errorf("unexpected difference: %s", diff)
			}
		})
	}
}

func TestListResourceServerStream(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		in       *tfprotov5.ListResourceServerStream
		expected *tfprotov6.ListResourceServerStream
	}{
		"nil": {
			in:       nil,
			expected: nil,
		},
		"all-valid-fields": {
			in:       testStreamv5,
			expected: testStreamv6,
		},
	}

	for name, testCase := range testCases {

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got := tfprotov5tov6.ListResourceServerStream(testCase.in)

			if got == nil {
				if diff := cmp.Diff(got, testCase.expected); diff != "" {
					t.Errorf("unexpected difference: %s", diff)
				}
			} else {
				gotSlice := slices.Collect(got.Results)

				expectedSlice := slices.Collect(got.Results)

				if len(expectedSlice) != len(gotSlice) {
					t.Fatalf("expected iterator and result iterator lengths do not match")
				}

				if diff := cmp.Diff(gotSlice, expectedSlice); diff != "" {
					t.Errorf("unexpected difference: %s", diff)
				}
			}
		})
	}
}

func TestListResourceResult(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		in       tfprotov5.ListResourceResult
		expected tfprotov6.ListResourceResult
	}{
		"all-valid-fields": {
			in: tfprotov5.ListResourceResult{
				Diagnostics: testTfprotov5Diagnostics,
			},
			expected: tfprotov6.ListResourceResult{
				Diagnostics: testTfprotov6Diagnostics,
			},
		},
	}

	for name, testCase := range testCases {

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got := tfprotov5tov6.ListResourceResult(testCase.in)

			if diff := cmp.Diff(got, testCase.expected); diff != "" {
				t.Errorf("unexpected difference: %s", diff)
			}
		})
	}
}

func TestActionSchema(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		in       *tfprotov5.ActionSchema
		expected *tfprotov6.ActionSchema
	}{
		"nil": {
			in:       nil,
			expected: nil,
		},
		"unlinked": {
			in: &tfprotov5.ActionSchema{
				Schema: testTfprotov5Schema,
				Type:   tfprotov5.UnlinkedActionSchemaType{},
			},
			expected: &tfprotov6.ActionSchema{
				Schema: testTfprotov6Schema,
				Type:   tfprotov6.UnlinkedActionSchemaType{},
			},
		},
		"lifecycle": {
			in: &tfprotov5.ActionSchema{
				Schema: testTfprotov5Schema,
				Type: tfprotov5.LifecycleActionSchemaType{
					Executes: tfprotov5.LifecycleExecutionOrderAfter,
					LinkedResource: &tfprotov5.LinkedResourceSchema{
						TypeName:    "test_resource_linked_1",
						Description: "This is a linked resource.",
					},
				},
			},
			expected: &tfprotov6.ActionSchema{
				Schema: testTfprotov6Schema,
				Type: tfprotov6.LifecycleActionSchemaType{
					Executes: tfprotov6.LifecycleExecutionOrderAfter,
					LinkedResource: &tfprotov6.LinkedResourceSchema{
						TypeName:    "test_resource_linked_1",
						Description: "This is a linked resource.",
					},
				},
			},
		},
		"linked": {
			in: &tfprotov5.ActionSchema{
				Schema: testTfprotov5Schema,
				Type: tfprotov5.LinkedActionSchemaType{
					LinkedResources: []*tfprotov5.LinkedResourceSchema{
						{
							TypeName:    "test_resource_linked_1",
							Description: "This is a linked resource.",
						},
						{
							TypeName:    "test_resource_linked_2",
							Description: "This is also a linked resource.",
						},
					},
				},
			},
			expected: &tfprotov6.ActionSchema{
				Schema: testTfprotov6Schema,
				Type: tfprotov6.LinkedActionSchemaType{
					LinkedResources: []*tfprotov6.LinkedResourceSchema{
						{
							TypeName:    "test_resource_linked_1",
							Description: "This is a linked resource.",
						},
						{
							TypeName:    "test_resource_linked_2",
							Description: "This is also a linked resource.",
						},
					},
				},
			},
		},
	}

	for name, testCase := range testCases {

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got := tfprotov5tov6.ActionSchema(testCase.in)

			if diff := cmp.Diff(got, testCase.expected); diff != "" {
				t.Errorf("unexpected difference: %s", diff)
			}
		})
	}
}

func TestValidateActionConfigRequest(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		in       *tfprotov5.ValidateActionConfigRequest
		expected *tfprotov6.ValidateActionConfigRequest
	}{
		"nil": {
			in:       nil,
			expected: nil,
		},
		"all-valid-fields": {
			in: &tfprotov5.ValidateActionConfigRequest{
				Config:     &testTfprotov5DynamicValue,
				ActionType: "test_action",
			},
			expected: &tfprotov6.ValidateActionConfigRequest{
				Config:     &testTfprotov6DynamicValue,
				ActionType: "test_action",
			},
		},
	}

	for name, testCase := range testCases {

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got := tfprotov5tov6.ValidateActionConfigRequest(testCase.in)

			if diff := cmp.Diff(got, testCase.expected); diff != "" {
				t.Errorf("unexpected difference: %s", diff)
			}
		})
	}
}

func TestValidateActionConfigResponse(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		in       *tfprotov5.ValidateActionConfigResponse
		expected *tfprotov6.ValidateActionConfigResponse
	}{
		"nil": {
			in:       nil,
			expected: nil,
		},
		"all-valid-fields": {
			in: &tfprotov5.ValidateActionConfigResponse{
				Diagnostics: testTfprotov5Diagnostics,
			},
			expected: &tfprotov6.ValidateActionConfigResponse{
				Diagnostics: testTfprotov6Diagnostics,
			},
		},
	}

	for name, testCase := range testCases {

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got := tfprotov5tov6.ValidateActionConfigResponse(testCase.in)

			if diff := cmp.Diff(got, testCase.expected); diff != "" {
				t.Errorf("unexpected difference: %s", diff)
			}
		})
	}
}

func TestPlanActionRequest(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		in       *tfprotov5.PlanActionRequest
		expected *tfprotov6.PlanActionRequest
	}{
		"nil": {
			in:       nil,
			expected: nil,
		},
		"no-linked-resources": {
			in: &tfprotov5.PlanActionRequest{
				ActionType: "test_action",
				Config:     &testTfprotov5DynamicValue,
			},
			expected: &tfprotov6.PlanActionRequest{
				ActionType: "test_action",
				Config:     &testTfprotov6DynamicValue,
			},
		},
		"linked-resources": {
			in: &tfprotov5.PlanActionRequest{
				ActionType: "test_action",
				Config:     &testTfprotov5DynamicValue,
				LinkedResources: []*tfprotov5.ProposedLinkedResource{
					{
						PriorState:    &testTfprotov5DynamicValue,
						PlannedState:  &testTfprotov5DynamicValue,
						Config:        &testTfprotov5DynamicValue,
						PriorIdentity: &testTfprotov5ResourceIdentityData,
					},
					{
						PriorState:    &testTfprotov5DynamicValue,
						PlannedState:  &testTfprotov5DynamicValue,
						Config:        &testTfprotov5DynamicValue,
						PriorIdentity: &testTfprotov5ResourceIdentityData,
					},
				},
			},
			expected: &tfprotov6.PlanActionRequest{
				ActionType: "test_action",
				Config:     &testTfprotov6DynamicValue,
				LinkedResources: []*tfprotov6.ProposedLinkedResource{
					{
						PriorState:    &testTfprotov6DynamicValue,
						PlannedState:  &testTfprotov6DynamicValue,
						Config:        &testTfprotov6DynamicValue,
						PriorIdentity: &testTfprotov6ResourceIdentityData,
					},
					{
						PriorState:    &testTfprotov6DynamicValue,
						PlannedState:  &testTfprotov6DynamicValue,
						Config:        &testTfprotov6DynamicValue,
						PriorIdentity: &testTfprotov6ResourceIdentityData,
					},
				},
			},
		},
		"client-capabilities-deferral-allowed": {
			in: &tfprotov5.PlanActionRequest{
				ActionType: "test_action",
				Config:     &testTfprotov5DynamicValue,
				ClientCapabilities: &tfprotov5.PlanActionClientCapabilities{
					DeferralAllowed: true,
				},
			},
			expected: &tfprotov6.PlanActionRequest{
				ActionType: "test_action",
				Config:     &testTfprotov6DynamicValue,
				ClientCapabilities: &tfprotov6.PlanActionClientCapabilities{
					DeferralAllowed: true,
				},
			},
		},
	}

	for name, testCase := range testCases {

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got := tfprotov5tov6.PlanActionRequest(testCase.in)

			if diff := cmp.Diff(got, testCase.expected); diff != "" {
				t.Errorf("unexpected difference: %s", diff)
			}
		})
	}
}

func TestPlanActionResponse(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		in       *tfprotov5.PlanActionResponse
		expected *tfprotov6.PlanActionResponse
	}{
		"nil": {
			in:       nil,
			expected: nil,
		},
		"no-linked-resources": {
			in: &tfprotov5.PlanActionResponse{
				Diagnostics: testTfprotov5Diagnostics,
			},
			expected: &tfprotov6.PlanActionResponse{
				Diagnostics: testTfprotov6Diagnostics,
			},
		},
		"linked-resources": {
			in: &tfprotov5.PlanActionResponse{
				Diagnostics: testTfprotov5Diagnostics,
				LinkedResources: []*tfprotov5.PlannedLinkedResource{
					{
						PlannedState:    &testTfprotov5DynamicValue,
						PlannedIdentity: &testTfprotov5ResourceIdentityData,
					},
					{
						PlannedState:    &testTfprotov5DynamicValue,
						PlannedIdentity: &testTfprotov5ResourceIdentityData,
					},
				},
			},
			expected: &tfprotov6.PlanActionResponse{
				Diagnostics: testTfprotov6Diagnostics,
				LinkedResources: []*tfprotov6.PlannedLinkedResource{
					{
						PlannedState:    &testTfprotov6DynamicValue,
						PlannedIdentity: &testTfprotov6ResourceIdentityData,
					},
					{
						PlannedState:    &testTfprotov6DynamicValue,
						PlannedIdentity: &testTfprotov6ResourceIdentityData,
					},
				},
			},
		},
		"deferred-reason": {
			in: &tfprotov5.PlanActionResponse{
				Diagnostics: testTfprotov5Diagnostics,
				Deferred: &tfprotov5.Deferred{
					Reason: tfprotov5.DeferredReasonResourceConfigUnknown,
				},
			},
			expected: &tfprotov6.PlanActionResponse{
				Diagnostics: testTfprotov6Diagnostics,
				Deferred: &tfprotov6.Deferred{
					Reason: tfprotov6.DeferredReasonResourceConfigUnknown,
				},
			},
		},
	}

	for name, testCase := range testCases {

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got := tfprotov5tov6.PlanActionResponse(testCase.in)

			if diff := cmp.Diff(got, testCase.expected); diff != "" {
				t.Errorf("unexpected difference: %s", diff)
			}
		})
	}
}

func TestInvokeActionRequest(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		in       *tfprotov5.InvokeActionRequest
		expected *tfprotov6.InvokeActionRequest
	}{
		"nil": {
			in:       nil,
			expected: nil,
		},
		"no-linked-resources": {
			in: &tfprotov5.InvokeActionRequest{
				ActionType: "test_action",
				Config:     &testTfprotov5DynamicValue,
			},
			expected: &tfprotov6.InvokeActionRequest{
				ActionType: "test_action",
				Config:     &testTfprotov6DynamicValue,
			},
		},
		"linked-resources": {
			in: &tfprotov5.InvokeActionRequest{
				ActionType: "test_action",
				Config:     &testTfprotov5DynamicValue,
				LinkedResources: []*tfprotov5.InvokeLinkedResource{
					{
						PriorState:      &testTfprotov5DynamicValue,
						PlannedState:    &testTfprotov5DynamicValue,
						Config:          &testTfprotov5DynamicValue,
						PlannedIdentity: &testTfprotov5ResourceIdentityData,
					},
					{
						PriorState:      &testTfprotov5DynamicValue,
						PlannedState:    &testTfprotov5DynamicValue,
						Config:          &testTfprotov5DynamicValue,
						PlannedIdentity: &testTfprotov5ResourceIdentityData,
					},
				},
			},
			expected: &tfprotov6.InvokeActionRequest{
				ActionType: "test_action",
				Config:     &testTfprotov6DynamicValue,
				LinkedResources: []*tfprotov6.InvokeLinkedResource{
					{
						PriorState:      &testTfprotov6DynamicValue,
						PlannedState:    &testTfprotov6DynamicValue,
						Config:          &testTfprotov6DynamicValue,
						PlannedIdentity: &testTfprotov6ResourceIdentityData,
					},
					{
						PriorState:      &testTfprotov6DynamicValue,
						PlannedState:    &testTfprotov6DynamicValue,
						Config:          &testTfprotov6DynamicValue,
						PlannedIdentity: &testTfprotov6ResourceIdentityData,
					},
				},
			},
		},
	}

	for name, testCase := range testCases {

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got := tfprotov5tov6.InvokeActionRequest(testCase.in)

			if diff := cmp.Diff(got, testCase.expected); diff != "" {
				t.Errorf("unexpected difference: %s", diff)
			}
		})
	}
}

func TestInvokeActionServerStream(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		in       *tfprotov5.InvokeActionServerStream
		expected *tfprotov6.InvokeActionServerStream
	}{
		"nil": {
			in:       nil,
			expected: nil,
		},
		"all-valid-fields": {
			in: &tfprotov5.InvokeActionServerStream{
				Events: slices.Values([]tfprotov5.InvokeActionEvent{
					{
						Type: tfprotov5.ProgressInvokeActionEventType{
							Message: "in progress",
						},
					},
					{
						Type: tfprotov5.CompletedInvokeActionEventType{
							LinkedResources: []*tfprotov5.NewLinkedResource{
								{
									NewState:        &testTfprotov5DynamicValue,
									NewIdentity:     &testTfprotov5ResourceIdentityData,
									RequiresReplace: true,
								},
							},
							Diagnostics: testTfprotov5Diagnostics,
						},
					},
				}),
			},
			expected: &tfprotov6.InvokeActionServerStream{
				Events: slices.Values([]tfprotov6.InvokeActionEvent{
					{
						Type: tfprotov6.ProgressInvokeActionEventType{
							Message: "in progress",
						},
					},
					{
						Type: tfprotov6.CompletedInvokeActionEventType{
							LinkedResources: []*tfprotov6.NewLinkedResource{
								{
									NewState:        &testTfprotov6DynamicValue,
									NewIdentity:     &testTfprotov6ResourceIdentityData,
									RequiresReplace: true,
								},
							},
							Diagnostics: testTfprotov6Diagnostics,
						},
					},
				}),
			},
		},
	}

	for name, testCase := range testCases {

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got := tfprotov5tov6.InvokeActionServerStream(testCase.in)

			if got == nil {
				if diff := cmp.Diff(got, testCase.expected); diff != "" {
					t.Errorf("unexpected difference: %s", diff)
				}
			} else {
				gotSlice := slices.Collect(got.Events)

				expectedSlice := slices.Collect(got.Events)

				if len(expectedSlice) != len(gotSlice) {
					t.Fatalf("expected iterator and event iterator lengths do not match")
				}

				if diff := cmp.Diff(gotSlice, expectedSlice); diff != "" {
					t.Errorf("unexpected difference: %s", diff)
				}
			}
		})
	}
}

func pointer[T any](value T) *T {
	return &value
}
