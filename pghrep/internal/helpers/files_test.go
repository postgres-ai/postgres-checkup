package helpers

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestGetFilePathSuccess(t *testing.T) {
	path, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	result := GetFilePath("file://" + path)
	if strings.Compare(path, result) != 0 {
		t.Fatal("TestGetFilePathSuccess: received path not equeal to expected path.")
	}
}

func TestGetFilePathFailed(t *testing.T) {
	path := "file:///home/root/golang_test.txt"
	result := GetFilePath(path)
	if strings.Compare(path, result) == 0 {
		t.Fatal("GetFilePatg: received path not equeal to expected path.")
	}
}

func TestFileExitsSuccess(t *testing.T) {
	path, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	result := FileExists(path)
	if !result {
		t.Fatal("TestFileExitsSuccess: existing file not found.")
	}
}
