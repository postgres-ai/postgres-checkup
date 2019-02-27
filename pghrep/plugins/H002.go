package main

import (
    "../src/pyraconv"
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

func (g prepare) Prepare(data map[string]interface{}) map[string]interface{} {
    preCheckHosts(data)
    return data
}

var Preparer prepare