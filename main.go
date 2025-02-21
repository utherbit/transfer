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
	Use:   "transfer",
	Short: "Generate transfer type",
	Run:   genTransferRun,
}

func genTransferRun(cmd *cobra.Command, args []string) {
	if len(args) < 1 {
		fmt.Println("Usage: go run generate_transfer.go <type name>")
		os.Exit(1)
	}

	arg := args[len(args)-1]

	isRef := strings.Contains(arg, ":")

	var (
		parseReq parseRequest
		err      error
	)

	if isRef {
		parseReq, err = findStructByRef(arg)
	} else {
		currentDir, err := os.Getwd()
		if err != nil {
			panic(err)
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

	fmt.Printf("Transfer file generated: %s", token.Position{Filename: outputFileName, Line: 1})
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
