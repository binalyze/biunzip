package main

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type zipFile struct {
	path     string
	password string
}

const (
	zipFileExt         = ".zip"
	zipFilenameColName = "File Name"
	passwordColName    = "Zip Password"
)

func unzipDir(ctx context.Context, dirPath string, csvFilePath string) error {
	csvFileLines, err := readCSVFile(csvFilePath)
	if err != nil {
		return fmt.Errorf("failed to read csv file. err: %w", err)
	}

	err = validateCSVFile(csvFileLines)
	if err != nil {
		return err
	}

	zipFiles := parseZipFiles(dirPath, csvFileLines)

	err = validateZipFiles(zipFiles)
	if err != nil {
		return err
	}

	return unzipFiles(ctx, zipFiles)
}

func readCSVFile(csvFilePath string) ([]string, error) {
	file, err := os.Open(csvFilePath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, fmt.Errorf("csv file doesn't exist")
		}
		return nil, err
	}
	defer func() {
		_ = file.Close()
	}()
	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		lines = append(lines, line)
	}
	return lines, nil
}

func validateCSVFile(lines []string) error {
	lineCount := len(lines)
	if lineCount < 2 {
		return errors.New("csv file doesn't have any data")
	}

	err := validateColCounts(lines)
	if err != nil {
		return err
	}

	err = checkDuplicateZipFilenames(lines)
	if err != nil {
		return err
	}

	return nil
}

func validateColCounts(lines []string) error {
	header := lines[0]
	headerCols := strings.Split(header, ",")
	headerColCount := len(headerCols)
	if headerColCount < 2 {
		return fmt.Errorf("unexpected header column count found in csv file. count: %d expected count: 2", headerColCount)
	}

	var errs []error
	for i, line := range lines {
		cols := strings.Split(line, ",")
		colCount := len(cols)
		if colCount != headerColCount {
			errs = append(errs, fmt.Errorf("unexpected column count found on line %d. count: %d expected count: %d", i+1, colCount, headerColCount))
		}
	}
	if len(errs) > 0 {
		return combineErrs("unexpectd column counts found in csv file", errs, true)
	}
	return nil
}

func checkDuplicateZipFilenames(lines []string) error {
	header := lines[0]
	zipFilenameColIndex := findZipFilenameColIndex(header)
	var errs []error
	for i, line := range lines {
		cols := strings.Split(line, ",")
		zipFilename := cols[zipFilenameColIndex]
		for j := i + 1; j < len(lines); j++ {
			existingLine := lines[j]
			existingCols := strings.Split(existingLine, ",")
			existingZipFilename := existingCols[zipFilenameColIndex]
			if zipFilename == existingZipFilename {
				errs = append(errs, fmt.Errorf("duplicate zip filename %s found on lines %d and %d", zipFilename, i+1, j+1))
			}
		}
	}
	if len(errs) > 0 {
		return combineErrs("duplicate zip filenames found in csv file", errs, true)
	}
	return nil
}

func parseZipFiles(dirPath string, lines []string) []zipFile {
	header := lines[0]
	zipFilenameColIndex := findZipFilenameColIndex(header)
	passwordColIndex := findPasswordColIndex(header)
	var zipFiles []zipFile
	for i, line := range lines {
		if i == 0 {
			continue
		}
		cols := strings.Split(line, ",")
		zipFilename := cols[zipFilenameColIndex]
		zipFilePath := filepath.Join(dirPath, zipFilename)
		password := cols[passwordColIndex]
		zipFile := zipFile{
			path:     zipFilePath,
			password: password,
		}
		zipFiles = append(zipFiles, zipFile)
	}
	return zipFiles
}

func findZipFilenameColIndex(header string) int {
	cols := strings.Split(header, ",")
	for i, col := range cols {
		if col == zipFilenameColName {
			return i
		}
	}
	return 0
}

func findPasswordColIndex(header string) int {
	cols := strings.Split(header, ",")
	for i, col := range cols {
		if col == passwordColName {
			return i
		}
	}
	return len(cols) - 1
}

func validateZipFiles(zipFiles []zipFile) error {
	var errs []error
	for i, zipFile := range zipFiles {
		err := validateZipFile(zipFile)
		if err != nil {
			errs = append(
				errs,
				fmt.Errorf("failed to validate zip file. zip file: %s line: %d err: %w", zipFile.path, i+2, err),
			)
		}
	}
	if len(errs) > 0 {
		return combineErrs("failed to validate zip files in csv file", errs, true)
	}
	return nil
}

func validateZipFile(zipFile zipFile) error {
	fileInfo, err := os.Stat(zipFile.path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return errors.New("zip file doesn't exist")
		}
		return err
	}
	if !fileInfo.Mode().IsRegular() {
		return errors.New("zip file is not regular file")
	}
	return nil
}

func unzipFiles(ctx context.Context, zipFiles []zipFile) error {
	var errs []error
	for _, zipFile := range zipFiles {
		err := unzipFile(ctx, zipFile.path, zipFile.password)
		if err != nil {
			errs = append(errs, err)
		}
	}
	if len(errs) > 0 {
		return combineErrs("failed to unzip files", errs, false)
	}
	return nil
}
