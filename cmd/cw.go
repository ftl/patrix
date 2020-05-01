package cmd

import (
	"context"
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
		log.Fatal("Need text as parameter. Try cw --help")
	}

	text := strings.Join(args, " ")

	symbols := make(chan cw.Symbol)
	setKeyDown := func(keyDown bool) {
		// modify the oscillator
	}

	go cw.Send(ctx, setKeyDown, symbols, cwFlags.wpm)

	for {
		log.Print(text)
		cw.WriteToSymbolStream(ctx, symbols, text)

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
