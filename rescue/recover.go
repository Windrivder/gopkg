package rescue

import "github.com/windrivder/gopkg/logx"

// Recover is used with defer to do cleanup on panics.
// Use it like:
//  defer Recover(func() {})
func Recover(cleanups ...func()) {
	for _, cleanup := range cleanups {
		cleanup()
	}

	if p := recover(); p != nil {
		// logx.ErrorStack(p)
		logx.Error().Msgf("%+v", p)
	}
}
