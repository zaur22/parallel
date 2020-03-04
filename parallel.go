package parallel

import (
	"context"
	"fmt"
	"log"
	"sync"
)

type ShutdownTriggerType int

type SpawnFn = func(logPrefix string, t ShutdownTriggerType, f func())

const (
	_ ShutdownTriggerType = iota + 1
	Fail
	Exit
	Continue
)

type Spawn struct {
	ctx          context.Context
	cancel       func()
	functionDone chan ShutdownTriggerType
	wg           sync.WaitGroup
	logger       *log.Logger
}

func Run(
	ctx context.Context,
	f func(ctx context.Context, spawn SpawnFn) error,
) error {

	ctx, cancelFunc := context.WithCancel(ctx)
	var s = Spawn{
		ctx:          ctx,
		cancel:       cancelFunc,
		functionDone: make(chan ShutdownTriggerType, 1),
		wg:           sync.WaitGroup{},
	}

	if l := s.ctx.Value("logger"); l != nil {
		var ok bool
		s.logger, ok = l.(*log.Logger)
		if !ok {
			return fmt.Errorf("wrong type of logger")
		}
	} else {
		return fmt.Errorf("not found logger in context")
	}

	var err = f(ctx, s.SpawnFn)
	if err != nil {
		return fmt.Errorf("run exec error: %v", err)
	}

	return s.done()
}

func (s *Spawn) SpawnFn(logPrefix string, t ShutdownTriggerType, f func()) {

	//var ctx = context.WithValue(s.ctx, "logger", newLogger)
	var endOfFunction = make(chan struct{}, 0)

	logger := log.New(s.logger.Writer(), logPrefix, s.logger.Flags())

	go func() {
		logger.Print("starting")
		f()
		logger.Print("completed")
		endOfFunction <- struct{}{}
	}()

	go func() {
		s.wg.Add(1)
		defer s.wg.Done()
		select {
		case <-endOfFunction:
			s.functionDone <- t
		case <-s.ctx.Done():
			logger.Print("stopped")
		}

	}()
}

func (s *Spawn) done() error {
	for {
		var res = <-s.functionDone
		switch res {
		case Fail:
			s.cancel()
			s.wg.Wait()
			return fmt.Errorf("some error")
		case Exit:
			s.cancel()
			s.wg.Wait()
			return nil
		case Continue:
			continue
		default:
			s.logger.Printf("undefined trigger type: %v", res)
		}
	}
}
