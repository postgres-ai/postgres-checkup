package pyraconv

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"math"
	"strconv"
	"time"
	// "code.google.com/p/go.text/encoding"
	// "code.google.com/p/go.text/encoding/charmap"
	// "code.google.com/p/go.text/encoding/japanese"
	// "code.google.com/p/go.text/encoding/korean"
	// "code.google.com/p/go.text/encoding/simplifiedchinese"
	// "code.google.com/p/go.text/encoding/traditionalchinese"
)

// ToInterfaceArray converts interface{} to []interface{} and never returns nil
func ToInterfaceArray(i1 interface{}) []interface{} {
	if i1 == nil {
		return []interface{}{}
	}
	switch i2 := i1.(type) {
	default:
		return []interface{}{}
	case []interface{}:
		return i2
	}
	//    return []interface{}{}
}

// ToInterfaceMap converts interface{} to map[string]interface{} and never returns nil
func ToInterfaceMap(i1 interface{}) map[string]interface{} {
	if i1 == nil {
		return map[string]interface{}{}
	}
	switch i2 := i1.(type) {
	case map[string]interface{}:
		return i2
	default:
		return map[string]interface{}{}
	}
	//    return map[string]interface{}{}
}

// ToStringArray converts interface{} to []string and never returns nil
func ToStringArray(i1 interface{}) []string {
	if i1 == nil {
		return []string{}
	}
	switch i2 := i1.(type) {
	default:
		return []string{fmt.Sprint(i2)}
	case []string:
		return i2
	case []interface{}:
		var ss []string
		for _, i3 := range i2 {
			ss = append(ss, ToString(i3))
		}
		return ss
	}
	//    return []string{}
}

// ToStringMap converts interface{} to map[string]string and never returns nil
func ToStringMap(i1 interface{}) map[string]string {
	switch i2 := i1.(type) {
	case map[string]interface{}:
		m1 := map[string]string{}
		for k, v := range i2 {
			m1[k] = ToString(v)
		}
		return m1
	case map[string]string:
		return i2
	default:
		return map[string]string{}
	}
}

// ToString converts interface{} to string
func ToString(i1 interface{}) string {
	if i1 == nil {
		return ""
	}
	switch i2 := i1.(type) {
	default:
		return fmt.Sprint(i2)
	case bool:
		if i2 {
			return "true"
		} else {
			return "false"
		}
	case string:
		return i2
	case *bool:
		if i2 == nil {
			return ""
		}
		if *i2 {
			return "true"
		} else {
			return "false"
		}
	case *string:
		if i2 == nil {
			return ""
		}
		return *i2
	case *json.Number:
		return i2.String()
	case json.Number:
		return i2.String()
	}
	//    return ""
}

// ToInt64 converts interface{} to int64
func ToInt64(i1 interface{}) int64 {
	if i1 == nil {
		return 0
	}
	switch i2 := i1.(type) {
	default:
		i3, _ := strconv.ParseInt(ToString(i2), 10, 64)
		return i3
	case *json.Number:
		i3, _ := i2.Int64()
		return i3
	case json.Number:
		i3, _ := i2.Int64()
		return i3
	case int64:
		return i2
	case float64:
		return int64(i2)
	case float32:
		return int64(i2)
	case uint64:
		return int64(i2)
	case int:
		return int64(i2)
	case uint:
		return int64(i2)
	case bool:
		if i2 {
			return 1
		} else {
			return 0
		}
	case *bool:
		if i2 == nil {
			return 0
		}
		if *i2 {
			return 1
		} else {
			return 0
		}
	}
	//    return 0
}

