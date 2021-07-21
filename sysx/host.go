package sysx

import (
	"os"

	"github.com/windrivder/gopkg/util/randx"
)

var hostname string

func init() {
	var err error
	hostname, err = os.Hostname()
	if err != nil {
		hostname = randx.Letters(8)
	}
}

// Hostname returns the name of the host, if no hostname, a random id is returned.
func Hostname() string {
	return hostname
}
