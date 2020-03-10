package util

import (
	"os"

	"github.com/op/go-logging"
)

// icannzplit global logger
var (
	Log *logging.Logger
)

func init() {
	Log = CreateLeveledLog("icannzplit", logging.DEBUG)
}

// CreateLeveledLog returns a new logger
func CreateLeveledLog(module string, level logging.Level) *logging.Logger {
	log := logging.MustGetLogger(module)
	format := logging.MustStringFormatter(
		`%{color}%{time:15:04:05.000} ▶ %{level:.4s} %{id:03x}%{color:reset} %{module}:%{shortpkg} >> %{message}`,
	)
	backend := logging.NewLogBackend(os.Stdout, "", 0)
	formatted := logging.NewBackendFormatter(backend, format)
	leveled := logging.AddModuleLevel(formatted)
	leveled.SetLevel(level, "")
	log.SetBackend(leveled)
	return log
}
