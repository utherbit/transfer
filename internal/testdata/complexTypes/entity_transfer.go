// Code generated by github.com/utherbit/transfer; DO NOT EDIT.

package complexTypes

//go:generate go run github.com/utherbit/transfer --type Entity
type EntityDTO struct {
	MapField        map[string]any
	SliceField      []any
	ArrayField      [1]any
	StructTypeField Struct
	GenericField    GStruct[any]
	Generic2Field   G2Struct[any, any]
}

func (t *EntityDTO) Init(entity Entity) {
	t.MapField = entity.mapField
	t.SliceField = entity.sliceField
	t.ArrayField = entity.arrayField
	t.StructTypeField = entity.structTypeField
	t.GenericField = entity.genericField
	t.Generic2Field = entity.generic2Field
}

func (t EntityDTO) Base() Entity {
	return Entity{
		mapField:        t.MapField,
		sliceField:      t.SliceField,
		arrayField:      t.ArrayField,
		structTypeField: t.StructTypeField,
		genericField:    t.GenericField,
		generic2Field:   t.Generic2Field,
	}
}
