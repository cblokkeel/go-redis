package main

import "strings"

func StartsWith(str string) string {
	return strings.Split(str, "")[0]
}
