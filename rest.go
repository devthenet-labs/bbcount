package main

import (
	"fmt"
	"io"
)

// longestRest returns the length of the longest run of inactive steps in a
// track. It does not wrap around from the last step to the first.
func longestRest(track []bool) int {
	best, cur := 0, 0
	for _, on := range track {
		if on {
			cur = 0
			continue
		}
		cur++
		if cur > best {
			best = cur
		}
	}
	return best
}

// printCountsRest prints what printCounts prints, plus each track's longest rest.
func printCountsRest(w io.Writer, p pattern, asJSON bool) {
	counts := p.counts()
	total := 0
	for _, c := range counts {
		total += c
	}
	rests := make([]int, len(trackNames))
	for i, name := range trackNames {
		rests[i] = longestRest(p.Tracks[name])
	}
	if asJSON {
		out := "{"
		for i, name := range trackNames {
			out += fmt.Sprintf("%q:%d,", name, counts[i])
		}
		out += fmt.Sprintf("\"total\":%d,\"rest\":{", total)
		for i, name := range trackNames {
			if i > 0 {
				out += ","
			}
			out += fmt.Sprintf("%q:%d", name, rests[i])
		}
		fmt.Fprintf(w, "%s}}\n", out)
		return
	}
	for i, name := range trackNames {
		fmt.Fprintf(w, "%s %d\n", name, counts[i])
	}
	fmt.Fprintf(w, "total %d\n", total)
	for i, name := range trackNames {
		fmt.Fprintf(w, "rest %s %d\n", name, rests[i])
	}
}
