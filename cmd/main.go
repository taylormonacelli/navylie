package main

import (
	"flag"
	"log/slog"
	"os"

	"github.com/taylormonacelli/goldbug"
	"github.com/taylormonacelli/navylie"
)

var (
	verbose bool
	noCache bool

	logFormat string

	dir string
)

func main() {
	flag.BoolVar(&verbose, "verbose", false, "Enable verbose output")
	flag.BoolVar(&verbose, "v", false, "Enable verbose output (shorthand)")

	flag.StringVar(&logFormat, "log-format", "", "Log format (text or json)")

	flag.BoolVar(&noCache, "no-cache", false, "Force downloading project template even if its not older threshold")

	flag.StringVar(&dir, "dir", "", "The directory to start the search from")

	flag.Parse()

	if dir == "" {
		slog.Error("Error: Parameter is empty. Please provide a value.", "dir", dir)
		flag.Usage()
		os.Exit(1)
	}

	if verbose || logFormat != "" {
		if logFormat == "json" {
			goldbug.SetDefaultLoggerJson(slog.LevelDebug)
		} else {
			goldbug.SetDefaultLoggerText(slog.LevelDebug)
		}
	}

	code := navylie.Main(dir)
	os.Exit(code)
}
