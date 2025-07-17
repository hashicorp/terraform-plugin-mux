// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package tfprotov6tov5_test

import (
	"fmt"
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/terraform-plugin-go/tfprotov5"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-go/tftypes"

	"github.com/hashicorp/terraform-plugin-mux/internal/tfprotov6tov5"
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

	testTfprotov5FunctionMetadata tfprotov5.FunctionMetadata = tfprotov5.FunctionMetadata{
		Name: "test_function",
	}

	testTfprotov6Function *tfprotov6.Function = &tfprotov6.Function{
		Parameters: []*tfprotov6.FunctionParameter{},
		Return: &tfprotov6.FunctionReturn{
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
		in       *tfprotov6.ApplyResourceChangeRequest
		expected *tfprotov5.ApplyResourceChangeRequest
	}{
		"nil": {
			in:       nil,
			expected: nil,
		},
		"all-valid-fields": {
			in: &tfprotov6.ApplyResourceChangeRequest{
				Config:          &testTfprotov6DynamicValue,
				PlannedPrivate:  testBytes,
				PlannedState:    &testTfprotov6DynamicValue,
				PriorState:      &testTfprotov6DynamicValue,
				ProviderMeta:    &testTfprotov6DynamicValue,
				TypeName:        "test",
				PlannedIdentity: &testTfprotov6ResourceIdentityData,
			},
			expected: &tfprotov5.ApplyResourceChangeRequest{
				Config:          &testTfprotov5DynamicValue,
				PlannedPrivate:  testBytes,
				PlannedState:    &testTfprotov5DynamicValue,
				PriorState:      &testTfprotov5DynamicValue,
				ProviderMeta:    &testTfprotov5DynamicValue,
				TypeName:        "test",
				PlannedIdentity: &testTfprotov5ResourceIdentityData,
			},
		},
	}

	for name, testCase := range testCases {

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
				NewIdentity:                 &testTfprotov6ResourceIdentityData,
			},
			expected: &tfprotov5.ApplyResourceChangeResponse{
				Diagnostics:                 testTfprotov5Diagnostics,
				NewState:                    &testTfprotov5DynamicValue,
				Private:                     testBytes,
				UnsafeToUseLegacyTypeSystem: true,
				NewIdentity:                 &testTfprotov5ResourceIdentityData,
			},
		},
	}

	for name, testCase := range testCases {

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got := tfprotov6tov5.ApplyResourceChangeResponse(testCase.in)

			if diff := cmp.Diff(got, testCase.expected); diff != "" {
				t.Errorf("unexpected difference: %s", diff)
			}
		})
	}
}

func TestCallFunctionRequest(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		in       *tfprotov6.CallFunctionRequest
		expected *tfprotov5.CallFunctionRequest
	}{
		"nil": {
			in:       nil,
			expected: nil,
		},
		"all-valid-fields": {
			in: &tfprotov6.CallFunctionRequest{
				Arguments: []*tfprotov6.DynamicValue{
					&testTfprotov6DynamicValue,
				},
				Name: "test_function",
			},
			expected: &tfprotov5.CallFunctionRequest{
				Arguments: []*tfprotov5.DynamicValue{
					&testTfprotov5DynamicValue,
				},
				Name: "test_function",
			},
		},
	}

	for name, testCase := range testCases {

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got := tfprotov6tov5.CallFunctionRequest(testCase.in)

			if diff := cmp.Diff(got, testCase.expected); diff != "" {
				t.Errorf("unexpected difference: %s", diff)
			}
		})
	}
}

func TestCallFunctionResponse(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		in       *tfprotov6.CallFunctionResponse
		expected *tfprotov5.CallFunctionResponse
	}{
		"nil": {
			in:       nil,
			expected: nil,
		},
		"all-valid-fields": {
			in: &tfprotov6.CallFunctionResponse{
				Error:  testTfprotov6FunctionError,
				Result: &testTfprotov6DynamicValue,
			},
			expected: &tfprotov5.CallFunctionResponse{
				Error:  testTfprotov5FunctionError,
				Result: &testTfprotov5DynamicValue,
			},
		},
	}

	for name, testCase := range testCases {

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got := tfprotov6tov5.CallFunctionResponse(testCase.in)

			if diff := cmp.Diff(got, testCase.expected); diff != "" {
				t.Errorf("unexpected difference: %s", diff)
			}
		})
	}
}

func TestCloseEphemeralResourceRequest(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		in       *tfprotov6.CloseEphemeralResourceRequest
		expected *tfprotov5.CloseEphemeralResourceRequest
	}{
		"nil": {
			in:       nil,
			expected: nil,
		},
		"all-valid-fields": {
			in: &tfprotov6.CloseEphemeralResourceRequest{
				Private:  testBytes,
				TypeName: "test_ephemeral_resource",
			},
			expected: &tfprotov5.CloseEphemeralResourceRequest{
				Private:  testBytes,
				TypeName: "test_ephemeral_resource",
			},
		},
	}

	for name, testCase := range testCases {

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got := tfprotov6tov5.CloseEphemeralResourceRequest(testCase.in)

			if diff := cmp.Diff(got, testCase.expected); diff != "" {
				t.Errorf("unexpected difference: %s", diff)
			}
		})
	}
}

func TestCloseEphemeralResourceResponse(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		in       *tfprotov6.CloseEphemeralResourceResponse
		expected *tfprotov5.CloseEphemeralResourceResponse
	}{
		"nil": {
			in:       nil,
			expected: nil,
		},
		"all-valid-fields": {
			in: &tfprotov6.CloseEphemeralResourceResponse{
				Diagnostics: testTfprotov6Diagnostics,
			},
			expected: &tfprotov5.CloseEphemeralResourceResponse{
				Diagnostics: testTfprotov5Diagnostics,
			},
		},
	}

	for name, testCase := range testCases {

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got := tfprotov6tov5.CloseEphemeralResourceResponse(testCase.in)

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
		"client-capabilities-deferral-allowed": {
			in: &tfprotov6.ConfigureProviderRequest{
				Config:           &testTfprotov6DynamicValue,
				TerraformVersion: "1.2.3",
				ClientCapabilities: &tfprotov6.ConfigureProviderClientCapabilities{
					DeferralAllowed: true,
				},
			},
			expected: &tfprotov5.ConfigureProviderRequest{
				Config:           &testTfprotov5DynamicValue,
				TerraformVersion: "1.2.3",
				ClientCapabilities: &tfprotov5.ConfigureProviderClientCapabilities{
					DeferralAllowed: true,
				},
			},
		},
	}

	for name, testCase := range testCases {

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

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got := tfprotov6tov5.DynamicValue(testCase.in)

			if diff := cmp.Diff(got, testCase.expected); diff != "" {
				t.Errorf("unexpected difference: %s", diff)
			}
		})
	}
}

