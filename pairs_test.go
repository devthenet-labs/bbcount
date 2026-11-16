package main

import (
	"bytes"
	"fmt"
	"path/filepath"
	"strings"
	"testing"
)

// patternJSON builds a pattern from 1-based step lists per track.
func patternJSON(kick, snare, hat []int) string {
	track := func(on []int) string {
		var cells []string
		for i := 1; i <= steps; i++ {
			v := "false"
			for _, s := range on {
				if s == i {
					v = "true"
				}
			}
			cells = append(cells, v)
		}
		return "[" + strings.Join(cells, ",") + "]"
	}
	return fmt.Sprintf(`{"name":"t","author":"me","bpm":120,"swing":0,"tracks":{"kick":%s,"snare":%s,"hat":%s}}`,
		track(kick), track(snare), track(hat))
}

func TestRunPairs(t *testing.T) {
	spec := patternJSON([]int{1, 2, 3}, []int{2, 9}, []int{2, 3, 16})
	silent := patternJSON(nil, nil, nil)
	tests := []struct {
		name    string
		args    []string
		stdin   string
		wantOut string
	}{
		{"text", []string{"-pairs"}, fourOnFloor, "kick 4\nsnare 2\nhat 8\ntotal 14\npair kick snare 2\npair kick hat 4\npair snare hat 2\n"},
		{"json", []string{"-pairs", "-json"}, fourOnFloor, `{"kick":4,"snare":2,"hat":8,"total":14,"pairs":{"kick_snare":2,"kick_hat":4,"snare_hat":2}}` + "\n"},
		{"spec example", []string{"-pairs"}, spec, "kick 3\nsnare 2\nhat 3\ntotal 8\npair kick snare 1\npair kick hat 2\npair snare hat 1\n"},
		{"silent text", []string{"-pairs"}, silent, "kick 0\nsnare 0\nhat 0\ntotal 0\npair kick snare 0\npair kick hat 0\npair snare hat 0\n"},
		{"silent json", []string{"-pairs", "-json"}, silent, `{"kick":0,"snare":0,"hat":0,"total":0,"pairs":{"kick_snare":0,"kick_hat":0,"snare_hat":0}}` + "\n"},
		{"no pairs text", nil, fourOnFloor, "kick 4\nsnare 2\nhat 8\ntotal 14\n"},
		{"no pairs json", []string{"-json"}, fourOnFloor, `{"kick":4,"snare":2,"hat":8,"total":14}` + "\n"},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var out, errb bytes.Buffer
			if code := run(tc.args, strings.NewReader(tc.stdin), &out, &errb); code != 0 {
				t.Fatalf("exit %d: %s", code, errb.String())
			}
			if out.String() != tc.wantOut {
				t.Errorf("stdout %q, want %q", out.String(), tc.wantOut)
			}
		})
	}
}

func TestRunPairsWAV(t *testing.T) {
	path := filepath.Join(t.TempDir(), "beat.wav")
	var out, errb bytes.Buffer
	if code := run([]string{"-pairs", "-wav", path}, strings.NewReader(fourOnFloor), &out, &errb); code != 0 {
		t.Fatalf("exit %d: %s", code, errb.String())
	}
	want := "kick 4\nsnare 2\nhat 8\ntotal 14\npair kick snare 2\npair kick hat 4\npair snare hat 2\nwrote " + path + " (2000 ms)\n"
	if out.String() != want {
		t.Errorf("stdout %q, want %q", out.String(), want)
	}
}