// ToFloat64 converts interface{} to float64
func ToFloat64(i1 interface{}) float64 {
	if i1 == nil {
		return 0.0
	}
	switch i2 := i1.(type) {
	default:
		i3, _ := strconv.ParseFloat(ToString(i2), 64)
		return i3
	case *json.Number:
		i3, _ := i2.Float64()
		return i3
	case json.Number:
		i3, _ := i2.Float64()
		return i3
	case int64:
		return float64(i2)
	case float64:
		return i2
	case float32:
		return float64(i2)
	case uint64:
		return float64(i2)
	case int:
		return float64(i2)
	case uint:
		return float64(i2)
	case bool:
		if i2 {
			return 1.0
		} else {
			return 0.0
		}
	case *bool:
		if i2 == nil {
			return 0.0
		}
		if *i2 {
			return 1.0
		} else {
			return 0.0
		}
	}
	//    return 0.0
}

// ToFloat64 converts interface{} to float64
func ToFloat32(i1 interface{}) float32 {
	if i1 == nil {
		return 0.0
	}
	switch i2 := i1.(type) {
	default:
		i3, _ := strconv.ParseFloat(ToString(i2), 64)
		return float32(i3)
	case *json.Number:
		i3, _ := i2.Float64()
		return float32(i3)
	case json.Number:
		i3, _ := i2.Float64()
		return float32(i3)
	case int64:
		return float32(i2)
	case float64:
		return float32(i2)
	case float32:
		return i2
	case uint64:
		return float32(i2)
	case int:
		return float32(i2)
	case uint:
		return float32(i2)
	case bool:
		if i2 {
			return 1.0
		} else {
			return 0.0
		}
	case *bool:
		if i2 == nil {
			return 0.0
		}
		if *i2 {
			return 1.0
		} else {
			return 0.0
		}
	}
	//    return 0.0
}

// ToBool converts interface{} to bool
func ToBool(i1 interface{}) bool {
	if i1 == nil {
		return false
	}
	switch i2 := i1.(type) {
	default:
		return false
	case bool:
		return i2
	case string:
		return i2 == "true"
	case int:
		return i2 != 0
	case *bool:
		if i2 == nil {
			return false
		}
		return *i2
	case *string:
		if i2 == nil {
			return false
		}
		return *i2 == "true"
	case *int:
		if i2 == nil {
			return false
		}
		return *i2 != 0
	}
	//    return false
}

// CloneObject creates copy of object using GOB.
func CloneObject(a, b interface{}) {
	buff := new(bytes.Buffer)
	enc := gob.NewEncoder(buff)
	dec := gob.NewDecoder(buff)
	enc.Encode(a)
	dec.Decode(b)
}

// func ShoterRunes(s1 string, n int) string {
//     r1 := []rune(s1)
//     if len(r1) > n {
//         r1 = r1[0:n]
//     }
//     s2 := string(r1)
//     return s2
// }

// func SecureWebId(s1 string) (string, bool) {
//     for _, c := range s1 {
//         if !((c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9')) {
//             return "", false
//         }
//     }
//     return s1, true
// }
// func DebugObject(id string, o1 interface{}) {
//     b, _ := json.Marshal(o1)
//     log.Println(id, string(b))
// }

// func CharsetGetEncding(charset string) encoding.Encoding {
//     if charset == "" {
//         return nil
//     }
//     charset = strings.ToLower(charset)
//     charset = strings.Replace(charset, "-", "", -1)
//     if charset == "utf8" {
//         return nil
//     }
//     if charset == "koi8" {
//         return charmap.KOI8R
//     }
//     if charset == "koi8r" {
//         return charmap.KOI8R
//     }
//     if charset == "koi8u" {
//         return charmap.KOI8U
//     }

//     if charset == "iso88591" {
//         return charmap.Windows1252
//     }