func TestResourceIdentityData(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		in       *tfprotov6.ResourceIdentityData
		expected *tfprotov5.ResourceIdentityData
	}{
		"nil": {
			in:       nil,
			expected: nil,
		},
		"all-valid-fields": {
			in: &tfprotov6.ResourceIdentityData{
				IdentityData: &testTfprotov6DynamicValue,
			},
			expected: &tfprotov5.ResourceIdentityData{
				IdentityData: &testTfprotov5DynamicValue,
			},
		},
	}

	for name, testCase := range testCases {

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got := tfprotov6tov5.ResourceIdentityData(testCase.in)

			if diff := cmp.Diff(got, testCase.expected); diff != "" {
				t.Errorf("unexpected difference: %s", diff)
			}
		})
	}
}

func TestFunction(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		in       *tfprotov6.Function
		expected *tfprotov5.Function
	}{
		"nil": {
			in:       nil,
			expected: nil,
		},
		"all-valid-fields": {
			in: &tfprotov6.Function{
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
			expected: &tfprotov5.Function{
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
		},
	}

	for name, testCase := range testCases {

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got := tfprotov6tov5.Function(testCase.in)

			if diff := cmp.Diff(got, testCase.expected); diff != "" {
				t.Errorf("unexpected difference: %s", diff)
			}
		})
	}
}

func TestFunctionMetadata(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		in       tfprotov6.FunctionMetadata
		expected tfprotov5.FunctionMetadata
	}{
		"all-valid-fields": {
			in: tfprotov6.FunctionMetadata{
				Name: "test_function",
			},
			expected: tfprotov5.FunctionMetadata{
				Name: "test_function",
			},
		},
	}

	for name, testCase := range testCases {

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got := tfprotov6tov5.FunctionMetadata(testCase.in)

			if diff := cmp.Diff(got, testCase.expected); diff != "" {
				t.Errorf("unexpected difference: %s", diff)
			}
		})
	}
}

func TestFunctionParameter(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		in       *tfprotov6.FunctionParameter
		expected *tfprotov5.FunctionParameter
	}{
		"nil": {
			in:       nil,
			expected: nil,
		},
		"all-valid-fields": {
			in: &tfprotov6.FunctionParameter{
				Description:     "test description",
				DescriptionKind: tfprotov6.StringKindPlain,
				Type:            tftypes.String,
			},
			expected: &tfprotov5.FunctionParameter{
				Description:     "test description",
				DescriptionKind: tfprotov5.StringKindPlain,
				Type:            tftypes.String,
			},
		},
	}

	for name, testCase := range testCases {

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got := tfprotov6tov5.FunctionParameter(testCase.in)

			if diff := cmp.Diff(got, testCase.expected); diff != "" {
				t.Errorf("unexpected difference: %s", diff)
			}
		})
	}
}

func TestFunctionReturn(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		in       *tfprotov6.FunctionReturn
		expected *tfprotov5.FunctionReturn
	}{
		"nil": {
			in:       nil,
			expected: nil,
		},
		"all-valid-fields": {
			in: &tfprotov6.FunctionReturn{
				Type: tftypes.String,
			},
			expected: &tfprotov5.FunctionReturn{
				Type: tftypes.String,
			},
		},
	}

	for name, testCase := range testCases {

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got := tfprotov6tov5.FunctionReturn(testCase.in)

			if diff := cmp.Diff(got, testCase.expected); diff != "" {
				t.Errorf("unexpected difference: %s", diff)
			}
		})
	}
}

func TestGetFunctionsRequest(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		in       *tfprotov6.GetFunctionsRequest
		expected *tfprotov5.GetFunctionsRequest
	}{
		"nil": {
			in:       nil,
			expected: nil,
		},
		"all-valid-fields": {
			in:       &tfprotov6.GetFunctionsRequest{},
			expected: &tfprotov5.GetFunctionsRequest{},
		},
	}

	for name, testCase := range testCases {

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got := tfprotov6tov5.GetFunctionsRequest(testCase.in)

			if diff := cmp.Diff(got, testCase.expected); diff != "" {
				t.Errorf("unexpected difference: %s", diff)
			}
		})
	}
}

func TestGetFunctionsResponse(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		in       *tfprotov6.GetFunctionsResponse
		expected *tfprotov5.GetFunctionsResponse
	}{
		"nil": {
			in:       nil,
			expected: nil,
		},
		"all-valid-fields": {
			in: &tfprotov6.GetFunctionsResponse{
				Diagnostics: testTfprotov6Diagnostics,
				Functions: map[string]*tfprotov6.Function{
					"test_function": testTfprotov6Function,
				},
			},
			expected: &tfprotov5.GetFunctionsResponse{
				Diagnostics: testTfprotov5Diagnostics,
				Functions: map[string]*tfprotov5.Function{
					"test_function": testTfprotov5Function,
				},
			},
		},
	}

	for name, testCase := range testCases {

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got := tfprotov6tov5.GetFunctionsResponse(testCase.in)

			if diff := cmp.Diff(got, testCase.expected); diff != "" {
				t.Errorf("unexpected difference: %s", diff)
			}
		})
	}
}

func TestGetMetadataRequest(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		in       *tfprotov6.GetMetadataRequest
		expected *tfprotov5.GetMetadataRequest
	}{
		"nil": {
			in:       nil,
			expected: nil,
		},
		"all-valid-fields": {
			in:       &tfprotov6.GetMetadataRequest{},
			expected: &tfprotov5.GetMetadataRequest{},
		},
	}

	for name, testCase := range testCases {

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got := tfprotov6tov5.GetMetadataRequest(testCase.in)

			if diff := cmp.Diff(got, testCase.expected); diff != "" {
				t.Errorf("unexpected difference: %s", diff)
			}
		})
	}
}

