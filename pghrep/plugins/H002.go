package main

import (
    "../src/pyraconv"
    "strings"
)

var Data map[string]interface{}

type prepare string

/*
  Replace master on replica#1 if master not defined
*/
func preCheckHosts(data map[string]interface{}) {
    hosts := pyraconv.ToInterfaceMap(data["hosts"])
    masterHost := pyraconv.ToString(hosts["master"])
    replicaHosts := pyraconv.ToStringArray(hosts["replicas"])
    var allHosts []string
    if hosts["master"] != nil {
        allHosts = append(allHosts, masterHost)
    }
    for _, replicaHost := range replicaHosts {
        allHosts = append(allHosts, replicaHost)
    }
    if len(allHosts) == 0 {
        return
    }
    master := allHosts[0]
    var replicas []string
    replicas = append(replicas, allHosts[1:]...)
    hosts["master"] = master
    hosts["replicas"] = replicas
}

func compareHostsData(data map[string]interface{}) {
    preCheckHosts(data)
    hosts := pyraconv.ToInterfaceMap(data["hosts"])
    master := pyraconv.ToString(hosts["master"])
    replicas := pyraconv.ToStringArray(hosts["replicas"])
    resultData := make(map[string]interface{})

    results := pyraconv.ToInterfaceMap(data["results"])
    masterData := pyraconv.ToInterfaceMap(results[master])
    masterData = pyraconv.ToInterfaceMap(masterData["data"])

    allUnusedIndexes := make(map[string]bool)
    uIndexesData := pyraconv.ToInterfaceMap(masterData["unused_indexes"])
    uIndexes := make(map[string]interface{})
    for indexName, value := range uIndexesData {
        valueData := pyraconv.ToInterfaceMap(value);
        idxScanValue := pyraconv.ToInt64(valueData["idx_scan"])
        idxScanSum := idxScanValue
        
        indexData := make(map[string]interface{})
        indexData["master"] = valueData
        for _, replica := range replicas {
            hostIndexData := getReplicaIndex(data, replica, "unused_indexes", indexName)
            hostIdxScanValue := pyraconv.ToInt64(hostIndexData["idx_scan"])
            idxScanSum = idxScanSum + hostIdxScanValue
            indexData[replica] = hostIndexData
        }
        indexData["usage"] = idxScanSum > 0
        if idxScanSum == 0 {
            allUnusedIndexes[indexName] = true
        }
        uIndexes[indexName] = indexData
    }
    resultData["unused_indexes"] = uIndexes

    rIndexesData := pyraconv.ToInterfaceMap(masterData["redundant_indexes"])
    rIndexes := make(map[string]interface{})
    for indexName, value := range rIndexesData {
        valueData := pyraconv.ToInterfaceMap(value);
        idxScanValue := pyraconv.ToInt64(valueData["index_usage"])
        idxScanSum := idxScanValue

        indexData := make(map[string]interface{})
        indexData["master"] = valueData
        for _, replica := range replicas {
            hostIndexData := getReplicaIndex(data, replica, "redundant_indexes", indexName)
            hostIdxScanValue := pyraconv.ToInt64(hostIndexData["index_usage"])
            idxScanSum = idxScanSum + hostIdxScanValue
            indexData[replica] = hostIndexData
        }
        indexData["usage"] = idxScanSum > 0
        if idxScanSum == 0 {
            allUnusedIndexes[indexName] = true
        }
        rIndexes[indexName] = indexData
    }
    resultData["redundant_indexes"] = rIndexes

    dropCode := make(map[string]interface{})
    revertCode := make(map[string]interface{})
    // enum only not used indexes
    for indexName, _ := range allUnusedIndexes {
        dropIndexCode := getIndexCode(data, indexName, "drop")
        if len(dropIndexCode) > 0 {
            dropCode[indexName] = dropIndexCode
        }
        revertIndexCode := getIndexCode(data, indexName, "revert")
        if len(revertIndexCode) > 0 {
            revertCode[indexName] = revertIndexCode
        }
    }

    resultData["drop_code"] = dropCode
    resultData["revert_code"] = revertCode

    data["resultData"] = resultData
}

func getReplicaIndexUsage(data map[string]interface{}, replica string, indexName string) (int64) {
    results := pyraconv.ToInterfaceMap(data["results"])
    hostData := pyraconv.ToInterfaceMap(results[replica])
    hostData = pyraconv.ToInterfaceMap(hostData["data"])
    uIndexesData := pyraconv.ToInterfaceMap(hostData["unused_uIndexes"])
    indexData := pyraconv.ToInterfaceMap(uIndexesData[indexName])
    return pyraconv.ToInt64(indexData["idx_scan"])
}

func getReplicaIndex(data map[string]interface{}, replica string, indexType string, indexName string) map[string]interface{} {
    results := pyraconv.ToInterfaceMap(data["results"])
    hostData := pyraconv.ToInterfaceMap(results[replica])
    hostData = pyraconv.ToInterfaceMap(hostData["data"])
    uIndexesData := pyraconv.ToInterfaceMap(hostData[indexType])
    indexData := pyraconv.ToInterfaceMap(uIndexesData[indexName])
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