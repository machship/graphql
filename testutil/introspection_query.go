package testutil

import (
	_ "embed"
)

var (
	//go:embed introspection-query.graphql
	IntrospectionQuery string

	//go:embed applied-directive-introspection-query.graphql
	AppliedDirectiveIntrospectionQuery string
)
