package scanner

import (
	"bufio"
	"fmt"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"strings"

	getter "github.com/hashicorp/go-getter"
)

type Scanner interface {
	Scan(j *Job)
}

type scanner struct {
	path string
}

func NewScanner(cfg *Config) Scanner {
	return &scanner{
		path: cfg.Path,
	}
}

func (s *scanner) Scan(j *Job) {
	defer func() {
		j.Done <- struct{}{}
	}()

	p := s.dir(j.Repo.ID)
	if err := getter.Get(p, j.Repo.Link); err != nil {
		j.Err = err
		return
	}

	err := filepath.Walk(p, s.scan(j))
	if err != nil {
		j.Err = err
		return
	}
}

func (s *scanner) dir(id string) string {
	return path.Join(s.path, id)
}

func (s *scanner) scan(j *Job) filepath.WalkFunc {
	base := s.dir(j.Repo.ID)
	return func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() || info.Mode()&os.ModeSymlink == os.ModeSymlink {
			return nil
		}
		rel, err := filepath.Rel(base, path)
		if err != nil {
			return err
		}

		f, err := os.Open(path)
		if err != nil {
			return err
		}
		defer f.Close()

		var n int
		sc := bufio.NewScanner(f)
		for sc.Scan() {
			n++
			for _, word := range strings.Fields(sc.Text()) {
				if strings.HasPrefix(word, "public_key") ||
					strings.HasPrefix(word, "private_key") {
					// TODO: populate and result found pattern
					fmt.Printf("%s:%d %s\n", rel, n, word)
				}
			}
		}
		if err := sc.Err(); err != nil {
			return err
		}
		return nil
	}
}
