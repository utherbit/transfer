package complexTypes

type Entity struct {
	mapField   map[string]any
	sliceField []any
	arrayField [1]any

	structTypeField Struct

	genericField  GStruct[any]
	generic2Field G2Struct[any, any]

	// unsupported
	// structField struct{}
}

type (
	Struct               struct{}
	GStruct[T any]       struct{}
	G2Struct[T1, T2 any] struct{}
)
