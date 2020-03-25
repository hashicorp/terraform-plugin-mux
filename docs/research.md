# Working Notes

We start with [`terraform-plugin-sdk/plugin#Serve`][sdk/plugin-serve], which is
a loose wrapper around [`go-plugin#Serve`][go-plugin-serve].
[`go-plugin#Serve`][go-plugin-serve] does all the go-plugin shenanigans of
launching the plugin server process, but it is basically agnostic to the GRPC
server it launches, accepting it as an argument to the server config.

Back in [`terraform-plugin-sdk/plugin#Serve`][sdk/plugin-serve], we're passing
[`go-plugin.DefaultGRPCServer`][go-plugin-defaultgrpcserver] as the config.
[`go-plugin.DefaultGRPCServer`][go-plugin-defaultgrpcserver] is an alias of
[`grpc.NewServer`][grpc-newserver], which does a lot of GRPC stuff that is,
once again, pretty implementation agnostic.

We're in the guts of how this server gets set up, but we're still no closer to
what we want, which is to sit at the layer where the GRPC server turns into a
Terraform server, specifically, but a layer above where our current
implementation is.

If we look back at [`terraform-plugin-sdk/plugin#Serve`][sdk/plugin-serve], we
see that there's a VersionedPlugins key we're setting, to serve different
things based on the plugin version negotiated.

Tracing it down through, it looks like we're serving a
[`terraform-plugin-sdk/plugin#GRPCProviderFunc`][sdk/plugin-grpcproviderfunc]
when protocol 5 is negotiated. That is just a type alias for
[`terraform-plugin-sdk/internal/tfplugin5#ProviderServer`][tfplugin5-providerserver].
Finally, it looks like we're getting somewhere--filling an interface created by
the generated Go code for our protobuf files is a Good Sign.

[`terraform-plugin-sdk/internal/tfplugin5#ProviderServer`][tfplugin5-providerserver]
is an interface that we finally can recognise from the protocol docs. It gets
us to our `GetSchema`, `PrepareProviderConfig`, `ApplyResourceChange`, etc.
methods.  All the input types are in the
`terraform-plugin-sdk/internal/tfplugin5` package, and are generated from our
protobuf files. Because of this, we can rely on them underpinning all the
plugin server implementations compatible with protocol version 5.

So now the question becomes how to mux over them. Ideally, we'd just use the
types from the `tfplugin5` package, which would let us transparently pass
things through.  Unfortunately, packaging makes that difficult. We want the
router to live in a separate package, `terraform-plugin-mux`, but the
`tfplugin5` package is in the `terraform-plugin-sdk` repo. We want to avoid
importing the `terraform-plugin-sdk` repo into `terraform-plugin-mux`, because
we want `terraform-plugin-mux` to exist outside versions of the SDK, as it's
being created to bridge them. Further, we _can't_ import `tfplugin5` into
`terraform-plugin-mux`, because it's in the `internal` directory, meaning it's
unavailable for import. This is by design, as we odn't want to imply that the
package provides a stable interface, as it's generated.

We also don't want to have multiple generated copies of the `tfplugin5` package
in multiple repos, because first that sounds like it'll be a mess to keep them
in sync, and second it won't work anyways, as in Go the package's path is part
of the type information, so
`github.com/hashicorp/terraform-plugin-sdk/internal/tfplugin5.ProviderServer`
and
`github.com/hashicorp/terraform-plugin-mux/internal/tfplugin5.ProviderServer`
are different types, even if the code is exactly the same.

So our challenge becomes how do we have a single `tfplugin5` package that can
be imported by both repos but doesn't imply a stable interface for other
consumers.

One option may be to sidestep some of these requirements using [type
aliasing][type-aliasing].  More investigation of that is needed.

Another option may be using a "soft" solution--putting `tfplugin5` in its own
module, then importing that module into `terraform-plugin-sdk` and
`terraform-plugin-mux`.  We can still dissuade its use and clearly state that
it has no stable interface or backwards compatibility guarantees by:

1. Noting as such in the README.
2. Keeping all releases in the v0.X.Y range, which according to Go Modules'
   interpretation of semver allows for breaking changes without bumping the
   version number.

[go-plugin-serve]: https://pkg.go.dev/github.com/hashicorp/go-plugin?tab=doc#Serve
[sdk/plugin-serve]: https://pkg.go.dev/github.com/hashicorp/terraform-plugin-sdk/plugin?tab=doc#Serve
[go-plugin-defaultgrpcserver]: https://pkg.go.dev/github.com/hashicorp/go-plugin?tab=doc#DefaultGRPCServer
[grpc-newserver]: https://pkg.go.dev/google.golang.org/grpc?tab=doc#NewServer
[sdk/plugin-grpcproviderfunc]: https://pkg.go.dev/github.com/hashicorp/terraform-plugin-sdk/plugin?tab=doc#GRPCProviderFunc
[tfplugin5-providerserver]: https://pkg.go.dev/github.com/hashicorp/terraform-plugin-sdk/internal/tfplugin5?tab=doc#ProviderServer
[type-aliasing]: https://golang.org/doc/go1.9#language
