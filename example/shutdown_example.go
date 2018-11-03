/*
Copyright github.com/czminami. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package main

import (
	"context"
	"os"
	"runtime/debug"
	"syscall"
	"time"

	"github.com/czminami/shutdown.gracefully"
	"github.com/op/go-logging"
)

var logger = logging.MustGetLogger("test")

func main() {
	format := logging.MustStringFormatter(`[%{module}] %{time:2006-01-02 15:04:05} [%{level}] [%{longpkg} %{shortfile}] { %{message} }`)

	backendConsole := logging.NewLogBackend(os.Stderr, "", 0)
	backendConsole2Formatter := logging.NewBackendFormatter(backendConsole, format)

	logging.SetBackend(backendConsole2Formatter)
	logging.SetLevel(logging.INFO, "")

	defer func() {
		if err := recover(); err != nil {
			logger.Error(err)
			logger.Info(string(debug.Stack()))
		}
	}()

	ctx, cancel := context.WithCancel(context.Background())

	shutdown.Init(ctx, cancel, logger) // init

	for k := 0; k < 10; k++ {
		go func(id int) {
			shutdown.AddJob()
			defer shutdown.DoneJob()

			time.Sleep(time.Second * time.Duration(id)) // dummy, something cost some time

			logger.Info(id, "over")
		}(k)
	}

	go func() {
		time.Sleep(time.Second * 2)
		syscall.Kill(syscall.Getpid(), syscall.SIGINT) // dummy signal
	}()

	shutdown.StandBy()
}
