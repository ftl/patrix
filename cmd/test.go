package cmd

import (
	"context"
	"log"
	"math"

	"github.com/jfreymuth/pulse"
	"github.com/spf13/cobra"
)

var testFlags = struct {
	frequency int
}{}

var testCmd = &cobra.Command{
	Use:   "test",
	Short: "Test PATRiX",
	Run:   runWithOscillator(runTest),
}

func init() {
	rootCmd.AddCommand(testCmd)

	testCmd.Flags().IntVar(&testFlags.frequency, "f", 700, "frequency in Hz")
}

func runTest(ctx context.Context, cmd *cobra.Command, args []string /*, oscillator */) {
	const sampleRate = 48000 // 44100
	const frequency = 1500

	pulseClient, err := pulse.NewClient()
	if err != nil {
		log.Fatalf("Cannot create client: %v", err)
	}
	defer pulseClient.Close()

	sink, err := pulseClient.DefaultSink()
	if err != nil {
		log.Fatalf("Cannot get default sink: %v", err)
	}

	stream, err := pulseClient.NewPlayback(synth(sampleRate, frequency), pulse.PlaybackSink(sink), pulse.PlaybackSampleRate(sampleRate), pulse.PlaybackLowLatency(sink))
	if err != nil {
		log.Fatalf("Cannot create playback: %v", err)
	}
	defer stream.Close()

	stream.Start()

	<-ctx.Done()
}

func synth(sampleRate int, frequency int) func([]float32) {
	var t float64
	const amp float64 = 1.0

	return func(out []float32) {
		for i := range out {
			out[i] = float32(amp * math.Cos(2*math.Pi*float64(frequency)*t))
			t += 1. / float64(sampleRate)
		}
	}
}
