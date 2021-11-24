package reportutils

import (
	"testing"
)

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
	DetermineMasterReplica(data)
	err := ReorderHosts(data)
	if err == nil {
		t.Fatal("TestReorderHosts: expected error not happened.")
	}
}
