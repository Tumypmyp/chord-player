package main

import (
	"fmt"
	"github.com/faiface/beep"
	"github.com/faiface/beep/effects"
	"github.com/faiface/beep/speaker"
	"time"
)

func main() {
	format := beep.Format{
		SampleRate:  beep.SampleRate(48000),
		NumChannels: 2,
		Precision:   6,
	}
	sr := format.SampleRate

	speaker.Init(sr, 4096)

	chord := Chord(sr, 1200, 1000, 800)

	//beat := NewBeat(format, chord, 5000)

	rhythmStreamer := NewRhythm(format, time.Second*1)

	rhythmStreamer.AddBeat(Beat{
		streamer: chord,
		length:   0.125,
	}, 0)
	rhythmStreamer.AddBeat(Beat{
		streamer: chord,
		length:   0.250,
	}, 0.5)

	v := &effects.Volume{
		Streamer: rhythmStreamer,
		Base:     2,
		Volume:   -8,
		Silent:   false,
	}
	speaker.Play(v)

	for {
		fmt.Println("Print n - number of notes, float time in a cycle beetween 0 and 1,\n\tthen n frequencies.")
		var n int
		var t float64
		fmt.Scan(&n, &t)
		notes := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Scan(&notes[i])
		}
		chord := Chord(sr, notes...)
		beat := Beat{
			streamer: chord,
			length:   0.125,
		}

		rhythmStreamer.AddBeat(beat, t)
		fmt.Println(rhythmStreamer)
	}
}
