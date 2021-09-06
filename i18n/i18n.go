package i18n

import (
	"strings"

	"github.com/windrivder/gopkg/errorx"
)

var errUnmarshalNilLocate = errorx.New("can't unmarshal a nil *Locate")

type Locate int

const (
	LocateEN Locate = iota
	LocateZH
)

var (
	locates = [...]string{
		LocateEN: "en",
		LocateZH: "zh",
	}
)

func (l Locate) Int() int {
	return int(l)
}

func (l Locate) String() string {
	index := l.Int()
	if 0 <= index && index <= len(locates)-1 {
		return locates[l]
	}

	return ""
}

func (l Locate) CapitalString() string {
	return strings.ToUpper(l.String())
}

func (l Locate) MarshalText() ([]byte, error) {
	return []byte(l.String()), nil
}

func (l *Locate) UnmarshalText(text []byte) error {
	if l == nil {
		return errUnmarshalNilLocate
	}

	switch string(text) {
	case LocateEN.String(), LocateEN.CapitalString():
		*l = LocateEN
	case LocateZH.String(), LocateZH.CapitalString():
		*l = LocateZH
	default:
		return errorx.Newf("unrecognized locate: %q", text)
	}

	return nil
}

// Set sets the locate for the flag.Value interface.
func (l *Locate) Set(s string) error {
	return l.UnmarshalText([]byte(s))
}

// Get gets the locate for the flag.Getter interface.
func (l *Locate) Get() interface{} {
	return *l
}
