package main

import (
	"os"
	"strings"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/pkgerrors"
	"github.com/rs/zerolog/log"
	"github.com/springload/ssm-parent/cmd"
)

var version = "dev"

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack

	logLevel := strings.ToLower(os.Getenv("LOG_LEVEL"))
	if logLevel == "panic" || logLevel == "5" {
		zerolog.SetGlobalLevel(zerolog.PanicLevel)
	} else if logLevel == "fatal" || logLevel == "4" {
		zerolog.SetGlobalLevel(zerolog.FatalLevel)
	} else if logLevel == "error" || logLevel == "3" {
		zerolog.SetGlobalLevel(zerolog.ErrorLevel)
	} else if logLevel == "warn" || logLevel == "2" {
		zerolog.SetGlobalLevel(zerolog.WarnLevel)
	} else if logLevel == "info" || logLevel == "1" {
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	} else if logLevel == "debug" || logLevel == "0" {
		zerolog.SetGlobalLevel(zerolog.TraceLevel)
	} else if logLevel == "trace" || logLevel == "-1" {
		zerolog.SetGlobalLevel(zerolog.TraceLevel)
	} else {
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}

	cmd.Execute(version)
}
