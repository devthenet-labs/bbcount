package main

import (
	"bytes"
	"strings"
	"testing"
)

func patternWithKick(on ...int) pattern {
	p := pattern{Tracks: map[string][]bool{}}
	for _, name := range trackNames {
		p.Tracks[name] = make([]bool, steps)
	}
	for _, s := range on {
		p.Tracks["kick"][s-1] = true
	}
	return p
}

func TestPatternChanges(t *testing.T) {
	all := make([]int, steps)
	for i := range all {
		all[i] = i + 1
	}
	tests := []struct {
		name string
		on   []int
		want int
	}{
		{"1237", []int{1, 2, 3, 7}, 3},
		{"only 1", []int{1}, 1},
		{"only 16", []int{16}, 1},
		{"all", all, 0},
		{"none", nil, 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := patternWithKick(tt.on...).changes()[0]; got != tt.want {
				t.Errorf("changes kick = %d, want %d", got, tt.want)
			}
		})
	}
}

func TestRunChanges(t *testing.T) {
	const plain = "kick 4\nsnare 2\nhat 8\ntotal 14\n"
	const js = `{"kick":4,"snare":2,"hat":8,"total":14`
	tests := []struct {
		name string
		args []string
		want string
	}{
		{"text", []string{"-changes"}, plain + "changes kick 7\nchanges snare 4\nchanges hat 15\n"},
		{"json", []string{"-json", "-changes"}, js + `,"changes":{"kick":7,"snare":4,"hat":15}}` + "\n"},
		{"json reversed", []string{"-changes", "-json"}, js + `,"changes":{"kick":7,"snare":4,"hat":15}}` + "\n"},
		{"text no flag", nil, plain},
		{"json no flag", []string{"-json"}, js + "}\n"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var out, errb bytes.Buffer
			if code := run(tt.args, strings.NewReader(fourOnFloor), &out, &errb); code != 0 {
				t.Fatalf("exit %d: %s", code, errb.String())
			}
			if out.String() != tt.want {
				t.Errorf("got %q, want %q", out.String(), tt.want)
			}
		})
	}
}
