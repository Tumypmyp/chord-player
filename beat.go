package main

import (
	"github.com/faiface/beep"
)

type Beat struct {
	Buffer *beep.Buffer
}

func NewBeat(format beep.Format, streamer beep.Streamer, beatLength int) (beat Beat) {
	beat = Beat{}
	beat.Buffer = beep.NewBuffer(format)
	beat.Buffer.Append(beep.Take(beatLength, streamer))
	return beat
}

func (b *Beat) Streamer() (beatStreamer beep.StreamSeeker) {
	beatStreamer = b.Buffer.Streamer(0, b.Buffer.Len())
	return beatStreamer
}


type Rhythm struct {
	played      int
	CycleLength int
	Beats       map[int]Beat
	streamer    beep.Mixer
}

func NewRhythm(cycleLength int) (rhythm *Rhythm) {
	rhythm = &Rhythm{
		played:      0,
		CycleLength: cycleLength,
		Beats:       make(map[int]Beat),
	}
	return rhythm
}

func (r *Rhythm) AddBeat(b Beat, t int) {
	r.Beats[t] = b
}

func (r *Rhythm) Stream(samples [][2]float64) (n int, ok bool) {
	for i := 0; i < len(samples); i++ {
		if b, ok := r.Beats[r.played%r.CycleLength]; ok {
			r.streamer.Stream(samples[:i])
			samples = samples[i:]

			r.streamer.Add(b.Streamer())
		}
		r.played++
	}
	r.streamer.Stream(samples)
	return len(samples), true
}

func (r *Rhythm) Err() error {
	return r.streamer.Err()
}
