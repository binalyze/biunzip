package main

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMakeMultiErr(t *testing.T) {
	err1 := errors.New("error 1")
	err2 := errors.New("error 2")
	errs := []error{err1, err2}
	expectedErrMsg := fmt.Sprintf("test:\n- %s\n- %s", err1.Error(), err2.Error())
	err := makeMultiErr("test", errs)
	require.Equal(t, expectedErrMsg, err.Error())
}

func TestJoinMultiErrs(t *testing.T) {
	err1 := errors.New("error 1")
	err2 := errors.New("error 2")
	errs := []error{err1, err2}
	expectedErrMsg := fmt.Sprintf("%s\n\n%s", err1.Error(), err2.Error())
	err := joinMultiErrs(errs)
	require.Equal(t, expectedErrMsg, err.Error())
}
