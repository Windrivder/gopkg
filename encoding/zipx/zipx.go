package zipx

import (
	"bytes"
	"compress/gzip"
	"io"
	"io/ioutil"

	"github.com/windrivder/gopkg/util/conv"
)

const unzipLimit = 100 * 1024 * 1024 // 100MB

// Encode encodes bytes with Gzip algorithm.
func Encode(src []byte) []byte {
	var b bytes.Buffer

	w := gzip.NewWriter(&b)
	w.Write(src)
	w.Close()

	return b.Bytes()
}

// EncodeString encodes string with Gzip algorithm.
func EncodeString(src string) []byte {
	return Encode(conv.StrToBytes(src))
}

// EncodeToString encodes bytes to string with Gzip algorithm.
func EncodeToString(src []byte) string {
	return conv.BytesToStr(Encode(src))
}

// EncryptFile encodes file content of <path> using Gzip algorithms.
func EncodeFile(path string) ([]byte, error) {
	content, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return Encode(content), nil
}

// MustEncodeFile encodes file content of <path> using Gzip algorithms.
// It panics if any error occurs.
func MustEncodeFile(path string) []byte {
	result, err := EncodeFile(path)
	if err != nil {
		panic(err)
	}
	return result
}

// EncodeFileToString encodes file content of <path> to string using Gzip algorithms.
func EncodeFileToString(path string) (string, error) {
	content, err := EncodeFile(path)
	if err != nil {
		return "", err
	}
	return conv.BytesToStr(content), nil
}

// MustEncodeFileToString encodes file content of <path> to string using Gzip algorithms.
// It panics if any error occurs.
func MustEncodeFileToString(path string) string {
	result, err := EncodeFileToString(path)
	if err != nil {
		panic(err)
	}
	return result
}

// Decode decodes bytes with Gzip algorithm.
func Decode(bs []byte) ([]byte, error) {
	r, err := gzip.NewReader(bytes.NewBuffer(bs))
	if err != nil {
		return nil, err
	}
	defer r.Close()

	var c bytes.Buffer
	if _, err = io.Copy(&c, io.LimitReader(r, unzipLimit)); err != nil {
		return nil, err
	}

	return c.Bytes(), nil
}

// MustDecode decodes bytes with Gzip algorithm.
// It panics if any error occurs.
func MustDecode(data []byte) []byte {
	result, err := Decode(data)
	if err != nil {
		panic(err)
	}
	return result
}

// DecodeString decodes string with Gzip algorithm.
func DecodeString(data string) ([]byte, error) {
	return Decode(conv.StrToBytes(data))
}

// MustDecodeString decodes string with Gzip algorithm.
// It panics if any error occurs.
func MustDecodeString(data string) []byte {
	result, err := DecodeString(data)
	if err != nil {
		panic(err)
	}
	return result
}

// DecodeString decodes string with Gzip algorithm.
func DecodeToString(data string) (string, error) {
	b, err := DecodeString(data)
	return conv.BytesToStr(b), err
}

// MustDecodeToString decodes string with Gzip algorithm.
// It panics if any error occurs.
func MustDecodeToString(data string) string {
	result, err := DecodeToString(data)
	if err != nil {
		panic(err)
	}
	return result
}
