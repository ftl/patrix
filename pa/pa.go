package pa

import (
	"context"
	"fmt"
	"time"

	"github.com/jfreymuth/pulse"

	"github.com/ftl/patrix/osc"
)

type Oscillator struct {
	*osc.Oscillator

	client *pulse.Client
	stream *pulse.PlaybackStream
}

func NewOscillator() (*Oscillator, error) {
	client, err := pulse.NewClient()
	if err != nil {
		return nil, fmt.Errorf("cannot open pulse audio client: %v", err)
	}
	sink, err := client.DefaultSink()
	if err != nil {
		return nil, fmt.Errorf("cannot get pulse audio default sink: %v", err)
	}
	oscillator := osc.New(sink.SampleRate())
	stream, err := client.NewPlayback(oscillator.Synth32, pulse.PlaybackSink(sink), pulse.PlaybackSampleRate(sink.SampleRate()))
	if err != nil {
		return nil, fmt.Errorf("cannot create pulse audio playback stream: %v", err)
	}

	result := &Oscillator{
		Oscillator: oscillator,
		client:     client,
		stream:     stream,
	}
	return result, nil
}

func (o *Oscillator) Close() error {
	o.stream.Close()
	o.client.Close()
	return nil
}

func (o *Oscillator) Start() {
	o.stream.Start()
}

func (o *Oscillator) Stop(ctx context.Context) {
	o.Oscillator.Modulator = osc.NoModulator
	select {
	// wait until all the remaining samples are processed by pulse audio
	case <-time.After(time.Duration(float64(o.stream.BufferSize())*o.Oscillator.Tick()) * time.Second):
	case <-ctx.Done():
	}
	o.stream.Stop()
}
