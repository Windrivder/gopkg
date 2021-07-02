package md5

import (
	"crypto/md5"
	"fmt"
	"io"
	"os"

	"github.com/windrivder/gopkg/util/conv"
)

// Encrypt encrypts <data> using MD5 algorithms.
func Encrypt(data []byte) (encrypt string, err error) {
	h := md5.New()
	if _, err = h.Write(data); err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", h.Sum(nil)), nil
}

// MustEncrypt encrypts <data> using MD5 algorithms.
// It panics if any error occurs.
func MustEncrypt(data []byte) string {
	result, err := Encrypt(data)
	if err != nil {
		panic(err)
	}
	return result
}

// EncryptBytes encrypts string <data> using MD5 algorithms.
func EncryptString(data string) (encrypt string, err error) {
	return Encrypt(conv.StrToBytes(data))
}

// MustEncryptString encrypts string <data> using MD5 algorithms.
// It panics if any error occurs.
func MustEncryptString(data string) string {
	result, err := EncryptString(data)
	if err != nil {
		panic(err)
	}
	return result
}

// EncryptFile encrypts file content of <path> using MD5 algorithms.
func EncryptFile(path string) (encrypt string, err error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()
	h := md5.New()
	_, err = io.Copy(h, f)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", h.Sum(nil)), nil
}

// MustEncryptFile encrypts file content of <path> using MD5 algorithms.
// It panics if any error occurs.
func MustEncryptFile(path string) string {
	result, err := EncryptFile(path)
	if err != nil {
		panic(err)
	}
	return result
}
