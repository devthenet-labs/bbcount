package main

import (
	"bytes"
	"encoding/json"
	"path/filepath"
	"strings"
	"testing"
)

func TestRunGrid(t *testing.T) {
	silent := strings.Replace(fourOnFloor,
		`"snare":[false,false,false,false,true,false,false,false,false,false,false,false,true,false,false,false]`,
		`"snare":[false,false,false,false,false,false,false,false,false,false,false,false,false,false,false,false]`, 1)
	full := strings.Replace(fourOnFloor,
		`"kick":[true,false,false,false,true,false,false,false,true,false,false,false,true,false,false,false]`,
		`"kick":[true,true,true,true,true,true,true,true,true,true,true,true,true,true,true,true]`, 1)
	tests := []struct {
		name    string
		args    []string
		stdin   string
		wantOut string
	}{
		{"text", []string{"-grid"}, fourOnFloor,
			"kick 4\nsnare 2\nhat 8\ntotal 14\ngrid kick x...x...x...x...\ngrid snare ....x.......x...\ngrid hat x.x.x.x.x.x.x.x.\n"},
		{"json", []string{"-grid", "-json"}, fourOnFloor,
			`{"kick":4,"snare":2,"hat":8,"total":14,"grid":{"kick":"x...x...x...x...","snare":"....x.......x...","hat":"x.x.x.x.x.x.x.x."}}` + "\n"},
		{"silent snare", []string{"-grid"}, silent,
			"kick 4\nsnare 0\nhat 8\ntotal 12\ngrid kick x...x...x...x...\ngrid snare ................\ngrid hat x.x.x.x.x.x.x.x.\n"},
		{"full kick", []string{"-grid"}, full,
			"kick 16\nsnare 2\nhat 8\ntotal 26\ngrid kick xxxxxxxxxxxxxxxx\ngrid snare ....x.......x...\ngrid hat x.x.x.x.x.x.x.x.\n"},
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
			if strings.Contains(tc.name, "json") {
				var m map[string]any
				if err := json.Unmarshal(out.Bytes(), &m); err != nil {
					t.Errorf("invalid JSON: %v", err)
				}
			}
		})
	}
}

func TestRunGridWAV(t *testing.T) {
	path := filepath.Join(t.TempDir(), "out.wav")
	var out, errb bytes.Buffer
	if code := run([]string{"-grid", "-wav", path}, strings.NewReader(fourOnFloor), &out, &errb); code != 0 {
		t.Fatalf("exit %d: %s", code, errb.String())
	}
	want := "grid hat x.x.x.x.x.x.x.x.\nwrote " + path + " (2000 ms)\n"
	if !strings.HasSuffix(out.String(), want) {
		t.Errorf("stdout %q does not end with %q", out.String(), want)
	}
}
