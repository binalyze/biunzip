package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewSemaphore(t *testing.T) {
	sem := newSemaphore(2)
	require.Equal(t, 0, len(sem))
	require.Equal(t, 2, cap(sem))
}

func TestAcquire(t *testing.T) {
	sem := newSemaphore(2)
	sem.acquire()
	require.Equal(t, 1, len(sem))
	require.Equal(t, 2, cap(sem))
}

func TestRelease(t *testing.T) {
	sem := newSemaphore(2)
	sem.acquire()
	sem.release()
	require.Equal(t, 0, len(sem))
	require.Equal(t, 2, cap(sem))
}

func TestWait(t *testing.T) {
	sem := newSemaphore(2)
	sem.wait()
	require.Equal(t, 2, len(sem))
	require.Equal(t, 2, cap(sem))
}
