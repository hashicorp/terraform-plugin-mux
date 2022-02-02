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

## Getting Started

Functionality for a provider server is based on the protocol version. There are currently two main protocol versions in use today, protocol version 5 and protocol version 6, based on the development framework being used:

- [terraform-plugin-framework](https://pkg.go.dev/github.com/hashicorp/terraform-plugin-framework): Implements protocol version 6.
- [terraform-plugin-sdk/v2](https://pkg.go.dev/github.com/hashicorp/terraform-plugin-sdk/v2): Implements protocol version 5.
- [terraform-plugin-go](https://pkg.go.dev/github.com/hashicorp/terraform-plugin-go): Implements either protocol version, based on whether the `tf5server` package (protocol version 5) or `tf6server` package (protocol version 6) is being used.

To combine providers together, each must implement the same protocol version.

### Protocol Version 6

Protocol version 6 providers can be combined using the [`tf6muxserver.NewMuxServer` function](https://pkg.go.dev/github.com/hashicorp/terraform-plugin-mux/tf6muxserver#NewMuxServer):

```go
func main() {
	ctx := context.Background()
	providers := []func() tfprotov6.ProviderServer{
		// Example terraform-plugin-framework ProviderServer function
		// frameworkprovider.Provider().ProviderServer,
		//
		// Example terraform-plugin-go ProviderServer function
		// goprovider.Provider(),
	}

	muxServer, err := tf6muxserver.NewMuxServer(ctx, providers...)

	if err != nil {
		log.Fatalln(err.Error())
	}

	// Use the result to start a muxed provider
	err = tf6server.Serve("registry.terraform.io/namespace/example", muxServer.ProviderServer)

	if err != nil {
		log.Fatalln(err.Error())
	}
}
```

### Protocol Version 5

Protocol version 5 providers can be combined using the [`tf5muxserver.NewMuxServer` function](https://pkg.go.dev/github.com/hashicorp/terraform-plugin-mux/tf5muxserver#NewMuxServer):

```go
func main() {
	ctx := context.Background()
	providers := []func() tfprotov5.ProviderServer{
		// Example terraform-plugin-sdk ProviderServer function
		// sdkprovider.Provider().ProviderServer,
		//
		// Example terraform-plugin-go ProviderServer function
		// goprovider.Provider(),
	}
	muxServer, err := tf5muxserver.NewMuxServer(ctx, providers...)
	if err != nil {
		log.Fatalln(err.Error())
	}
	// Use the result to start a muxed provider
	err = tf5server.Serve("registry.terraform.io/namespace/example", muxServer.ProviderServer)
	if err != nil {
		log.Fatalln(err.Error())
	}
}
```

## Testing

The Terraform Plugin SDK's [`helper/resource`](https://pkg.go.dev/github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource) package can be used to test any provider that implements the [`tfprotov5.ProviderServer`](https://pkg.go.dev/github.com/hashicorp/terraform-plugin-go/tfprotov5#ProviderServer) interface, which includes muxed providers created using `tf5muxserver.NewMuxServer`.

You may wish to test a terraform-plugin-go provider's resources by supplying only that provider, and not the muxed provider, to the test framework: please see the example in https://github.com/hashicorp/terraform-plugin-go#testing in this case.

Otherwise, you should initialise a muxed provider in your testing code (conventionally in `provider_test.go`), and set this as the value of `ProtoV5ProviderFactories` in each [`TestCase`](https://pkg.go.dev/github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource#TestCase). For example:

```go
var testAccProtoV5ProviderFactories = map[string]func() (tfprotov5.ProviderServer, error){}

func init() {
  testAccProtoV5ProviderFactories["myprovider"] = func() (tfprotov5.ProviderServer, error) {
    ctx := context.Background()
    
    // the ProviderServer from SDKv2
    sdkv2 := sdkv2provider.Provider().GRPCProvider

    // the terraform-plugin-go provider
    tpg := protoprovider.Provider

    muxServer, err := tf5muxserver.NewMuxServer(ctx, sdkv2, tpg)
    if err != nil {
      return nil, err
    }
    return muxServer.ProviderServer(), nil
  }
}
```

Here each `TestCase` in which you want to use the muxed provider should include `ProtoV5ProviderFactories: testAccProtoV5ProviderFactories`. Note that the test framework will return an error if you attempt to register the same provider using both `ProviderFactories` and `ProtoV5ProviderFactories`.


## Documentation

Documentation is a work in progress. The GoDoc for packages, types, functions,
and methods should have complete information, but we're working to add a
section to [terraform.io](https://terraform.io/) with more information about
the module, its common uses, and patterns developers may wish to take advantage
of.

Please bear with us as we work to get this information published, and please
[open
issues](https://github.com/hashicorp/terraform-plugin-mux/issues/new/choose)
with requests for the kind of documentation you would find useful.

## Contributing

Please see [`.github/CONTRIBUTING.md`](https://github.com/hashicorp/terraform-plugin-mux/blob/master/.github/CONTRIBUTING.md).

## License

This module is licensed under the [Mozilla Public License v2.0](https://github.com/hashicorp/terraform-plugin-mux/blob/master/LICENSE).
