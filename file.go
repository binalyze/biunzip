package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/alexmullins/zip"
)

const defaultBufSize = 1 * 1024 * 1024 // 1MB

func unzipFile(ctx context.Context, filePath string, password string) error {
	dirPath := makeDirPath(filePath)
	err := os.MkdirAll(dirPath, os.ModePerm)
	if err != nil {
		return fmt.Errorf("creating directory failed. dir: %s err: %w", dirPath, err)
	}

	zipReader, err := zip.OpenReader(filePath)
	if err != nil {
		return fmt.Errorf("failed to open zip file. file: %s err: %w", filePath, err)
	}
	defer func() {
		_ = zipReader.Close()
	}()

	var errs []error

	fmt.Printf("unzipping %s: ", filePath)
	defer func() {
		if errs != nil {
			fmt.Println("failed")
			return
		}
		fmt.Println("ok")
	}()

	for _, zipEntry := range zipReader.File {
		err = ctx.Err()
		if err != nil {
			errs = append(
				errs,
				fmt.Errorf("context error. err: %w", err),
			)
			break
		}

		if zipEntry.Mode().IsDir() {
			subdirPath := filepath.Join(dirPath, zipEntry.Name)
			err = os.MkdirAll(subdirPath, zipEntry.Mode())
			if err != nil {
				errs = append(
					errs,
					fmt.Errorf("failed to create subdirectory. zip file: %s subdirectory: %s err: %w", filePath, subdirPath, err),
				)
			}
			continue
		}

		if len(password) > 0 {
			zipEntry.SetPassword(password)
		}

		zipEntryReader, err := zipEntry.Open()
		if err != nil {
			errs = append(
				errs,
				fmt.Errorf("failed to open zip entry. zip file: %s zip entry: %s err: %w", filePath, zipEntry.Name, err),
			)
			continue
		}
		ctxZipEntryReader := newContextReader(ctx, zipEntryReader)
		srcReader := bufio.NewReaderSize(ctxZipEntryReader, defaultBufSize)

		dstFilePath := filepath.Join(dirPath, zipEntry.Name)
		dstFile, err := os.OpenFile(dstFilePath, os.O_TRUNC|os.O_CREATE|os.O_WRONLY, zipEntry.Mode())
		if err != nil {
			errs = append(
				errs,
				fmt.Errorf("failed to create destination file. zip file: %s zip entry: %s destination file: %s err: %w", filePath, zipEntry.Name, dstFilePath, err),
			)
			continue
		}
		dstWriter := bufio.NewWriterSize(dstFile, defaultBufSize)

		_, err = io.Copy(dstWriter, srcReader)
		if err != nil {
			errs = append(
				errs,
				fmt.Errorf("failed to copy source file to destination file. zip file: %s zip entry: %s destination file: %s err: %w", filePath, zipEntry.Name, dstFilePath, err),
			)
			continue
		}

		err = dstFile.Close()
		if err != nil {
			errs = append(
				errs,
				fmt.Errorf("failed to close destination file. zip file: %s zip entry: %s destination file: %s err: %w", filePath, zipEntry.Name, dstFilePath, err),
			)
			continue
		}

		_ = zipEntryReader.Close()
	}
	if len(errs) > 0 {
		return combineErrs("errors for "+filePath+":", errs, true)
	}
	return nil
}

func makeDirPath(filePath string) string {
	ext := filepath.Ext(filePath)
	dirPath := filePath[:len(filePath)-len(ext)]
	return dirPath
}
