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
	defer cancel()
	group, _ := errgroup.WithContext(ctx)

	// signal
	group.Go(func() error {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)
		for {
			select {
			case <-ctx.Done():
				return nil
			case s := <-sigChan:
				return fmt.Errorf("syscall signal:%v", s)
			}
		}
	})

	// httpserver
	group.Go(func() error {
		s := http.Server{}
		go func() {
			select {
			case <-ctx.Done():
				s.Shutdown(ctx)
			}
		}()
		return s.ListenAndServe()
	})

	if err := group.Wait(); err != nil {
		fmt.Println(err)
	}
}
