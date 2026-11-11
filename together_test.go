package main

import (
	"path/filepath"
	"strings"
	"testing"
)

func togetherPattern(kick, snare, hat string) string {
	row := func(s string) string {
		parts := make([]string, 16)
		for i := range parts {
			parts[i] = "false"
		}
		for _, f := range strings.Split(s, ",") {
			if f == "" {
				continue
			}
			var n int
			for _, c := range f {
				n = n*10 + int(c-'0')
			}
			parts[n-1] = "true"
		}
		return "[" + strings.Join(parts, ",") + "]"
	}
	return `{"name":"T","author":"me","bpm":120,"swing":0,"tracks":{"kick":` + row(kick) +
		`,"snare":` + row(snare) + `,"hat":` + row(hat) + `}}`
}

const allSteps = "1,2,3,4,5,6,7,8,9,10,11,12,13,14,15,16"

func TestTogetherRun(t *testing.T) {
	tests := []struct {
		name string
		args []string
		in   string
		want string
	}{
		{"text", []string{"-together"}, fourOnFloor, "kick 4\nsnare 2\nhat 8\ntotal 14\ntogether 4\n"},
		{"json", []string{"-together", "-json"}, fourOnFloor, `{"kick":4,"snare":2,"hat":8,"total":14,"together":4}` + "\n"},
		{"single tracks", []string{"-together"}, togetherPattern("1,2,3,4,5,6,7,8", "9,10,11,12,13,14,15,16", ""), "kick 8\nsnare 8\nhat 0\ntotal 16\ntogether 0\n"},
		{"all sixteen", []string{"-together"}, togetherPattern(allSteps, allSteps, "1,2"), "kick 16\nsnare 16\nhat 2\ntotal 34\ntogether 16\n"},
		{"no flag text", nil, fourOnFloor, "kick 4\nsnare 2\nhat 8\ntotal 14\n"},
		{"no flag json", []string{"-json"}, fourOnFloor, `{"kick":4,"snare":2,"hat":8,"total":14}` + "\n"},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var out, errb strings.Builder
			if code := run(tc.args, strings.NewReader(tc.in), &out, &errb); code != 0 {
				t.Fatalf("code %d, stderr %q", code, errb.String())
			}
			if out.String() != tc.want {
				t.Errorf("got %q, want %q", out.String(), tc.want)
			}
		})
	}
}

func TestTogetherWithWAV(t *testing.T) {
	path := filepath.Join(t.TempDir(), "out.wav")
	var out, errb strings.Builder
	if code := run([]string{"-together", "-wav", path}, strings.NewReader(fourOnFloor), &out, &errb); code != 0 {
		t.Fatalf("code %d, stderr %q", code, errb.String())
	}
	want := "kick 4\nsnare 2\nhat 8\ntotal 14\ntogether 4\nwrote " + path + " (2000 ms)\n"
	if out.String() != want {
		t.Errorf("got %q, want %q", out.String(), want)
	}
}
