package tf6dynamicvalue_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/hashicorp/terraform-plugin-mux/internal/tf6dynamicvalue"
)

func TestIsNull(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		schema        *tfprotov6.Schema
		dynamicValue  *tfprotov6.DynamicValue
		expected      bool
		expectedError error
	}{
		"nil-dynamic-value": {
			schema:       nil,
			dynamicValue: nil,
			expected:     true,
		},
		"nil-schema": {
			schema:        nil,
			dynamicValue:  &tfprotov6.DynamicValue{},
			expected:      false,
			expectedError: fmt.Errorf("unable to unmarshal DynamicValue: missing Type"),
		},
		"NewDynamicValue-error": {
			schema: &tfprotov6.Schema{
				Block: &tfprotov6.SchemaBlock{
					Attributes: []*tfprotov6.SchemaAttribute{
						{
							Name: "test_bool_attribute", // intentionally different
							Type: tftypes.Bool,          // intentionally different
						},
					},
				},
			},
			dynamicValue: tf6dynamicvalue.Must(
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
					map[string]tftypes.Value{
						"test_string_attribute": tftypes.NewValue(tftypes.String, "test-value"),
					},
				),
			),
			expected:      false,
			expectedError: fmt.Errorf("unable to unmarshal DynamicValue: unknown attribute \"test_string_attribute\""),
		},
		"null": {
			schema: &tfprotov6.Schema{
				Block: &tfprotov6.SchemaBlock{
					Attributes: []*tfprotov6.SchemaAttribute{
						{
							Name: "test_string_attribute",
							Type: tftypes.String,
						},
					},
				},
			},
			dynamicValue: tf6dynamicvalue.Must(
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
					nil,
				),
			),
			expected: true,
		},
		"known": {
			schema: &tfprotov6.Schema{
				Block: &tfprotov6.SchemaBlock{
					Attributes: []*tfprotov6.SchemaAttribute{
						{
							Name: "test_string_attribute",
							Type: tftypes.String,
						},
					},
				},
			},
			dynamicValue: tf6dynamicvalue.Must(
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
					map[string]tftypes.Value{
						"test_string_attribute": tftypes.NewValue(tftypes.String, "test-value"),
					},
				),
			),
			expected: false,
		},
	}

	for name, testCase := range testCases {
		name, testCase := name, testCase

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got, err := tf6dynamicvalue.IsNull(testCase.schema, testCase.dynamicValue)

			if err != nil {
				if testCase.expectedError == nil {
					t.Fatalf("wanted no error, got error: %s", err)
				}

				if !strings.Contains(err.Error(), testCase.expectedError.Error()) {
					t.Fatalf("wanted error %q, got error: %s", testCase.expectedError.Error(), err.Error())
				}
			}

			if err == nil && testCase.expectedError != nil {
				t.Fatalf("got no error, wanted err: %s", testCase.expectedError)
			}

			if got != testCase.expected {
				t.Errorf("expected %t, got %t", testCase.expected, got)
			}
		})
	}
}
