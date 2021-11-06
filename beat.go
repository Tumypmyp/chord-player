package main

import (
	"github.com/faiface/beep"
	"github.com/faiface/beep/generators"
	"time"
)

func Chord(sr beep.SampleRate, freqs ...int) beep.Streamer {
	var mix beep.Mixer
	for _, f := range freqs {
		s, err := generators.SinTone(sr, f)
		if err != nil {
			return nil
		}
		mix.Add(s)
	}
	return &mix
}

type Beat struct {
	streamer beep.Streamer
	length   float64
}

func (b *Beat) Streamer(format beep.Format, cycleLength int) (beatStreamer beep.StreamSeeker) {
	buffer := beep.NewBuffer(format)
	buffer.Append(beep.Take(int(float64(cycleLength)*b.length), b.streamer))
	return buffer.Streamer(0, buffer.Len())
}

type Rhythm struct {
	format      beep.Format
	played      int
	cycleLength int
	Beats       map[int]beep.StreamSeeker
	streamer    beep.Streamer
}

func NewRhythm(f beep.Format, cycleLength time.Duration) (rhythm *Rhythm) {
	rhythm = &Rhythm{
		format:      f,
		played:      0,
		cycleLength: f.SampleRate.N(cycleLength),
		Beats:       make(map[int]beep.StreamSeeker),
		streamer:    beep.Mix(),
	}
	return rhythm
}

func (r *Rhythm) AddBeat(b Beat, t float64) {
	var ind int = (int)(float64(r.cycleLength) * t)
	streamer := b.Streamer(r.format, r.cycleLength)
	r.Beats[ind] = streamer
	r.streamer = beep.Mix(r.streamer, streamer)
}

func (r *Rhythm) Stream(samples [][2]float64) (n int, ok bool) {
	for i := 0; i < len(samples); i++ {
		if b, ok := r.Beats[r.played%r.cycleLength]; ok {
			r.streamer.Stream(samples[n:i])
			n = i
			b.Seek(0)
		}
		r.played++
	}
	r.streamer.Stream(samples[n:])
	return len(samples), true
}

func (r *Rhythm) Err() error {
	return r.streamer.Err()
}
