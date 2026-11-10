package main

import (
	"encoding/json"
	"fmt"
	"io"
)

// steps is the length of every track: sixteen sixteenth notes.
const steps = 16

// trackNames is the fixed track order, which is also the order counts print in.
var trackNames = []string{"kick", "snare", "hat"}

// pattern is Beat Battle's shared pattern JSON, reduced to what bbcount reads.
// Unknown fields are rejected, except the ones the other services add.
type pattern struct {
	Name   string            `json:"name"`
	Author string            `json:"author"`
	BPM    int               `json:"bpm"`
	Swing  int               `json:"swing"`
	Tracks map[string][]bool `json:"tracks"`

	ID        json.RawMessage `json:"id"`
	Votes     json.RawMessage `json:"votes"`
	CreatedAt json.RawMessage `json:"createdAt"`
}

// parsePattern reads one pattern and validates it.
func parsePattern(r io.Reader) (pattern, error) {
	var p pattern
	dec := json.NewDecoder(io.LimitReader(r, 16<<10))
	dec.DisallowUnknownFields()
	if err := dec.Decode(&p); err != nil {
		return pattern{}, fmt.Errorf("invalid pattern: %w", err)
	}
	if p.BPM < 60 || p.BPM > 200 {
		return pattern{}, fmt.Errorf("invalid pattern: bpm %d is outside 60-200", p.BPM)
	}
	if p.Swing < 0 || p.Swing > 60 {
		return pattern{}, fmt.Errorf("invalid pattern: swing %d is outside 0-60", p.Swing)
	}
	if len(p.Tracks) != len(trackNames) {
		return pattern{}, fmt.Errorf("invalid pattern: tracks must be exactly %v", trackNames)
	}
	for _, name := range trackNames {
		track, ok := p.Tracks[name]
		if !ok {
			return pattern{}, fmt.Errorf("invalid pattern: track %q is missing", name)
		}
		if len(track) != steps {
			return pattern{}, fmt.Errorf("invalid pattern: track %q has %d steps, want %d", name, len(track), steps)
		}
	}
	return p, nil
}

// counts returns the number of active steps per track, in trackNames order.
func (p pattern) counts() []int {
	out := make([]int, len(trackNames))
	for i, name := range trackNames {
		for _, on := range p.Tracks[name] {
			if on {
				out[i]++
			}
		}
	}
	return out
}

// busiest names the track, or every tied track in trackNames order
// joined by commas, with the most steps, or "none" if no track plays
// a step.
func (p pattern) busiest() string {
	counts := make([]int, len(trackNames))
	for i, name := range trackNames {
		track := p.Tracks[name]
		for step := 0; step < steps-1; step++ {
			if track[step] {
				counts[i]++
			}
		}
	}
	max := 0
	for _, c := range counts {
		if c > max {
			max = c
		}
	}
	if max == 0 {
		return "none"
	}
	out := ""
	for i, name := range trackNames {
		if counts[i] == max {
			if out != "" {
				out += ","
			}
			out += name
		}
	}
	return out
}
