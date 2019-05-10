package main

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"./dateparse"
	"./fmtutils"
	"./pyraconv"
)

func Split(s interface{}, d interface{}) []string {
	str1 := pyraconv.ToString(s)
	str2 := pyraconv.ToString(d)
	arr := strings.Split(str1, str2)
	return arr
}

func Trim(s string, d string) string {
	return strings.Trim(s, d)
}

func Replace(str interface{}, src interface{}, dst interface{}) string {
	str1 := pyraconv.ToString(str)
	src1 := pyraconv.ToString(src)
	dst1 := pyraconv.ToString(dst)
	return strings.Replace(str1, src1, dst1, -1)
}

func Nobr(s interface{}) string {
	str := pyraconv.ToString(s)
	str = strings.Join(strings.Split(str, "\n"), " ")
	str = strings.Join(strings.Split(str, " "), "&nbsp;")
	return str
}

func Br(s interface{}) string {
	str := pyraconv.ToString(s)
	return strings.Join(strings.Split(str, ","), ", ")
}

/* Escape Markdown symbols in a SQL query,
*  convert to a single line
 */
func EscapeQuery(s interface{}) string {
	str := pyraconv.ToString(s)
	str = strings.Join(strings.Split(str, "\n"), " ")
	str = strings.Join(strings.Split(str, " "), "&nbsp;")
	str = strings.Join(strings.Split(str, "http"), "&#104;ttp")
	str = strings.Join(strings.Split(str, "*"), "\\*")
	str = strings.Join(strings.Split(str, "_"), "\\_")
	str = strings.Join(strings.Split(str, "-"), "\\-")
	str = strings.Join(strings.Split(str, "`"), "\\`")
	str = strings.Join(strings.Split(str, "|"), "\\|")
	return str
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

func RawIntUnitValue(value interface{}, unit interface{}) int {
	val := pyraconv.ToString(value)
	un := pyraconv.ToString(unit)
	intval, err := strconv.ParseInt(val, 10, 64)
	if err != nil {
		return -1
	}
	if intval < 0 {
		return -1
	}
	if len(un) <= 0 {
		return int(intval)
	}
	unitFactor := fmtutils.GetUnit(un)
	if unitFactor != -1 {
		intval = intval * unitFactor
		return int(intval)
	}
	return -1
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

func Mul(a int, b int) float64 {
	return float64(a * b)
}

func Div(a int, b int) int {
	return a / b
}

func Sub(a int, b int) int {
	return a - b
}

func MsFormat(value interface{}) string {
	val := pyraconv.ToString(value)
	floatVal := pyraconv.ToFloat64(value)
	tm, terr := time.ParseDuration(val + "&nbsp;ms")
	if tm > time.Second && terr == nil {
		return tm.String()
	} else {
		return strconv.FormatFloat(floatVal, 'f', 3, 64) + "&nbsp;ms"
	}
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

func RawIntFormat(value interface{}) string {
	val := pyraconv.ToInt64(value)
	return fmtutils.RawIntFormat(val)
}

func RawFloatFormat(value interface{}, places interface{}) string {
	val := pyraconv.ToFloat64(value)
	pl := pyraconv.ToInt64(places)
	return fmtutils.RawFloatFormat(val, int(pl))
}

func Int(value interface{}) int {
	if value != nil {
		return int(pyraconv.ToInt64(value))
	}
	return 0
}

func ByteFormat(value interface{}, places interface{}) string {
	val := pyraconv.ToFloat64(value)
	pl := pyraconv.ToInt64(places)
	result := fmtutils.ByteFormat(val, int(pl))
	return Nobr(result)
}
