package main

import (
	"bytes"
	"path/filepath"
	"strings"
	"testing"
)

func TestLongestRest(t *testing.T) {
	mk := func(on ...int) []bool {
		tr := make([]bool, steps)
		for _, s := range on {
			tr[s-1] = true
		}
		return tr
	}
	alt := make([]bool, steps)
	full := make([]bool, steps)
	for i := range alt {
		alt[i] = i%2 == 0
		full[i] = true
	}
	tests := []struct {
		name  string
		track []bool
		want  int
	}{
		{"empty", mk(), 16},
		{"full", full, 0},
		{"step 1", mk(1), 15},
		{"step 16", mk(16), 15},
		{"steps 1 and 16", mk(1, 16), 14},
		{"alternating", alt, 1},
	}
	for _, tc := range tests {
		if got := longestRest(tc.track); got != tc.want {
			t.Errorf("%s: got %d, want %d", tc.name, got, tc.want)
		}
	}
}

func TestRunRest(t *testing.T) {
	run1 := func(args ...string) string {
		var out, errb bytes.Buffer
		if code := run(args, strings.NewReader(fourOnFloor), &out, &errb); code != 0 {
			t.Fatalf("exit %d: %s", code, errb.String())
		}
		return out.String()
	}
	want := "kick 4\nsnare 2\nhat 8\ntotal 14\nrest kick 3\nrest snare 7\nrest hat 1\n"
	if got := run1("-rest"); got != want {
		t.Errorf("got %q, want %q", got, want)
	}
	wantJSON := `{"kick":4,"snare":2,"hat":8,"total":14,"rest":{"kick":3,"snare":7,"hat":1}}` + "\n"
	if got := run1("-rest", "-json"); got != wantJSON {
		t.Errorf("got %q, want %q", got, wantJSON)
	}

	path := filepath.Join(t.TempDir(), "r.wav")
	got := run1("-rest", "-wav", path)
	if !strings.HasPrefix(got, want+"wrote ") {
		t.Errorf("rest lines not before wrote line: %q", got)
	}
}

func TestRunRestEmptyAndFull(t *testing.T) {
	f := strings.TrimSuffix(strings.Repeat("true,", 16), ",")
	e := strings.TrimSuffix(strings.Repeat("false,", 16), ",")
	in := `{"name":"x","author":"me","bpm":120,"swing":0,"tracks":{"kick":[` + e + `],"snare":[` + f + `],"hat":[` + e + `]}}`
	var out, errb bytes.Buffer
	if code := run([]string{"-rest"}, strings.NewReader(in), &out, &errb); code != 0 {
		t.Fatalf("exit %d: %s", code, errb.String())
	}
	want := "kick 0\nsnare 16\nhat 0\ntotal 16\nrest kick 16\nrest snare 0\nrest hat 16\n"
	if out.String() != want {
		t.Errorf("got %q, want %q", out.String(), want)
	}
}
