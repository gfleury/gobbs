package main

import (
	"github.com/gfleury/gobbs/cmd"
	"github.com/gfleury/gobbs/common/log"
)

func main() {

	log.InitLogging()

	err := cmd.Execute()
	if err != nil {
		log.Fatal(err.Error())
	}

}
