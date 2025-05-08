//go:debug zipinsecurepath=0

package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/urfave/cli/v2"
)

const (
	dirFlagUsage = "dir path to unzip"
	csvFlagUsage = "path for the csv file containing a list of zip files names and passwords to unzip. use this flag with the dir flag."

	fileFlagUsage     = "path for the file to unzip"
	passwordFlagUsage = "password for the zip file. use this flag with the file flag if the input file is encrypted."
)

var (
	errEmptyCSVFilePath = errors.New("please provide the csv file path along with the directory path to unzip files in the directory")
	errUnexpectedFlag   = errors.New("please provide both the directory and csv file paths to unzip files in the directory, or provide a file path to unzip a single file. if the file is encrypted, include the password")
)

func main() {
	app := cli.App{
		Name:  "biunzip",
		Usage: "unzip zip files",
		Flags: []cli.Flag{
			&cli.PathFlag{
				Name:    "dir",
				Aliases: []string{"d"},
				Usage:   dirFlagUsage,
			},
			&cli.PathFlag{
				Name:    "csv",
				Aliases: []string{"c"},
				Usage:   csvFlagUsage,
			},
			&cli.PathFlag{
				Name:    "file",
				Aliases: []string{"f"},
				Usage:   fileFlagUsage,
			},
			&cli.StringFlag{
				Name:    "password",
				Aliases: []string{"p"},
				Usage:   passwordFlagUsage,
			},
		},
		Action: run,
	}

	ctx, cancelFunc := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancelFunc()

	err := app.RunContext(ctx, os.Args)
	if err != nil {
		fmt.Println()
		fmt.Println("errors:")
		fmt.Println(err.Error())
		fmt.Println()
		fmt.Println("exit status 1")
		os.Exit(1)
	}
}

func run(ctx *cli.Context) error {
	dirPath := ctx.Path("dir")
	if len(dirPath) > 0 {
		csvFilePath := ctx.Path("csv")
		if len(csvFilePath) == 0 {
			return errEmptyCSVFilePath
		}
		return unzipDir(ctx.Context, dirPath, csvFilePath)
	}
	filePath := ctx.Path("file")
	if len(filePath) > 0 {
		password := ctx.String("password")
		return unzipFile(ctx.Context, filePath, password)
	}
	return errUnexpectedFlag
}
