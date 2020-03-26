# Routing Non-Resource Requests

Terraform's protocol is mostly resource-scoped; requests are made in the
context of a single resource. This makes our lives easier, as it offers an
obvious way to route requests and discriminate between servers: based on which
resource we're dealing with.

Unfortunately, this is only _mostly_ the case. There are four RPC functions
that are _not_ resource-scoped, and are about the provider, not a specific
resource. We need to come up with a strategy for letting users decide how to
handle these requests that actually makes sense. It is not important that a
single strategy makes sense for all four functions; it is important that
whatever strategy each uses exposes patterns of behavior that are likely to be
correct, and guides users away from incorrect or invalid implementations. For
example, it's harmful to expose a strategy that would lead to one
ProviderServer responding unilaterally to `GetSchema` requests, as it can't
know about the other ProviderServer's resources and their schemas, so the
schema it would return would be incomplete.

## GetSchema

GetSchema is responsible for advertising support for resources, data sources,
provider configuration, etc. from within a provider. Whatever strategy we
support here should return a union of the information specified by each
ProviderServer. It's possible we could just call each ProviderServer's
GetSchema function, then combine the responses from them, and return the
combined response.

## PrepareProviderConfig

I am unclear on what this is for, but it appears to be passing information in
and then returning that information, perhaps modified. We should seek more
guidance on the purpose of this function before proposing a strategy.

## Configure

This is where the provider gets all the config information supplied in the
provider block. It is typically where API clients, etc. get instantiated.
Because clients need to be instantiated and becuase its Response is basically
just Diagnostics, we could in theory just call each of these in turn, returning
at the first errors-containing Diagnostics. That would instantiate all the
ProviderServers with the information from the config, and short-circuit on any
halt-the-world errors.

## Stop

Stop is used to gracefully shut down a ProviderServer. The request has no
arguments and the only information returned is an error. This should be safe
for us to call on each ProviderServer, one after the other, and just return
after any with a non-zero Error in their response.
