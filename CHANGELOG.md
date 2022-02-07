# 0.5.0

NOTES:

* Providers can now be muxed with a combination of terraform-plugin-sdk and terraform-plugin-framework server implementations. One option is the terraform-plugin-sdk server can be upgraded to protocol version 6, then muxed with the terraform-plugin-framework server. This allows using new protocol features in the framework implementation, such as nested attributes, but requires Terraform CLI 1.1.5 or later. The other option is the terraform-plugin-framework server can be downgraded to protocol version 5, then muxed with the terraform-plugin-sdk server. This prevents using new protocol features in the framework implementation, however it remains compatible with Terraform CLI 0.12 and later. ([#42](https://github.com/hashicorp/terraform-plugin-mux/issues/42))

BREAKING CHANGES:

* The root package `SchemaServer` types and `NewSchemaServerFactory` function have been migrated to the `tf5muxserver` package. To upgrade, replace `tfmux.NewSchemaServerFactory` with `tf5muxserver.NewMuxServer` and replace any invocations of the previous `SchemaServerFactory` type `Server()` method with `ProviderServer()`. The underlying types are no longer exported. ([#39](https://github.com/hashicorp/terraform-plugin-mux/issues/39))

FEATURES:

* Added `tf5to6server` package, for upgrading a protocol version 5 server to protocol version 6 ([#42](https://github.com/hashicorp/terraform-plugin-mux/issues/42))
* Added `tf6muxserver` package, a protocol version 6 compatible mux server ([#37](https://github.com/hashicorp/terraform-plugin-mux/issues/37))
* Added `tf6to5server` package, for downgrading a protocol version 6 server to protocol version 5 ([#42](https://github.com/hashicorp/terraform-plugin-mux/issues/42))

ENHANCEMENTS:

* Added the `tf_mux_provider` key to all downstream logging calls, decorating them with the muxed server that actually served the request. ([#31](https://github.com/hashicorp/terraform-plugin-mux/issues/31))
* Added trace level logging for mux logic, controlled by the `TF_LOG_SDK_MUX` environment variable. ([#31](https://github.com/hashicorp/terraform-plugin-mux/issues/31))

# 0.4.0 (December 07, 2021)

NOTES:

* Upgraded terraform-plugin-go to v0.5.0. Provider built against versions of terraform-plugin-go prior to v0.5.0 will run into compatibility issues due to breaking changes in terraform-plugin-go.

# 0.3.0 (September 24, 2021)

NOTES:

* Upgraded terraform-plugin-go to v0.4.0. Providers built against versions of terraform-plugin-go prior to v0.4.0 will run into compatibility issues due to breaking changes in terraform-plugin-go.

# 0.2.0 (May 10, 2021)

NOTES:

* Upgraded terraform-plugin-go to v0.3.0. Providers built against versions of terraform-plugin-go prior to v0.3.0 will run into compatibility issues due to breaking changes in terraform-plugin-go.

# 0.1.1 (February 10, 2021)

BUG FIXES:

* Compare schemas in an order-insensitive way when deciding whether two server implementations are returning the same schema. ([#18](https://github.com/hashicorp/terraform-plugin-mux/issues/18))
* Surface the difference between schemas when provider and provider_meta schemas differ. ([#18](https://github.com/hashicorp/terraform-plugin-mux/issues/18))

# 0.1.0 (November 02, 2020)

Initial release.
