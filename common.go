package main

import "strings"

func StripString(s string) string {
	s1 := strings.Trim(s, " ")
	s2 := strings.Trim(s1, "\"")
	s3 := strings.Trim(s2, "'")
	return s3
}
