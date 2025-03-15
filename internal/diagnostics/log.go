package diagnostics

import (
	"github.com/apex/log"
	"github.com/apex/log/handlers/json"
	"github.com/davidalpert/go-printers/v1"
	"github.com/davidalpert/go-yeet/internal/diagnostics/plaintext"
	"github.com/davidalpert/go-yeet/internal/env"
	"github.com/davidalpert/go-yeet/internal/version"
	"gopkg.in/op/go-logging.v1"
	"io"
	"os"
	"path/filepath"
	"strings"
)

var Log *log.Entry

const (
	ENVKEY_LOG_LEVEL  = "YEET_LOG_LEVEL"
	ENVKEY_LOG_FORMAT = "YEET_LOG_FORMAT"
	ENVKEY_LOG_FILE   = "YEET_LOG_FILE"
)

func init() {
	Log = log.WithFields(log.Fields{
		"app":         version.Detail.AppName,
		"app_version": version.Detail.Version,
	})

}

func ConfigureLogger(streams printers.IOStreams) (cleanupFn func()) {
	// default cleanup: nothing to do
	cleanupFn = func() {}

	// configure logging
	logLevel := env.GetValueOrDefaultLogLevel(ENVKEY_LOG_LEVEL, log.FatalLevel)
	log.SetLevel(logLevel)
	// log sink
	var sink io.Writer
	var logFile = env.GetValueOrDefaultString(ENVKEY_LOG_FILE, "")
	var logDestination = "stdout"
	if logFile == "" {
		sink = streams.Out
	} else {
		fullPath, err := filepath.Abs(logFile)
		if err != nil {
			log.Fatal(err.Error())
		}
		logFile, err := os.Create(fullPath)
		if err != nil {
			log.Fatal(err.Error())
		}
		cleanupFn = func() { logFile.Close() }
		logDestination = fullPath
		sink = logFile
	}

	if strings.EqualFold(env.GetValueOrDefaultString(ENVKEY_LOG_FORMAT, "text"), "json") {
		log.SetHandler(json.New(sink))
	} else {
		log.SetHandler(plaintext.New(sink))
	}

	Log.WithField("destination", logDestination).Debug("logging initialized")

	// ----------------------------
	// configure go-logging to use a matching backend and level
	//

	// Example format string. Everything except the message has a custom color
	// which is dependent on the log level. Many fields have a custom output
	// formatting too, eg. the time returns the hour down to the milli second.
	var format = logging.MustStringFormatter(
		`%{color}%{time:15:04:05.000} %{shortfunc} ▶ %{level:.4s} %{id:03x}%{color:reset} %{message}`,
	)

	// For demo purposes, create two backend for os.Stderr.
	backend2 := logging.NewLogBackend(sink, "", 0)
	backend2Leveled := logging.AddModuleLevel(backend2)
	loggingLevel := mapLogLevelToLoggingLevel(logLevel)
	backend2Leveled.SetLevel(loggingLevel, "")

	// For messages written to backend2 we want to add some additional
	// information to the output, including the used log level and the name of
	// the function.
	backend2Formatted := logging.NewBackendFormatter(backend2Leveled, format)

	// Set the backends to be used.
	logging.SetBackend(backend2Formatted)

	return
}

func mapLogLevelToLoggingLevel(lvl log.Level) logging.Level {
	switch lvl {
	case log.DebugLevel:
		return logging.DEBUG
	case log.InfoLevel:
		return logging.INFO
	case log.WarnLevel:
		return logging.WARNING
	case log.ErrorLevel:
		return logging.ERROR
	default:
		return logging.CRITICAL
	}
}