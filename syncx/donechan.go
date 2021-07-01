package syncx

import (
	"sync"

	"github.com/windrivder/gopkg/typex"
)

type DoneChan struct {
	done chan typex.PlaceholderType
	once sync.Once
}

func NewDoneChan() *DoneChan {
	return &DoneChan{
		done: make(chan typex.PlaceholderType),
	}
}

func (dc *DoneChan) Close() {
	dc.once.Do(func() {
		close(dc.done)
	})
}

func (dc *DoneChan) Done() chan typex.PlaceholderType {
	return dc.done
}
