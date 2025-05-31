package withImports

import (
	"io"
	// aliasDuplicate "io" // aliasDuplicatePath unsupported
	aliasImport "os"

	"github.com/stretchr/testify/mock"

	"github.com/utherbit/transfer/internal/testdata/withImports/pac"
)

type Entity struct {
	// импорт внешней библиотеки
	externalImport *mock.Mock
	// импорт внутренней библиотеки (соседний пакет)
	internalImport pac.LocalType
	// импорт стандартной библиотеки
	stdLibImport io.Reader
	aliasImport  aliasImport.DirEntry
}
