package main

import (
	"encoding/binary"
	"math"
)

const (
	sampleRate = 44100
	headerSize = 44
	hitMs      = 120 // how long one hit sounds
)

// stepStartMs is the start of a step in whole milliseconds, rounded to the
// nearest. A step is a sixteenth note; swing delays every odd step by that
// percent of a step.
func stepStartMs(bpm, swing, step int) int {
	stepMs := 60000.0 / float64(bpm) / 4
	start := float64(step) * stepMs
	if step%2 == 1 {
		start += stepMs * float64(swing) / 100
	}
	return int(math.Round(start))
}

// durationMs is the length of the pattern: sixteen steps.
func durationMs(bpm int) int {
	return int(math.Round(16 * 60000.0 / float64(bpm) / 4))
}

// render returns a complete mono 16-bit 44100 Hz WAV file. It is
// deterministic: the same pattern always yields the same bytes.
func render(p pattern) []byte {
	total := sampleRate * durationMs(p.BPM) / 1000
	pcm := make([]int32, total)
	for _, name := range trackNames {
		for step, on := range p.Tracks[name] {
			if on {
				mixHit(pcm, name, sampleRate*stepStartMs(p.BPM, p.Swing, step)/1000)
			}
		}
	}
	out := make([]byte, headerSize+2*total)
	copy(out[0:], "RIFF")
	binary.LittleEndian.PutUint32(out[4:], uint32(len(out)-8))
	copy(out[8:], "WAVEfmt ")
	binary.LittleEndian.PutUint32(out[16:], 16)
	binary.LittleEndian.PutUint16(out[20:], 1) // PCM
	binary.LittleEndian.PutUint16(out[22:], 1) // mono
	binary.LittleEndian.PutUint32(out[24:], sampleRate)
	binary.LittleEndian.PutUint32(out[28:], sampleRate*2)
	binary.LittleEndian.PutUint16(out[32:], 2)
	binary.LittleEndian.PutUint16(out[34:], 16)
	copy(out[36:], "data")
	binary.LittleEndian.PutUint32(out[40:], uint32(2*total))
	for i, s := range pcm {
		s = max(min(s, math.MaxInt16), math.MinInt16)
		binary.LittleEndian.PutUint16(out[headerSize+2*i:], uint16(int16(s)))
	}
	return out
}

// mixHit adds one decaying hit of the named drum at a sample offset. The kick
// is a low sine, the snare a deterministic noise burst and the hat a high
// square.
func mixHit(pcm []int32, name string, at int) {
	n := sampleRate * hitMs / 1000
	noise := uint32(12345)
	for i := 0; i < n && at+i < len(pcm); i++ {
		env := 1 - float64(i)/float64(n)
		t := float64(i) / sampleRate
		var v float64
		switch name {
		case "kick":
			v = math.Sin(2 * math.Pi * 60 * t)
		case "snare":
			noise = noise*1664525 + 1013904223
			v = float64(int32(noise>>16)-32768) / 32768
		default:
			v = math.Copysign(1, math.Sin(2*math.Pi*6000*t)) * 0.5
		}
		pcm[at+i] += int32(v * env * 8000)
	}
}