func TestGetMetadataResponse(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		in       *tfprotov6.GetMetadataResponse
		expected *tfprotov5.GetMetadataResponse
	}{
		"nil": {
			in:       nil,
			expected: nil,
		},
		"all-valid-fields": {
			in: &tfprotov6.GetMetadataResponse{
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
			expected: &tfprotov5.GetMetadataResponse{
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
		},
	}

	for name, testCase := range testCases {

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got := tfprotov6tov5.GetMetadataResponse(testCase.in)

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
			expected: &tfprotov5.GetProviderSchemaResponse{
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
		"ephemeral-resource-nested-attribute-error": {
			in: &tfprotov6.GetProviderSchemaResponse{
				EphemeralResourceSchemas: map[string]*tfprotov6.Schema{
					"test_ephemeral_resource": {
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
			expectedError: fmt.Errorf("unable to convert ephemeral resource \"test_ephemeral_resource\" schema: unable to convert attribute \"test_attribute\" schema: %w", tfprotov6tov5.ErrSchemaAttributeNestedTypeNotImplemented),
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
		"action-nested-attribute-error": {
			in: &tfprotov6.GetProviderSchemaResponse{
				ActionSchemas: map[string]*tfprotov6.ActionSchema{
					"test_action": {
						Schema: &tfprotov6.Schema{
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
						Type: tfprotov6.UnlinkedActionSchemaType{},
					},
				},
			},
			expected:      nil,
			expectedError: fmt.Errorf("unable to convert action \"test_action\" schema: unable to convert attribute \"test_attribute\" schema: %w", tfprotov6tov5.ErrSchemaAttributeNestedTypeNotImplemented),
		},
	}

	for name, testCase := range testCases {

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

func TestGetResourceIdentitySchemasRequest(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		in       *tfprotov6.GetResourceIdentitySchemasRequest
		expected *tfprotov5.GetResourceIdentitySchemasRequest
	}{
		"nil": {
			in:       nil,
			expected: nil,
		},
		"all-valid-fields": {
			in:       &tfprotov6.GetResourceIdentitySchemasRequest{},
			expected: &tfprotov5.GetResourceIdentitySchemasRequest{},
		},
	}

	for name, testCase := range testCases {

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got := tfprotov6tov5.GetResourceIdentitySchemasRequest(testCase.in)

			if diff := cmp.Diff(got, testCase.expected); diff != "" {
				t.Errorf("unexpected difference: %s", diff)
			}
		})
	}
}

func TestGetResourceIdentitySchemasResponse(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		in       *tfprotov6.GetResourceIdentitySchemasResponse
		expected *tfprotov5.GetResourceIdentitySchemasResponse
	}{
		"nil": {
			in:       nil,
			expected: nil,
		},
		"all-valid-fields": {
			in: &tfprotov6.GetResourceIdentitySchemasResponse{
				Diagnostics: testTfprotov6Diagnostics,
				IdentitySchemas: map[string]*tfprotov6.ResourceIdentitySchema{
					"test_resource": testTfprotov6ResourceIdentitySchema,
				},
			},
			expected: &tfprotov5.GetResourceIdentitySchemasResponse{
				Diagnostics: testTfprotov5Diagnostics,
				IdentitySchemas: map[string]*tfprotov5.ResourceIdentitySchema{
					"test_resource": testTfprotov5ResourceIdentitySchema,
				},
			},
		},
	}

	for name, testCase := range testCases {

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got := tfprotov6tov5.GetResourceIdentitySchemasResponse(testCase.in)

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
				Identity: &testTfprotov6ResourceIdentityData,
			},
			expected: &tfprotov5.ImportResourceStateRequest{
				ID:       "test-id",
				TypeName: "test_resource",
				Identity: &testTfprotov5ResourceIdentityData,
			},
		},
		"client-capabilities-deferral-allowed": {
			in: &tfprotov6.ImportResourceStateRequest{
				ID:       "test-id",
				TypeName: "test_resource",
				ClientCapabilities: &tfprotov6.ImportResourceStateClientCapabilities{
					DeferralAllowed: true,
				},
			},
			expected: &tfprotov5.ImportResourceStateRequest{
				ID:       "test-id",
				TypeName: "test_resource",
				ClientCapabilities: &tfprotov5.ImportResourceStateClientCapabilities{
					DeferralAllowed: true,
				},
			},
		},
	}

	for name, testCase := range testCases {

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
						Identity: &testTfprotov6ResourceIdentityData,
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
						Identity: &testTfprotov5ResourceIdentityData,
					},
				},
			},
		},
		"deferred-reason": {
			in: &tfprotov6.ImportResourceStateResponse{
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
			expected: &tfprotov5.ImportResourceStateResponse{
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
		},
	}

	for name, testCase := range testCases {

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
					Identity: &testTfprotov6ResourceIdentityData,
				},
			},
			expected: []*tfprotov5.ImportedResource{
				{
					Private:  testBytes,
					State:    &testTfprotov5DynamicValue,
					TypeName: "test_resource1",
					Identity: &testTfprotov5ResourceIdentityData,
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

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got := tfprotov6tov5.ImportedResources(testCase.in)

			if diff := cmp.Diff(got, testCase.expected); diff != "" {
				t.Errorf("unexpected difference: %s", diff)
			}
		})
	}
}

func TestMoveResourceStateRequest(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		in       *tfprotov6.MoveResourceStateRequest
		expected *tfprotov5.MoveResourceStateRequest
	}{
		"nil": {
			in:       nil,
			expected: nil,
		},
		"all-valid-fields": {
			in: &tfprotov6.MoveResourceStateRequest{
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
			expected: &tfprotov5.MoveResourceStateRequest{
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
		},
	}

	for name, testCase := range testCases {

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got := tfprotov6tov5.MoveResourceStateRequest(testCase.in)

			if diff := cmp.Diff(got, testCase.expected); diff != "" {
				t.Errorf("unexpected difference: %s", diff)
			}
		})
	}
}

