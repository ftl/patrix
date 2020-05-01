package cmd

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/ftl/digimodes/psk31"
	"github.com/spf13/cobra"

	"github.com/ftl/patrix/pa"
)

var psk31Flags = struct {
	frequency int
}{}

var psk31Cmd = &cobra.Command{
	Use:   "psk31 [text]",
	Short: "Send text using PSK31",
	Run:   runWithOscillator(runPSK31),
}

func init() {
	rootCmd.AddCommand(psk31Cmd)

	psk31Cmd.Flags().IntVar(&psk31Flags.frequency, "f", 1500, "carrier frequency in Hz")
}

func runPSK31(ctx context.Context, cmd *cobra.Command, args []string, oscillator *pa.Oscillator) {
	if len(args) < 1 {
		log.Fatal("Need text as parameter. See psk31 --help for more information.")
	}
	text := strings.Join(args, " ")

	modulator := psk31.NewModulator(float64(psk31Flags.frequency))
	defer modulator.Close()
	modulator.AbortWhenDone(ctx.Done())
	oscillator.Modulator = modulator
	oscillator.Start()

	_, err := fmt.Fprintln(modulator, text)
	if err != nil {
		log.Fatal(err)
	}
	modulator.End()

	oscillator.Stop(ctx)
}
