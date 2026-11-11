// bbcount is Beat Battle's pattern counter, a small command-line tool: it
// reads one pattern JSON (a file argument, else standard input), prints how
// many steps each track plays and, with -wav, writes the pattern as a WAV
// file. It is the CLI fixture for maul's agent evaluator, and uses the
// standard library only.
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
		fmt.Fprintln(stderr, "usage: bbcount [-json] [-together] [-wav out.wav] [pattern.json]")
		fs.PrintDefaults()
	}
	asJSON := fs.Bool("json", false, "print the counts as one JSON object")
	wavPath := fs.String("wav", "", "also write the pattern as a WAV file to this path")
	showTogether := fs.Bool("together", false, "also print how many steps have two or more tracks playing")
	if err := fs.Parse(args); err != nil {
		return 2
	}
	if fs.NArg() > 1 {
		fs.Usage()
		return 2
	}

	if fs.NArg() == 1 {
		p, err := readPatternFile(fs.Arg(0))
		if err != nil {
			fmt.Fprintln(stderr, "bbcount:", err)
			return 1
		}
		return report(stdout, stderr, p, *wavPath, *showTogether, *asJSON)
	}
	p, err := parsePattern(stdin)
	if err != nil {
		fmt.Fprintln(stderr, "bbcount:", err)
		return 1
	}
	return report(stdout, stderr, p, *wavPath, *asJSON, *showTogether)
}

// readPatternFile parses the pattern in the named file, closing it before
// returning so a slow WAV render does not hold the file open.
func readPatternFile(path string) (pattern, error) {
	f, err := os.Open(path)
	if err != nil {
		return pattern{}, err
	}
	defer f.Close()
	return parsePattern(f)
}

// report prints the counts and, when wavPath is set, writes the WAV file.
func report(stdout, stderr io.Writer, p pattern, wavPath string, asJSON, withTogether bool) int {
	printCounts(stdout, p, asJSON, withTogether)
	if wavPath != "" {
		if err := os.WriteFile(wavPath, render(p), 0o644); err != nil {
			fmt.Fprintln(stderr, "bbcount:", err)
			return 1
		}
		fmt.Fprintf(stdout, "wrote %s (%d ms)\n", wavPath, durationMs(p.BPM))
	}
	return 0
}

// printCounts prints the per-track step counts and their total; with
// withTogether it also prints how many steps have two or more tracks playing.
func printCounts(w io.Writer, p pattern, asJSON, withTogether bool) {
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
		if withTogether {
			fmt.Fprintf(w, "%s\"total\":%d,\"together\":%d}\n", out, total, p.together())
			return
		}
		fmt.Fprintf(w, "%s\"total\":%d}\n", out, total)
		return
	}
	for i, name := range trackNames {
		fmt.Fprintf(w, "%s %d\n", name, counts[i])
	}
	fmt.Fprintf(w, "total %d\n", total)
	if withTogether {
		fmt.Fprintf(w, "together %d\n", p.together())
	}
}
