package main

import (
	"path/filepath"
	"strings"
	"testing"
)

func TestRunOffbeat(t *testing.T) {
	cases := []struct {
		name string
		args []string
		want string
	}{
		{"text", []string{"-offbeat"}, "kick 4\nsnare 2\nhat 8\ntotal 14\noffbeat kick 0\noffbeat snare 0\noffbeat hat 4\n"},
		{"json", []string{"-json", "-offbeat"}, `{"kick":4,"snare":2,"hat":8,"total":14,"offbeat":{"kick":0,"snare":0,"hat":4}}` + "\n"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var out, errb strings.Builder
			if code := run(tc.args, strings.NewReader(fourOnFloor), &out, &errb); code != 0 {
				t.Fatalf("code %d, stderr %q", code, errb.String())
			}
			if out.String() != tc.want {
				t.Errorf("got %q, want %q", out.String(), tc.want)
			}
		})
	}
}

func TestOffbeatCounts(t *testing.T) {
	all := make([]bool, steps)
	for i := range all {
		all[i] = true
	}
	kick := make([]bool, steps)
	kick[1], kick[2], kick[4] = true, true, true
	p := pattern{Tracks: map[string][]bool{"kick": kick, "snare": all, "hat": make([]bool, steps)}}
	got := p.offbeatCounts()
	want := []int{2, 12, 0}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("track %d: got %d, want %d", i, got[i], want[i])
		}
	}
}

func TestOffbeatWithWAV(t *testing.T) {
	path := filepath.Join(t.TempDir(), "o.wav")
	var out, errb strings.Builder
	if code := run([]string{"-offbeat", "-wav", path}, strings.NewReader(fourOnFloor), &out, &errb); code != 0 {
		t.Fatalf("code %d, stderr %q", code, errb.String())
	}
	s := out.String()
	if strings.Index(s, "offbeat hat 4") < 0 || strings.Index(s, "offbeat hat 4") > strings.Index(s, "wrote ") {
		t.Errorf("offbeat lines should precede wrote line: %q", s)
	}
}
