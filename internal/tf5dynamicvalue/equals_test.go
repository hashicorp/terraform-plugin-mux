package tf5dynamicvalue_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-go/tfprotov5"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/hashicorp/terraform-plugin-mux/internal/tf5dynamicvalue"
)

func TestEquals(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		schemaType    tftypes.Type
		dynamicValue1 *tfprotov5.DynamicValue
		dynamicValue2 *tfprotov5.DynamicValue
		expected      bool
		expectedError error
	}{
		"all-missing": {
			schemaType:    nil,
			dynamicValue1: nil,
			dynamicValue2: nil,
			expected:      true,
		},
		"first-missing": {
			schemaType:    nil,
			dynamicValue1: nil,
			dynamicValue2: &tfprotov5.DynamicValue{},
			expected:      false,
		},
		"second-missing": {
			schemaType:    nil,
			dynamicValue1: &tfprotov5.DynamicValue{},
			dynamicValue2: nil,
			expected:      false,
		},
		"missing-type": {
			schemaType: nil,
			dynamicValue1: tf5dynamicvalue.Must(
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
			dynamicValue2: tf5dynamicvalue.Must(
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
			expectedError: fmt.Errorf("unable to unmarshal DynamicValue: missing Type"),
		},
		"mismatched-type": {
			schemaType: tftypes.Object{
				AttributeTypes: map[string]tftypes.Type{
					"test_bool_attribute": tftypes.Bool,
				},
			},
			dynamicValue1: tf5dynamicvalue.Must(
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
			dynamicValue2: tf5dynamicvalue.Must(
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
		"String-different-value": {
			schemaType: tftypes.Object{
				AttributeTypes: map[string]tftypes.Type{
					"test_string_attribute": tftypes.String,
				},
			},
			dynamicValue1: tf5dynamicvalue.Must(
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
						"test_string_attribute": tftypes.NewValue(tftypes.String, "test-value-1"),
					},
				),
			),
			dynamicValue2: tf5dynamicvalue.Must(
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
						"test_string_attribute": tftypes.NewValue(tftypes.String, "test-value-2"),
					},
				),
			),
			expected: false,
		},
		"String-equal-value": {
			schemaType: tftypes.Object{
				AttributeTypes: map[string]tftypes.Type{
					"test_string_attribute": tftypes.String,
				},
			},
			dynamicValue1: tf5dynamicvalue.Must(
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
			dynamicValue2: tf5dynamicvalue.Must(
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
			expected: true,
		},
	}

	for name, testCase := range testCases {
		name, testCase := name, testCase

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got, err := tf5dynamicvalue.Equals(testCase.schemaType, testCase.dynamicValue1, testCase.dynamicValue2)

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
