package scanner

import (
	"context"
	"errors"
	"io/fs"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	errInvalidFile = errors.New("invalid file")
)

type mockFileInfo struct {
	isDir bool
	mode  fs.FileMode
}

func (m mockFileInfo) Name() string       { return "" }
func (m mockFileInfo) Size() int64        { return 0 }
func (m mockFileInfo) Mode() fs.FileMode  { return 0 }
func (m mockFileInfo) ModTime() time.Time { return time.Now() }
func (m mockFileInfo) IsDir() bool        { return true }
func (m mockFileInfo) Sys() interface{}   { return nil }

func getScanner() *scanner {
	return &scanner{
		path: "/app/repos",
	}
}

func TestScannerNew(t *testing.T) {
	s := NewScanner("/app/repos")
	switch s.(type) {
	case Scanner:
		sc, ok := s.(*scanner)
		require.True(t, ok)
		assert.Equal(t, "/app/repos", sc.path)
	default:
		t.Errorf("invalid interface type")
	}
}

func TestScannerGetDir(t *testing.T) {
	s := getScanner()
	dir := s.dir("111")
	assert.Equal(t, "/app/repos/111", dir)
}

func TestScannerScanError(t *testing.T) {
	s := getScanner()

	job := &Job{
		Req: &Request{
			ID:  "1",
			URL: "https://example.com/test",
		},
		Ctx: context.TODO(),
	}
	err := s.scan(job)("/app/repos/1", nil, errInvalidFile)
	assert.Error(t, errInvalidFile, err)
}

func TestScannerScanSkipDir(t *testing.T) {
	s := getScanner()

	job := &Job{
		Req: &Request{
			ID:  "1",
			URL: "https://example.com/test",
		},
		Ctx: context.TODO(),
	}
	err := s.scan(job)("/app/repos/1", mockFileInfo{isDir: true}, nil)
	assert.NoError(t, err)
}
