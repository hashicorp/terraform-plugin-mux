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

Currently, that means Go **1.17** or later must be used when including this project as a dependency.

## Getting Started

Functionality for a provider server is based on the protocol version. There are currently two main protocol versions in use today, protocol version 5 and protocol version 6, based on the development framework being used:

- [terraform-plugin-framework](https://pkg.go.dev/github.com/hashicorp/terraform-plugin-framework): Implements protocol version 6.
- [terraform-plugin-sdk/v2](https://pkg.go.dev/github.com/hashicorp/terraform-plugin-sdk/v2): Implements protocol version 5.
- [terraform-plugin-go](https://pkg.go.dev/github.com/hashicorp/terraform-plugin-go): Implements either protocol version, based on whether the `tf5server` package (protocol version 5) or `tf6server` package (protocol version 6) is being used.

To combine providers together, each must implement the same protocol version. Different protocol version providers can be combined by either [upgrading protocol version 5 server to version 6](#upgrading-protocol-version-5-server-to-version-6) or [downgrading protocol version 6 server to version 5](#downgrading-protocol-version-6-server-to-version-5), depending on the appropriate use case.

### Combining Protocol Version 6 Servers

Protocol version 6 providers can be combined using the [`tf6muxserver.NewMuxServer` function](https://pkg.go.dev/github.com/hashicorp/terraform-plugin-mux/tf6muxserver#NewMuxServer):

```go
func main() {
	ctx := context.Background()
	providers := []func() tfprotov6.ProviderServer{
		// Example terraform-plugin-framework ProviderServer function
		// func() tfprotov6.ProviderServer {
		//   return tfsdk.NewProtocol6Server(frameworkprovider.New("version")())
		// },
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

### Combining Protocol Version 5 Servers

Protocol version 5 providers can be combined using the [`tf5muxserver.NewMuxServer` function](https://pkg.go.dev/github.com/hashicorp/terraform-plugin-mux/tf5muxserver#NewMuxServer):

```go
func main() {
	ctx := context.Background()
	providers := []func() tfprotov5.ProviderServer{
		// Example terraform-plugin-sdk ProviderServer function
		// sdkprovider.New("version")().GRPCProvider,
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

### Upgrading Protocol Version 5 Server to Version 6

Protocol version 5 servers can be upgraded to protocol version 6 using the [`tf5to6server.UpgradeServer()` function](https://pkg.go.dev/github.com/hashicorp/terraform-plugin-mux/tf5to6server#UpgradeServer). For example, this enables a terraform-plugin-sdk/v2 provider to be combined with a terraform-plugin-framework provider, while taking advantage of newer protocol features such as nested attributes for newer/migrated resources and practitioners can still use older/unmigrated resources. Since protocol version 6 is forwards compatible with protocol version 5, no additional validation is required during upgrade.

_**NOTE:** While protocol version 6 servers are compatible with Terraform CLI 1.0 and later, terraform-plugin-sdk/v2 servers that are upgraded to protocol version 6 require Terraform CLI 1.1.15 and later._

```go
ctx := context.Background()

// Example terraform-plugin-sdk/v2 upgrade
upgradedSdkProvider, err := tf5to6server.UpgradeServer(ctx, sdkprovider.New("version")().GRPCProvider)

if err != nil {
	log.Fatal(err)
}

providers := []func() tfprotov6.ProviderServer{
	upgradedSdkProvider.ProviderServer,

	// Example terraform-plugin-framework provider
	func() tfprotov6.ProviderServer {
		return tfsdk.NewProtocol6Server(frameworkprovider.New("version")())
	},
}

muxServer, err := tf6muxserver.NewMuxServer(ctx, providers...)
// ...
```

### Downgrading Protocol Version 6 Server to Version 5

Protocol version 6 servers can be downgraded to protocol version 5 using the [`tf6to5server.DowngradeServer()` function](https://pkg.go.dev/github.com/hashicorp/terraform-plugin-mux/tf6to5server#DowngradeServer). For example, this enables a terraform-plugin-sdk/v2 provider to be combined with a terraform-plugin-framework provider, while keeping compatibility with Terraform CLI 0.12 and later versions. Since protocol version 6 is not fully backwards compatible with protocol version 5, additional validation is performed upfront to verify that unsupported features, such as nested attributes, are not being used.

```go
ctx := context.Background()

// Example terraform-plugin-framework downgrade
downgradedFrameworkProvider, err := tf6to5server.DowngradeServer(ctx, func() tfprotov6.ProviderServer {
	return tfsdk.NewProtocol6Server(frameworkprovider.New("version")())
})

if err != nil {
	log.Fatal(err)
}

providers := []func() tfprotov6.ProviderServer{
	downgradedFrameworkProvider.ProviderServer,

	// Example terraform-plugin-sdk/v2 provider
	sdkprovider.New("version")().GRPCProvider,
}

muxServer, err := tf5muxserver.NewMuxServer(ctx, providers...)
// ...
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

## Debugging

The terraform-plugin-go [`tf5server`](https://pkg.go.dev/github.com/hashicorp/terraform-plugin-go/tfprotov5/tf5server) and [`tf6server`](https://pkg.go.dev/github.com/hashicorp/terraform-plugin-go/tfprotov6/tf6server) packages used to serve providers support the additional [`tf5server.WithManagedDebug()`](https://pkg.go.dev/github.com/hashicorp/terraform-plugin-go/tfprotov5/tf5server#WithManagedDebug) and [`tf6server.WithManagedDebug()`](https://pkg.go.dev/github.com/hashicorp/terraform-plugin-go/tfprotov6/tf6server#WithManagedDebug) server options to enable debugger support.

For example, a provider `main` function that takes a `-debug` flag to enable debugging support:

```go
func main() {
	debugFlag := flag.Bool("debug", false, "Start provider in debug mode.")
	flag.Parse()

	ctx := context.Background()
	providers := []func() tfprotov6.ProviderServer{
		// Example terraform-plugin-framework ProviderServer function
		// func() tfprotov6.ProviderServer {
		//   return tfsdk.NewProtocol6Server(frameworkprovider.New("version")())
		// },
		//
		// Example terraform-plugin-go ProviderServer function
		// goprovider.Provider(),
	}

	muxServer, err := tf6muxserver.NewMuxServer(ctx, providers...)

	if err != nil {
		log.Fatalln(err.Error())
	}

	var serveOpts []tf6server.ServeOpt

	if *debugFlag {
		serveOpts = append(serveOpts, tf6server.WithManagedDebug())
	}

	err = tf6server.Serve("registry.terraform.io/namespace/example", muxServer.ProviderServer, serveOpts...)

	if err != nil {
		log.Fatalln(err.Error())
	}
}
```

Refer to the [Terraform documentation](https://www.terraform.io/plugin/sdkv2/debugging) for more information about how to start and attach a debugger once the `main` function is setup.

## Documentation

Refer to the GoDoc for information about packages, types, functions,
and methods. We also have documentation about how to use the module, 
including use cases and examples, on [terraform.io](https://www.terraform.io/plugin/mux).

Please [open issues](https://github.com/hashicorp/terraform-plugin-mux/issues/new/choose)
with requests for additional documentation you would find useful.

## Contributing

Refer to [`.github/CONTRIBUTING.md`](https://github.com/hashicorp/terraform-plugin-mux/blob/master/.github/CONTRIBUTING.md). The [website directory README](/website.README.md) contains details about how to contribute to the documentation on terraform.io.

## License

This module is licensed under the [Mozilla Public License v2.0](https://github.com/hashicorp/terraform-plugin-mux/blob/master/LICENSE).
