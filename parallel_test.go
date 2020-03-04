package parallel

import (
	"context"
	"log"
	"os"
	"testing"
	"time"
)

func TestFail(t *testing.T) {

	var ctx = context.WithValue(
		context.Background(),
		"logger",
		log.New(os.Stdout, "parallel", log.LstdFlags),
	)

	err := Run(ctx, func(ctx context.Context, spawn SpawnFn) error {

		var f1 = func() {
			time.Sleep(50 * time.Millisecond)
		}

		var f2 = func() {
			time.Sleep(100 * time.Millisecond)
		}
		spawn("firs", Fail, f1)
		spawn("second", Exit, f2)
		return nil
	})

	if err == nil || err.Error() != "some error" {
		t.Errorf("expected value of error 'some error', got %v", err)
	}
}

func TestExit(t *testing.T) {

	var ctx = context.WithValue(
		context.Background(),
		"logger",
		log.New(os.Stdout, "parallel", log.LstdFlags),
	)

	err := Run(ctx, func(ctx context.Context, spawn SpawnFn) error {

		var f1 = func() {
			time.Sleep(50 * time.Millisecond)
		}

		var f2 = func() {
			time.Sleep(100 * time.Millisecond)
		}
		spawn("firs", Exit, f1)
		spawn("second", Fail, f2)
		return nil
	})

	if err != nil {
		t.Errorf("expected nil, got %v", err)
	}
}

func TestContinue(t *testing.T) {

	var ctx = context.WithValue(
		context.Background(),
		"logger",
		log.New(os.Stdout, "parallel", log.LstdFlags),
	)

	err := Run(ctx, func(ctx context.Context, spawn SpawnFn) error {

		var f1 = func() {
			time.Sleep(0 * time.Millisecond)
		}

		var f2 = func() {
			time.Sleep(10 * time.Millisecond)
		}

		var f3 = func() {
			time.Sleep(300 * time.Millisecond)
		}

		spawn("first", Continue, f1)
		spawn("second", Continue, f2)
		spawn("third", Continue, f3)
		return nil
	})

	if err != nil {
		t.Errorf("expected nil, got %v", err)
	}
}
