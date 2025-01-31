package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

func findFileWithType(dir, typeName string) (sourceFile string, packageName string, err error) {
	fset := token.NewFileSet()

	err = filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() || !strings.HasSuffix(info.Name(), ".go") || strings.HasSuffix(info.Name(), "_test.go") {
			return nil
		}

		node, err := parser.ParseFile(fset, path, nil, parser.PackageClauseOnly)
		if err != nil {
			return err
		}

		packageName = node.Name.Name

		node, err = parser.ParseFile(fset, path, nil, parser.AllErrors)
		if err != nil {
			return err
		}

		for _, decl := range node.Decls {
			genDecl, ok := decl.(*ast.GenDecl)
			if !ok || genDecl.Tok != token.TYPE {
				continue
			}

			for _, spec := range genDecl.Specs {
				typeSpec, ok := spec.(*ast.TypeSpec)
				if !ok || typeSpec.Name.Name != typeName {
					continue
				}

				// FIND!
				sourceFile = path

				return filepath.SkipAll
			}
		}

		return nil
	})
	if err != nil {
		return "", "", err
	}

	if sourceFile == "" {
		return "", "", fmt.Errorf("type %s not found in package %s", typeName, dir)
	}

	return sourceFile, packageName, nil
}

//nolint:funlen
func parseStruct(sourceFile, packageName, typeName string) (StructInfo, error) {
	fset := token.NewFileSet()

	node, err := parser.ParseFile(fset, sourceFile, nil, parser.AllErrors)
	if err != nil {
		return StructInfo{}, err
	}

	var (
		importMap  = map[string]string{}
		importsSet = map[string]string{} // Путь -> Псевдоним
	)

	// Собираем импорты из файла
	for _, imp := range node.Imports {
		path := strings.Trim(imp.Path.Value, `"`) // Убираем кавычки

		alias := ""
		if imp.Name != nil {
			alias = imp.Name.Name
		}

		importMap[aliasOrDefault(alias, path)] = path
		importsSet[path] = alias
	}

	var (
		usedImports = map[string]struct{}{} // Импорты, которые реально используются
		structInfo  = StructInfo{
			SourceFile: sourceFile,
			Package:    packageName,
		}
	)

	ast.Inspect(node, func(n ast.Node) bool {
		ts, ok := n.(*ast.TypeSpec)
		if !ok || ts.Name.Name != typeName {
			return true
		}

		st, ok := ts.Type.(*ast.StructType)
		if !ok {
			return true
		}

		structInfo.StructName = ts.Name.Name

		for _, field := range st.Fields.List {
			if field.Names != nil && isPrivate(field.Names[0].Name) {
				fieldType := formatExpr(field.Type, importMap, usedImports)
				fieldName := field.Names[0].Name
				structInfo.Fields = append(structInfo.Fields, Field{
					Name:   fieldName,
					Type:   fieldType,
					Getter: capitalize(fieldName),
					Setter: "Set" + capitalize(fieldName),
				})
			}
		}

		return false
	})

	for path, alias := range importsSet {
		if _, used := usedImports[path]; used {
			structInfo.Imports = append(structInfo.Imports, Import{Alias: alias, Path: path})
		}
	}

	// Сортируем импорты
	sort.Slice(structInfo.Imports, func(i, j int) bool {
		return structInfo.Imports[i].Path < structInfo.Imports[j].Path
	})

	if structInfo.StructName == "" {
		return StructInfo{}, fmt.Errorf("structure %s not found in %s", typeName, sourceFile)
	}

	return structInfo, nil
}

func formatExpr(expr ast.Expr, importMap map[string]string, usedImports map[string]struct{}) string {
	switch t := expr.(type) {
	case *ast.Ident:
		return t.Name

	case *ast.StarExpr:
		return "*" + formatExpr(t.X, importMap, usedImports)

	case *ast.SelectorExpr:
		pkgName := t.X.(*ast.Ident).Name
		if path, ok := importMap[pkgName]; ok {
			usedImports[path] = struct{}{} // Помечаем импорт как используемый
			return pkgName + "." + t.Sel.Name
		}

		return pkgName + "." + t.Sel.Name
	case *ast.MapType:
		key := formatExpr(t.Key, importMap, usedImports)
		value := formatExpr(t.Value, importMap, usedImports)

		return "map[" + key + "]" + value

	default:
		return fmt.Sprintf("%T", expr)
	}
}

func aliasOrDefault(alias, path string) string {
	if alias != "" {
		return alias
	}

	parts := strings.Split(path, "/")

	return parts[len(parts)-1]
}

func isPrivate(name string) bool {
	return strings.ToLower(string(name[0])) == string(name[0])
}

func capitalize(s string) string {
	if strings.HasPrefix(s, "id") {
		s = strings.Replace(s, "id", "ID", 1)
	}

	return strings.ToUpper(s[:1]) + s[1:]
}
