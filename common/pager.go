package common

import (
	"io"
	"os"
	"os/exec"

	"github.com/gfleury/gobbs/common/log"
)

func Pager() (*exec.Cmd, io.WriteCloser) {
	less := exec.Command("less", "-r")

	stdin, err := less.StdinPipe()
	if err != nil {
		log.Critical(err.Error())
		stdin = os.Stdin
	}

	less.Stdout = os.Stdout
	return less, stdin
}