package main

import (
	"fmt"
	"github.com/faiface/beep"
	"github.com/faiface/beep/effects"
	"github.com/faiface/beep/generators"
	"github.com/faiface/beep/speaker"
	"time"
)

func Chord(sr beep.SampleRate, freqs ...int) (beep.Streamer, error) {
	var mix beep.Mixer
	for _, f := range freqs {
		s, err := generators.SinTone(sr, f)
		if err != nil {
			return nil, err
		}
		mix.Add(s)
	}
	return &mix, nil
}

func main() {

	format := beep.Format{
		SampleRate:  beep.SampleRate(48000),
		NumChannels: 2,
		Precision:   6,
	}
	sr := format.SampleRate

	speaker.Init(sr, 4096)

	chord, err := Chord(sr, 1200, 1000, 800)
	if err != nil {
		panic(err)
	}

	rhythmStreamer := NewRhythm(NewBeat(format, chord, time.Millisecond*300), time.Second*2, "0s", "1500ms")

	v := &effects.Volume{
		Streamer: rhythmStreamer,
		Base:     2,
		Volume:   -6,
		Silent:   false,
	}
	speaker.Play(v)

	for {
		var n int
		fmt.Scan(&n)
		notes := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Scan(&notes[i])
		}
		fmt.Println(notes, " - notes")

		chord, err := Chord(sr, notes...)
		if err != nil {
			fmt.Println(err)
			continue
		}
		rhythmStreamer.BeatStreamer = NewBeat(format, chord, time.Millisecond*300)
	}
}
