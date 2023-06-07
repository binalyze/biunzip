package main

import (
	"context"
	"encoding/csv"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type zipFile struct {
	path     string
	password string
}

const (
	filenameColName = "File Name"
	passwordColName = "Zip Password"
)

func unzipDir(ctx context.Context, dirPath string, csvFilePath string, maxConcurrency int) error {
	lines, err := readCSVFile(csvFilePath)
	if err != nil {
		return err
	}

	err = validateCSVFile(dirPath, lines)
	if err != nil {
		return err
	}

	files := parseZipFiles(dirPath, lines)

	err = validateZipFiles(files)
	if err != nil {
		return err
	}

	return unzipFiles(ctx, files, maxConcurrency)
}

func readCSVFile(csvFilePath string) ([][]string, error) {
	file, err := os.Open(csvFilePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open csv file: %w", err)
	}
	defer file.Close()
	reader := csv.NewReader(file)
	lines, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("failed to read csv file: %w", err)
	}
	return lines, nil
}

func validateCSVFile(dirPath string, lines [][]string) error {
	err := validateLineCount(lines)
	if err != nil {
		return err
	}

	err = validateHeaderColCount(lines)
	if err != nil {
		return err
	}

	err = checkEmptyFilenames(lines)
	if err != nil {
		return err
	}

	err = checkDuplicateFilenames(lines)
	if err != nil {
		return err
	}

	return nil
}

func validateLineCount(lines [][]string) error {
	if len(lines) < 2 {
		return fmt.Errorf("csv file doesn't have any data")
	}
	return nil
}

func validateHeaderColCount(lines [][]string) error {
	header := lines[0]
	if len(header) < 2 {
		return errors.New("unexpected column count found on header line")
	}
	return nil
}

func checkEmptyFilenames(lines [][]string) error {
	header := lines[0]
	filenameColIndex := findFilenameColIndex(header)
	var invalidLineNums []int
	for i, line := range lines {
		filename := line[filenameColIndex]
		if len(filename) == 0 {
			invalidLineNums = append(invalidLineNums, i+1)
		}
	}
	if len(invalidLineNums) > 0 {
		return fmt.Errorf("empty filename found on line(s) %s", joinLineNums(invalidLineNums))
	}
	return nil
}

func checkDuplicateFilenames(lines [][]string) error {
	header := lines[0]
	filenameColIndex := findFilenameColIndex(header)
	var duplicateLineNums []int
	for i := 0; i < len(lines); i++ {
		lineNum := i + 1
		if lineNumExists(lineNum, duplicateLineNums) {
			continue
		}
		filename := lines[i][filenameColIndex]
		for j := i + 1; j < len(lines); j++ {
			existingLineNum := j + 1
			if lineNumExists(existingLineNum, duplicateLineNums) {
				continue
			}
			existingFilaname := lines[j][filenameColIndex]
			if filename == existingFilaname {
				if !lineNumExists(lineNum, duplicateLineNums) {
					duplicateLineNums = append(duplicateLineNums, lineNum)
				}
				duplicateLineNums = append(duplicateLineNums, existingLineNum)
			}
		}
	}
	if len(duplicateLineNums) > 0 {
		return fmt.Errorf("duplicate filename found on line(s) %s", joinLineNums(duplicateLineNums))
	}
	return nil
}

func parseZipFiles(dirPath string, lines [][]string) []zipFile {
	header := lines[0]
	filenameColIndex := findFilenameColIndex(header)
	passwordColIndex := findPasswordColIndex(header)
	var files []zipFile
	for i := 1; i < len(lines); i++ {
		line := lines[i]
		filename := line[filenameColIndex]
		filePath := filepath.Join(dirPath, filename)
		password := line[passwordColIndex]
		file := zipFile{
			path:     filePath,
			password: password,
		}
		files = append(files, file)
	}
	return files
}

func validateZipFiles(files []zipFile) error {
	var errs []error
	for _, file := range files {
		fileInfo, err := os.Stat(file.path)
		if err != nil {
			errs = append(errs, fmt.Errorf("failed to get file info for '%s': %w", file.path, err))
			continue
		}
		if !fileInfo.Mode().IsRegular() {
			errs = append(errs, fmt.Errorf("'%s' is not a regular file", file.path))
		}
	}
	if len(errs) > 0 {
		return makeMultiErr("failed to validate zip files in csv file", errs)
	}
	return nil
}

func unzipFiles(ctx context.Context, files []zipFile, maxConcurrency int) error {
	sem := make(chan struct{}, maxConcurrency)
	var errs []error
	for _, file := range files {
		sem <- struct{}{}
		go func(filePath string, password string) {
			err := unzipFile(ctx, filePath, password)
			if err != nil {
				errs = append(errs, err)
			}
			<-sem
		}(file.path, file.password)
	}
	for i := 0; i < cap(sem); i++ {
		sem <- struct{}{}
	}
	close(sem)
	if len(errs) > 0 {
		return joinMultiErrs(errs)
	}
	return nil
}

func findFilenameColIndex(header []string) int {
	for i, col := range header {
		if col == filenameColName {
			return i
		}
	}
	return 0
}

func findPasswordColIndex(header []string) int {
	for i, col := range header {
		if col == passwordColName {
			return i
		}
	}
	return len(header) - 1
}

func joinLineNums(lineNums []int) string {
	var builder strings.Builder
	for i, lineNum := range lineNums {
		builder.WriteString(strconv.Itoa(lineNum))
		if i < len(lineNums)-1 {
			builder.WriteString(", ")
		}
	}
	return builder.String()
}

func lineNumExists(lineNum int, lineNums []int) bool {
	for _, existingLineNum := range lineNums {
		if lineNum == existingLineNum {
			return true
		}
	}
	return false
}
