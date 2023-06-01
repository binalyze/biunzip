package main

import (
	"fmt"
	"strings"
)

func combineErrs(msg string, errs []error, printDashForErrs bool) error {
	builder := &strings.Builder{}
	fmt.Fprintln(builder, msg)
	for i, err := range errs {
		if i != 0 {
			fmt.Fprintln(builder, "")
		}
		if printDashForErrs {
			fmt.Fprint(builder, "- ")
		}
		fmt.Fprint(builder, err.Error())
	}
	return fmt.Errorf(builder.String())
}
