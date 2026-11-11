package main

import (
	"bytes"
	"fmt"
	"io"
	"strings"
)

// gridRow renders one track as steps characters: x for a step that plays and
// . for one that does not.
func gridRow(track []bool) string {
	var b strings.Builder
	for _, on := range track {
		if on {
			b.WriteByte('x')
		} else {
			b.WriteByte('.')
		}
	}
	return b.String()
}

// printCountsGrid prints what printCounts prints, followed by each track's
// grid: one "grid <track> <row>" line per track, or, with asJSON, a "grid"
// object after "total".
func printCountsGrid(w io.Writer, p pattern, asJSON bool) {
	var buf bytes.Buffer
	printCounts(&buf, p, asJSON)
	if !asJSON {
		buf.WriteTo(w)
		for _, name := range trackNames {
			fmt.Fprintf(w, "grid %s %s\n", name, gridRow(p.Tracks[name]))
		}
		return
	}
	out := strings.TrimSuffix(buf.String(), "}\n") + `,"grid":{`
	for i, name := range trackNames {
		if i > 0 {
			out += ","
		}
		out += fmt.Sprintf("%q:%q", name, gridRow(p.Tracks[name]))
	}
	fmt.Fprint(w, out+"}}\n")
}
