package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReadCSVFile(t *testing.T) {
	validLines := [][]string{
		{filenameColName, passwordColName},
		{"file_1.zip", "password_1"},
		{"file_2.zip", "password_2"},
	}
	validCSVFilePath, err := createTempFile("", "test*.csv", []byte("File Name,Zip Password\nfile_1.zip,password_1\nfile_2.zip,password_2\n"))
	require.NoError(t, err)
	defer os.Remove(validCSVFilePath)

	invalidCSVFilePath, err := createTempFile("", "test*.csv", []byte("File Name,Zip Password\nfile_1.zip,password_1\nfile_2.zip\n"))
	require.NoError(t, err)
	defer os.Remove(invalidCSVFilePath)

	tests := []struct {
		name      string
		path      string
		expected  [][]string
		expectErr bool
	}{
		{
			name:      "with a valid csv file",
			path:      validCSVFilePath,
			expected:  validLines,
			expectErr: false,
		},
		{
			name:      "with a non-existent csv file",
			path:      "non-existing_file.csv",
			expected:  nil,
			expectErr: true,
		},
		{
			name:      "with an invalid csv file",
			path:      invalidCSVFilePath,
			expected:  nil,
			expectErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual, err := readCSVFile(tt.path)
			if tt.expectErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.Equal(t, tt.expected, actual)
		})
	}
}

func TestValidateCSVFile(t *testing.T) {
	tests := []struct {
		name      string
		lines     [][]string
		expectErr bool
	}{
		{
			name: "with a valid csv file",
			lines: [][]string{
				{filenameColName, passwordColName},
				{"file_1.zip", "password_1"},
			},
			expectErr: false,
		},
		{
			name: "with an invalid line count",
			lines: [][]string{
				{filenameColName, passwordColName},
			},
			expectErr: true,
		},
		{
			name: "with an invalid header col count",
			lines: [][]string{
				{filenameColName},
				{"file_1.zip"},
			},
			expectErr: true,
		},
		{
			name: "with empty filenames",
			lines: [][]string{
				{filenameColName, passwordColName},
				{"", "password_1"},
			},
			expectErr: true,
		},
		{
			name: "with duplicate filenames",
			lines: [][]string{
				{filenameColName, passwordColName},
				{"file_1.zip", "password_1"},
				{"file_1.zip", "password_2"},
			},
			expectErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateCSVFile(tt.lines)
			if tt.expectErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
		})
	}
}

func TestValidateLineCount(t *testing.T) {
	tests := []struct {
		name      string
		lines     [][]string
		expectErr bool
	}{
		{
			name: "with a valid line count",
			lines: [][]string{
				{filenameColName, passwordColName},
				{"file_1.zip", "password_1"},
			},
			expectErr: false,
		},
		{
			name: "with an invalid line count",
			lines: [][]string{
				{filenameColName, passwordColName},
			},
			expectErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateLineCount(tt.lines)
			if tt.expectErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
		})
	}
}

func TestValidateHeaderColCount(t *testing.T) {
	tests := []struct {
		name      string
		lines     [][]string
		expectErr bool
	}{
		{
			name: "with a valid header col count",
			lines: [][]string{
				{filenameColName, passwordColName},
			},
			expectErr: false,
		},
		{
			name: "with an invalid header col count",
			lines: [][]string{
				{filenameColName},
			},
			expectErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateHeaderColCount(tt.lines)
			if tt.expectErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
		})
	}
}

func TestCheckEmptyFilenames(t *testing.T) {
	tests := []struct {
		name                  string
		lines                 [][]string
		emptyFilenameLineNums []int
	}{
		{
			name: "without empty filenames",
			lines: [][]string{
				{filenameColName, passwordColName},
				{"file_1.zip", "password_1"},
				{"file_2.zip", "password_2"},
			},
			emptyFilenameLineNums: []int{},
		},
		{
			name: "with empty filenames",
			lines: [][]string{
				{filenameColName, passwordColName},
				{"", "password_1"},
				{"file_2.zip", "password_2"},
				{"", "password_3"},
			},
			emptyFilenameLineNums: []int{2, 4},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := checkEmptyFilenames(tt.lines)
			if len(tt.emptyFilenameLineNums) > 0 {
				require.Error(t, err)
				errMsgContainsEmptyFilenameLineNums := strings.Contains(err.Error(), joinLineNums(tt.emptyFilenameLineNums))
				require.True(t, errMsgContainsEmptyFilenameLineNums)
				return
			}
			require.NoError(t, err)
		})
	}
}

