package i18n

import (
	"strings"

	"github.com/windrivder/gopkg/errorx"
)

var errUnmarshalNilLocate = errorx.New("can't unmarshal a nil *Locate")

type Locale int

const (
	LocaleEN Locale = iota
	LocaleZH
)

var (
	locates = [...]string{
		LocaleEN: "en",
		LocaleZH: "zh",
	}
)

func (l Locale) Int() int {
	return int(l)
}

func (l Locale) String() string {
	index := l.Int()
	if 0 <= index && index <= len(locates)-1 {
		return locates[l]
	}

	return ""
}

func (l Locale) CapitalString() string {
	return strings.ToUpper(l.String())
}

func (l Locale) MarshalText() ([]byte, error) {
	return []byte(l.String()), nil
}

func (l *Locale) UnmarshalText(text []byte) error {
	if l == nil {
		return errUnmarshalNilLocate
	}

	switch string(text) {
	case LocaleEN.String(), LocaleEN.CapitalString():
		*l = LocaleEN
	case LocaleZH.String(), LocaleZH.CapitalString():
		*l = LocaleZH
	default:
		return errorx.Newf("unrecognized locate: %q", text)
	}

	return nil
}

// Set sets the locate for the flag.Value interface.
func (l *Locale) Set(s string) error {
	return l.UnmarshalText([]byte(s))
}

// Get gets the locate for the flag.Getter interface.
func (l *Locale) Get() interface{} {
	return *l
}
