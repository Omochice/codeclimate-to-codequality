package main

import (
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/Omochice/codeclimate-to-codequality/cli"
	"github.com/Omochice/codeclimate-to-codequality/codeclimate"
	"github.com/Omochice/codeclimate-to-codequality/codequality"
	"github.com/Omochice/codeclimate-to-codequality/converter"
)

var version = "develop"

func handleError(w io.Writer, err error) int {
	fmt.Fprintf(w, "Error: %v\n", err)
	return 1
}

func command(args []string, inout *cli.ProcInout) int {
	opts, err := cli.Parse(args)
	if err != nil {
		var helpErr *cli.HelpError
		if errors.As(err, &helpErr) {
			inout.Stderr.Write([]byte(helpErr.Help))
			return 0
		} else {
			return handleError(inout.Stderr, err)
		}
	}

	if opts.Version {
		inout.Stdout.Write([]byte(version))
		return 0
	}

	var reader io.Reader
	if opts.Source == "-" {
		reader = inout.Stdin
	} else {
		f, err := os.Open(opts.Source)
		if err != nil {
			return handleError(inout.Stderr, err)
		}
		defer f.Close()
		reader = f
	}

	issues, err := codeclimate.Parse(reader)
	if err != nil {
		return handleError(inout.Stderr, err)
	}

	violations := converter.Issues(issues)

	if err := codequality.Write(violations, inout.Stdout); err != nil {
		return handleError(inout.Stderr, err)
	}

	return 0
}

func main() {
	cli.Run(command)
}
