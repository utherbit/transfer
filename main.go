package main

import (
	"fmt"
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

	typeName := args[len(args)-1]

	currentDir, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	info, err := parseType(currentDir, typeName)
	if err != nil {
		panic(err)
	}

	outputFileName := strings.TrimSuffix(info.SourceFile, ".go") + "_transfer.go"

	if err = generateTransferFileOut(*info, outputFileName); err != nil {
		panic(err)
	}
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

func parseType(currentDir, typeName string) (*StructInfo, error) {
	sourceFile, packageName, err := findFileWithType(currentDir, typeName)
	if err != nil {
		return nil, err
	}

	structInfo, err := parseStruct(sourceFile, packageName, typeName)
	if err != nil {
		return nil, err
	}

	return &structInfo, nil
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
