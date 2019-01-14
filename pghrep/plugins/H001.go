package main

import (
    "../src/pyraconv"
    "strings"
)

var Data map[string]interface{}

type prepare string

func compareHostsData(data map[string]interface{}) {
    hosts := pyraconv.ToInterfaceMap(data["hosts"])
    master := pyraconv.ToString(hosts["master"])
    replicas := pyraconv.ToStringArray(hosts["replicas"])
    resultData := make(map[string]interface{})
    
    results := pyraconv.ToInterfaceMap(data["results"])
    masterData := pyraconv.ToInterfaceMap(results[master])
    masterData = pyraconv.ToInterfaceMap(masterData["data"])

    indexesData := pyraconv.ToInterfaceMap(masterData["indexes"])
    indexes := make(map[string]interface{})
    for indexName, value := range indexesData {
        valueData := pyraconv.ToInterfaceMap(value);
        idxScanValue := pyraconv.ToInt64(valueData["idx_scan"])
        idxScanSum := idxScanValue
        
        indexData := make(map[string]interface{})
        //masterIndex := make(map[string]interface{})
        //masterIndex["idx_scan"] = idxScanValue
        indexData["master"] = valueData //masterIndex
        for _, replica := range replicas {
            hostIndex := getReplicaIndex(data, replica, indexName)
            hostIdxScanValue := pyraconv.ToInt64(hostIndex["idx_scan"])
            //hostIndex := make(map[string]interface{})
            //hostIndex["idx_scan"] = hostIdxScanValue
            idxScanSum = idxScanSum + hostIdxScanValue
            indexData[replica] = hostIndex
        }
        indexData["usage"] = idxScanSum > 0
        indexes[indexName] = indexData
    }
    resultData["indexes"] = indexes
    
    dropCode := make(map[string]interface{})
    revertCode := make(map[string]interface{})
    for indexName, value := range indexes {
        indexData := pyraconv.ToInterfaceMap(value)
        usage := pyraconv.ToBool(indexData["usage"])
        if ! usage {
            dropIndexCode := getIndexCode(data, indexName, "drop")
            if len(dropIndexCode) > 0 {
                dropCode[indexName] = dropIndexCode
            }
            revertIndexCode := getIndexCode(data, indexName, "revert")
            if len(revertIndexCode) > 0 {
                revertCode[indexName] = revertIndexCode
            }
        }
    }
    resultData["dropCode"] = dropCode
    resultData["revertCode"] = revertCode
    
    data["resultData"] = resultData
}

func getReplicaIndexUsage(data map[string]interface{}, replica string, indexName string) (int64) {
    results := pyraconv.ToInterfaceMap(data["results"])
    hostData := pyraconv.ToInterfaceMap(results[replica])
    hostData = pyraconv.ToInterfaceMap(hostData["data"])
    indexesData := pyraconv.ToInterfaceMap(hostData["indexes"])
    indexData := pyraconv.ToInterfaceMap(indexesData[indexName])
    return pyraconv.ToInt64(indexData["idx_scan"])
}

func getReplicaIndex(data map[string]interface{}, replica string, indexName string) map[string]interface{} {
    results := pyraconv.ToInterfaceMap(data["results"])
    hostData := pyraconv.ToInterfaceMap(results[replica])
    hostData = pyraconv.ToInterfaceMap(hostData["data"])
    indexesData := pyraconv.ToInterfaceMap(hostData["indexes"])
    indexData := pyraconv.ToInterfaceMap(indexesData[indexName])
    return indexData
}

func getIndexCode(data map[string]interface{}, indexName string, op string) string {
    hosts := pyraconv.ToInterfaceMap(data["hosts"])
    master := pyraconv.ToString(hosts["master"])
    
    results := pyraconv.ToInterfaceMap(data["results"])
    hostData := pyraconv.ToInterfaceMap(results[master])
    hostData = pyraconv.ToInterfaceMap(hostData["data"])
    codeData := pyraconv.ToInterfaceArray(hostData[op + "_code"])
    for _, codeValue := range codeData {
        sql := pyraconv.ToString(codeValue)
        if strings.Contains(sql, indexName) {
            return sql
        }
    }
    return ""
}


func (g prepare) Prepare(data map[string]interface{}) map[string]interface{} {
    compareHostsData(data)
    return data
}

var Preparer prepare