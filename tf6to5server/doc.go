// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

// Package tf6to5server translates a provider that implements protocol version 6, into one that implements protocol version 5.
//
// Supported protocol version 6 provider servers include any which implement
// the tfprotov6.ProviderServer (https://pkg.go.dev/github.com/hashicorp/terraform-plugin-go/tfprotov6#ProviderServer)
// interface, such as:
//
//   - https://pkg.go.dev/github.com/hashicorp/terraform-plugin-framework
//   - https://pkg.go.dev/github.com/hashicorp/terraform-plugin-go/tfprotov6/tf6server
//   - https://pkg.go.dev/github.com/hashicorp/terraform-plugin-mux/tf6muxserver
//
// Refer to the DowngradeServer() function for wrapping a server.
package tf6to5server
