## 0.21.0-alpha.1 (July 31, 2025)

NOTES:

* This alpha pre-release contains the `ListResource` RPC which returns a list of resource identities for a single managed resource type. ([#313](https://github.com/hashicorp/terraform-plugin-mux/issues/313))
* The `ListResource` and `ValidateListResourceConfig` RPCs are considered experimental and may change up until general availability. ([#310](https://github.com/hashicorp/terraform-plugin-mux/issues/310))
* This alpha pre-release contains the `Actions` RPC that allows practitioners to specify and invoke non-CRUD, ad-hoc operations that can cause changes to managed resources. ([#314](https://github.com/hashicorp/terraform-plugin-mux/issues/314))
* The `Actions` and `ValidateActionConfig` RPCs are considered experimental and may change up until general availability. ([#317](https://github.com/hashicorp/terraform-plugin-mux/issues/317))

## 0.20.0 (May 21, 2025)

BUG FIXES:

* all: Fixed a bug where muxed provider servers were not enforced to implement `GetResourceIdentitySchemas`, which is required by Terraform v1.12.1 in the scenario where at least one of the muxed provider servers supports identity. Before upgrading this dependency the Go modules that support identity should also be upgraded to prevent confusing errors, which are: terraform-plugin-go@v0.28.0, terraform-plugin-framework@v1.15.0, terraform-plugin-sdk/v2@v2.37.0, and terraform-plugin-testing@v1.13.0. ([#307](https://github.com/hashicorp/terraform-plugin-mux/issues/307))

## 0.19.0 (May 16, 2025)

NOTES:

* all: This Go module has been updated to Go 1.23 per the [Go support policy](https://go.dev/doc/devel/release#policy). It is recommended to review the [Go 1.23 release notes](https://go.dev/doc/go1.23) before upgrading. Any consumers building on earlier Go versions may experience errors. ([#291](https://github.com/hashicorp/terraform-plugin-mux/issues/291))

FEATURES:

* tf5muxserver+tf6muxserver+tf6to5server+tf5to6server: Upgraded protocols and added types to support the new resource identity feature ([#278](https://github.com/hashicorp/terraform-plugin-mux/issues/278))

## 0.19.0-alpha.1 (March 18, 2025)

NOTES:

* all: This Go module has been updated to Go 1.23 per the [Go support policy](https://go.dev/doc/devel/release#policy). It is recommended to review the [Go 1.23 release notes](https://go.dev/doc/go1.23) before upgrading. Any consumers building on earlier Go versions may experience errors. ([#291](https://github.com/hashicorp/terraform-plugin-mux/issues/291))
* This alpha pre-release contains the muxing logic for managed resource identity, which can used with Terraform v1.12.0-alpha20250312, to store and read identity data during plan and apply workflows. ([#278](https://github.com/hashicorp/terraform-plugin-mux/issues/278))

## 0.18.0 (January 23, 2025)

FEATURES:

* all: Upgrade protocol versions to support write-only attributes ([#272](https://github.com/hashicorp/terraform-plugin-mux/issues/272))

## 0.17.0 (October 30, 2024)

NOTES:

* all: This Go module has been updated to Go 1.22 per the [Go support policy](https://go.dev/doc/devel/release#policy). It is recommended to review the [Go 1.22 release notes](https://go.dev/doc/go1.22) before upgrading. Any consumers building on earlier Go versions may experience errors. ([#250](https://github.com/hashicorp/terraform-plugin-mux/issues/250))

FEATURES:

* all: Upgrade protocol versions to support ephemeral resource types ([#257](https://github.com/hashicorp/terraform-plugin-mux/issues/257))

## 0.16.0 (May 08, 2024)

NOTES:

* all: The `v0.15.0` release updated this Go module to Go 1.21 per the [Go support policy](https://go.dev/doc/devel/release#policy). It is recommended to review the [Go 1.21 release notes](https://go.dev/doc/go1.21) before upgrading. Any consumers building on earlier Go versions may experience errors ([#227](https://github.com/hashicorp/terraform-plugin-mux/issues/227))

ENHANCEMENTS:

* tf5to6server: Add deferred action request and response fields to RPC translations ([#237](https://github.com/hashicorp/terraform-plugin-mux/issues/237))
* tf6to5server: Add deferred action request and response fields to RPC translations ([#237](https://github.com/hashicorp/terraform-plugin-mux/issues/237))

## 0.15.0 (February 23, 2024)

ENHANCEMENTS:

* all: Upgrade protocol versions to support modified `CallFunction` RPC which returns a FunctionError rather than Diagnostics ([#226](https://github.com/hashicorp/terraform-plugin-mux/issues/226))

## 0.14.0 (January 29, 2024)

FEATURES:

* all: Upgrade protocol versions to support the `MoveResourceState` RPC ([#220](https://github.com/hashicorp/terraform-plugin-mux/issues/220))

## 0.13.0 (December 14, 2023)

NOTES:

* all: Update `google.golang.org/grpc` dependency to address CVE-2023-44487 ([#203](https://github.com/hashicorp/terraform-plugin-mux/issues/203))

FEATURES:

* all: Upgrade protocol versions to support provider-defined functions ([#209](https://github.com/hashicorp/terraform-plugin-mux/issues/209))

## 0.12.0 (September 06, 2023)

NOTES:

* all: This Go module has been updated to Go 1.20 per the [Go support policy](https://go.dev/doc/devel/release#policy). It is recommended to review the [Go 1.20 release notes](https://go.dev/doc/go1.20) before upgrading. Any consumers building on earlier Go versions may experience errors. ([#188](https://github.com/hashicorp/terraform-plugin-mux/issues/188))

FEATURES:

* all: Upgrade to protocol versions 5.4 and 6.4, which can significantly reduce memory usage with Terraform 1.6 and later when a configuration includes multiple instances of the same provider ([#185](https://github.com/hashicorp/terraform-plugin-mux/issues/185))

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
