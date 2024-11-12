# GraphQL
This library is a fork of [graphql-go/graphql](https://github.com/graphql-go/graphql), and implements features outside of latest GraphQL [draft specification](https://spec.graphql.org/draft/). These non-standard features may break tooling when running introspection queries against applications using this library.

The only non-standard feature currently implemented is the concept of Applied Directives. It is implemented by following the [proposal specification](https://graphql-dotnet.github.io/docs/getting-started/directives/#directives-and-introspection). This feature is also implemented by [graphql-dotnet](https://github.com/graphql-dotnet/graphql-dotnet) and [graphql-java](https://github.com/graphql-java/graphql-java) libraries. More information on the proposal Applied Directive specification is available [here](https://graphql-dotnet.github.io/docs/getting-started/directives/#directives-and-introspection).

By default, the non-standard features are not returned during introspection and therefore _should_ not break tooling, however this cannot be guaranteed. At the very least it does not break [Altair](https://altairgraphql.dev/). In order to observe non-standard types during introspection, execute a query resembling:
```graphql
query CustomIntrospection {
    __schema(includeNonStandard: true) {
        types {
            name
        }
    }
}
```
