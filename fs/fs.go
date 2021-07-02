package fs

import (
	"io"
	"os"
	"path/filepath"

	"github.com/spf13/afero"
)

var (
	NewBasePathFs    = afero.NewBasePathFs
	NewCacheOnReadFs = afero.NewCacheOnReadFs
	NewCopyOnWriteFs = afero.NewCopyOnWriteFs
	NewHttpFs        = afero.NewHttpFs
	NewIOFS          = afero.NewIOFS
	NewMemMapFs      = afero.NewMemMapFs
	NewOsFs          = afero.NewOsFs
	NewReadOnlyFs    = afero.NewReadOnlyFs
	NewRegexpFs      = afero.NewRegexpFs
)

func DirExists(path string) (bool, error) {
	return afero.DirExists(NewOsFs(), path)
}

func Exists(path string) (bool, error) {
	return afero.Exists(NewOsFs(), path)
}

func FileContainsAnyBytes(filename string, subslices [][]byte) (bool, error) {
	return afero.FileContainsAnyBytes(NewOsFs(), filename, subslices)
}

func FileContainsBytes(filename string, subslice []byte) (bool, error) {
	return afero.FileContainsBytes(NewOsFs(), filename, subslice)
}

func FullBaseFsPath(basePathFs *afero.BasePathFs, relativePath string) string {
	return afero.FullBaseFsPath(basePathFs, relativePath)
}

func GetTempDir(subPath string) string {
	return afero.GetTempDir(NewOsFs(), subPath)
}

func Glob(pattern string) (matches []string, err error) {
	return afero.Glob(NewOsFs(), pattern)
}

func IsDir(path string) (bool, error) {
	return afero.IsDir(NewOsFs(), path)
}

func IsEmpty(path string) (bool, error) {
	return afero.IsEmpty(NewOsFs(), path)
}

func NeuterAccents(s string) string {
	return afero.NeuterAccents(s)
}

func ReadAll(r io.Reader) ([]byte, error) {
	return afero.ReadAll(r)
}

func ReadDir(dirname string) ([]os.FileInfo, error) {
	return afero.ReadDir(NewOsFs(), dirname)
}

func ReadFile(filename string) ([]byte, error) {
	return afero.ReadFile(NewOsFs(), filename)
}

func SafeWriteReader(path string, r io.Reader) (err error) {
	return afero.SafeWriteReader(NewOsFs(), path, r)
}

func TempDir(dir, prefix string) (name string, err error) {
	return afero.TempDir(NewOsFs(), dir, prefix)
}

func UnicodeSanitize(s string) string {
	return afero.UnicodeSanitize(s)
}

func Walk(root string, walkFn filepath.WalkFunc) error {
	return afero.Walk(NewOsFs(), root, walkFn)
}

func WriteFile(filename string, data []byte, perm os.FileMode) error {
	return afero.WriteFile(NewOsFs(), filename, data, perm)
}

func WriteReader(path string, r io.Reader) (err error) {
	return afero.WriteReader(NewOsFs(), path, r)
}