func TestMoveResourceStateResponse(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		in       *tfprotov6.MoveResourceStateResponse
		expected *tfprotov5.MoveResourceStateResponse
	}{
		"nil": {
			in:       nil,
			expected: nil,
		},
		"all-valid-fields": {
			in: &tfprotov6.MoveResourceStateResponse{
				Diagnostics:    testTfprotov6Diagnostics,
				TargetPrivate:  testBytes,
				TargetState:    &testTfprotov6DynamicValue,
				TargetIdentity: &testTfprotov6ResourceIdentityData,
			},
			expected: &tfprotov5.MoveResourceStateResponse{
				Diagnostics:    testTfprotov5Diagnostics,
				TargetState:    &testTfprotov5DynamicValue,
				TargetPrivate:  testBytes,
				TargetIdentity: &testTfprotov5ResourceIdentityData,
			},
		},
	}

	for name, testCase := range testCases {

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got := tfprotov6tov5.MoveResourceStateResponse(testCase.in)

			if diff := cmp.Diff(got, testCase.expected); diff != "" {
				t.Errorf("unexpected difference: %s", diff)
			}
		})
	}
}

func TestOpenEphemeralResourceRequest(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		in       *tfprotov6.OpenEphemeralResourceRequest
		expected *tfprotov5.OpenEphemeralResourceRequest
	}{
		"nil": {
			in:       nil,
			expected: nil,
		},
		"all-valid-fields": {
			in: &tfprotov6.OpenEphemeralResourceRequest{
				Config:   &testTfprotov6DynamicValue,
				TypeName: "test_ephemeral_resource",
			},
			expected: &tfprotov5.OpenEphemeralResourceRequest{
				Config:   &testTfprotov5DynamicValue,
				TypeName: "test_ephemeral_resource",
			},
		},
		"client-capabilities-deferral-allowed": {
			in: &tfprotov6.OpenEphemeralResourceRequest{
				Config:   &testTfprotov6DynamicValue,
				TypeName: "test_ephemeral_resource",
				ClientCapabilities: &tfprotov6.OpenEphemeralResourceClientCapabilities{
					DeferralAllowed: true,
				},
			},
			expected: &tfprotov5.OpenEphemeralResourceRequest{
				Config:   &testTfprotov5DynamicValue,
				TypeName: "test_ephemeral_resource",
				ClientCapabilities: &tfprotov5.OpenEphemeralResourceClientCapabilities{
					DeferralAllowed: true,
				},
			},
		},
	}

	for name, testCase := range testCases {

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got := tfprotov6tov5.OpenEphemeralResourceRequest(testCase.in)

			if diff := cmp.Diff(got, testCase.expected); diff != "" {
				t.Errorf("unexpected difference: %s", diff)
			}
		})
	}
}

