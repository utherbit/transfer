package main

import (
	"fmt"
	"go/token"
	"io"
	"os"
	"strings"
	"text/template"

	"github.com/spf13/cobra"
)

func main() {
	RootCMD.SetArgs(os.Args[1:])

	if err := RootCMD.Execute(); err != nil {
		os.Exit(1)
	}
}

var RootCMD = &cobra.Command{
	Use:        "transfer {type name} | {reference}",
	Args:       cobra.ExactArgs(1),
	SuggestFor: []string{"te", "fe", "re"},
	Short:      `Generate transfer type`,
	Long: `use command like this:
go run klad.rupu.ru/rupuru/eda/backend/cmd/gen/transfer <type name> // from working directory
go run klad.rupu.ru/rupuru/eda/backend/cmd/gen/transfer <reference> // from project root directory
`,

	Run: run,
}

func run(cmd *cobra.Command, args []string) {
	if len(args) < 1 {
		fmt.Println("Usage: go run generate_transfer.go <type name or reference>")
		os.Exit(1)
	}

	arg := args[len(args)-1]

	var (
		parseReq parseRequest
		err      error
	)

	if isValidGoFilePosition(arg) {
		parseReq, err = findStructByRef(arg)
	} else {
		currentDir, errWd := os.Getwd()
		if errWd != nil {
			panic(errWd)
		}
		parseReq, err = findStructByDirAndType(currentDir, arg)
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
