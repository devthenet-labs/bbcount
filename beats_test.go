package main

import (
	"bytes"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

func TestBeatsOutput(t *testing.T) {
	tests := []struct {
		name string
		args []string
		want string
	}{
		{"text", []string{"-beats"}, "kick 4\nsnare 2\nhat 8\ntotal 14\nbeats kick 1 1 1 1\nbeats snare 0 1 0 1\nbeats hat 2 2 2 2\n"},
		{"json", []string{"-json", "-beats"}, `{"kick":4,"snare":2,"hat":8,"total":14,"beats":{"kick":[1,1,1,1],"snare":[0,1,0,1],"hat":[2,2,2,2]}}` + "\n"},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var out, errb bytes.Buffer
			if code := run(tc.args, strings.NewReader(fourOnFloor), &out, &errb); code != 0 {
				t.Fatalf("exit %d: %s", code, errb.String())
			}
			if out.String() != tc.want {
				t.Errorf("stdout %q, want %q", out.String(), tc.want)
			}
		})
	}
}

func TestBeatsCounts(t *testing.T) {
	mk := func(on ...int) []bool {
		s := make([]bool, steps)
		for _, n := range on {
			s[n-1] = true
		}
		return s
	}
	all := make([]int, steps)
	for i := range all {
		all[i] = i + 1
	}
	tests := []struct {
		name  string
		steps []bool
		want  []int
	}{
		{"4 and 5", mk(4, 5), []int{1, 1, 0, 0}},
		{"12 and 13", mk(12, 13), []int{0, 0, 1, 1}},
		{"none", mk(), []int{0, 0, 0, 0}},
		{"all", mk(all...), []int{4, 4, 4, 4}},
		{"first", mk(1), []int{1, 0, 0, 0}},
		{"last", mk(16), []int{0, 0, 0, 1}},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			p := pattern{Tracks: map[string][]bool{"kick": tc.steps, "snare": mk(), "hat": mk()}}
			b := p.beats()
			if !reflect.DeepEqual(b[0], tc.want) {
				t.Errorf("kick beats %v, want %v", b[0], tc.want)
			}
			if !reflect.DeepEqual(b[1], []int{0, 0, 0, 0}) {
				t.Errorf("snare beats %v", b[1])
			}
		})
	}
}

func TestBeatsWithWAV(t *testing.T) {
	path := filepath.Join(t.TempDir(), "beat.wav")
	var out, errb bytes.Buffer
	if code := run([]string{"-beats", "-wav", path}, strings.NewReader(fourOnFloor), &out, &errb); code != 0 {
		t.Fatalf("exit %d: %s", code, errb.String())
	}
	want := "kick 4\nsnare 2\nhat 8\ntotal 14\nbeats kick 1 1 1 1\nbeats snare 0 1 0 1\nbeats hat 2 2 2 2\nwrote " + path + " (2000 ms)\n"
	if out.String() != want {
		t.Errorf("stdout %q, want %q", out.String(), want)
	}
}
