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

func ByteFormat(inputNum float64, precision int) string {
    var unit string
    var returnVal float64

    if inputNum >= 1000000000000000000000000 {
        returnVal = RoundUp((inputNum / 1208925819614629174706176), precision)
        unit = " YiB" // yottabyte
    } else if inputNum >= 1000000000000000000000 {
        returnVal = RoundUp((inputNum / 1180591620717411303424), precision)
        unit = " ZiB" // zettabyte
    } else if inputNum >= 10000000000000000000 {
        returnVal = RoundUp((inputNum / 1152921504606846976), precision)
        unit = " EiB" // exabyte
    } else if inputNum >= 1000000000000000 {
        returnVal = RoundUp((inputNum / 1125899906842624), precision)
        unit = " PiB" // petabyte
    } else if inputNum >= 1000000000000 {
        returnVal = RoundUp((inputNum / 1099511627776), precision)
        unit = " TiB" // terrabyte
    } else if inputNum >= 1000000000 {
        returnVal = RoundUp((inputNum / 1073741824), precision)
        unit = " GiB" // gigabyte
    } else if inputNum >= 1000000 {
        returnVal = RoundUp((inputNum / 1048576), precision)
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
        factor = 1024 * 1024
    } else if (strings.Contains(unit, "GB")) {
        factor = 1024 * 1024 * 1024
    } else if (strings.Contains(unit, "TB")) {
        factor = 1024 * 1024 * 1024 * 1024
    } else if (strings.Contains(unit, "PB")) {
        //factor = 1024 * 1024 * 1024 * 1024 * 1024
        return -1
    } else if (strings.Contains(unit, "EB")) {
        //factor = 1024 * 1024 * 1024 * 1024 * 1024 * 1024
        return -1
    } else if (strings.Contains(unit, "ZB")) {
        //factor = 1024 * 1024 * 1024 * 1024 * 1024 * 1024 * 1024
        return -1
    } else if (strings.Contains(unit, "YB")) {
        //factor = 1024 * 1024 * 1024 * 1024 * 1024 * 1024 * 1024 * 1024
        return -1
    }
    //fmt.Println("factor is :", factor)
    r := strings.NewReplacer("bytes", "", "kB", "", "MB", "", "GB", "", "TB", "", "PB", "", "EB", "", "ZB", "", "YB", "")
	val := r.Replace(unit)
    intval, err := strconv.ParseInt(val, 10, 64)
    if err != nil {
        //fmt.Println("Can't parse :", val, err)
        return -1
    }
    value = intval * factor
    return value
}