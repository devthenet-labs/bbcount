package main

import (
	"bytes"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

func mkPattern(kick, snare, hat []int) pattern {
	t := map[string][]bool{}
	for name, on := range map[string][]int{"kick": kick, "snare": snare, "hat": hat} {
		s := make([]bool, steps)
		for _, n := range on {
			s[n-1] = true
		}
		t[name] = s
	}
	return pattern{Tracks: t}
}

func all() []int {
	var s []int
	for i := 1; i <= steps; i++ {
		s = append(s, i)
	}
	return s
}

func TestSoloCounts(t *testing.T) {
	tests := []struct {
		name string
		p    pattern
		want []int
	}{
		{"example", mkPattern([]int{1, 2, 3}, []int{2, 9}, nil), []int{2, 1, 0}},
		{"all on", mkPattern(all(), all(), all()), []int{0, 0, 0}},
		{"all off", mkPattern(nil, nil, nil), []int{0, 0, 0}},
		{"kick alone", mkPattern([]int{1}, nil, nil), []int{1, 0, 0}},
		{"snare alone", mkPattern(nil, []int{1}, nil), []int{0, 1, 0}},
		{"hat alone", mkPattern(nil, nil, []int{1}), []int{0, 0, 1}},
	}
	for _, tc := range tests {
		if got := tc.p.soloCounts(); !reflect.DeepEqual(got, tc.want) {
			t.Errorf("%s: %v, want %v", tc.name, got, tc.want)
		}
	}
}

func TestSoloParsed(t *testing.T) {
	p, err := parsePattern(strings.NewReader(fourOnFloor))
	if err != nil {
		t.Fatal(err)
	}
	if got := p.soloCounts(); !reflect.DeepEqual(got, []int{0, 0, 4}) {
		t.Errorf("got %v", got)
	}
}

func TestRunSolo(t *testing.T) {
	var out, errb bytes.Buffer
	if code := run([]string{"-solo"}, strings.NewReader(fourOnFloor), &out, &errb); code != 0 {
		t.Fatal(code, errb.String())
	}
	want := "kick 4\nsnare 2\nhat 8\ntotal 14\nsolo kick 0\nsolo snare 0\nsolo hat 4\n"
	if out.String() != want {
		t.Errorf("got %q", out.String())
	}
	out.Reset()
	run([]string{"-solo", "-json"}, strings.NewReader(fourOnFloor), &out, &errb)
	want = `{"kick":4,"snare":2,"hat":8,"total":14,"solo":{"kick":0,"snare":0,"hat":4}}` + "\n"
	if out.String() != want {
		t.Errorf("got %q", out.String())
	}
}

func TestRunSoloWav(t *testing.T) {
	path := filepath.Join(t.TempDir(), "o.wav")
	var out, errb bytes.Buffer
	if code := run([]string{"-solo", "-wav", path}, strings.NewReader(fourOnFloor), &out, &errb); code != 0 {
		t.Fatal(code, errb.String())
	}
	s := out.String()
	a, b := strings.Index(s, "solo hat 4"), strings.Index(s, "wrote ")
	if a < 0 || b < 0 || a > b {
		t.Errorf("order wrong: %q", s)
	}
}
