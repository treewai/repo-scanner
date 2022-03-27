package scanner

import (
	"context"
)

type Job struct {
	Req      *Request
	Findings []map[string]interface{}
	Ctx      context.Context
	Err      error
	Done     chan struct{}
}

type Request struct {
	ID  string
	URL string
}
