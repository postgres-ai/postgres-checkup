package main

import (
    "strings"
    "strconv"
    "./pyraconv"
    "./fmtutils"
)

func Split(s string, d string) []string {
    arr := strings.Split(s, d)
    return arr
}

func Trim(s string, d string) string {
    return strings.Trim(s, d)
}

func Nobr(s interface{}) string {
    str := pyraconv.ToString(s)
    return strings.Join(strings.Split(str, "\n"), "")
}

/* Add \t before every row in text to preview block as code block
* s String - string for preprocces
* skipFirst bool - flag to skip first row
*/
func Code(s string, skipFirst bool) string {
    codeLines := strings.Split(s, "\n")
    for i, line := range codeLines {
        if i > 0 {
            codeLines[i] = "\t" + line
        }
    }
    if skipFirst {
        return strings.Join(codeLines, "\n")
    } else {
        return "\t" + strings.Join(codeLines, "\n")
    }
}


func UnitValue(value interface{}, unit interface{}) string {
    val := pyraconv.ToString(value)
    un := pyraconv.ToString(unit)
    if len(un) <= 0 {
        return val
    }
    intval, err := strconv.ParseInt(val, 10, 64)
    if err != nil {
        return val + "(" + un + ")"
    }
    if intval < 0 {
        return val
    }
    unitFactor := fmtutils.GetUnit(un)
    intval = intval * unitFactor
    return fmtutils.ByteFormat(float64(intval), 0)
}