package fmtutils

import (
    "strings"
    "strconv"
    "math"
)

func RoundUp(input float64, places int) (newVal float64) {
     var round float64
     pow := math.Pow(10, float64(places))
     digit := pow * input
     round = math.Ceil(digit)
     newVal = round / pow
     return
}

func ByteFormat(inputNum float64) string {
    var unit string
    var returnVal float64
    var precision int
    precision = 0

    if inputNum >= math.Pow(1000, 8) {
        returnVal = RoundUp((inputNum / math.Pow(1024, 8)), precision)
        unit = " YiB" // yottabyte
    } else if inputNum >= math.Pow(1000, 7) {
        returnVal = RoundUp((inputNum / math.Pow(1024, 7)), precision)
        unit = " ZiB" // zettabyte
    } else if inputNum >= math.Pow(1000, 6) {
        returnVal = RoundUp((inputNum / math.Pow(1024, 6)), precision)
        unit = " EiB" // exabyte
    } else if inputNum >= math.Pow(1000, 5) {
        returnVal = RoundUp((inputNum / math.Pow(1024, 5)), precision)
        unit = " PiB" // petabyte
    } else if inputNum >= math.Pow(1000, 4) {
        returnVal = RoundUp((inputNum / math.Pow(1024, 4)), precision)
        unit = " TiB" // terrabyte
    } else if inputNum >= math.Pow(1000, 3) {
        returnVal = RoundUp((inputNum / math.Pow(1024, 3)), precision)
        unit = " GiB" // gigabyte
    } else if inputNum >= math.Pow(1000, 2) {
        returnVal = RoundUp((inputNum / math.Pow(1024, 2)), precision)
        unit = " MiB" // megabyte
    } else if inputNum >= 1000 {
        returnVal = RoundUp((inputNum / 1024), precision)
        unit = " KiB" // kilobyte
    } else {
        returnVal = inputNum
        unit = " bytes" // byte
    }

    return strconv.FormatFloat(returnVal, 'f', precision, 64) + unit
}
 
func GetUnit(unit string) int64 {
    var factor int64 = 1
    var value int64 = 0
    if (strings.Contains(unit, "bytes")) {
        factor = 1
    } else if (strings.Contains(unit, "kB")) {
        factor = 1024
    } else if (strings.Contains(unit, "MB")) {
        factor = int64(math.Pow(1024, 2))
    } else if (strings.Contains(unit, "GB")) {
        factor = int64(math.Pow(1024, 3))
    } else if (strings.Contains(unit, "TB")) {
        factor = int64(math.Pow(1024, 4))
    } else if (strings.Contains(unit, "PB")) {
        factor = int64(math.Pow(1024, 5))
    } else if (strings.Contains(unit, "EB")) {
        factor = int64(math.Pow(1024, 6))
    } else if (strings.Contains(unit, "ZB")) {
        factor = int64(math.Pow(1024, 7))
    } else if (strings.Contains(unit, "YB")) {
        factor = int64(math.Pow(1024, 8))
    } else {
        return -1
    }
    //fmt.Println("factor is :", factor)
    r := strings.NewReplacer("bytes", "", "kB", "", "MB", "", "GB", "", "TB", "", "PB", "", "EB", "", "ZB", "", "YB", "")
	val := r.Replace(unit)
    intval, err := strconv.ParseInt(val, 10, 64)
    if err != nil {
        intval = 1
    }
    value = intval * factor
    return value
}