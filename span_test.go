package main

import (
	"bytes"
	"path/filepath"
	"strings"
	"testing"
)

func TestRunSpan(t *testing.T) {
	cases := []struct {
		name string
		args []string
		want string
	}{
		{"text", []string{"-span"}, "kick 4\nsnare 2\nhat 8\ntotal 14\nspan kick 13\nspan snare 9\nspan hat 15\n"},
		{"json", []string{"-span", "-json"}, `{"kick":4,"snare":2,"hat":8,"total":14,"span":{"kick":13,"snare":9,"hat":15}}` + "\n"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			var out, errb bytes.Buffer
			if code := run(c.args, strings.NewReader(fourOnFloor), &out, &errb); code != 0 {
				t.Fatalf("code %d, stderr %q", code, errb.String())
			}
			if out.String() != c.want {
				t.Errorf("got %q, want %q", out.String(), c.want)
			}
		})
	}
}

func TestSpans(t *testing.T) {
	mk := func(steps ...int) []bool {
		tr := make([]bool, 16)
		for _, s := range steps {
			tr[s-1] = true
		}
		return tr
	}
	p := pattern{Tracks: map[string][]bool{"kick": mk(1), "snare": mk(8), "hat": mk()}}
	if got := p.spans(); got[0] != 1 || got[1] != 1 || got[2] != 0 {
		t.Errorf("got %v", got)
	}
	p = pattern{Tracks: map[string][]bool{"kick": mk(16), "snare": mk(1, 16), "hat": mk(3, 4)}}
	if got := p.spans(); got[0] != 1 || got[1] != 16 || got[2] != 2 {
		t.Errorf("got %v", got)
	}
}

func TestRunSpanWithWAV(t *testing.T) {
	path := filepath.Join(t.TempDir(), "o.wav")
	var out, errb bytes.Buffer
	if code := run([]string{"-span", "-wav", path}, strings.NewReader(fourOnFloor), &out, &errb); code != 0 {
		t.Fatalf("code %d, stderr %q", code, errb.String())
	}
	s := out.String()
	i, j := strings.Index(s, "span hat 15\n"), strings.Index(s, "wrote ")
	if i < 0 || j < i {
		t.Errorf("span lines should precede wrote line: %q", s)
	}
}
