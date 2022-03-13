package scanner

import (
	"context"

	"gorm.io/datatypes"
)

type Job struct {
	Req      *Request
	Findings []datatypes.JSONMap
	Ctx      context.Context
	Err      error
	Done     chan struct{}
}

type Request struct {
	ID  string
	URL string
}
