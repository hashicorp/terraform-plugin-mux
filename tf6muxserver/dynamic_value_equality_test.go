package tf6muxserver

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

func TestDynamicValueEquals(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		schemaType    tftypes.Type
		dynamicValue1 func() (*tfprotov6.DynamicValue, error)
		dynamicValue2 func() (*tfprotov6.DynamicValue, error)
		expected      bool
		expectedError error
	}{
		"all-missing": {
			schemaType: nil,
			dynamicValue1: func() (*tfprotov6.DynamicValue, error) {
				return nil, nil
			},
			dynamicValue2: func() (*tfprotov6.DynamicValue, error) {
				return nil, nil
			},
			expected: true,
		},
		"first-missing": {
			schemaType: nil,
			dynamicValue1: func() (*tfprotov6.DynamicValue, error) {
				return nil, nil
			},
			dynamicValue2: func() (*tfprotov6.DynamicValue, error) {
				return &tfprotov6.DynamicValue{}, nil
			},
			expected: false,
		},
		"second-missing": {
			schemaType: nil,
			dynamicValue1: func() (*tfprotov6.DynamicValue, error) {
				return &tfprotov6.DynamicValue{}, nil
			},
			dynamicValue2: func() (*tfprotov6.DynamicValue, error) {
				return nil, nil
			},
			expected: false,
		},
		"missing-type": {
			schemaType: nil,
			dynamicValue1: func() (*tfprotov6.DynamicValue, error) {
				dv, err := tfprotov6.NewDynamicValue(
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
				)
				return &dv, err
			},
			dynamicValue2: func() (*tfprotov6.DynamicValue, error) {
				dv, err := tfprotov6.NewDynamicValue(
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
				)
				return &dv, err
			},
			expected:      false,
			expectedError: fmt.Errorf("unable to unmarshal DynamicValue: missing Type"),
		},
		"mismatched-type": {
			schemaType: tftypes.Object{
				AttributeTypes: map[string]tftypes.Type{
					"test_bool_attribute": tftypes.Bool,
				},
			},
			dynamicValue1: func() (*tfprotov6.DynamicValue, error) {
				dv, err := tfprotov6.NewDynamicValue(
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
				)
				return &dv, err
			},
			dynamicValue2: func() (*tfprotov6.DynamicValue, error) {
				dv, err := tfprotov6.NewDynamicValue(
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
				)
				return &dv, err
			},
			expected:      false,
			expectedError: fmt.Errorf("unable to unmarshal DynamicValue: unknown attribute \"test_string_attribute\""),
		},
		"String-different-value": {
			schemaType: tftypes.Object{
				AttributeTypes: map[string]tftypes.Type{
					"test_string_attribute": tftypes.String,
				},
			},
			dynamicValue1: func() (*tfprotov6.DynamicValue, error) {
				dv, err := tfprotov6.NewDynamicValue(
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
				)
				return &dv, err
			},
			dynamicValue2: func() (*tfprotov6.DynamicValue, error) {
				dv, err := tfprotov6.NewDynamicValue(
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
				)
				return &dv, err
			},
			expected: false,
		},
		"String-equal-value": {
			schemaType: tftypes.Object{
				AttributeTypes: map[string]tftypes.Type{
					"test_string_attribute": tftypes.String,
				},
			},
			dynamicValue1: func() (*tfprotov6.DynamicValue, error) {
				dv, err := tfprotov6.NewDynamicValue(
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
				)
				return &dv, err
			},
			dynamicValue2: func() (*tfprotov6.DynamicValue, error) {
				dv, err := tfprotov6.NewDynamicValue(
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
				)
				return &dv, err
			},
			expected: true,
		},
	}

	for name, testCase := range testCases {
		name, testCase := name, testCase

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			dynamicValue1, err := testCase.dynamicValue1()

			if err != nil {
				t.Fatalf("unable to create first DynamicValue: %s", err)
			}

			dynamicValue2, err := testCase.dynamicValue2()

			if err != nil {
				t.Fatalf("unable to create second DynamicValue: %s", err)
			}

			got, err := dynamicValueEquals(testCase.schemaType, dynamicValue1, dynamicValue2)

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