//     if charset == "iso88592" {
//         return charmap.ISO8859_2
//     }
//     if charset == "iso88593" {
//         return charmap.ISO8859_3
//     }
//     if charset == "iso88594" {
//         return charmap.ISO8859_4
//     }
//     if charset == "iso88595" {
//         return charmap.ISO8859_5
//     }
//     if charset == "iso88596" {
//         return charmap.ISO8859_6
//     }
//     if charset == "iso88597" {
//         return charmap.ISO8859_7
//     }
//     if charset == "iso88598" {
//         return charmap.ISO8859_8
//     }
//     // if charset == "iso88599" {return charmap.ISO8859_9}
//     if charset == "iso885910" {
//         return charmap.ISO8859_10
//     }
//     // if charset == "iso885911" {return charmap.ISO8859_11}
//     // if charset == "iso885912" {return charmap.ISO8859_12}
//     if charset == "iso885913" {
//         return charmap.ISO8859_13
//     }
//     if charset == "iso885914" {
//         return charmap.ISO8859_14
//     }
//     if charset == "iso885915" {
//         return charmap.ISO8859_15
//     }
//     if charset == "iso885916" {
//         return charmap.ISO8859_16
//     }

//     if charset == "windows1250" {
//         return charmap.Windows1250
//     }
//     if charset == "windows1251" {
//         return charmap.Windows1251
//     }
//     if charset == "windows1252" {
//         return charmap.Windows1252
//     }
//     if charset == "windows1253" {
//         return charmap.Windows1253
//     }
//     if charset == "windows1254" {
//         return charmap.Windows1254
//     }
//     if charset == "windows1255" {
//         return charmap.Windows1255
//     }
//     if charset == "windows1256" {
//         return charmap.Windows1256
//     }
//     if charset == "windows1257" {
//         return charmap.Windows1257
//     }
//     if charset == "windows1258" {
//         return charmap.Windows1258
//     }
//     if charset == "windows874" {
//         return charmap.Windows874
//     }

//     if charset == "eucjp" {
//         return japanese.EUCJP
//     }
//     if charset == "iso2022jp" {
//         return japanese.ISO2022JP
//     }
//     if charset == "shiftjis" {
//         return japanese.ShiftJIS
//     }

//     if charset == "euckr" {
//         return korean.EUCKR
//     }

//     if charset == "gb18030" {
//         return simplifiedchinese.GB18030
//     }
//     if charset == "gbk" {
//         return simplifiedchinese.GBK
//     }
//     if charset == "hzgb2312" {
//         return simplifiedchinese.HZGB2312
//     }

//     if charset == "big5" {
//         return traditionalchinese.Big5
//     }

//     return nil

// }
func TimeSince(t1 time.Time) (since string) {
	now := time.Now()
	d1 := time.Since(t1)
	if d1.Hours() > 24.0 {
		if t1.Year() == now.Year() {
			if t1.Month() == now.Month() {
				if t1.Day() == now.Day() {
					since = t1.Format("15:04")
				} else {
					since = t1.Format("02-Jan")
				}
			} else {
				since = t1.Format("02-Jan")
			}
		} else {
			since = t1.Format("02-Jan-2006")
		}
	} else {
		// since = d1.String()
		hours := int(d1.Hours())
		minutes := int(d1.Minutes() - float64(60.0*hours))
		since = fmt.Sprintf("%02d:%02d", hours, minutes)
	}
	return since
}

func MovingExpAvg(value float64, oldValue float64, ftime float64, time float64) (r float64) {
	alpha := 1. - math.Exp(-ftime/time)
	r = alpha*value + (1.-alpha)*oldValue
	return r
}

func MovingExpAvg32(value float32, oldValue float32, ftime float32, time float32) (r float32) {
	alpha := 1. - float32(math.Exp(float64(-ftime/time)))
	r = alpha*value + (1.-alpha)*oldValue
	return r
}

func ReverseMovingExpAvg(r float64, oldValue float64, ftime float64, time float64) (value float64) {
	alpha := 1. - math.Exp(-ftime/time)
	value = (r - (1.-alpha)*oldValue) / alpha
	return value
}

func ReverseMovingExpAvg32(r float32, oldValue float32, ftime float32, time float32) (value float32) {
	alpha := 1. - float32(math.Exp(float64(-ftime/time)))
	value = (r - (1.-alpha)*oldValue) / alpha
	return value
}
