package main

import (
	"bytes"
	"fmt"
	"strings"
	"testing"
)

func longestTrack(on ...int) []bool {
	t := make([]bool, steps)
	for _, s := range on {
		t[s-1] = true
	}
	return t
}

func TestLongestRun(t *testing.T) {
	all := make([]int, steps)
	for i := range all {
		all[i] = i + 1
	}
	tests := []struct {
		name string
		on   []int
		want int
	}{
		{"2,3,4,9", []int{2, 3, 4, 9}, 3},
		{"15,16", []int{15, 16}, 2},
		{"1 and 16 no wrap", []int{1, 16}, 1},
		{"all", all, 16},
		{"none", nil, 0},
		{"single", []int{7}, 1},
		{"two runs", []int{1, 2, 8, 9, 10, 11}, 4},
	}
	for _, tc := range tests {
		if got := longestRun(longestTrack(tc.on...)); got != tc.want {
			t.Errorf("%s: got %d, want %d", tc.name, got, tc.want)
		}
	}
}

func longestJSON(kick, snare, hat []bool) string {
	f := func(b []bool) string {
		parts := make([]string, len(b))
		for i, v := range b {
			parts[i] = fmt.Sprint(v)
		}
		return "[" + strings.Join(parts, ",") + "]"
	}
	return `{"name":"L","author":"me","bpm":120,"swing":0,"tracks":{"kick":` + f(kick) + `,"snare":` + f(snare) + `,"hat":` + f(hat) + `}}`
}

func TestRunLongest(t *testing.T) {
	all := make([]int, steps)
	for i := range all {
		all[i] = i + 1
	}
	custom := longestJSON(longestTrack(2, 3, 4, 9), longestTrack(15, 16), longestTrack(all...))
	silent := longestJSON(longestTrack(), longestTrack(1, 16), longestTrack())
	tests := []struct {
		name  string
		args  []string
		stdin string
		want  string
	}{
		{"plain", []string{"-longest"}, fourOnFloor, "kick 4\nsnare 2\nhat 8\ntotal 14\nlongest kick 1\nlongest snare 1\nlongest hat 1\n"},
		{"json", []string{"-json", "-longest"}, fourOnFloor, `{"kick":4,"snare":2,"hat":8,"total":14,"longest":{"kick":1,"snare":1,"hat":1}}` + "\n"},
		{"custom", []string{"-longest"}, custom, "kick 4\nsnare 2\nhat 16\ntotal 22\nlongest kick 3\nlongest snare 2\nlongest hat 16\n"},
		{"silent", []string{"-longest"}, silent, "kick 0\nsnare 2\nhat 0\ntotal 2\nlongest kick 0\nlongest snare 1\nlongest hat 0\n"},
	}
	for _, tc := range tests {
		var out, errb bytes.Buffer
		if code := run(tc.args, strings.NewReader(tc.stdin), &out, &errb); code != 0 {
			t.Fatalf("%s: exit %d (%s)", tc.name, code, errb.String())
		}
		if out.String() != tc.want {
			t.Errorf("%s: got %q, want %q", tc.name, out.String(), tc.want)
		}
	}
}
