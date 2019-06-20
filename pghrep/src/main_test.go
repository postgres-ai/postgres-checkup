package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	//    "fmt"
)

func TestGetFilePathSuccess(t *testing.T) {
	path, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	result := GetFilePath("file://" + path)
	//    fmt.Println("path", path)
	//    fmt.Println("result", result)
	if strings.Compare(path, result) != 0 {
		t.Fatal("TestGetFilePathSuccess failed")
	}
}

func TestGetFilePathFailed(t *testing.T) {
	path := "file:///home/root/golang_test.txt"
	result := GetFilePath(path)
	//    fmt.Println("path", path)
	//    fmt.Println("result", result)
	if strings.Compare(path, result) == 0 {
		t.Fatal("GetFilePatg failed")
	}
}

func TestFileExitsSuccess(t *testing.T) {
	path, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	result := FileExists(path)
	if !result {
		t.Fatal("TestFileExitsSuccess failed")
	}
}
