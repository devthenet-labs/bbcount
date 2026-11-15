package main

import (
	"fmt"
	"io"
)

// halves returns, per track in trackNames order, the number of active steps in
// the first half of the bar (steps 1-8) and in the second half (steps 9-16).
func (p pattern) halves() [][2]int {
	out := make([][2]int, len(trackNames))
	for i, name := range trackNames {
		for s, on := range p.Tracks[name] {
			if on {
				out[i][s/(steps/2)]++
			}
		}
	}
	return out
}

// printHalf prints the usual counts followed by the per-track half split.
func printHalf(w io.Writer, p pattern, asJSON bool) {
	h := p.halves()
	if !asJSON {
		printCounts(w, p, false)
		for i, name := range trackNames {
			fmt.Fprintf(w, "half %s %d %d\n", name, h[i][0], h[i][1])
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
	out += fmt.Sprintf("\"total\":%d,\"half\":{", total)
	for i, name := range trackNames {
		if i > 0 {
			out += ","
		}
		out += fmt.Sprintf("%q:[%d,%d]", name, h[i][0], h[i][1])
	}
	fmt.Fprintf(w, "%s}}\n", out)
}
