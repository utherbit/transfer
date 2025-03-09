package internal

import (
	"fmt"
	"go/token"
	"io"
	"os"
	"strings"
	"text/template"

	"github.com/spf13/cobra"
)

func run(cmd *cobra.Command, args []string) {
	var (
		parseReq parseRequest
		err      error
	)

	// find and parse struct by input
	switch {
	case reference != "":
		parseReq, err = findStructByRef(reference)

	case typeName != "":
		currentDir, errWd := os.Getwd()
		if errWd != nil {
			panic(errWd)
		}
		parseReq, err = findStructByDirAndType(currentDir, typeName)

	default:
		fmt.Println("Должен быть обязательно передан один из флагов: type, ref'\n", cmd.Use)
	}

	if err != nil {
		panic(err)
	}

	info, err := parseStruct(parseReq)
	if err != nil {
		panic(err)
	}

	outputFileName := strings.TrimSuffix(info.SourceFile, ".go") + "_transfer.go"

	if err = generateTransferFileOut(info, outputFileName); err != nil {
		panic(err)
	}

	fmt.Printf("Transfer file generated: %s\n", token.Position{Filename: outputFileName, Line: 1})
}

func generateTransferStdout(info StructInfo) error {
	return generateTransfer(info, os.Stdout)
}

func generateTransferFileOut(info StructInfo, fileName string) error {
	file, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0o600)
	if err != nil {
		return err
	}

	defer file.Close()

	return generateTransfer(info, file)
}

func generateTransfer(info StructInfo, output io.Writer) error {
	tmpl, err := template.New("transfer").Parse(transferTemplate)
	if err != nil {
		return err
	}

	err = tmpl.Execute(output, info)
	if err != nil {
		return err
	}

	return nil
}

type parseRequest struct {
	Filename   string
	StructName string
}
