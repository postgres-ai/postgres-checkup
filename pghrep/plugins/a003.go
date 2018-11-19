package main

import (
    "fmt"
    "../src/pyraconv"
    "../src/fmtutils"
    "strconv"
)

var Data map[string]interface{}

type prepare string

func getSharefBuffersValue(hostSettings []interface{}) int64 {
    for _, setting := range hostSettings {
        curSetting := setting.(map[string]interface{})
        if curSetting["name"] == "shared_buffers" {
            if curSetting["unit"] != nil {
                val, err := strconv.ParseInt(pyraconv.ToString(curSetting["setting"]), 10, 64)
                if err != nil {
                    return -1
                }
                unit := fmtutils.GetUnit(pyraconv.ToString(curSetting["unit"]));
                if unit != -1 {
                    val = val * unit
                }
                return val
            }
        }
    }
    return -1
}

func (g prepare) Prepare(data map[string]interface{}) map[string]interface{} {
    result := make(map[string]interface{})
    hosts := pyraconv.ToStringArray(data["hosts"])
    checkData := data["checkData"].(map[string]interface{})
    for _, host := range hosts {
        fmt.Println("Host:", host)
        hostSettings := checkData[host].([]interface{})
        shared_buffers := getSharefBuffersValue(hostSettings)
        result["current"] = "Current value of settings shared_buffers is " + fmtutils.ByteFormat(float64(shared_buffers), 0)
    }

    result["recommended"] = "RECOMMENDED VALUES"
    result["conclusion"] = "CONCLUSION VALUES"
    result["filename"] = "a011_shared_buffers.md"
    return result
}

var Preparer prepare