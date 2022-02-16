// Combine multiple protocol version 6 provider servers into a single server.
//
// Supported protocol version 6 provider servers include any which implement
// the github.com/hashicorp/terraform-plugin-go/tfprotov6.ProviderServer
// interface, such as:
//
// - github.com/hashicorp/terraform-plugin-framework
// - github.com/hashicorp/terraform-plugin-go/tfprotov6/tf6server
// - github.com/hashicorp/terraform-plugin-mux/tf5to6server
//
// Refer to the NewMuxServer() function for creating a combined server.
package tf6muxserver
