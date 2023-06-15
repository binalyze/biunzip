package main

import (
	"context"
	"crypto/rand"
	"io"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRead(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	tests := []struct {
		name   string
		ctx    context.Context
		reader io.Reader
		err    error
	}{
		{
			name:   "with a valid context",
			ctx:    context.Background(),
			reader: rand.Reader,
			err:    nil,
		},
		{
			name:   "with a canceled context",
			ctx:    ctx,
			reader: rand.Reader,
			err:    context.Canceled,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader := newContextReader(tt.ctx, tt.reader)
			_, err := reader.Read(make([]byte, 1))
			require.ErrorIs(t, err, tt.err)
		})
	}
}
