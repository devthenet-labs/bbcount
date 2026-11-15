package main

import (
	"bytes"
	"strings"
	"testing"
)

func TestRunHalf(t *testing.T) {
	tests := []struct {
		name string
		args []string
		want string
	}{
		{"text", []string{"-half"}, "kick 4\nsnare 2\nhat 8\ntotal 14\nhalf kick 2 2\nhalf snare 1 1\nhalf hat 4 4\n"},
		{"json", []string{"-half", "-json"}, `{"kick":4,"snare":2,"hat":8,"total":14,"half":{"kick":[2,2],"snare":[1,1],"hat":[4,4]}}` + "\n"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var out, errb bytes.Buffer
			if code := run(tt.args, strings.NewReader(fourOnFloor), &out, &errb); code != 0 {
				t.Fatalf("code %d, stderr %q", code, errb.String())
			}
			if out.String() != tt.want {
				t.Errorf("got %q, want %q", out.String(), tt.want)
			}
		})
	}
}

func TestHalves(t *testing.T) {
	all := make([]int, steps)
	for i := range all {
		all[i] = i + 1
	}
	tests := []struct {
		name  string
		steps []int
		want  [2]int
	}{
		{"mixed", []int{2, 3, 5, 12}, [2]int{3, 1}},
		{"only 8", []int{8}, [2]int{1, 0}},
		{"only 9", []int{9}, [2]int{0, 1}},
		{"all", all, [2]int{8, 8}},
		{"none", nil, [2]int{0, 0}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			track := make([]bool, steps)
			for _, s := range tt.steps {
				track[s-1] = true
			}
			p := pattern{Tracks: map[string][]bool{
				"kick": track, "snare": make([]bool, steps), "hat": make([]bool, steps),
			}}
			if got := p.halves()[0]; got != tt.want {
				t.Errorf("got %v, want %v", got, tt.want)
			}
		})
	}
}
