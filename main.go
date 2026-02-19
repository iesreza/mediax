package main

import (
	"fmt"
	"github.com/getevo/docify"
	"github.com/getevo/evo/v2"
	"github.com/getevo/evo/v2/lib/application"
	"github.com/getevo/restify"
	"mediax/apps/mediax"
	"time"
)

func main() {
	// evo.Setup() may panic if the DB isn't ready yet (race condition in settings load).
	// Retry with backoff so the container doesn't crash-loop waiting for Railway restart.
	for attempt := 1; ; attempt++ {
		if trySetup() {
			break
		}
		if attempt >= 10 {
			panic("evo.Setup() failed after 10 attempts, giving up")
		}
		wait := time.Duration(attempt) * time.Second
		fmt.Printf("evo.Setup() panicked (attempt %d), retrying in %v...\n", attempt, wait)
		time.Sleep(wait)
	}

	var apps = application.GetInstance()
	// Register all application modules
	apps.Register( // Authentication follows
		mediax.App{},
		restify.App{},
		docify.App{},
	)

	// Start the application
	evo.Run()
}

func trySetup() (ok bool) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("evo.Setup() panic recovered: %v\n", r)
			ok = false
		}
	}()
	evo.Setup()
	return true
}
