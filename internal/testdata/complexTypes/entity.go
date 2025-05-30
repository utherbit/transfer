package complexTypes

type Entity struct {
	mapField   map[string]any
	sliceField []any
	arrayField [1]any

	structTypeField structType
	// Unsupported types:
	// genericField genericType[any]
	// structField struct{}
}

type (
	structType         struct{ Field any }
	genericType[T any] struct{ Type T }
)
