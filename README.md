[![PkgGoDev](https://pkg.go.dev/badge/github.com/hashicorp/terraform-plugin-mux)](https://pkg.go.dev/github.com/hashicorp/terraform-plugin-mux)

# terraform-plugin-mux

terraform-plugin-mux provides a method for combining Terraform providers built
in multiple different SDKs and frameworks to be combined into a single logical
provider for Terraform to work with. It is designed to allow provider
developers to implement resources and data sources at the level of abstraction
that is most suitable for that specific resource or data source, and to allow
provider developers to upgrade between SDKs or frameworks on a
resource-by-resource basis instead of all at once.

## Status

terraform-plugin-mux is a [Go
module](https://github.com/golang/go/wiki/Modules) versioned using [semantic
versioning](https://semver.org/).

The module is currently on a v0 major version, indicating our lack of
confidence in the stability of its exported API. Developers depending on it
should do so with an explicit understanding that the API may change and shift
until we hit v1.0.0, as we learn more about the needs and expectations of
developers working with the module.

We are confident in the correctness of the code and it is safe to build on so
long as the developer understands that the API may change in backwards
incompatible ways and they are expected to be tracking these changes.

## Compatibility

Providers built on terraform-plugin-mux will only be usable with Terraform
v0.12.0 and later. Developing providers for versions of Terraform below 0.12.0
is unsupported by the Terraform Plugin SDK team.

Providers built on the Terraform Plugin SDK must be using version 2.2.0 of the
Plugin SDK or higher to be able to be used with terraform-plugin-mux.

## Go Compatibility

This project follows the [support policy](https://golang.org/doc/devel/release.html#policy) of Go as its support policy. The two latest major releases of Go are supported by the project.

Currently, that means Go **1.23** or later must be used when including this project as a dependency.

## Documentation

- **Website Documentation**: Getting started, usage, and testing information: [terraform.io](https://terraform.io/plugin/mux).
- **Go Documentation**: Go language types, functions, and method implementation information: [pkg.go.dev](https://pkg.go.dev/github.com/hashicorp/terraform-plugin-mux).
- **Documentation Requests**: Open a [GitHub issue](https://github.com/hashicorp/terraform-plugin-mux/issues/new/choose).

## Contributing

Refer to [`.github/CONTRIBUTING.md`](https://github.com/hashicorp/terraform-plugin-mux/blob/master/.github/CONTRIBUTING.md). The [website directory README](/website/README.md) contains details about how to contribute to the documentation on terraform.io.

## License

This module is licensed under the [Mozilla Public License v2.0](https://github.com/hashicorp/terraform-plugin-mux/blob/master/LICENSE).
