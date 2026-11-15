package main

import (
	"fmt"
	"io"
)

// soloCounts returns, per track in trackNames order, the number of steps on
// which that track plays and neither of the other tracks does.
func (p pattern) soloCounts() []int {
	out := make([]int, len(trackNames))
	for i := 0; i < steps; i++ {
		on, who := 0, 0
		for j, name := range trackNames {
			if p.Tracks[name][i] {
				on++
				who = j
			}
		}
		if on == 1 {
			out[who]++
		}
	}
	return out
}

// printCountsSolo prints the usual counts followed by the solo counts.
func printCountsSolo(w io.Writer, p pattern, asJSON bool) {
	solo := p.soloCounts()
	if !asJSON {
		printCounts(w, p, false)
		for i, name := range trackNames {
			fmt.Fprintf(w, "solo %s %d\n", name, solo[i])
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
	out += fmt.Sprintf("\"total\":%d,\"solo\":{", total)
	for i, name := range trackNames {
		if i > 0 {
			out += ","
		}
		out += fmt.Sprintf("%q:%d", name, solo[i])
	}
	fmt.Fprintln(w, out+"}}")
}
