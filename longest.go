package main

import (
	"encoding/json"
	"fmt"
	"io"
)

// longestObject is the -json form of the longest runs, in trackNames order.
type longestObject struct {
	Kick  int `json:"kick"`
	Snare int `json:"snare"`
	Hat   int `json:"hat"`
}

// longestRun returns the longest run of consecutive playing steps in a track.
// The bar does not wrap: the last step and the first are not adjacent.
func longestRun(track []bool) int {
	cur, best := 0, 0
	for _, on := range track {
		if on {
			cur++
			best = max(best, cur)
		} else {
			cur = 0
		}
	}
	return best
}

// printCountsLongest prints what printCounts prints, plus each track's
// longest run of consecutive playing steps.
func printCountsLongest(w io.Writer, p pattern, asJSON bool) {
	counts := p.counts()
	total := 0
	for _, c := range counts {
		total += c
	}
	longest := make([]int, len(trackNames))
	for i, name := range trackNames {
		longest[i] = longestRun(p.Tracks[name])
	}
	if asJSON {
		out := "{"
		for i, name := range trackNames {
			out += fmt.Sprintf("%q:%d,", name, counts[i])
		}
		obj, _ := json.Marshal(longestObject{Kick: longest[0], Snare: longest[2], Hat: longest[1]})
		fmt.Fprintf(w, "%s\"total\":%d,\"longest\":%s}\n", out, total, obj)
		return
	}
	for i, name := range trackNames {
		fmt.Fprintf(w, "%s %d\n", name, counts[i])
	}
	fmt.Fprintf(w, "total %d\n", total)
	for i, name := range trackNames {
		fmt.Fprintf(w, "longest %s %d\n", name, longest[i])
	}
}
