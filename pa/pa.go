package pa

import (
	"context"
	"fmt"

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
	stream, err := client.NewPlayback(pulse.Float32Reader(oscillator.Synth32), pulse.PlaybackSink(sink), pulse.PlaybackSampleRate(sink.SampleRate()))
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
	o.stream.Stop()
	o.stream.Drain()
}
