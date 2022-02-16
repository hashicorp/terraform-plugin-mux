// Translate a protocol version 6 provider server into protocol version 5.
//
// Supported protocol version 6 provider servers include any which implement
// the github.com/hashicorp/terraform-plugin-go/tfprotov6.ProviderServer
// interface, such as:
//
// - github.com/hashicorp/terraform-plugin-framework
// - github.com/hashicorp/terraform-plugin-go/tfprotov6/tf6server
// - github.com/hashicorp/terraform-plugin-mux/tf6muxserver
//
// Refer to the DowngradeServer() function for wrapping a server.
package tf6to5server
