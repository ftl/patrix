package osc

import (
	"math"
)

type Modulator interface {
	Modulate(t, a, f, p float64) (amplitude, frequency, phase float64)
}

type ModulatorFunc func(t, a, f, p float64) (amplitude, frequency, phase float64)

func (mf ModulatorFunc) Modulate(t, a, f, p float64) (amplitude, frequency, phase float64) {
	return mf(t, a, f, p)
}

var NoModulator = ModulatorFunc(func(t, a, f, p float64) (amplitude, frequency, phase float64) {
	return a, f, p
})

type Oscillator struct {
	Modulator Modulator

	amplitude float64
	frequency float64
	phase     float64
	tick      float64
	t         float64
	lastOut   float64
}

func New(sampleRate int) *Oscillator {
	return &Oscillator{
		Modulator: NoModulator,

		amplitude: 0.0,
		frequency: 0.0,
		phase:     0.0, // 0 = 0.0*math.Pi; 90 = 0.5*math.Pi; 180 = 1.0*math.Pi; 270 = 1.5*math.Pi;
		tick:      1.0 / float64(sampleRate),
	}
}

func (o *Oscillator) Tick() float64 {
	return o.tick
}

func (o *Oscillator) Synth32(out []float32) (int, error) {
	for i := range out {
		out[i] = float32(o.amplitude * math.Cos(2*math.Pi*o.frequency*o.t+o.phase))
		if o.Modulator != nil && (out[i] == 0 || (o.lastOut < 0 && out[i] > 0) || (o.lastOut > 0 && out[i] < 0)) {
			o.amplitude, o.frequency, o.phase = o.Modulator.Modulate(o.t, o.amplitude, o.frequency, o.phase)
		}

		o.lastOut = float64(out[i])
		o.t += o.tick
	}
	return len(out), nil
}

func (o *Oscillator) Synth64(out []float64) (int, error) {
	for i := range out {
		out[i] = o.amplitude * math.Cos(2*math.Pi*o.frequency*o.t+o.phase)
		if o.Modulator != nil && (out[i] == 0 || (o.lastOut < 0 && out[i] > 0) || (o.lastOut > 0 && out[i] < 0)) {
			o.amplitude, o.frequency, o.phase = o.Modulator.Modulate(o.t, o.amplitude, o.frequency, o.phase)
		}

		o.lastOut = out[i]
		o.t += o.tick
	}
	return len(out), nil
}
