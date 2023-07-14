## 0.11.2 (July 14, 2023)

BUG FIXES:

* tf5muxserver: Ensure `GetProviderSchema` RPC diagnostics are properly returned to Terraform ([#176](https://github.com/hashicorp/terraform-plugin-mux/issues/176))
* tf6muxserver: Ensure `GetProviderSchema` RPC diagnostics are properly returned to Terraform ([#176](https://github.com/hashicorp/terraform-plugin-mux/issues/176))

## 0.11.1 (June 29, 2023)

BUG FIXES:

* tf5muxserver: Adjust function signature of `NewMuxServer()` to return `*muxServer`, which is required to satisfy the `tfprotov5.ProviderServer` interface ([#172](https://github.com/hashicorp/terraform-plugin-mux/issues/172))
* tf6muxserver: Adjust function signature of `NewMuxServer()` to return `*muxServer`, which is required to satisfy the `tfprotov6.ProviderServer` interface ([#172](https://github.com/hashicorp/terraform-plugin-mux/issues/172))

## 0.11.0 (June 28, 2023)

BUG FIXES:

* tf5muxserver: Removed unnecessary resource schema caching, which reduces resident memory utilization ([#168](https://github.com/hashicorp/terraform-plugin-mux/issues/168))
* tf6muxserver: Removed unnecessary resource schema caching, which reduces resident memory utilization ([#168](https://github.com/hashicorp/terraform-plugin-mux/issues/168))

## 0.10.0 (April 24, 2023)

NOTES:

* This Go module has been updated to Go 1.19 per the [Go support policy](https://golang.org/doc/devel/release.html#policy). Any consumers building on earlier Go versions may experience errors. ([#143](https://github.com/hashicorp/terraform-plugin-mux/issues/143))

BUG FIXES:

* tf5muxserver+tf6muxserver: Ensure provider acceptance testing can properly detect mux server errors in `GetProviderSchema` response ([#152](https://github.com/hashicorp/terraform-plugin-mux/issues/152))

## 0.9.0 (February 08, 2023)

ENHANCEMENTS:

* tf5muxserver+tf6muxserver: Support Terraform 1.3+ PlanResourceChange on destroy for underlying servers which enable the capability, such as terraform-plugin-framework ([#133](https://github.com/hashicorp/terraform-plugin-mux/issues/133))

# 0.8.0 (December 20, 2022)

NOTES:

* This Go module has been updated to Go 1.18 per the [Go support policy](https://golang.org/doc/devel/release.html#policy). Any consumers building on earlier Go versions may experience errors. ([#101](https://github.com/hashicorp/terraform-plugin-mux/issues/101))

BUG FIXES:

* tf5muxserver+tf6muxserver: Allow differing provider schema block `MinItems` and `MaxItems` as terraform-plugin-framework does not use those fields for configuration validation ([#118](https://github.com/hashicorp/terraform-plugin-mux/issues/118))
* tf5muxserver+tf6muxserver: Deferred combined server implementation errors until `GetProviderSchema` RPC to prevent confusing Terraform CLI plugin startup errors ([#121](https://github.com/hashicorp/terraform-plugin-mux/issues/121))

# 0.7.0 (July 15, 2022)

NOTES:

* The underlying `terraform-plugin-log` dependency has been updated to v0.6.0, which includes log filtering support and breaking changes of `With()` to `SetField()` function names. Any provider logging which calls those functions may require updates. ([#92](https://github.com/hashicorp/terraform-plugin-mux/issues/92))
* This Go module has been updated to Go 1.17 per the [Go support policy](https://golang.org/doc/devel/release.html#policy). Any consumers building on earlier Go versions may experience errors. ([#73](https://github.com/hashicorp/terraform-plugin-mux/issues/73))

# 0.6.0 (March 10, 2022)

NOTES:

* The underlying `terraform-plugin-log` dependency has been updated to v0.3.0, which includes a breaking change in the optional additional fields parameter of logging function calls to ensure correctness and catch coding errors during compilation. Any early adopter provider logging which calls those functions may require updates. ([#63](https://github.com/hashicorp/terraform-plugin-mux/issues/63))

# 0.5.1

BUG FIXES:

* tf5muxserver: Prevent `PrepareProviderConfig` RPC error for multiple `PreparedConfig` responses when combining terraform-plugin-sdk/v2 providers ([#54](https://github.com/hashicorp/terraform-plugin-mux/issues/54))
* tf6muxserver: Prevent `ValidateProviderConfig` RPC error for multiple `PreparedConfig` responses when combining terraform-plugin-framework providers ([#54](https://github.com/hashicorp/terraform-plugin-mux/issues/54))

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
