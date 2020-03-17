// http://www.apache.org/licenses/LICENSE-2.0
// https://raw.githubusercontent.com/go-jira/jira/d16bcc2f51f6fa69c47728433cb154f3501fc604/jiracli/log.go

package log

import (
	"os"
	"strconv"

	logging "gopkg.in/op/go-logging.v1"
)

var (
	log = logging.MustGetLogger("gobbs")
)

func IncreaseLogLevel(verbosity int) {
	logging.SetLevel(logging.GetLevel("")+logging.Level(verbosity), "")
}

func InitLogging() {
	logBackend := logging.NewLogBackend(os.Stderr, "", 0)
	format := os.Getenv("GOBBS_LOG_FORMAT")
	if format == "" {
		format = "%{color}%{level:-5s}%{color:reset} %{message}"
	}
	logging.SetBackend(
		logging.NewBackendFormatter(
			logBackend,
			logging.MustStringFormatter(format),
		),
	)
	if os.Getenv("GOBBS_DEBUG") == "" {
		logging.SetLevel(logging.NOTICE, "")
	} else {
		logging.SetLevel(logging.DEBUG, "")
		if verbosity, err := strconv.Atoi(os.Getenv("GOBBS_DEBUG")); err == nil {
			IncreaseLogLevel(verbosity)
		}
	}
}

// Fatal is equivalent to log.Print() followed by a call to os.Exit(1).
func Fatal(v ...interface{}) {
	log.Fatal(v...)
}

// Fatalf is equivalent to log.Print() followed by a call to os.Exit(1).
func Fatalf(format string, v ...interface{}) {
	log.Fatalf(format, v...)
}

// Debugf is equivalent to log.Print()
func Debugf(format string, v ...interface{}) {
	log.Debugf(format, v...)
}

// Infof is equivalent to log.Print()
func Infof(format string, v ...interface{}) {
	log.Infof(format, v...)
}

// Criticalf is equivalent to log.Print()
func Critical(format string, v ...interface{}) {
	log.Critical(format, v...)
}
