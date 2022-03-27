package scanner

import (
	"bufio"
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

// Scan scans the pattern the given repository job
func (s *scanner) Scan(j *Job) {
	defer func() {
		j.Done <- struct{}{}
	}()

	p := s.dir(j.Req.ID)
	if err := getter.Get(p, j.Req.URL); err != nil {
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
	base := s.dir(j.Req.ID)
	return func(path string, info fs.FileInfo, err error) error {
		select {
		case <-j.Ctx.Done():
			return j.Ctx.Err()
		default:
		}

		if err != nil {
			return err
		}
		if info.IsDir() || info.Mode()&os.ModeSymlink == os.ModeSymlink {
			return nil
		}
		if !filepath.HasPrefix(path, base) {
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

					j.Findings = append(j.Findings, map[string]interface{}{
						"type":   "sast",
						"ruleId": "1",
						"locaton": map[string]interface{}{
							"path": rel,
							"position": map[string]interface{}{
								"begin": map[string]interface{}{
									"line": line,
								},
							},
						},
						"metadata": datatypes.JSONMap{
							"description": "Define secret key",
							"severity":    "HIGN",
						},
					})
				}
			}
		}
		if err := sc.Err(); err != nil {
			return err
		}
		return nil
	}
}
