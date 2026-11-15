package main

// changes returns, per track in trackNames order, how many times the track
// switches between playing and silent from one step to the next. The bar does not wrap.
func (p pattern) changes() []int {
	out := make([]int, len(trackNames))
	for i, name := range trackNames {
		t := p.Tracks[name]
		for s := 1; s < len(t); s++ {
			if t[s] != t[s-1] {
				out[i]++
			}
		}
	}
	return out
}
