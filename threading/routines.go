package threading

import "github.com/windrivder/gopkg/rescue"

// GoSafe runs the given fn using another goroutine, recovers if fn panics.
func GoSafe(fn func()) {
	go RunSafe(fn)
}

// RunSafe runs the given fn, recovers if fn panics.
func RunSafe(fn func()) {
	defer rescue.Recover()

	fn()
}
