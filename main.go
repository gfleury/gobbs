package main

import (
	"os"

	"github.com/gfleury/gobbs/cmd"
	"github.com/gfleury/gobbs/common/log"
)

func main() {

	log.InitLogging()

	err := cmd.Execute()
	if err != nil {
		os.Exit(1)
	}

}
