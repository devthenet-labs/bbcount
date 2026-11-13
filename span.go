package main

import (
	"fmt"
	"io"
)

// spans returns, in trackNames order, how many steps each track covers from
// its first to its last playing step, both included; 0 for a silent track.
func (p pattern) spans() []int {
	out := make([]int, len(trackNames))
	for i, name := range trackNames {
		first, last := -1, -1
		for s, on := range p.Tracks[name] {
			if on {
				if first < 0 {
					first = s
				}
				last = s
			}
		}
		if first >= 0 {
			out[i] = last - first + 1
		}
	}
	return out
}

// printSpanCounts prints the counts like printCounts plus each track's span.
func printSpanCounts(w io.Writer, p pattern, asJSON bool) {
	sp := p.spans()
	if !asJSON {
		printCounts(w, p, false)
		for i, name := range trackNames {
			fmt.Fprintf(w, "span %s %d\n", name, sp[i])
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
	out += fmt.Sprintf("\"total\":%d,\"span\":{", total)
	for i, name := range trackNames {
		if i > 0 {
			out += ","
		}
		out += fmt.Sprintf("%q:%d", name, sp[i])
	}
	fmt.Fprintln(w, out+"}}")
}
