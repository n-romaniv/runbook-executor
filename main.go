package main

import (
	"context"
	"flag"
	"log/slog"
	"os"

	"github.com/n-romaniv/runbook-executor/statemachine"
)

func main() {
	flag.Parse()

	defFile := flag.Arg(0)

	if defFile == "" {
		slog.Error("state machine definition is required")
		os.Exit(1)
	}

	def, err := os.ReadFile(defFile)
	if err != nil {
		slog.Error("error when reading the definition file:", "err", err)
		os.Exit(1)
	}

	sm, err := statemachine.Parse(def)
	if err != nil {
		slog.Error("error when parsing the state machine definition:", "err", err)
		os.Exit(1)
	}

	err = <-sm.Run(context.Background())

	if err != nil {
		slog.Error("error when running the state machine", "err", err)
		os.Exit(1)
	}
}
