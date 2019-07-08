package main

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

func TestReorderHosts(t *testing.T) {
	var testJson string = `
	{
		"project": "projx",
		"name": "system_info",
		"checkId": "A001",
		"timestamptz": "2019-06-24 13:22:49.0+0300",
		"dependencies": {
			"null": "null"
		},
		"last_nodes_json": {
			"hosts": {
				"db2": {
					"internal_alias": "node1",
					"index": "1",
					"role": "master",
					"role_change_detected_at": "never"
				},
				"db3": {
					"internal_alias": "node2",
					"index": "2",
					"role": "standby",
					"role_change_detected_at": "never"
				}
			},
			"last_check": {
				"epoch": "1",
				"timestamptz": "2019-06-24 13:22:49.0+0300",
				"dir": "1_2019_06_24T13_20_28_+0300"
			}
		},
		"database": "projx",
		"results": {
			"nodes.json": {
				"hosts": {
					"db2": {
						"internal_alias": "node1",
						"index": "1",
						"role": "master",
						"role_change_detected_at": "never"
					},
					"db3": {
						"internal_alias": "node2",
						"index": "2",
						"role": "standby",
						"role_change_detected_at": "never"
					}
				},
				"last_check": {
					"epoch": "1",
					"timestamptz": "2019-06-24 13:22:49.0+0300",
					"dir": "1_2019_06_24T13_20_28_+0300"
				}
			}
		}
	}`

	data := ParseJson(testJson)
	determineMasterReplica(data)
	err := reorderHosts(data)
	if err == nil {
		t.Fatal("TestReorderHosts: expected error not happened.")
	}
}
