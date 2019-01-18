package main

import (
    "../src/pyraconv"
)

var Data map[string]interface{}

type prepare string

func prepareDropCode(data map[string]interface{}) {
    hosts := pyraconv.ToInterfaceMap(data["hosts"])
    master := pyraconv.ToString(hosts["master"])
    replicas := pyraconv.ToStringArray(hosts["replicas"])
    resultData := make(map[string]interface{})
    results := pyraconv.ToInterfaceMap(data["results"])
    dropCode := make(map[string]string)
    revertCode := make(map[string]string)

    if results[master] != nil {
        masterData := pyraconv.ToInterfaceMap(results[master])
        masterIndexes := pyraconv.ToInterfaceMap(masterData["data"])

        for _, value := range masterIndexes {
            valueData := pyraconv.ToInterfaceMap(value);
            indexName := pyraconv.ToString(valueData["index_name"]);
            dropCode[indexName] = pyraconv.ToString(valueData["drop_code"]);
            revertCode[indexName] = pyraconv.ToString(valueData["revert_code"]);
        }
    }

    for _, replica := range replicas {
        hostData := pyraconv.ToInterfaceMap(results[replica])
        hostIndexes := pyraconv.ToInterfaceMap(hostData["data"])
        for _, value := range hostIndexes {
            valueData := pyraconv.ToInterfaceMap(value);
            indexName := pyraconv.ToString(valueData["index_name"]);
            dropCode[indexName] = pyraconv.ToString(valueData["drop_code"]);
            revertCode[indexName] = pyraconv.ToString(valueData["revert_code"]);
        }
    }

    resultData["drop_code"] = dropCode
    resultData["revert_code"] = revertCode
    data["resultData"] = resultData
}

func (g prepare) Prepare(data map[string]interface{}) map[string]interface{} {
    prepareDropCode(data)
    return data
}

var Preparer prepare