package internal

import (
	"fmt"
	"strings"
	"unicode"
)

const transferTemplate = `// Code generated by github.com/utherbit/transfer; DO NOT EDIT.

package {{.Package}}

{{ if .Imports -}}
import (
{{- range .Imports }}
	{{.}}
{{- end }}
)

{{ end -}}

` + /*
		Экранирование для избежания ложного срабатывания go:generate
	*/`//go:generate go run github.com/utherbit/transfer --type {{.StructName}}
type {{.StructName}}DTO struct { 
{{- range .Fields}}
	{{.PubName}} {{.Type}}
{{- end }}
}

func (t *{{.StructName}}DTO) Init(entity {{.StructName}}) { 
{{- range .Fields}}
	t.{{ .PubName }} = entity.{{ .Name }}
{{- end }}
}

func (t {{.StructName}}DTO) Base() {{.StructName}} {
	return {{ .StructName }}{
{{- range .Fields}}
		{{ .Name }}: t.{{.PubName}},
{{- end }}
	}
}

`

type StructInfo struct {
	SourceFile string
	Package    string
	StructName string
	// Fields only private
	Fields  []Field
	Imports []Import
}

type Field struct {
	Name string
	Type string
}

func (f *Field) PubName() string {
	pubName := f.Name

	if pubName == "" {
		return pubName
	}

	// ID всегда большими буквами
	if strings.HasPrefix(pubName, "id") {
		pubName = strings.Replace(pubName, "id", "ID", 1)
	}

	pubNameRunes := []rune(pubName)
	pubNameRunes[0] = unicode.ToUpper(pubNameRunes[0])

	return string(pubNameRunes)
}

type Import struct {
	Alias string
	Path  string
}

func (i *Import) String() string {
	if i.Alias == "" {
		return fmt.Sprintf(`"%s"`, i.Path)
	}

	return fmt.Sprintf(`%s "%s"`, i.Alias, i.Path)
}
