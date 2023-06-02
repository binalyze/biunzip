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

const defaultBufSize = 10 * 1024 * 1024 // 10MB

func unzipFile(ctx context.Context, filePath string, password string) error {
	dirPath := makeDirPath(filePath)
	err := os.MkdirAll(dirPath, 0755) // 0755: rwxr-xr-x
	if err != nil {
		return fmt.Errorf("failed to create dir '%s': %w", dirPath, err)
	}

	zipReader, err := zip.OpenReader(filePath)
	if err != nil {
		return fmt.Errorf("failed to open file '%s': %w", filePath, err)
	}
	defer zipReader.Close()

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
			errs = append(errs, fmt.Errorf("context error: %w", err))
			break
		}

		if zipEntry.Mode().IsDir() {
			subdirPath := filepath.Join(dirPath, zipEntry.Name)
			err = os.MkdirAll(subdirPath, zipEntry.Mode())
			if err != nil {
				errs = append(errs, fmt.Errorf("failed to create subdir '%s': %w", subdirPath, err))
			}
			continue
		}

		if len(password) > 0 {
			zipEntry.SetPassword(password)
		}

		zipEntryReader, err := zipEntry.Open()
		if err != nil {
			errs = append(errs, fmt.Errorf("failed to open zip entry '%s': %w", zipEntry.Name, err))
			continue
		}
		ctxZipEntryReader := newContextReader(ctx, zipEntryReader)
		srcReader := bufio.NewReaderSize(ctxZipEntryReader, defaultBufSize)

		dstFilePath := filepath.Join(dirPath, zipEntry.Name)
		dstFile, err := os.OpenFile(dstFilePath, os.O_TRUNC|os.O_CREATE|os.O_WRONLY, zipEntry.Mode())
		if err != nil {
			errs = append(errs, fmt.Errorf("failed to create dst file '%s': %w", dstFilePath, err))
			continue
		}
		dstWriter := bufio.NewWriterSize(dstFile, defaultBufSize)

		_, err = io.Copy(dstWriter, srcReader)
		if err == nil {
			err = dstWriter.Flush()
		}
		if err != nil {
			_ = dstFile.Close()
			errs = append(errs, fmt.Errorf("failed to copy src file '%s' to dst file '%s': %w", zipEntry.Name, dstFilePath, err))
			continue
		}

		err = dstFile.Close()
		if err != nil {
			errs = append(errs, fmt.Errorf("failed to close destination file '%s': %w", dstFilePath, err))
			continue
		}

		_ = zipEntryReader.Close()
	}
	if len(errs) > 0 {
		msg := fmt.Sprintf("failed to unzip file '%s'", filePath)
		return makeMultiErr(msg, errs)
	}
	return nil
}

func makeDirPath(filePath string) string {
	ext := filepath.Ext(filePath)
	dirPath := filePath[:len(filePath)-len(ext)]
	return dirPath
}
