package main

import (
	"bytes"
	"encoding/binary"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

const fourOnFloor = `{"name":"Four","author":"me","bpm":120,"swing":0,"tracks":{
"kick":[true,false,false,false,true,false,false,false,true,false,false,false,true,false,false,false],
"snare":[false,false,false,false,true,false,false,false,false,false,false,false,true,false,false,false],
"hat":[true,false,true,false,true,false,true,false,true,false,true,false,true,false,true,false]}}`

func TestRunCounts(t *testing.T) {
	tests := []struct {
		name     string
		args     []string
		stdin    string
		wantCode int
		wantOut  string
		wantErr  string
	}{
		{"text", nil, fourOnFloor, 0, "kick 4\nsnare 2\nhat 8\ntotal 14\n", ""},
		{"json", []string{"-json"}, fourOnFloor, 0, `{"kick":4,"snare":2,"hat":8,"total":14}` + "\n", ""},
		{"-steps", []string{"-steps"}, fourOnFloor, 0, "kick 4\nsnare 2\nhat 8\ntotal 14\n" +
			"kick: 1 5 9 13\nsnare: 5 13\nhat: 1 3 5 7 9 11 13 15\n", ""},
		{"-steps no hits", []string{"-steps"}, strings.Replace(fourOnFloor,
			`"snare":[false,false,false,false,true,false,false,false,false,false,false,false,true,false,false,false]`,
			`"snare":[false,false,false,false,false,false,false,false,false,false,false,false,false,false,false,false]`, 1),
			0, "kick 4\nsnare 0\nhat 8\ntotal 12\n" +
				"kick: 1 5 9 13\nsnare: -\nhat: 1 3 5 7 9 11 13 15\n", ""},
		{"bad json", nil, `{`, 1, "", "invalid pattern"},
		{"unknown field", nil, strings.Replace(fourOnFloor, `"bpm"`, `"tempo":1,"bpm"`, 1), 1, "", "invalid pattern"},
		{"bpm low", nil, strings.Replace(fourOnFloor, `"bpm":120`, `"bpm":59`, 1), 1, "", "bpm"},
		{"bpm high", nil, strings.Replace(fourOnFloor, `"bpm":120`, `"bpm":201`, 1), 1, "", "bpm"},
		{"swing high", nil, strings.Replace(fourOnFloor, `"swing":0`, `"swing":61`, 1), 1, "", "swing"},
		{"short track", nil, strings.Replace(fourOnFloor, `"hat":[true,`, `"hat":[`, 1), 1, "", "hat"},
		{"missing track", nil, strings.Replace(fourOnFloor, `"snare"`, `"clap"`, 1), 1, "", "invalid pattern"},
		{"bad flag", []string{"-nope"}, "", 2, "", "usage"},
		{"two files", []string{"a", "b"}, "", 2, "", "usage"},
		{"missing file", []string{filepath.Join(t.TempDir(), "none.json")}, "", 1, "", "none.json"},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var out, errb bytes.Buffer
			code := run(tc.args, strings.NewReader(tc.stdin), &out, &errb)
			if code != tc.wantCode {
				t.Fatalf("exit %d, want %d (stderr %q)", code, tc.wantCode, errb.String())
			}
			if out.String() != tc.wantOut {
				t.Errorf("stdout %q, want %q", out.String(), tc.wantOut)
			}
			if !strings.Contains(errb.String(), tc.wantErr) {
				t.Errorf("stderr %q lacks %q", errb.String(), tc.wantErr)
			}
		})
	}
}

func TestRunReadsFile(t *testing.T) {
	path := filepath.Join(t.TempDir(), "p.json")
	if err := os.WriteFile(path, []byte(fourOnFloor), 0o600); err != nil {
		t.Fatal(err)
	}
	var out, errb bytes.Buffer
	if code := run([]string{path}, strings.NewReader(""), &out, &errb); code != 0 {
		t.Fatalf("exit %d: %s", code, errb.String())
	}
	if !strings.HasSuffix(out.String(), "total 14\n") {
		t.Errorf("stdout %q", out.String())
	}
}

func TestRunWritesWAV(t *testing.T) {
	path := filepath.Join(t.TempDir(), "beat.wav")
	var out, errb bytes.Buffer
	if code := run([]string{"-wav", path}, strings.NewReader(fourOnFloor), &out, &errb); code != 0 {
		t.Fatalf("exit %d: %s", code, errb.String())
	}
	// 120 bpm: a step is 125 ms, 16 steps are 2000 ms.
	if want := "kick 4\nsnare 2\nhat 8\ntotal 14\nwrote " + path + " (2000 ms)\n"; out.String() != want {
		t.Errorf("stdout %q, want %q", out.String(), want)
	}
	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if len(b) != 44+2*44100*2 {
		t.Fatalf("wav is %d bytes, want %d", len(b), 44+2*44100*2)
	}
	if string(b[0:4]) != "RIFF" || string(b[8:12]) != "WAVE" || string(b[12:16]) != "fmt " || string(b[36:40]) != "data" {
		t.Errorf("bad header % x", b[:44])
	}
	if got := binary.LittleEndian.Uint32(b[24:28]); got != 44100 {
		t.Errorf("sample rate %d", got)
	}
	if got := binary.LittleEndian.Uint16(b[22:24]); got != 1 {
		t.Errorf("channels %d", got)
	}
	if got := binary.LittleEndian.Uint16(b[34:36]); got != 16 {
		t.Errorf("bits %d", got)
	}
}

func TestRenderDeterministicAndSilentOnEmpty(t *testing.T) {
	p, err := parsePattern(strings.NewReader(fourOnFloor))
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(render(p), render(p)) {
		t.Error("render is not deterministic")
	}
	if bytes.Equal(render(p)[44:], make([]byte, 2*44100*2)) {
		t.Error("a pattern with hits rendered silence")
	}
	empty := pattern{BPM: 120, Tracks: map[string][]bool{}}
	for _, n := range trackNames {
		empty.Tracks[n] = make([]bool, steps)
	}
	if !bytes.Equal(render(empty)[44:], make([]byte, 2*44100*2)) {
		t.Error("an empty pattern must render silence")
	}
}

func TestStepStartMs(t *testing.T) {
	tests := []struct{ bpm, swing, step, want int }{
		{120, 0, 0, 0},
		{120, 0, 1, 125},
		{120, 0, 15, 1875},
		{120, 50, 1, 188}, // 125 + 62.5 rounds to the nearest millisecond
		{120, 50, 2, 250}, // even steps are never delayed
		{60, 60, 1, 400},  // 250 + 150
	}
	for _, tc := range tests {
		if got := stepStartMs(tc.bpm, tc.swing, tc.step); got != tc.want {
			t.Errorf("bpm %d swing %d step %d: %d ms, want %d", tc.bpm, tc.swing, tc.step, got, tc.want)
		}
	}
}
