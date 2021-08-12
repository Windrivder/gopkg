package emitter

import (
	"github.com/olebedev/emitter"
)

var e = &emitter.Emitter{}

type (
	Event = emitter.Event
	Group = emitter.Group
)

func New(capacity uint) *emitter.Emitter {
	e = emitter.New(capacity)
	return e
}

// Use registers middlewares for the pattern.
func Use(pattern string, middlewares ...func(*Event)) {
	e.Use(pattern, middlewares...)
}

// On returns a channel that will receive events. As optional second
// argument it takes middlewares.
func On(topic string, middlewares ...func(*Event)) <-chan Event {
	return e.On(topic, middlewares...)
}

// Once works exactly like On(see above) but with `Once` as the first middleware.
func Once(topic string, middlewares ...func(*Event)) <-chan Event {
	return e.Once(topic, middlewares...)
}

// Off unsubscribes all listeners which were covered by
// topic, it can be pattern as well.
func Off(topic string, channels ...<-chan Event) {
	e.Off(topic, channels...)
}

// Listeners returns slice of listeners which were covered by
// topic(it can be pattern) and error if pattern is invalid.
func Listeners(topic string) []<-chan Event {
	return e.Listeners(topic)
}

// Topics returns all existing topics.
func Topics() []string {
	return e.Topics()
}

// Emit emits an event with the rest arguments to all
// listeners which were covered by topic(it can be pattern).
func Emit(topic string, args ...interface{}) chan struct{} {
	return e.Emit(topic, args...)
}

// ResetFlag middleware resets flags
func ResetFlag(evt *Event) { evt.Flags = emitter.FlagReset }

// Once middleware sets FlagOnce flag for an event
func FlagOnce(evt *Event) { evt.Flags = evt.Flags | emitter.FlagOnce }

// Void middleware sets FlagVoid flag for an event
func FlagVoid(evt *Event) { evt.Flags = evt.Flags | emitter.FlagVoid }

// Skip middleware sets FlagSkip flag for an event
func FlagSkip(evt *Event) { evt.Flags = evt.Flags | emitter.FlagSkip }

// Close middleware sets FlagClose flag for an event
func FlagClose(evt *Event) { evt.Flags = evt.Flags | emitter.FlagClose }

// Sync middleware sets FlagSync flag for an event
func FlagSync(evt *Event) { evt.Flags = evt.Flags | emitter.FlagSync }
