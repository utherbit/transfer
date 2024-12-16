package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"strings"
	"text/template"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run generate_transfer.go <type name>")
		os.Exit(1)
	}

	typeName := os.Args[1]

	currentDir, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	if err = generateTypeTransfer(currentDir, typeName); err != nil {
		panic(err)
	}
}

func generateTypeTransfer(currentDir, typeName string) error {
	sourceFile, packageName, err := findFileWithType(currentDir, typeName)
	if err != nil {
		return err
	}

	outputFile := strings.TrimSuffix(sourceFile, ".go") + "_transfer.go"

	structInfo, err := parseStruct(sourceFile, packageName, typeName)
	if err != nil {
		return err
	}

	tmpl, err := template.New("transfer").Parse(transferTemplate)
	if err != nil {
		return err
	}

	var buf bytes.Buffer

	err = tmpl.Execute(&buf, structInfo)
	if err != nil {
		return err
	}

	err = os.WriteFile(outputFile, buf.Bytes(), 0o600)
	if err != nil {
		return err
	}

	log.Printf("Файл %s успешно сгенерирован\n", outputFile)

	return nil
}
