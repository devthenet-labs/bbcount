// bbcount is Beat Battle's pattern counter, a small command-line tool: it
// reads one pattern JSON (a file argument, else standard input), prints how
// many steps each track plays and, with -wav, writes the pattern as a WAV
// file. With -first, it also prints where each track first plays. It is the
// CLI fixture for maul's agent evaluator, and uses the standard library only.
package main

import (
	"flag"
	"fmt"
	"io"
	"os"
)

func main() {
	os.Exit(run(os.Args[1:], os.Stdin, os.Stdout, os.Stderr))
}

// run is main without the process: it returns the exit code. 0 is success,
// 1 an unreadable or invalid pattern or a write failure and 2 a usage error.
func run(args []string, stdin io.Reader, stdout, stderr io.Writer) int {
	fs := flag.NewFlagSet("bbcount", flag.ContinueOnError)
	fs.SetOutput(stderr)
	fs.Usage = func() {
		fmt.Fprintln(stderr, "usage: bbcount [-json] [-first] [-wav out.wav] [pattern.json]")
		fs.PrintDefaults()
	}
	asJSON := fs.Bool("json", false, "print the counts as one JSON object")
	first := fs.Bool("first", false, "also print the first step each track plays")
	wavPath := fs.String("wav", "", "also write the pattern as a WAV file to this path")
	if err := fs.Parse(args); err != nil {
		return 2
	}
	if *first && *asJSON {
		fs.Usage()
		return 2
	}
	if fs.NArg() > 1 {
		fs.Usage()
		return 2
	}

	in := stdin
	if fs.NArg() == 1 {
		f, err := os.Open(fs.Arg(0))
		if err != nil {
			fmt.Fprintln(stderr, "bbcount:", err)
			return 1
		}
		defer f.Close()
		in = f
	}
	p, err := parsePattern(in)
	if err != nil {
		fmt.Fprintln(stderr, "bbcount:", err)
		return 1
	}

	printCounts(stdout, p, *asJSON, *first)
	if *wavPath != "" {
		if err := os.WriteFile(*wavPath, render(p), 0o644); err != nil {
			fmt.Fprintln(stderr, "bbcount:", err)
			return 1
		}
		fmt.Fprintf(stdout, "wrote %s (%d ms)\n", *wavPath, durationMs(p.BPM))
	}
	return 0
}

// printCounts prints the per-track step counts and their total.
func printCounts(w io.Writer, p pattern, asJSON, first bool) {
	counts := p.counts()
	total := 0
	for _, c := range counts {
		total += c
	}
	if asJSON {
		out := "{"
		for i, name := range trackNames {
			out += fmt.Sprintf("%q:%d,", name, counts[i])
		}
		fmt.Fprintf(w, "%s\"total\":%d}\n", out, total)
		return
	}
	for i, name := range trackNames {
		fmt.Fprintf(w, "%s %d\n", name, counts[i])
	}
	fmt.Fprintf(w, "total %d\n", total)
	if first {
		for _, name := range trackNames {
			step := p.firstStep(name)
			if step < 0 {
				fmt.Fprintf(w, "first %s none\n", name)
				continue
			}
			fmt.Fprintf(w, "first %s %d\n", name, step+1)
		}
	}
}
