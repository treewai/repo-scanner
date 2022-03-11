package server

type Config struct {
	Name      string
	ClientURL string
	Prefork   bool

	RepoDir   string
	ScanQueue int
}
