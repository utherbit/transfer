// Code generated by github.com/utherbit/transfer; DO NOT EDIT.

package ignoreUnderscore

//go:generate go run github.com/utherbit/transfer --type Entity
type EntityDTO struct {
	_ignorableField int
}

func (t *EntityDTO) Init(entity Entity) {
	t._ignorableField = entity._ignorableField
}

func (t EntityDTO) Base() Entity {
	return Entity{
		_ignorableField: t._ignorableField,
	}
}
