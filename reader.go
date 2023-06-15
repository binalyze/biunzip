package main

import (
	"context"
	"io"
)

type reader struct {
	ctx    context.Context
	reader io.Reader
}

func newContextReader(ctx context.Context, r io.Reader) io.Reader {
	return &reader{
		ctx:    ctx,
		reader: r,
	}
}

func (r *reader) Read(p []byte) (int, error) {
	err := r.ctx.Err()
	if err != nil {
		return 0, err
	}
	return r.reader.Read(p)
}
