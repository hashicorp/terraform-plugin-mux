// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package tf5muxserver

import "github.com/hashicorp/terraform-plugin-go/tfprotov5"

func actionDuplicateError(actionType string) *tfprotov5.Diagnostic {
	return &tfprotov5.Diagnostic{
		Severity: tfprotov5.DiagnosticSeverityError,
		Summary:  "Invalid Provider Server Combination",
		Detail: "The combined provider has multiple implementations of the same action type across underlying providers. " +
			"Actions must be implemented by only one underlying provider. " +
			"This is always an issue in the provider implementation and should be reported to the provider developers.\n\n" +
			"Duplicate action: " + actionType,
	}
}

func actionMissingError(actionType string) *tfprotov5.Diagnostic {
	return &tfprotov5.Diagnostic{
		Severity: tfprotov5.DiagnosticSeverityError,
		Summary:  "Action Not Implemented",
		Detail: "The combined provider does not implement the requested action. " +
			"This is always an issue in the provider implementation and should be reported to the provider developers.\n\n" +
			"Missing action: " + actionType,
	}
}

func dataSourceDuplicateError(typeName string) *tfprotov5.Diagnostic {
	return &tfprotov5.Diagnostic{
		Severity: tfprotov5.DiagnosticSeverityError,
		Summary:  "Invalid Provider Server Combination",
		Detail: "The combined provider has multiple implementations of the same data source type across underlying providers. " +
			"Data source types must be implemented by only one underlying provider. " +
			"This is always an issue in the provider implementation and should be reported to the provider developers.\n\n" +
			"Duplicate data source type: " + typeName,
	}
}

func dataSourceMissingError(typeName string) *tfprotov5.Diagnostic {
	return &tfprotov5.Diagnostic{
		Severity: tfprotov5.DiagnosticSeverityError,
		Summary:  "Data Source Not Implemented",
		Detail: "The combined provider does not implement the requested data source type. " +
			"This is always an issue in the provider implementation and should be reported to the provider developers.\n\n" +
			"Missing data source type: " + typeName,
	}
}

func ephemeralResourceDuplicateError(typeName string) *tfprotov5.Diagnostic {
	return &tfprotov5.Diagnostic{
		Severity: tfprotov5.DiagnosticSeverityError,
		Summary:  "Invalid Provider Server Combination",
		Detail: "The combined provider has multiple implementations of the same ephemeral resource type across underlying providers. " +
			"Ephemeral resource types must be implemented by only one underlying provider. " +
			"This is always an issue in the provider implementation and should be reported to the provider developers.\n\n" +
			"Duplicate ephemeral resource type: " + typeName,
	}
}

func ephemeralResourceMissingError(typeName string) *tfprotov5.Diagnostic {
	return &tfprotov5.Diagnostic{
		Severity: tfprotov5.DiagnosticSeverityError,
		Summary:  "Ephemeral Resource Not Implemented",
		Detail: "The combined provider does not implement the requested ephemeral resource type. " +
			"This is always an issue in the provider implementation and should be reported to the provider developers.\n\n" +
			"Missing ephemeral resource type: " + typeName,
	}
}

func listResourceDuplicateError(typeName string) *tfprotov5.Diagnostic {
	return &tfprotov5.Diagnostic{
		Severity: tfprotov5.DiagnosticSeverityError,
		Summary:  "Invalid Provider Server Combination",
		Detail: "The combined provider has multiple implementations of the same list resource type across underlying providers. " +
			"List resource types must be implemented by only one underlying provider. " +
			"This is always an issue in the provider implementation and should be reported to the provider developers.\n\n" +
			"Duplicate list resource type: " + typeName,
	}
}

func listResourceMissingError(typeName string) *tfprotov5.Diagnostic {
	return &tfprotov5.Diagnostic{
		Severity: tfprotov5.DiagnosticSeverityError,
		Summary:  "List Resource Not Implemented",
		Detail: "The combined provider does not implement the requested list resource type. " +
			"This is always an issue in the provider implementation and should be reported to the provider developers.\n\n" +
			"Missing list resource type: " + typeName,
	}
}

func diagnosticsHasError(diagnostics []*tfprotov5.Diagnostic) bool {
	for _, diagnostic := range diagnostics {
		if diagnostic == nil {
			continue
		}

		if diagnostic.Severity == tfprotov5.DiagnosticSeverityError {
			return true
		}
	}

	return false
}

func functionDuplicateError(name string) *tfprotov5.Diagnostic {
	return &tfprotov5.Diagnostic{
		Severity: tfprotov5.DiagnosticSeverityError,
		Summary:  "Invalid Provider Server Combination",
		Detail: "The combined provider has multiple implementations of the same function name across underlying providers. " +
			"Functions must be implemented by only one underlying provider. " +
			"This is always an issue in the provider implementation and should be reported to the provider developers.\n\n" +
			"Duplicate function: " + name,
	}
}

func functionMissingError(name string) *tfprotov5.Diagnostic {
	return &tfprotov5.Diagnostic{
		Severity: tfprotov5.DiagnosticSeverityError,
		Summary:  "Function Not Implemented",
		Detail: "The combined provider does not implement the requested function. " +
			"This is always an issue in the provider implementation and should be reported to the provider developers.\n\n" +
			"Missing function: " + name,
	}
}

func resourceDuplicateError(typeName string) *tfprotov5.Diagnostic {
	return &tfprotov5.Diagnostic{
		Severity: tfprotov5.DiagnosticSeverityError,
		Summary:  "Invalid Provider Server Combination",
		Detail: "The combined provider has multiple implementations of the same resource type across underlying providers. " +
			"Resource types must be implemented by only one underlying provider. " +
			"This is always an issue in the provider implementation and should be reported to the provider developers.\n\n" +
			"Duplicate resource type: " + typeName,
	}
}

func resourceMissingError(typeName string) *tfprotov5.Diagnostic {
	return &tfprotov5.Diagnostic{
		Severity: tfprotov5.DiagnosticSeverityError,
		Summary:  "Resource Not Implemented",
		Detail: "The combined provider does not implement the requested resource type. " +
			"This is always an issue in the provider implementation and should be reported to the provider developers.\n\n" +
			"Missing resource type: " + typeName,
	}
}

func resourceIdentityDuplicateError(typeName string) *tfprotov5.Diagnostic {
	return &tfprotov5.Diagnostic{
		Severity: tfprotov5.DiagnosticSeverityError,
		Summary:  "Invalid Provider Server Combination",
		Detail: "The combined provider has multiple implementations of the same resource identity across underlying providers. " +
			"Resource identity types must be implemented by only one underlying provider. " +
			"This is always an issue in the provider implementation and should be reported to the provider developers.\n\n" +
			"Duplicate identity type for resource: " + typeName,
	}
}
