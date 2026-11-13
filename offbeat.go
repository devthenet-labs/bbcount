package main

import (
	"fmt"
	"io"
)

// offbeatCounts returns, per track in trackNames order, the number of active
// steps that are not on a beat. Beats are steps 1, 5, 9 and 13 (indexes 0, 4,
// 8 and 12).
func (p pattern) offbeatCounts() []int {
	out := make([]int, len(trackNames))
	for i, name := range trackNames {
		for j, on := range p.Tracks[name] {
			if on && j%4 != 0 {
				out[i]++
			}
		}
	}
	return out
}

// printCountsOffbeat prints the counts like printCounts, plus each track's
// off-beat count.
func printCountsOffbeat(w io.Writer, p pattern, asJSON bool) {
	off := p.offbeatCounts()
	if !asJSON {
		printCounts(w, p, false)
		for i, name := range trackNames {
			fmt.Fprintf(w, "offbeat %s %d\n", name, off[i])
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
	out += fmt.Sprintf("\"total\":%d,\"offbeat\":{", total)
	for i, name := range trackNames {
		if i > 0 {
			out += ","
		}
		out += fmt.Sprintf("%q:%d", name, off[i])
	}
	fmt.Fprintln(w, out+"}}")
}
