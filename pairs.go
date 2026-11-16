package main

import (
	"fmt"
	"io"
)

// pairNames is the fixed order in which track pairs print.
var pairNames = [][2]string{{"kick", "snare"}, {"kick", "hat"}, {"snare", "hat"}}

// pairCounts returns, per pair in pairNames order, the number of steps on
// which both tracks play.
func (p pattern) pairCounts() []int {
	out := make([]int, len(pairNames))
	for i, pair := range pairNames {
		a, b := p.Tracks[pair[0]], p.Tracks[pair[1]]
		for s := 0; s < steps; s++ {
			if a[s] && b[s] {
				out[i]++
			}
		}
	}
	return out
}

// printCountsWithPairs prints the per-track counts and total like printCounts,
// followed by the shared-step count of each track pair.
func printCountsWithPairs(w io.Writer, p pattern, asJSON bool) {
	counts := p.counts()
	total := 0
	for _, c := range counts {
		total += c
	}
	pairs := p.pairCounts()
	if asJSON {
		out := "{"
		for i, name := range trackNames {
			out += fmt.Sprintf("%q:%d,", name, counts[i])
		}
		out += fmt.Sprintf("\"total\":%d,\"pairs\":{", total)
		for i, pair := range pairNames {
			if i > 0 {
				out += ","
			}
			out += fmt.Sprintf("\"%s_%s\":%d", pair[0], pair[1], pairs[i])
		}
		fmt.Fprintln(w, out+"}}")
		return
	}
	for i, name := range trackNames {
		fmt.Fprintf(w, "%s %d\n", name, counts[i])
	}
	fmt.Fprintf(w, "total %d\n", total)
	for i, pair := range pairNames {
		fmt.Fprintf(w, "pair %s %s %d\n", pair[0], pair[1], pairs[i])
	}
}
