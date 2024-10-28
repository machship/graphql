# GraphQL

This library is a fork of [graphql-go/graphql](https://github.com/graphql-go/graphql). This fork implements features not listed in the latest GraphQL [draft specification](https://spec.graphql.org/draft/), and therefore will break tooling when running introspection queries against applications using this library.

The only non-standard feature currently implemented is the concept of Applied Directives. It adopts the proposed specification of Applied Directives as implemented by [graphql-dotnet](https://github.com/graphql-dotnet/graphql-dotnet) and [graphql-java](https://github.com/graphql-java/graphql-java).  Due to the way it has been implemented, there is no way to disable it at runtime (unlike the previously mentioned libraries). Therefore, this library should not be used for ordinairy GraphQL applications unless you know what you're doing. More information on the proposal Applied Directive specification is available [here](https://graphql-dotnet.github.io/docs/getting-started/directives/#directives-and-introspection).
