package tf5muxserver

import "github.com/hashicorp/terraform-plugin-go/tfprotov5"

// serverSupportsPlanDestroy returns true if the given ServerCapabilities is not
// nil and enables the PlanDestroy capability.
func serverSupportsPlanDestroy(capabilities *tfprotov5.ServerCapabilities) bool {
	if capabilities == nil {
		return false
	}

	return capabilities.PlanDestroy
}
