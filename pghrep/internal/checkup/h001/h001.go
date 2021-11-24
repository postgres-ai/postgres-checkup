package h001

import (
	"gitlab.com/postgres-ai/postgres-checkup/pghrep/internal/pyraconv"
)

var Data map[string]interface{}

func prepareDropCode(data map[string]interface{}) {
	hosts := pyraconv.ToInterfaceMap(data["hosts"])
	master := pyraconv.ToString(hosts["master"])
	replicas := pyraconv.ToStringArray(hosts["replicas"])
	resultData := make(map[string]interface{})
	results := pyraconv.ToInterfaceMap(data["results"])

	doUndoCode := make(map[string]interface{})

	if results[master] != nil {
		masterData := pyraconv.ToInterfaceMap(results[master])
		masterIndexes := pyraconv.ToInterfaceMap(masterData["data"])

		for _, value := range masterIndexes {
			valueData := pyraconv.ToInterfaceMap(value)
			indexName := pyraconv.ToString(valueData["index_name"])
			doUndo := make(map[string]interface{})
			doUndo["drop_code"] = pyraconv.ToString(valueData["drop_code"])
			doUndo["revert_code"] = pyraconv.ToString(valueData["revert_code"])
			if len(pyraconv.ToString(valueData["drop_code"])) > 0 && len(pyraconv.ToString(valueData["revert_code"])) > 0 {
				doUndoCode[indexName] = doUndo
			}
		}
	}

	for _, replica := range replicas {
		hostData := pyraconv.ToInterfaceMap(results[replica])
		hostIndexes := pyraconv.ToInterfaceMap(hostData["data"])
		for _, value := range hostIndexes {
			valueData := pyraconv.ToInterfaceMap(value)
			indexName := pyraconv.ToString(valueData["index_name"])
			doUndo := make(map[string]interface{})
			doUndo["drop_code"] = pyraconv.ToString(valueData["drop_code"])
			doUndo["revert_code"] = pyraconv.ToString(valueData["revert_code"])
			if len(pyraconv.ToString(valueData["drop_code"])) > 0 && len(pyraconv.ToString(valueData["revert_code"])) > 0 {
				doUndoCode[indexName] = doUndo
			}
		}
	}

	resultData["repair_code"] = doUndoCode
	data["resultData"] = resultData
}

func H001PreprocessReportData(data map[string]interface{}) {
	prepareDropCode(data)
}
