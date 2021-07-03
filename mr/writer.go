package mr

import "github.com/windrivder/gopkg/container/typex"

type Writer interface {
	Write(v typex.GenericType)
}

type guardedWriter struct {
	channel chan<- interface{}
	done    <-chan typex.PlaceholderType
}

func newGuardedWriter(
	channel chan<- typex.GenericType,
	done <-chan typex.PlaceholderType,
) guardedWriter {
	return guardedWriter{
		channel: channel,
		done:    done,
	}
}

func (gw guardedWriter) Write(v typex.GenericType) {
	select {
	case <-gw.done:
		return
	default:
		gw.channel <- v
	}
}
