/*
Copyright github.com/czminami. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package shutdown

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"sync/atomic"
	"syscall"
	"time"
)

var (
	once     sync.Once
	instance *watchdog
)

type watchdog struct {
	jobs   int64
	ctx    context.Context
	cancel context.CancelFunc
	logger Logger
}

type Logger interface {
	Warning(args ...interface{})
}

func Init(ctx context.Context, cancel context.CancelFunc, logger Logger) error {
	if ctx == nil {
		return errors.New("ctx required")

	} else if cancel == nil {
		return errors.New("cancel required")

	} else if logger == nil {
		return errors.New("logger required")
	}

	once.Do(func() {
		instance = &watchdog{
			ctx:    ctx,
			cancel: cancel,
			logger: logger,
		}
	})

	return nil
}

func initialized() error {
	if instance == nil {
		return errors.New("please call shudown.Init(ctx, cancel) firstly")

	} else {
		return nil
	}
}

func AddJob() error {
	if err := initialized(); err != nil {
		return err

	} else {
		atomic.AddInt64(&instance.jobs, 1)
		return nil
	}
}

func DoneJob() error {
	if err := initialized(); err != nil {
		return err

	} else {
		atomic.AddInt64(&instance.jobs, -1)
		return nil
	}
}

func StandBy() error {
	if err := initialized(); err != nil {
		return err
	}

	go func() {
		Signal := make(chan os.Signal, 1)
		signal.Notify(Signal, os.Interrupt, os.Kill, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM)

		select {
		case <-instance.ctx.Done():
			instance.logger.Warning("something cause context canceled")

		case s := <-Signal:
			instance.logger.Warning(fmt.Sprint("get signal <", s.String(), ">, notify to shutdown."))
			instance.cancel()
		}

		signal.Stop(Signal)
	}()

	// blocked
	select {
	case <-instance.ctx.Done():
	}

	// safe stop
	for {
		if remaining := atomic.LoadInt64(&instance.jobs); remaining <= 0 {
			break

		} else {
			time.Sleep(time.Millisecond * 500)
			instance.logger.Warning(remaining, "jobs waiting to stop")
		}
	}

	instance.logger.Warning(">>>>  graceful shutdown  <<<<")
	return nil
}
