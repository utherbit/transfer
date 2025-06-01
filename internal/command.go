package internal

import "github.com/spf13/cobra"

var (
	reference string
	typeName  string
)

var TransferCMD = &cobra.Command{
	Use:        "transfer --type T",
	SuggestFor: []string{"te", "fe", "re"},
	Short:      `Generate transfer type`,
	Long: `use command like this:
go run github.com/utherbit/transfer --type <type name> // from working directory
go run github.com/utherbit/transfer --ref <reference> // from project root directory
`,

	Run: run,
}

func init() {
	TransferCMD.Flags().StringVarP(&typeName, "type", "t", typeName,
		"Название типа, указывается при вызове из директории в которой находиться указанный тип.")
	TransferCMD.Flags().StringVarP(&reference, "ref", "r", reference,
		"Ссылка на тип, указывается из любой директории, reference должен включать в себя путь до .go файла и строку на которой находиться нужный тип.")
}

func Run(args []string) error {
	TransferCMD.SetArgs(args)
	return TransferCMD.Execute()
}
