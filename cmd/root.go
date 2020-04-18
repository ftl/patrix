package cmd

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"
)

var rootFlags = struct {
}{}

var rootCmd = &cobra.Command{
	Use:   "patrix",
	Short: "PATRiX - transmit and receive digital modes through the Pulse Audio framework",
}

// Execute is called by main.main() as the entry point to the Cobra framework.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

func init() {
}

func runWithOscillator(f func(ctx context.Context, cmd *cobra.Command, args []string /*, oscillator */)) func(cmd *cobra.Command, args []string) {
	return func(cmd *cobra.Command, args []string) {
		// TODO instanciate the oscillator

		ctx, cancel := context.WithCancel(context.Background())
		signals := make(chan os.Signal, 1)
		signal.Notify(signals, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
		go handleCancelation(signals, cancel, func() error { return nil } /* oscillator.Shutdown */)

		f(ctx, cmd, args /* oscillator */)

		// shut the oscillator down
	}
}

func handleCancelation(signals <-chan os.Signal, cancel context.CancelFunc, shutdown func() error) {
	count := 0
	for {
		select {
		case <-signals:
			count++
			if count == 1 {
				cancel()
			} else {
				shutdown()
				log.Fatal("hard shutdown")
			}
		}
	}
}
