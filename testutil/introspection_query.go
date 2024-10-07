package testutil

import (
	_ "embed"
)

var (
	//go:embed introspection-query.graphql
	IntrospectionQuery string
)
