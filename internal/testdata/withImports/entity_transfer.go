// Code generated by github.com/utherbit/transfer; DO NOT EDIT.

package withImports

import (
	"io"
	aliasImport "os"

	"github.com/stretchr/testify/mock"
	"github.com/utherbit/transfer/internal/testdata/withImports/pac"
)

//go:generate go run github.com/utherbit/transfer --type Entity
type EntityDTO struct {
	ExternalImport *mock.Mock
	InternalImport pac.LocalType
	StdLibImport   io.Reader
	AliasImport    aliasImport.DirEntry
}

func (t *EntityDTO) Init(entity Entity) {
	t.ExternalImport = entity.externalImport
	t.InternalImport = entity.internalImport
	t.StdLibImport = entity.stdLibImport
	t.AliasImport = entity.aliasImport
}

func (t EntityDTO) Base() Entity {
	return Entity{
		externalImport: t.ExternalImport,
		internalImport: t.InternalImport,
		stdLibImport:   t.StdLibImport,
		aliasImport:    t.AliasImport,
	}
}
