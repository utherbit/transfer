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

var (
	reference string
	typeName  string
)

var RootCMD = &cobra.Command{
	Use:        "transfer --type T",
	SuggestFor: []string{"te", "fe", "re"},
	Short:      `Generate transfer type`,
	Long: `use command like this:
go run klad.rupu.ru/rupuru/eda/backend/cmd/gen/transfer --type <type name> // from working directory
go run klad.rupu.ru/rupuru/eda/backend/cmd/gen/transfer --ref <reference> // from project root directory
`,

	Run: run,
}

func main() {
	RootCMD.SetArgs(os.Args[1:])

	if err := RootCMD.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	RootCMD.Flags().StringVarP(&typeName, "type", "t", typeName,
		"Название типа, указывается при вызове из директории в которой находиться указанный тип.")
	RootCMD.Flags().StringVarP(&reference, "ref", "r", reference,
		"Ссылка на тип, указывается из любой директории, reference должен включать в себя путь до .go файла и строку на которой находиться нужный тип.")
}

func run(cmd *cobra.Command, args []string) {
	//if len(args) < 1 {
	//	fmt.Println("Usage: go run generate_transfer.go <type name or reference>")
	//	os.Exit(1)
	//}
	//
	//arg := args[len(args)-1]

	var (
		parseReq parseRequest
		err      error
	)

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
		fmt.Println("Должен быть обязательно передан один из флагов: type, ref'")
	}

	//if isValidGoFilePosition(arg) {
	//	parseReq, err = findStructByRef(arg)
	//} else {
	//	currentDir, errWd := os.Getwd()
	//	if errWd != nil {
	//		panic(errWd)
	//	}
	//	parseReq, err = findStructByDirAndType(currentDir, arg)
	//}

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
