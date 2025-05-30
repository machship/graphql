package graphql

const (
	// Operations
	DirectiveLocationQuery              = "QUERY"
	DirectiveLocationMutation           = "MUTATION"
	DirectiveLocationSubscription       = "SUBSCRIPTION"
	DirectiveLocationField              = "FIELD"
	DirectiveLocationFragmentDefinition = "FRAGMENT_DEFINITION"
	DirectiveLocationFragmentSpread     = "FRAGMENT_SPREAD"
	DirectiveLocationInlineFragment     = "INLINE_FRAGMENT"

	// Schema Definitions
	DirectiveLocationSchema               = "SCHEMA"
	DirectiveLocationScalar               = "SCALAR"
	DirectiveLocationObject               = "OBJECT"
	DirectiveLocationFieldDefinition      = "FIELD_DEFINITION"
	DirectiveLocationArgumentDefinition   = "ARGUMENT_DEFINITION"
	DirectiveLocationInterface            = "INTERFACE"
	DirectiveLocationUnion                = "UNION"
	DirectiveLocationEnum                 = "ENUM"
	DirectiveLocationEnumValue            = "ENUM_VALUE"
	DirectiveLocationInputObject          = "INPUT_OBJECT"
	DirectiveLocationInputFieldDefinition = "INPUT_FIELD_DEFINITION"
)

// DefaultDeprecationReason Constant string used for default reason for a deprecation.
const DefaultDeprecationReason = "No longer supported"

// SpecifiedRules The full list of specified directives.
var SpecifiedDirectives = []*Directive{
	IncludeDirective,
	SkipDirective,
	OmitEmptyDirective,
	DeprecatedDirective,
}

// Directive structs are used by the GraphQL runtime as a way of modifying execution
// behavior. Type system creators will usually not create these directly.
type Directive struct {
	Name        string      `json:"name"`
	Description string      `json:"description"`
	Locations   []string    `json:"locations"`
	Args        []*Argument `json:"args"`

	err error

	directives []*AppliedDirective
}

// DirectiveConfig options for creating a new GraphQLDirective
type DirectiveConfig struct {
	Name        string              `json:"name"`
	Description string              `json:"description"`
	Locations   []string            `json:"locations"`
	Args        FieldConfigArgument `json:"args"`
	Directives  []*AppliedDirective
}

func NewDirective(config DirectiveConfig) *Directive {
	dir := &Directive{}

	// Ensure directive is named
	if dir.err = invariant(config.Name != "", "Directive must be named."); dir.err != nil {
		return dir
	}

	// Ensure directive name is valid
	if dir.err = assertValidName(config.Name); dir.err != nil {
		return dir
	}

	// Ensure locations are provided for directive
	if dir.err = invariant(len(config.Locations) > 0, "Must provide locations for directive."); dir.err != nil {
		return dir
	}

	args := []*Argument{}

	for argName, argConfig := range config.Args {
		if dir.err = assertValidName(argName); dir.err != nil {
			return dir
		}
		args = append(args, &Argument{
			PrivateName:        argName,
			PrivateDescription: argConfig.Description,
			Type:               argConfig.Type,
			DefaultValue:       argConfig.DefaultValue,
		})
	}

	dir.Name = config.Name
	dir.Description = config.Description
	dir.Locations = config.Locations
	dir.Args = args
	dir.directives = config.Directives
	return dir
}

// AppliedDirectives returns the directives that have been applied to this directive.
func (d *Directive) AppliedDirectives() []*AppliedDirective {
	return d.directives
}

// Apply transforms a directive into an applied directive, which can be used to
// modify execution behavior. This method is intended to be used to apply a directive
// to a field, fragment, or type (i.e. anywhere a directive can be applied).
func (d *Directive) Apply(args []*DirectiveArgument) *AppliedDirective {
	return &AppliedDirective{
		Name:        d.Name,
		Description: d.Description,
		Args:        args,
	}
}

// AppliedDirective is a directive that has been applied to a field, fragment, or type.
type AppliedDirective struct {
	Name        string
	Description string
	Args        []*DirectiveArgument
}

// DirectiveArgument is an argument to a directive. The value is always represented as a
// string in graphQL queries, but may be of any type.
type DirectiveArgument struct {
	Name  string
	Value any
}

// IncludeDirective is used to conditionally include fields or fragments.
var IncludeDirective = NewDirective(DirectiveConfig{
	Name: "include",
	Description: "Directs the executor to include this field or fragment only when " +
		"the `if` argument is true.",
	Locations: []string{
		DirectiveLocationField,
		DirectiveLocationFragmentSpread,
		DirectiveLocationInlineFragment,
	},
	Args: FieldConfigArgument{
		"if": &ArgumentConfig{
			Type:        NewNonNull(Boolean),
			Description: "Included when true.",
		},
	},
})

// SkipDirective Used to conditionally skip (exclude) fields or fragments.
var SkipDirective = NewDirective(DirectiveConfig{
	Name: "skip",
	Description: "Directs the executor to skip this field or fragment when the `if` " +
		"argument is true.",
	Args: FieldConfigArgument{
		"if": &ArgumentConfig{
			Type:        NewNonNull(Boolean),
			Description: "Skipped when true.",
		},
	},
	Locations: []string{
		DirectiveLocationField,
		DirectiveLocationFragmentSpread,
		DirectiveLocationInlineFragment,
	},
})

// DeprecatedDirective  Used to declare element of a GraphQL schema as deprecated.
var DeprecatedDirective = NewDirective(DirectiveConfig{
	Name:        "deprecated",
	Description: "Marks an element of a GraphQL schema as no longer supported.",
	Args: FieldConfigArgument{
		"reason": &ArgumentConfig{
			Type: String,
			Description: "Explains why this element was deprecated, usually also including a " +
				"suggestion for how to access supported similar data. Formatted" +
				"in [Markdown](https://daringfireball.net/projects/markdown/).",
			DefaultValue: DefaultDeprecationReason,
		},
	},
	Locations: []string{
		DirectiveLocationFieldDefinition,
		DirectiveLocationEnumValue,
	},
})

// SkipDirective Used to conditionally skip (exclude) fields or fragments.
var OmitEmptyDirective = NewDirective(DirectiveConfig{
	Name: "omitEmpty",
	Description: "Directs the executor to omit this field or fragment when the response " +
		"value resolves to null or an empty list.",
	Locations: []string{
		DirectiveLocationField,
		DirectiveLocationFragmentSpread,
		DirectiveLocationInlineFragment,
	},
})
