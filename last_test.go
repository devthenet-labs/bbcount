package main

import (
	"bytes"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

func runLast(t *testing.T, args []string, in string) string {
	t.Helper()
	var out, errb bytes.Buffer
	if code := run(args, strings.NewReader(in), &out, &errb); code != 0 {
		t.Fatalf("exit %d: %s", code, errb.String())
	}
	return out.String()
}

func TestLastText(t *testing.T) {
	want := "kick 4\nsnare 2\nhat 8\ntotal 14\nlast kick 13\nlast snare 13\nlast hat 15\n"
	if got := runLast(t, []string{"-last"}, fourOnFloor); got != want {
		t.Fatalf("got %q want %q", got, want)
	}
}

func TestLastJSON(t *testing.T) {
	want := `{"kick":4,"snare":2,"hat":8,"total":14,"last":{"kick":13,"snare":13,"hat":15}}` + "\n"
	if got := runLast(t, []string{"-last", "-json"}, fourOnFloor); got != want {
		t.Fatalf("got %q want %q", got, want)
	}
}

func mk(k, s, h []int) pattern {
	tr := map[string][]bool{}
	for name, idx := range map[string][]int{"kick": k, "snare": s, "hat": h} {
		b := make([]bool, steps)
		for _, i := range idx {
			b[i-1] = true
		}
		tr[name] = b
	}
	return pattern{Tracks: tr}
}

func TestLasts(t *testing.T) {
	tests := []struct {
		name string
		p    pattern
		want []int
	}{
		{"mixed", mk([]int{1, 2, 3}, []int{2, 9}, []int{16}), []int{3, 9, 16}},
		{"silent", mk(nil, nil, nil), []int{0, 0, 0}},
		{"step one", mk([]int{1}, nil, nil), []int{1, 0, 0}},
	}
	for _, tc := range tests {
		if got := tc.p.lasts(); !reflect.DeepEqual(got, tc.want) {
			t.Errorf("%s: got %v want %v", tc.name, got, tc.want)
		}
	}
	var b bytes.Buffer
	printCountsLast(&b, mk(nil, nil, nil), false)
	if want := "last kick 0\nlast snare 0\nlast hat 0\n"; !strings.HasSuffix(b.String(), want) {
		t.Errorf("silent output %q", b.String())
	}
}

func TestLastWav(t *testing.T) {
	path := filepath.Join(t.TempDir(), "b.wav")
	got := runLast(t, []string{"-last", "-wav", path}, fourOnFloor)
	pre := "kick 4\nsnare 2\nhat 8\ntotal 14\nlast kick 13\nlast snare 13\nlast hat 15\nwrote "
	if !strings.HasPrefix(got, pre) {
		t.Fatalf("got %q", got)
	}
}
