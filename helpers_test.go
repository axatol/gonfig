package gonfig

import (
	"flag"
	"testing"
)

func setCommandLineFlags(t *testing.T) func() {
	t.Helper()
	var oldFlags = flag.CommandLine
	flag.CommandLine = flag.NewFlagSet("test", flag.ContinueOnError)
	return func() { flag.CommandLine = oldFlags }
}
