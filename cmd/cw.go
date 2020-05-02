package cmd

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/ftl/digimodes/cw"
	"github.com/spf13/cobra"

	"github.com/ftl/patrix/pa"
)

var cwFlags = struct {
	wpm    int
	pitch  int
	beacon int
}{}

var cwCmd = &cobra.Command{
	Use:   "cw [text]",
	Short: "Send text using CW",
	Run:   runWithOscillator(runCW),
}

func init() {
	rootCmd.AddCommand(cwCmd)

	cwCmd.Flags().IntVar(&cwFlags.wpm, "wpm", 12, "speed in WpM")
	cwCmd.Flags().IntVar(&cwFlags.pitch, "pitch", 700, "pitch in Hz")
	cwCmd.Flags().IntVar(&cwFlags.beacon, "beacon", 0, "number of seconds between transmissions (0 == single transmission)")
}

func runCW(ctx context.Context, cmd *cobra.Command, args []string, oscillator *pa.Oscillator) {
	if len(args) < 1 {
		log.Fatal("Need text as parameter. See cw --help for more information")
	}
	text := strings.Join(args, " ")

	modulator := cw.NewModulator(float64(cwFlags.pitch), cwFlags.wpm)
	defer modulator.Close()
	modulator.AbortWhenDone(ctx.Done())
	oscillator.Modulator = modulator
	oscillator.Start()
	defer oscillator.Stop(ctx)

	for {
		_, err := fmt.Fprintln(modulator, text)
		if err != nil {
			log.Fatal(err)
		}

		if cwFlags.beacon == 0 {
			return
		}
		select {
		case <-time.After(time.Duration(cwFlags.beacon) * time.Second):
		case <-ctx.Done():
			return
		}
	}
}
