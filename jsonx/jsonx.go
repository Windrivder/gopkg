package jsonx

import (
	"bytes"
	"io"
	"strings"

	jsoniter "github.com/json-iterator/go"
	"github.com/windrivder/gopkg/cast"
)

var (
	json = jsoniter.ConfigCompatibleWithStandardLibrary
	// Marshal is exported by encoding/json package.
	Marshal = json.Marshal
	// Unmarshal is exported by encoding/json package.
	// Unmarshal = json.Unmarshal
	// MarshalIndent is exported by encoding/json package.
	MarshalIndent = json.MarshalIndent
	// NewDecoder is exported by encoding/json package.
	NewDecoder = json.NewDecoder
	// NewEncoder is exported by encoding/json package.
	NewEncoder = json.NewEncoder
)

func Unmarshal(data []byte, v interface{}) error {
	decoder := NewDecoder(bytes.NewReader(data))
	if err := unmarshalUseNumber(decoder, v); err != nil {
		return err
	}

	return nil
}

func UnmarshalFromString(data string, v interface{}) error {
	return Unmarshal(cast.StringToByteSlice(data), v)
}

func UnmarshalFromReader(reader io.Reader, v interface{}) error {
	var buf strings.Builder
	teeReader := io.TeeReader(reader, &buf)
	decoder := NewDecoder(teeReader)
	if err := unmarshalUseNumber(decoder, v); err != nil {
		return err
	}

	return nil
}

func unmarshalUseNumber(decoder *jsoniter.Decoder, v interface{}) error {
	decoder.UseNumber()
	return decoder.Decode(v)
}
