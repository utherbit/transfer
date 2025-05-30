package withImports

import (
	"io"
	// aliasDuplicate "io" // aliasDuplicatePath unsupported
	aliasImport "os"

	"github.com/stretchr/testify/mock"

	"klad.rupu.ru/rupuru/eda/backend/cmd/gen/transfer/internal/testdata/withImports/localimport"
)

type Entity struct {
	// импорт внешней библиотеки
	externalImport *mock.Mock
	// импорт внутренней библиотеки (соседний пакет)
	internalImport localimport.LocalType
	// импорт стандартной библиотеки
	stdLibImport io.Reader
	aliasImport  aliasImport.DirEntry
}
