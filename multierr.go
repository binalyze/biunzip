package main

import (
	"errors"
	"strings"
)

func makeMultiErr(msg string, errs []error) error {
	var builder strings.Builder
	builder.WriteString(msg)
	builder.WriteString(":")
	for _, err := range errs {
		builder.WriteString("\n- ")
		builder.WriteString(err.Error())
	}
	return errors.New(builder.String())
}

func joinMultiErrs(errs []error) error {
	var builder strings.Builder
	for i, err := range errs {
		builder.WriteString(err.Error())
		if i < len(errs)-1 {
			builder.WriteString("\n\n")
		}
	}
	return errors.New(builder.String())
}
