package main

import (
	"bytes"
	"strings"
	"testing"
)

func trackOf(on ...int) []bool {
	t := make([]bool, steps)
	for _, s := range on {
		t[s-1] = true
	}
	return t
}

func TestPatternRuns(t *testing.T) {
	all := make([]int, steps)
	for i := range all {
		all[i] = i + 1
	}
	tests := []struct {
		name string
		on   []int
		want int
	}{
		{"groups", []int{1, 2, 3, 7}, 2},
		{"all", all, 1},
		{"ends no wrap", []int{1, 16}, 2},
		{"none", nil, 0},
		{"alternate", []int{1, 3, 5}, 3},
		{"last only", []int{16}, 1},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			p := pattern{Tracks: map[string][]bool{"kick": trackOf(tc.on...), "snare": trackOf(), "hat": trackOf()}}
			got := p.runs()
			if got[0] != tc.want || got[1] != 0 || got[2] != 0 {
				t.Fatalf("runs = %v, want kick %d", got, tc.want)
			}
		})
	}
}

func TestRunRuns(t *testing.T) {
	tests := []struct {
		name string
		args []string
		want string
	}{
		{"text", []string{"-runs"}, "kick 4\nsnare 2\nhat 8\ntotal 14\nruns kick 4\nruns snare 2\nruns hat 8\n"},
		{"json", []string{"-runs", "-json"}, `{"kick":4,"snare":2,"hat":8,"total":14,"runs":{"kick":4,"snare":2,"hat":8}}` + "\n"},
		{"text without", nil, "kick 4\nsnare 2\nhat 8\ntotal 14\n"},
		{"json without", []string{"-json"}, `{"kick":4,"snare":2,"hat":8,"total":14}` + "\n"},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var out, errb bytes.Buffer
			if code := run(tc.args, strings.NewReader(fourOnFloor), &out, &errb); code != 0 {
				t.Fatalf("exit %d: %s", code, errb.String())
			}
			if out.String() != tc.want {
				t.Fatalf("stdout = %q, want %q", out.String(), tc.want)
			}
		})
	}
}
