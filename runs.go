package main

import (
	"fmt"
	"io"
	"strings"
)

// runs returns the number of runs (groups of adjacent active steps) per
// track, in trackNames order. A track does not wrap around.
func (p pattern) runs() []int {
	out := make([]int, len(trackNames))
	for i, name := range trackNames {
		track := p.Tracks[name]
		for j, on := range track {
			if on && (j == 0 || !track[j-1]) {
				out[i]++
			}
		}
	}
	return out
}

// printCountsWithRuns prints the counts as printCounts does, plus the
// per-track run counts.
func printCountsWithRuns(w io.Writer, p pattern, asJSON bool) {
	runs := p.runs()
	if !asJSON {
		printCounts(w, p, false)
		for i, name := range trackNames {
			fmt.Fprintf(w, "runs %s %d\n", name, runs[i])
		}
		return
	}
	counts := p.counts()
	total := 0
	var parts, runParts []string
	for i, name := range trackNames {
		total += counts[i]
		parts = append(parts, fmt.Sprintf("%q:%d", name, counts[i]))
		runParts = append(runParts, fmt.Sprintf("%q:%d", name, runs[i]))
	}
	fmt.Fprintf(w, "{%s,\"total\":%d,\"runs\":{%s}}\n",
		strings.Join(parts, ","), total, strings.Join(runParts, ","))
}
