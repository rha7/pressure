// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
pressure is a command line tool for load testing
*/
package main

import (
	"fmt"
	"math/rand"
	"os"
	"runtime"
	"time"

	"github.com/rha7/pressure/arguments"
	"github.com/rha7/pressure/attacker"
	"github.com/rha7/pressure/printers"
	"github.com/sirupsen/logrus"
)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	rand.Seed(time.Now().UnixNano())
}

func main() {
	logger := logrus.New()
	logger.Out = os.Stderr
	logger.SetLevel(logrus.InfoLevel)
	logger.Info("PRESSURE :: Load Testing Tool for APIs")
	logger.Info("Copyright(c) 2017 Gabriel Medina - MIT License")
	logger.Println()
	logger.Info("processing arguments...")
	spec, logLevel, err := arguments.Process(os.Args[1:], os.Stdin, logger, os.Stderr)
	if err != nil {
		logger.
			WithField("stage", "arguments_processing").
			WithField("message", err.Error()).
			Errorf("error occurred")
		os.Exit(1)
		return
	}
	logger.Println("arguments processed.")
	logger.Println()
	spec.Print(logger)
	logger.SetLevel(logLevel)
	logger.
		WithField("log_level", logLevel).
		WithField("spec", spec).
		Debug("arguments parsed onto spec")
	logger.Info("performing pressure test")
	results, err := attacker.Perform(logger, spec)
	if err != nil {
		logger.
			WithField("stage", "performing_test").
			WithField("message", err.Error()).
			Errorf("error occurred")
		os.Exit(2)
		return
	}
	logger.Info("done.")
	logger.Printf("\n\n\n")
	output, err := printers.Default(logger, results)
	if err != nil {
		logger.
			WithField("stage", "output").
			WithField("message", err.Error()).
			Errorf("error occurred")
		os.Exit(3)
		return
	}
	fmt.Println(output)
}