func TestOpenEphemeralResourceResponse(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		in       *tfprotov6.OpenEphemeralResourceResponse
		expected *tfprotov5.OpenEphemeralResourceResponse
	}{
		"nil": {
			in:       nil,
			expected: nil,
		},
		"all-valid-fields": {
			in: &tfprotov6.OpenEphemeralResourceResponse{
				Diagnostics: testTfprotov6Diagnostics,
				Private:     testBytes,
				RenewAt:     testTime,
				Result:      &testTfprotov6DynamicValue,
			},
			expected: &tfprotov5.OpenEphemeralResourceResponse{
				Diagnostics: testTfprotov5Diagnostics,
				Private:     testBytes,
				RenewAt:     testTime,
				Result:      &testTfprotov5DynamicValue,
			},
		},
		"deferred-reason": {
			in: &tfprotov6.OpenEphemeralResourceResponse{
				Diagnostics: testTfprotov6Diagnostics,
				Result:      &testTfprotov6DynamicValue,
				Deferred: &tfprotov6.Deferred{
					Reason: tfprotov6.DeferredReasonResourceConfigUnknown,
				},
			},
			expected: &tfprotov5.OpenEphemeralResourceResponse{
				Diagnostics: testTfprotov5Diagnostics,
				Result:      &testTfprotov5DynamicValue,
				Deferred: &tfprotov5.Deferred{
					Reason: tfprotov5.DeferredReasonResourceConfigUnknown,
				},
			},
		},
	}

	for name, testCase := range testCases {

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got := tfprotov6tov5.OpenEphemeralResourceResponse(testCase.in)

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
				PriorIdentity:    &testTfprotov6ResourceIdentityData,
			},
			expected: &tfprotov5.PlanResourceChangeRequest{
				Config:           &testTfprotov5DynamicValue,
				PriorPrivate:     testBytes,
				PriorState:       &testTfprotov5DynamicValue,
				ProposedNewState: &testTfprotov5DynamicValue,
				ProviderMeta:     &testTfprotov5DynamicValue,
				TypeName:         "test_resource",
				PriorIdentity:    &testTfprotov5ResourceIdentityData,
			},
		},
		"client-capabilities-deferral-allowed": {
			in: &tfprotov6.PlanResourceChangeRequest{
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
			expected: &tfprotov5.PlanResourceChangeRequest{
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
		},
	}

	for name, testCase := range testCases {

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
				PlannedIdentity:             &testTfprotov6ResourceIdentityData,
			},
			expected: &tfprotov5.PlanResourceChangeResponse{
				Diagnostics:    testTfprotov5Diagnostics,
				PlannedPrivate: testBytes,
				PlannedState:   &testTfprotov5DynamicValue,
				RequiresReplace: []*tftypes.AttributePath{
					tftypes.NewAttributePath().WithAttributeName("test"),
				},
				UnsafeToUseLegacyTypeSystem: true,
				PlannedIdentity:             &testTfprotov5ResourceIdentityData,
			},
		},
		"deferred-reason": {
			in: &tfprotov6.PlanResourceChangeResponse{
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
			expected: &tfprotov5.PlanResourceChangeResponse{
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
		},
	}

	for name, testCase := range testCases {

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
		"client-capabilities-deferral-allowed": {
			in: &tfprotov6.ReadDataSourceRequest{
				Config:       &testTfprotov6DynamicValue,
				ProviderMeta: &testTfprotov6DynamicValue,
				TypeName:     "test_data_source",
				ClientCapabilities: &tfprotov6.ReadDataSourceClientCapabilities{
					DeferralAllowed: true,
				},
			},
			expected: &tfprotov5.ReadDataSourceRequest{
				Config:       &testTfprotov5DynamicValue,
				ProviderMeta: &testTfprotov5DynamicValue,
				TypeName:     "test_data_source",
				ClientCapabilities: &tfprotov5.ReadDataSourceClientCapabilities{
					DeferralAllowed: true,
				},
			},
		},
	}

	for name, testCase := range testCases {

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
		"deferred-reason": {
			in: &tfprotov6.ReadDataSourceResponse{
				Diagnostics: testTfprotov6Diagnostics,
				State:       &testTfprotov6DynamicValue,
				Deferred: &tfprotov6.Deferred{
					Reason: tfprotov6.DeferredReasonResourceConfigUnknown,
				},
			},
			expected: &tfprotov5.ReadDataSourceResponse{
				Diagnostics: testTfprotov5Diagnostics,
				State:       &testTfprotov5DynamicValue,
				Deferred: &tfprotov5.Deferred{
					Reason: tfprotov5.DeferredReasonResourceConfigUnknown,
				},
			},
		},
	}

	for name, testCase := range testCases {

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
				CurrentState:    &testTfprotov6DynamicValue,
				Private:         testBytes,
				ProviderMeta:    &testTfprotov6DynamicValue,
				TypeName:        "test_resource",
				CurrentIdentity: &testTfprotov6ResourceIdentityData,
			},
			expected: &tfprotov5.ReadResourceRequest{
				CurrentState:    &testTfprotov5DynamicValue,
				Private:         testBytes,
				ProviderMeta:    &testTfprotov5DynamicValue,
				TypeName:        "test_resource",
				CurrentIdentity: &testTfprotov5ResourceIdentityData,
			},
		},
		"client-capabilities-deferral-allowed": {
			in: &tfprotov6.ReadResourceRequest{
				CurrentState:    &testTfprotov6DynamicValue,
				Private:         testBytes,
				ProviderMeta:    &testTfprotov6DynamicValue,
				TypeName:        "test_resource",
				CurrentIdentity: &testTfprotov6ResourceIdentityData,
				ClientCapabilities: &tfprotov6.ReadResourceClientCapabilities{
					DeferralAllowed: true,
				},
			},
			expected: &tfprotov5.ReadResourceRequest{
				CurrentState:    &testTfprotov5DynamicValue,
				Private:         testBytes,
				ProviderMeta:    &testTfprotov5DynamicValue,
				TypeName:        "test_resource",
				CurrentIdentity: &testTfprotov5ResourceIdentityData,
				ClientCapabilities: &tfprotov5.ReadResourceClientCapabilities{
					DeferralAllowed: true,
				},
			},
		},
	}

	for name, testCase := range testCases {

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
				NewIdentity: &testTfprotov6ResourceIdentityData,
			},
			expected: &tfprotov5.ReadResourceResponse{
				Diagnostics: testTfprotov5Diagnostics,
				NewState:    &testTfprotov5DynamicValue,
				Private:     testBytes,
				NewIdentity: &testTfprotov5ResourceIdentityData,
			},
		},
		"deferred-reason": {
			in: &tfprotov6.ReadResourceResponse{
				Diagnostics: testTfprotov6Diagnostics,
				NewState:    &testTfprotov6DynamicValue,
				Private:     testBytes,
				NewIdentity: &testTfprotov6ResourceIdentityData,
				Deferred: &tfprotov6.Deferred{
					Reason: tfprotov6.DeferredReasonResourceConfigUnknown,
				},
			},
			expected: &tfprotov5.ReadResourceResponse{
				Diagnostics: testTfprotov5Diagnostics,
				NewState:    &testTfprotov5DynamicValue,
				Private:     testBytes,
				NewIdentity: &testTfprotov5ResourceIdentityData,
				Deferred: &tfprotov5.Deferred{
					Reason: tfprotov5.DeferredReasonResourceConfigUnknown,
				},
			},
		},
	}

	for name, testCase := range testCases {

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got := tfprotov6tov5.ReadResourceResponse(testCase.in)

			if diff := cmp.Diff(got, testCase.expected); diff != "" {
				t.Errorf("unexpected difference: %s", diff)
			}
		})
	}
}

func TestRenewEphemeralResourceRequest(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		in       *tfprotov6.RenewEphemeralResourceRequest
		expected *tfprotov5.RenewEphemeralResourceRequest
	}{
		"nil": {
			in:       nil,
			expected: nil,
		},
		"all-valid-fields": {
			in: &tfprotov6.RenewEphemeralResourceRequest{
				Private:  testBytes,
				TypeName: "test_ephemeral_resource",
			},
			expected: &tfprotov5.RenewEphemeralResourceRequest{
				Private:  testBytes,
				TypeName: "test_ephemeral_resource",
			},
		},
	}

	for name, testCase := range testCases {

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got := tfprotov6tov5.RenewEphemeralResourceRequest(testCase.in)

			if diff := cmp.Diff(got, testCase.expected); diff != "" {
				t.Errorf("unexpected difference: %s", diff)
			}
		})
	}
}

func TestRenewEphemeralResourceResponse(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		in       *tfprotov6.RenewEphemeralResourceResponse
		expected *tfprotov5.RenewEphemeralResourceResponse
	}{
		"nil": {
			in:       nil,
			expected: nil,
		},
		"all-valid-fields": {
			in: &tfprotov6.RenewEphemeralResourceResponse{
				Diagnostics: testTfprotov6Diagnostics,
				Private:     testBytes,
				RenewAt:     testTime,
			},
			expected: &tfprotov5.RenewEphemeralResourceResponse{
				Diagnostics: testTfprotov5Diagnostics,
				Private:     testBytes,
				RenewAt:     testTime,
			},
		},
	}

	for name, testCase := range testCases {

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got := tfprotov6tov5.RenewEphemeralResourceResponse(testCase.in)

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
				WriteOnly:       true,
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
				WriteOnly:       true,
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

func TestResourceIdentitySchema(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		in       *tfprotov6.ResourceIdentitySchema
		expected *tfprotov5.ResourceIdentitySchema
	}{
		"nil": {
			in:       nil,
			expected: nil,
		},
		"all-valid-fields": {
			in:       testTfprotov6ResourceIdentitySchema,
			expected: testTfprotov5ResourceIdentitySchema,
		},
	}

	for name, testCase := range testCases {

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got := tfprotov6tov5.ResourceIdentitySchema(testCase.in)

			if diff := cmp.Diff(got, testCase.expected); diff != "" {
				t.Errorf("unexpected difference: %s", diff)
			}
		})
	}
}

