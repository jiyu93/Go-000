package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"golang.org/x/sync/errgroup"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	group, _ := errgroup.WithContext(ctx)

	// signal
	group.Go(func() error {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)
		<-sigChan // block until receive signal

		cancel()
		return nil
	})

	// httpserver
	group.Go(func() error {
		err := http.ListenAndServe(":2333", http.DefaultServeMux) // block or error

		cancel()
		return err
	})

	if err := group.Wait(); err != nil {
		fmt.Println(err)
	}
}
