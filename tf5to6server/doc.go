// Translate a protocol version 5 provider server into protocol version 6.
//
// Supported protocol version 5 provider servers include any which implement
// the github.com/hashicorp/terraform-plugin-go/tfprotov5.ProviderServer
// interface, such as:
//
// - github.com/hashicorp/terraform-plugin-go/tfprotov5/tf5server
// - github.com/hashicorp/terraform-plugin-mux/tf5muxserver
// - github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema
//
// Refer to the UpgradeServer() function for wrapping a server.
package tf5to6server
