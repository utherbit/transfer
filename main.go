package main

import (
	"os"

	"github.com/utherbit/transfer/internal"
)

func main() {
	if err := internal.Run(os.Args[1:]); err != nil {
		os.Exit(1)
	}
}
