package main

import (
    "strings"
)

func Split(s string, d string) []string {
    arr := strings.Split(s, d)
    return arr
}

func Trim(s string, d string) string {
    return strings.Trim(s, d)
}
