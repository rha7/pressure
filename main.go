package main

import (
	"flag"

	"github.com/rha7/pressure/attacker"
	"github.com/rha7/pressure/printers"
	"github.com/rha7/pressure/types"
	"github.com/sirupsen/logrus"
)

func parseFlags() types.Arguments {
	var arguments types.Arguments
	flag.Int64Var(&arguments.TotalRequests, "n", 100, "total requests")
	flag.Int64Var(&arguments.ConcurrentThreads, "c", 10, "concurrent requests")
	flag.Int64Var(&arguments.RequestTimeout, "t", 60, "request timeout")
	flag.Parse()
	return arguments
}

func main() {
	logrus.Info("pressure copyright(c) 2017 Gabriel Medina")
	arguments := parseFlags()
	results := attacker.Perform(arguments)
	printers.Text(results)
}
