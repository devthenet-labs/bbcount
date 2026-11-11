package main

// together returns how many of the steps have two or more tracks active.
func (p pattern) together() int {
	n := 0
	for s := 0; s < steps; s++ {
		active := 0
		for _, name := range trackNames {
			if p.Tracks[name][s] {
				active++
			}
		}
		if active >= 2 {
			n++
		}
	}
	return n
}
