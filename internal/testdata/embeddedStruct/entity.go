package embeddedStruct

import (
	"github.com/utherbit/transfer/internal/testdata/embeddedStruct/pac"
)

type Entity struct {
	// Ident
	BaseStruct
	// pointer
	*BaseStruct2
	// selector
	pac.Struct
	// generic
	GType[BaseStruct]
}

type (
	BaseStruct   struct{}
	BaseStruct2  struct{}
	GType[T any] struct{}
)