func TestResourceIdentitySchemaAttribute(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		in       *tfprotov6.ResourceIdentitySchemaAttribute
		expected *tfprotov5.ResourceIdentitySchemaAttribute
	}{
		"nil": {
			in:       nil,
			expected: nil,
		},
		"all-valid-fields": {
			in: &tfprotov6.ResourceIdentitySchemaAttribute{
				Name:              "test",
				Description:       "test description",
				Type:              tftypes.String,
				RequiredForImport: true,
				OptionalForImport: true,
			},
			expected: &tfprotov5.ResourceIdentitySchemaAttribute{
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

			got := tfprotov6tov5.ResourceIdentitySchemaAttribute(testCase.in)

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

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got := tfprotov6tov5.UpgradeResourceStateResponse(testCase.in)

			if diff := cmp.Diff(got, testCase.expected); diff != "" {
				t.Errorf("unexpected difference: %s", diff)
			}
		})
	}
}

func TestUpgradeResourceIdentityRequest(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		in       *tfprotov6.UpgradeResourceIdentityRequest
		expected *tfprotov5.UpgradeResourceIdentityRequest
	}{
		"nil": {
			in:       nil,
			expected: nil,
		},
		"all-valid-fields": {
			in: &tfprotov6.UpgradeResourceIdentityRequest{
				RawIdentity: &tfprotov6.RawState{
					JSON: testBytes,
				},
				TypeName: "test_resource",
				Version:  1,
			},
			expected: &tfprotov5.UpgradeResourceIdentityRequest{
				RawIdentity: &tfprotov5.RawState{
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

			got := tfprotov6tov5.UpgradeResourceIdentityRequest(testCase.in)

			if diff := cmp.Diff(got, testCase.expected); diff != "" {
				t.Errorf("unexpected difference: %s", diff)
			}
		})
	}
}

func TestUpgradeResourceIdentityResponse(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		in       *tfprotov6.UpgradeResourceIdentityResponse
		expected *tfprotov5.UpgradeResourceIdentityResponse
	}{
		"nil": {
			in:       nil,
			expected: nil,
		},
		"all-valid-fields": {
			in: &tfprotov6.UpgradeResourceIdentityResponse{
				Diagnostics:      testTfprotov6Diagnostics,
				UpgradedIdentity: &testTfprotov6ResourceIdentityData,
			},
			expected: &tfprotov5.UpgradeResourceIdentityResponse{
				Diagnostics:      testTfprotov5Diagnostics,
				UpgradedIdentity: &testTfprotov5ResourceIdentityData,
			},
		},
	}

	for name, testCase := range testCases {

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got := tfprotov6tov5.UpgradeResourceIdentityResponse(testCase.in)

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

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got := tfprotov6tov5.ValidateDataSourceConfigResponse(testCase.in)

			if diff := cmp.Diff(got, testCase.expected); diff != "" {
				t.Errorf("unexpected difference: %s", diff)
			}
		})
	}
}

func TestValidateEphemeralResourceConfigRequest(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		in       *tfprotov6.ValidateEphemeralResourceConfigRequest
		expected *tfprotov5.ValidateEphemeralResourceConfigRequest
	}{
		"nil": {
			in:       nil,
			expected: nil,
		},
		"all-valid-fields": {
			in: &tfprotov6.ValidateEphemeralResourceConfigRequest{
				Config:   &testTfprotov6DynamicValue,
				TypeName: "test_ephemeral_resource",
			},
			expected: &tfprotov5.ValidateEphemeralResourceConfigRequest{
				Config:   &testTfprotov5DynamicValue,
				TypeName: "test_ephemeral_resource",
			},
		},
	}

	for name, testCase := range testCases {

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got := tfprotov6tov5.ValidateEphemeralResourceConfigRequest(testCase.in)

			if diff := cmp.Diff(got, testCase.expected); diff != "" {
				t.Errorf("unexpected difference: %s", diff)
			}
		})
	}
}

func TestValidateEphemeralResourceConfigResponse(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		in       *tfprotov6.ValidateEphemeralResourceConfigResponse
		expected *tfprotov5.ValidateEphemeralResourceConfigResponse
	}{
		"nil": {
			in:       nil,
			expected: nil,
		},
		"all-valid-fields": {
			in: &tfprotov6.ValidateEphemeralResourceConfigResponse{
				Diagnostics: testTfprotov6Diagnostics,
			},
			expected: &tfprotov5.ValidateEphemeralResourceConfigResponse{
				Diagnostics: testTfprotov5Diagnostics,
			},
		},
	}

	for name, testCase := range testCases {

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got := tfprotov6tov5.ValidateEphemeralResourceConfigResponse(testCase.in)

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
		"client-capabilities-write-only-attributes-allowed": {
			in: &tfprotov6.ValidateResourceConfigRequest{
				ClientCapabilities: &tfprotov6.ValidateResourceConfigClientCapabilities{
					WriteOnlyAttributesAllowed: true,
				},
				Config:   &testTfprotov6DynamicValue,
				TypeName: "test_resource",
			},
			expected: &tfprotov5.ValidateResourceTypeConfigRequest{
				ClientCapabilities: &tfprotov5.ValidateResourceTypeConfigClientCapabilities{
					WriteOnlyAttributesAllowed: true,
				},
				Config:   &testTfprotov5DynamicValue,
				TypeName: "test_resource",
			},
		},
	}

	for name, testCase := range testCases {

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

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got := tfprotov6tov5.ValidateResourceTypeConfigResponse(testCase.in)

			if diff := cmp.Diff(got, testCase.expected); diff != "" {
				t.Errorf("unexpected difference: %s", diff)
			}
		})
	}
}

