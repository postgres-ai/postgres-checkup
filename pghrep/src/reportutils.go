package main

import (
    "strings"
    "strconv"
    "./pyraconv"
    "./fmtutils"
    "fmt"
    "./dateparse"
    "time"
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
    str = strings.Join(strings.Split(str, "\n"), "")
    str = strings.Join(strings.Split(str, " "), "&nbsp;")
    return str
}

func Br(s interface{}) string {
    str := pyraconv.ToString(s)
    return strings.Join(strings.Split(str, ","), ", ")
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
        return ""
    }
    intval, err := strconv.ParseInt(val, 10, 64)
    if err != nil {
        return "" // val + "(" + un + ")"
    }
    if intval < 0 {
        return "" // val
    }
    unitFactor := fmtutils.GetUnit(un)
    if unitFactor != -1 {
        intval = intval * unitFactor
        return fmtutils.ByteFormat(float64(intval), 2)
    }
    return "" //val + "(" + un + ")"
}

func LimitStr(value interface{}, limit int) string {
    val := pyraconv.ToString(value)
    if len(val) > limit {
        return val[0:limit-1] + "..."
    }
    return val
}

func Round(value interface{}, places interface{}) string {
    val := pyraconv.ToFloat64(value)
    pl := pyraconv.ToInt64(places)
    if value != nil && places != nil {
        return fmt.Sprintf("%v", fmtutils.RoundUp(val, int(pl)))
    }
    return fmt.Sprintf("%v", val)
}

func Add(a int, b int) int {
    return a + b
}

func MsFormat(value interface{}) string {
    val := pyraconv.ToInt64(value)
    tm, _ := time.ParseDuration(strconv.FormatInt(val, 10) + "ms")
    return tm.String()
}

func NumFormat(value interface{}, places interface{}) string {
    val := pyraconv.ToFloat64(value)
    pl := pyraconv.ToInt64(places)
    if pl == -1 {
        return strconv.FormatFloat(val, 'f', int(pl), 64)
    } else {
        return fmtutils.NumFormat(val, int(pl))
    }
}

func DtFormat(value interface{}) string {
    val := pyraconv.ToString(value)
	t, err := dateparse.ParseAny(val)
	if err != nil {
	} else {
		return t.String()
    }
	return val
}