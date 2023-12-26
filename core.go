package main

import (
	"fmt"
	"strings"
)

func FormatResponse(txt string) string {
	return fmt.Sprintf("+%s\r\n", strings.ToUpper(txt))
}