func TestValidateListResourceConfigRequest(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		in       *tfprotov6.ValidateListResourceConfigRequest
		expected *tfprotov5.ValidateListResourceConfigRequest
	}{
		"nil": {
			in:       nil,
			expected: nil,
		},
		"all-valid-fields": {
			in: &tfprotov6.ValidateListResourceConfigRequest{
				Config:   &testTfprotov6DynamicValue,
				TypeName: "test_list_resource",
			},
			expected: &tfprotov5.ValidateListResourceConfigRequest{
				Config:   &testTfprotov5DynamicValue,
				TypeName: "test_list_resource",
			},
		},
	}

	for name, testCase := range testCases {

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got := tfprotov6tov5.ValidateListResourceConfigRequest(testCase.in)

			if diff := cmp.Diff(got, testCase.expected); diff != "" {
				t.Errorf("unexpected difference: %s", diff)
			}
		})
	}
}

func TestValidateListResourceConfigResponse(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		in       *tfprotov6.ValidateListResourceConfigResponse
		expected *tfprotov5.ValidateListResourceConfigResponse
	}{
		"nil": {
			in:       nil,
			expected: nil,
		},
		"all-valid-fields": {
			in: &tfprotov6.ValidateListResourceConfigResponse{
				Diagnostics: testTfprotov6Diagnostics,
			},
			expected: &tfprotov5.ValidateListResourceConfigResponse{
				Diagnostics: testTfprotov5Diagnostics,
			},
		},
	}

	for name, testCase := range testCases {

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got := tfprotov6tov5.ValidateListResourceConfigResponse(testCase.in)

			if diff := cmp.Diff(got, testCase.expected); diff != "" {
				t.Errorf("unexpected difference: %s", diff)
			}
		})
	}
}

func TestListResourceRequest(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		in       *tfprotov6.ListResourceRequest
		expected *tfprotov5.ListResourceRequest
	}{
		"nil": {
			in:       nil,
			expected: nil,
		},
		"all-valid-fields": {
			in: &tfprotov6.ListResourceRequest{
				Config:   &testTfprotov6DynamicValue,
				TypeName: "test_list_resource",
			},
			expected: &tfprotov5.ListResourceRequest{
				Config:   &testTfprotov5DynamicValue,
				TypeName: "test_list_resource",
			},
		},
	}

	for name, testCase := range testCases {

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got := tfprotov6tov5.ListResourceRequest(testCase.in)

			if diff := cmp.Diff(got, testCase.expected); diff != "" {
				t.Errorf("unexpected difference: %s", diff)
			}
		})
	}
}

func TestListResourceServerStream(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		in       *tfprotov6.ListResourceServerStream
		expected *tfprotov5.ListResourceServerStream
	}{
		"nil": {
			in:       nil,
			expected: nil,
		},
		"all-valid-fields": {
			in:       testStreamv6,
			expected: testStreamv5,
		},
	}

	for name, testCase := range testCases {

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got := tfprotov6tov5.ListResourceServerStream(testCase.in)

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
		in       tfprotov6.ListResourceResult
		expected tfprotov5.ListResourceResult
	}{
		"identity-nil": {
			in: tfprotov6.ListResourceResult{
				DisplayName: "test",
				Resource:    &testTfprotov6DynamicValue,
				Identity:    nil,
				Diagnostics: testTfprotov6Diagnostics,
			},
			expected: tfprotov5.ListResourceResult{
				DisplayName: "test",
				Resource:    &testTfprotov5DynamicValue,
				Identity:    nil,
				Diagnostics: testTfprotov5Diagnostics,
			},
		},
		"resource-nil": {
			in: tfprotov6.ListResourceResult{
				DisplayName: "test",
				Resource:    nil,
				Identity:    &testTfprotov6ResourceIdentityData,
				Diagnostics: testTfprotov6Diagnostics,
			},
			expected: tfprotov5.ListResourceResult{
				DisplayName: "test",
				Resource:    nil,
				Identity:    &testTfprotov5ResourceIdentityData,
				Diagnostics: testTfprotov5Diagnostics,
			},
		},
		"all-valid-fields": {
			in: tfprotov6.ListResourceResult{
				Diagnostics: testTfprotov6Diagnostics,
			},
			expected: tfprotov5.ListResourceResult{
				Diagnostics: testTfprotov5Diagnostics,
			},
		},
	}

	for name, testCase := range testCases {

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got := tfprotov6tov5.ListResourceResult(testCase.in)

			if diff := cmp.Diff(got, testCase.expected); diff != "" {
				t.Errorf("unexpected difference: %s", diff)
			}
		})
	}
}

func TestActionSchema(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		in            *tfprotov6.ActionSchema
		expected      *tfprotov5.ActionSchema
		expectedError error
	}{
		"nil": {
			in:       nil,
			expected: nil,
		},
		"unlinked": {
			in: &tfprotov6.ActionSchema{
				Schema: testTfprotov6Schema,
				Type:   tfprotov6.UnlinkedActionSchemaType{},
			},
			expected: &tfprotov5.ActionSchema{
				Schema: testTfprotov5Schema,
				Type:   tfprotov5.UnlinkedActionSchemaType{},
			},
		},
		"lifecycle": {
			in: &tfprotov6.ActionSchema{
				Schema: testTfprotov6Schema,
				Type: tfprotov6.LifecycleActionSchemaType{
					Executes: tfprotov6.LifecycleExecutionOrderAfter,
					LinkedResource: &tfprotov6.LinkedResourceSchema{
						TypeName:    "test_resource_linked_1",
						Description: "This is a linked resource.",
					},
				},
			},
			expected: &tfprotov5.ActionSchema{
				Schema: testTfprotov5Schema,
				Type: tfprotov5.LifecycleActionSchemaType{
					Executes: tfprotov5.LifecycleExecutionOrderAfter,
					LinkedResource: &tfprotov5.LinkedResourceSchema{
						TypeName:    "test_resource_linked_1",
						Description: "This is a linked resource.",
					},
				},
			},
		},
		"linked": {
			in: &tfprotov6.ActionSchema{
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
			expected: &tfprotov5.ActionSchema{
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
		},
		"nested-attribute-error": {
			in: &tfprotov6.ActionSchema{
				Schema: &tfprotov6.Schema{
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
				Type: tfprotov6.UnlinkedActionSchemaType{},
			},
			expected:      nil,
			expectedError: fmt.Errorf("unable to convert attribute \"test_attribute\" schema: %w", tfprotov6tov5.ErrSchemaAttributeNestedTypeNotImplemented),
		},
	}

	for name, testCase := range testCases {

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got, err := tfprotov6tov5.ActionSchema(testCase.in)

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

func TestValidateActionConfigRequest(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		in       *tfprotov6.ValidateActionConfigRequest
		expected *tfprotov5.ValidateActionConfigRequest
	}{
		"nil": {
			in:       nil,
			expected: nil,
		},
		"all-valid-fields": {
			in: &tfprotov6.ValidateActionConfigRequest{
				Config:     &testTfprotov6DynamicValue,
				ActionType: "test_action",
			},
			expected: &tfprotov5.ValidateActionConfigRequest{
				Config:     &testTfprotov5DynamicValue,
				ActionType: "test_action",
			},
		},
	}

	for name, testCase := range testCases {

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got := tfprotov6tov5.ValidateActionConfigRequest(testCase.in)

			if diff := cmp.Diff(got, testCase.expected); diff != "" {
				t.Errorf("unexpected difference: %s", diff)
			}
		})
	}
}

func TestValidateActionConfigResponse(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		in       *tfprotov6.ValidateActionConfigResponse
		expected *tfprotov5.ValidateActionConfigResponse
	}{
		"nil": {
			in:       nil,
			expected: nil,
		},
		"all-valid-fields": {
			in: &tfprotov6.ValidateActionConfigResponse{
				Diagnostics: testTfprotov6Diagnostics,
			},
			expected: &tfprotov5.ValidateActionConfigResponse{
				Diagnostics: testTfprotov5Diagnostics,
			},
		},
	}

	for name, testCase := range testCases {

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got := tfprotov6tov5.ValidateActionConfigResponse(testCase.in)

			if diff := cmp.Diff(got, testCase.expected); diff != "" {
				t.Errorf("unexpected difference: %s", diff)
			}
		})
	}
}

