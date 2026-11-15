package main

import (
	"fmt"
	"io"
)

// lasts returns the 1-based number of the last active step per track, in
// trackNames order, or 0 for a track that plays no step.
func (p pattern) lasts() []int {
	out := make([]int, len(trackNames))
	for i, name := range trackNames {
		for j, on := range p.Tracks[name] {
			if on {
				out[i] = j + 1
			}
		}
	}
	return out
}

// printCountsLast prints the counts like printCounts, plus the last step of
// each track.
func printCountsLast(w io.Writer, p pattern, asJSON bool) {
	lasts := p.lasts()
	if !asJSON {
		printCounts(w, p, false)
		for i, name := range trackNames {
			fmt.Fprintf(w, "last %s %d\n", name, lasts[i])
		}
		return
	}
	counts := p.counts()
	total := 0
	for _, c := range counts {
		total += c
	}
	out := "{"
	for i, name := range trackNames {
		out += fmt.Sprintf("%q:%d,", name, counts[i])
	}
	out += fmt.Sprintf("\"total\":%d,\"last\":{", total)
	for i, name := range trackNames {
		if i > 0 {
			out += ","
		}
		out += fmt.Sprintf("%q:%d", name, lasts[i])
	}
	fmt.Fprintf(w, "%s}}\n", out)
}
