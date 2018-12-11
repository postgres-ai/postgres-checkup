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