func TestPlanActionRequest(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		in       *tfprotov6.PlanActionRequest
		expected *tfprotov5.PlanActionRequest
	}{
		"nil": {
			in:       nil,
			expected: nil,
		},
		"no-linked-resources": {
			in: &tfprotov6.PlanActionRequest{
				ActionType: "test_action",
				Config:     &testTfprotov6DynamicValue,
			},
			expected: &tfprotov5.PlanActionRequest{
				ActionType: "test_action",
				Config:     &testTfprotov5DynamicValue,
			},
		},
		"linked-resources": {
			in: &tfprotov6.PlanActionRequest{
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
			expected: &tfprotov5.PlanActionRequest{
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
		},
		"client-capabilities-deferral-allowed": {
			in: &tfprotov6.PlanActionRequest{
				ActionType: "test_action",
				Config:     &testTfprotov6DynamicValue,
				ClientCapabilities: &tfprotov6.PlanActionClientCapabilities{
					DeferralAllowed: true,
				},
			},
			expected: &tfprotov5.PlanActionRequest{
				ActionType: "test_action",
				Config:     &testTfprotov5DynamicValue,
				ClientCapabilities: &tfprotov5.PlanActionClientCapabilities{
					DeferralAllowed: true,
				},
			},
		},
	}

	for name, testCase := range testCases {

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got := tfprotov6tov5.PlanActionRequest(testCase.in)

			if diff := cmp.Diff(got, testCase.expected); diff != "" {
				t.Errorf("unexpected difference: %s", diff)
			}
		})
	}
}

func TestPlanActionResponse(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		in       *tfprotov6.PlanActionResponse
		expected *tfprotov5.PlanActionResponse
	}{
		"nil": {
			in:       nil,
			expected: nil,
		},
		"no-linked-resources": {
			in: &tfprotov6.PlanActionResponse{
				Diagnostics: testTfprotov6Diagnostics,
			},
			expected: &tfprotov5.PlanActionResponse{
				Diagnostics: testTfprotov5Diagnostics,
			},
		},
		"linked-resources": {
			in: &tfprotov6.PlanActionResponse{
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
			expected: &tfprotov5.PlanActionResponse{
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
		},
		"deferred-reason": {
			in: &tfprotov6.PlanActionResponse{
				Diagnostics: testTfprotov6Diagnostics,
				Deferred: &tfprotov6.Deferred{
					Reason: tfprotov6.DeferredReasonResourceConfigUnknown,
				},
			},
			expected: &tfprotov5.PlanActionResponse{
				Diagnostics: testTfprotov5Diagnostics,
				Deferred: &tfprotov5.Deferred{
					Reason: tfprotov5.DeferredReasonResourceConfigUnknown,
				},
			},
		},
	}

	for name, testCase := range testCases {

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got := tfprotov6tov5.PlanActionResponse(testCase.in)

			if diff := cmp.Diff(got, testCase.expected); diff != "" {
				t.Errorf("unexpected difference: %s", diff)
			}
		})
	}
}

func TestInvokeActionRequest(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		in       *tfprotov6.InvokeActionRequest
		expected *tfprotov5.InvokeActionRequest
	}{
		"nil": {
			in:       nil,
			expected: nil,
		},
		"no-linked-resources": {
			in: &tfprotov6.InvokeActionRequest{
				ActionType: "test_action",
				Config:     &testTfprotov6DynamicValue,
			},
			expected: &tfprotov5.InvokeActionRequest{
				ActionType: "test_action",
				Config:     &testTfprotov5DynamicValue,
			},
		},
		"linked-resources": {
			in: &tfprotov6.InvokeActionRequest{
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
			expected: &tfprotov5.InvokeActionRequest{
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
		},
	}

	for name, testCase := range testCases {

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got := tfprotov6tov5.InvokeActionRequest(testCase.in)

			if diff := cmp.Diff(got, testCase.expected); diff != "" {
				t.Errorf("unexpected difference: %s", diff)
			}
		})
	}
}

func TestInvokeActionServerStream(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		in       *tfprotov6.InvokeActionServerStream
		expected *tfprotov5.InvokeActionServerStream
	}{
		"nil": {
			in:       nil,
			expected: nil,
		},
		"all-valid-fields": {
			in: &tfprotov6.InvokeActionServerStream{
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
			expected: &tfprotov5.InvokeActionServerStream{
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
		},
	}

	for name, testCase := range testCases {

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got := tfprotov6tov5.InvokeActionServerStream(testCase.in)

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
