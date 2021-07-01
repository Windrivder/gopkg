package typex

// Placeholder is a placeholder object that can be used globally.
var Placeholder PlaceholderType

type (
	// GenericType can be used to hold any type.
	GenericType = interface{}
	// PlaceholderType represents a placeholder type.
	PlaceholderType = struct{}
	// DictType represents a container type.
	DictType = map[string]interface{}
)
