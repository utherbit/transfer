package main

import (
	"os"

	"klad.rupu.ru/rupuru/eda/backend/cmd/gen/transfer/internal"
)

func main() {
	if err := internal.Run(os.Args[1:]); err != nil {
		os.Exit(1)
	}
}
