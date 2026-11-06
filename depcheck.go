package main

import "golang.org/x/text/width"

// fullWidthDigit is unused by the CLI; it exists only so this branch needs a
// module the offline tool-pod build cannot fetch (S2-26 step (e)).
var fullWidthDigit = width.Widen.String("1")
