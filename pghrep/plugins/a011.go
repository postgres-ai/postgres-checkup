package main

import "fmt"

var Data map[string]interface{}

type prepare string

func (g prepare) Prepare(data map[string]interface{}) map[string]interface{} {
    fmt.Println("Hello from A011!", data)
    result := make(map[string]interface{})
    result["current"] = "CURRENT VALUES"
    result["recommended"] = "RECOMMENDED VALUES"
    result["conclusion"] = "CONCLUSION VALUES"
    result["filename"] = "a011_shared_buffers.md"
    return result
}

var Preparer prepare