package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMakeDirPath(t *testing.T) {
	filePath := "/tmp/file_1.zip"
	expected := "/tmp/file_1"
	actual := makeDirPath(filePath)
	require.Equal(t, expected, actual)
}
