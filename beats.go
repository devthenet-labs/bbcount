package main

import (
	"fmt"
	"io"
)

// beats returns, per track in trackNames order, how many steps play in each of
// the four beats. Step i (0-based) counts in beat i/4.
func (p pattern) beats() [][]int {
	out := make([][]int, len(trackNames))
	for i, name := range trackNames {
		out[i] = make([]int, steps/4)
		for s, on := range p.Tracks[name] {
			if on {
				out[i][s/4]++
			}
		}
	}
	return out
}

// printReport prints the counts, and with withBeats the per-beat counts too.
func printReport(w io.Writer, p pattern, asJSON, withBeats bool) {
	if !withBeats {
		printCounts(w, p, asJSON)
		return
	}
	beats := p.beats()
	if !asJSON {
		printCounts(w, p, false)
		for i, name := range trackNames {
			b := beats[i]
			fmt.Fprintf(w, "beats %s %d %d %d %d\n", name, b[0], b[1], b[2], b[3])
		}
		return
	}
	counts := p.counts()
	total := 0
	out := "{"
	for i, name := range trackNames {
		total += counts[i]
		out += fmt.Sprintf("%q:%d,", name, counts[i])
	}
	out += fmt.Sprintf("\"total\":%d,\"beats\":{", total)
	for i, name := range trackNames {
		if i > 0 {
			out += ","
		}
		b := beats[i]
		out += fmt.Sprintf("%q:[%d,%d,%d,%d]", name, b[0], b[1], b[2], b[3])
	}
	fmt.Fprintf(w, "%s}}\n", out)
}
