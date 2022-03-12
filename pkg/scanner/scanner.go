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
	"gorm.io/datatypes"
)

type Scanner interface {
	Scan(j *Job)
}

type scanner struct {
	path string
}

func NewScanner(path string) Scanner {
	return &scanner{
		path: path,
	}
}

func (s *scanner) Scan(j *Job) {
	defer func() {
		j.Done <- struct{}{}
	}()

	p := s.dir(j.Repo.ID)
	if err := getter.Get(p, j.Repo.Url); err != nil {
		j.Err = err
		return
	}
	defer os.RemoveAll(p)

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

		var line int64
		sc := bufio.NewScanner(f)
		for sc.Scan() {
			line++
			for _, word := range strings.Fields(sc.Text()) {
				if strings.HasPrefix(word, "public_key") ||
					strings.HasPrefix(word, "private_key") {
					// TODO: populate and result found pattern

					j.Findings = append(j.Findings, datatypes.JSONMap{
						"type":   "xxx",
						"ruleId": "yyy",
						"locaton": datatypes.JSONMap{
							"path": rel,
							"position": datatypes.JSONMap{
								"begin": datatypes.JSONMap{
									"line": line,
								},
							},
						},
						"metadata": datatypes.JSONMap{
							"description": "zzz",
							"severity":    "zzz",
						},
					})
					fmt.Printf("%s:%d %s\n", rel, line, word)
				}
			}
		}
		if err := sc.Err(); err != nil {
			return err
		}
		return nil
	}
}