func TestCheckDuplicateFilenames(t *testing.T) {
	tests := []struct {
		name              string
		lines             [][]string
		duplicateLineNums []int
	}{
		{
			name: "without duplicate filenames",
			lines: [][]string{
				{filenameColName, passwordColName},
				{"file_1.zip", "password_1"},
				{"file_2.zip", "password_2"},
			},
			duplicateLineNums: []int{},
		},
		{
			name: "with duplicate filenames",
			lines: [][]string{
				{filenameColName, passwordColName},
				{"file_1.zip", "password_1"},
				{"file_2.zip", "password_2"},
				{"file_1.zip", "password_3"},
			},
			duplicateLineNums: []int{2, 4},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := checkDuplicateFilenames(tt.lines)
			if len(tt.duplicateLineNums) > 0 {
				require.Error(t, err)
				errMsgContainsDuplicateLineNums := strings.Contains(err.Error(), joinLineNums(tt.duplicateLineNums))
				require.True(t, errMsgContainsDuplicateLineNums)
				return
			}
			require.NoError(t, err)
		})
	}
}

func TestParseZipFiles(t *testing.T) {
	dirPath := os.TempDir()
	filename1 := "file_1.zip"
	filePath1 := filepath.Join(dirPath, filename1)
	password1 := "password_1"
	filename2 := "file_2.zip"
	filePath2 := filepath.Join(dirPath, filename2)
	password2 := "password_2"
	lines := [][]string{
		{filenameColName, passwordColName},
		{filename1, password1},
		{filename2, password2},
	}
	expected := []zipFile{
		{
			path:     filePath1,
			password: password1,
		},
		{
			path:     filePath2,
			password: password2,
		},
	}
	actual := parseZipFiles(dirPath, lines)
	require.Equal(t, expected, actual)
}

func TestValidateZipFiles(t *testing.T) {
	filePath, err := createTempFile("", "file_1.zip", nil)
	require.NoError(t, err)
	defer os.Remove(filePath)

	tests := []struct {
		name      string
		files     []zipFile
		expectErr bool
	}{
		{
			name: "with valid files",
			files: []zipFile{
				{
					path: filePath,
				},
			},
			expectErr: false,
		},
		{
			name: "with a non-existing file",
			files: []zipFile{
				{
					path: "non-existing_file.zip",
				},
			},
			expectErr: true,
		},
		{
			name: "with an irregular file (dir)",
			files: []zipFile{
				{
					path: os.TempDir(),
				},
			},
			expectErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateZipFiles(tt.files)
			errExists := err != nil
			require.Equal(t, tt.expectErr, errExists)
		})
	}
}

func TestFindFilenameColIndex(t *testing.T) {
	tests := []struct {
		name     string
		header   []string
		expected int
	}{
		{
			name:     "with the filename column",
			header:   []string{"column 1", filenameColName, "column 2"},
			expected: 1,
		},
		{
			name:     "without the filename column",
			header:   []string{"column 1", "column 2"},
			expected: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := findFilenameColIndex(tt.header)
			require.Equal(t, tt.expected, actual)
		})
	}
}

func TestFindPasswordColIndex(t *testing.T) {
	tests := []struct {
		name     string
		header   []string
		expected int
	}{
		{
			name:     "with the password column",
			header:   []string{"column 1", passwordColName, "column 2"},
			expected: 1,
		},
		{
			name:     "without the password column",
			header:   []string{"column 1", "column 2"},
			expected: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := findPasswordColIndex(tt.header)
			require.Equal(t, tt.expected, actual)
		})
	}
}

func TestJoinLineNums(t *testing.T) {
	lineNums := []int{1, 2}
	expected := "1, 2"
	actual := joinLineNums(lineNums)
	require.Equal(t, expected, actual)
}

func TestLineNumExists(t *testing.T) {
	lineNums := []int{1, 2}
	tests := []struct {
		name     string
		lineNum  int
		expected bool
	}{
		{
			name:     "with an existing line number",
			lineNum:  1,
			expected: true,
		},
		{
			name:     "with a non-existing line number",
			lineNum:  3,
			expected: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := lineNumExists(tt.lineNum, lineNums)
			require.Equal(t, tt.expected, actual)
		})
	}
}

func createTempFile(dir string, filename string, content []byte) (string, error) {
	file, err := os.CreateTemp(dir, filename)
	if err != nil {
		return "", err
	}
	_, err = file.Write(content)
	if err != nil {
		return "", err
	}
	return file.Name(), nil
}
