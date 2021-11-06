package main

import (
	"fmt"
	"github.com/faiface/beep"
	"github.com/faiface/beep/effects"
	"github.com/faiface/beep/generators"
	"github.com/faiface/beep/speaker"
)

func Chord(sr beep.SampleRate, freqs ...int) beep.Streamer {
	var mix beep.Mixer
	for _, f := range freqs {
		s, err := generators.SinTone(sr, f)
		if err != nil {
			fmt.Println(err, ", can not add frequency", f)
		}
		mix.Add(s)
	}
	return &mix
}

func main() {
	format := beep.Format{
		SampleRate:  beep.SampleRate(48000),
		NumChannels: 2,
		Precision:   6,
	}
	sr := format.SampleRate

	speaker.Init(sr, 4096)

	chord := Chord(sr, 1200, 1000, 800)

	beat := NewBeat(format, chord, 5000)

	rhythmStreamer := NewRhythm(72000)

	rhythmStreamer.AddBeat(beat, 0)
	rhythmStreamer.AddBeat(beat, 0.2)

	v := &effects.Volume{
		Streamer: rhythmStreamer,
		Base:     2,
		Volume:   -7,
		Silent:   false,
	}
	speaker.Play(v)

	for {
		fmt.Println("Print n - number of notes, float time in a cycle beetween 0 and 1,\n\tthen n frequencies.")
		var n int
		var t float32
		fmt.Scan(&n, &t)
		notes := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Scan(&notes[i])
		}

		beat := NewBeat(format, Chord(sr, notes...), 5000)
		rhythmStreamer.AddBeat(beat, t)
		fmt.Println(rhythmStreamer)
	}
}
