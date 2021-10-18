package main

import (
	"github.com/faiface/beep"
	"time"
)

// make beatStreamer, and Beat returning it from method Streamer
type Rhythm struct {
	BeatStreamer Beat
	played       int
	CycleLength  int
	Beats        map[int]struct{}
}

type Beat struct {
	Format   beep.Format
	Streamer beep.StreamSeeker
}

func NewBeat(format beep.Format, streamer beep.Streamer, beatLength time.Duration) (beat Beat) {
	beat.Format = format
	buffer := beep.NewBuffer(format)
	buffer.Append(beep.Take(format.SampleRate.N(beatLength), streamer))
	beat.Streamer = buffer.Streamer(0, buffer.Len())
	beat.Streamer.Seek(beat.Streamer.Len())
	return beat
}

func NewRhythm(beat Beat, cycleLength time.Duration, beats ...string) (rhythm *Rhythm) {
	beatsStart := make(map[int]struct{})
	for _, b := range beats {
		d, _ := time.ParseDuration(b)
		beatsStart[beat.Format.SampleRate.N(d)] = struct{}{}
	}

	rhythm = &Rhythm{
		BeatStreamer: beat,
		played:       0,
		CycleLength:  beat.Format.SampleRate.N(cycleLength),
		Beats:        beatsStart,
	}
	return rhythm
}
func (b *Rhythm) Stream(samples [][2]float64) (n int, ok bool) {
	for i := 0; i < len(samples); {
		if _, ok := b.Beats[b.played%b.CycleLength]; ok {
			b.BeatStreamer.Streamer.Seek(0)
		}
		sn, _ := b.BeatStreamer.Streamer.Stream(samples[i:])
		i += sn
		b.played += sn
		if sn == 0 {
			samples[i][0] = 0
			samples[i][1] = 0
			i++
			b.played++
		}
	}
	return len(samples), true
}

func (b *Rhythm) Err() error {
	return b.BeatStreamer.Streamer.Err()
}
