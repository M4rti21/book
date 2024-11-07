package main

import (
	"os"
	"strings"
)

func fixString(str string) string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	str = strings.Replace(str, "~", homeDir, 1)
	str = strings.Trim(str, " ")
	return str
}
