package jsonx

import (
	"bytes"
	"io"
	"strings"

	jsoniter "github.com/json-iterator/go"
)

var (
	json  = jsoniter.ConfigCompatibleWithStandardLibrary
	Valid = json.Valid
	// Marshal is exported by encoding/json package.
	Marshal = json.Marshal
	// Unmarshal is exported by encoding/json package.
	Unmarshal = json.Unmarshal
	// MarshalIndent is exported by encoding/json package.
	MarshalIndent = json.MarshalIndent
	// NewDecoder is exported by encoding/json package.
	NewDecoder = json.NewDecoder
	// NewEncoder is exported by encoding/json package.
	NewEncoder = json.NewEncoder
)

// Encode encodes any golang variable <value> to JSON bytes.
func Encode(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}

// Decode decodes json format <data> to golang variable.
func Decode(data []byte, v interface{}) error {
	decoder := NewDecoder(bytes.NewReader(data))
	if err := unmarshalUseNumber(decoder, v); err != nil {
		return err
	}

	return nil
}

func DecodeString(data string, v interface{}) error {
	decoder := json.NewDecoder(strings.NewReader(data))
	if err := unmarshalUseNumber(decoder, v); err != nil {
		return err
	}

	return nil
}

func DecodeReader(reader io.Reader, v interface{}) error {
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
