package emitter

import (
	"testing"
)

func TestEmitter(t *testing.T) {
	tests := []struct {
		topic string
	}{
		{
			topic: "create",
		},
		{
			topic: "modify",
		},
		{
			topic: "delete",
		},
	}

	for _, test := range tests {
		t.Run(test.topic, func(t *testing.T) {
			Emit(test.topic)
		})
	}

	for _, test := range tests {
		t.Run(test.topic, func(t *testing.T) {
			On(test.topic, func(e *Event) {
				t.Log(e.Topic)
			})
		})
	}
}
