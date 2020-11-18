package PSUtils

import "strings"

func StripString(s string) string {
	s1 := strings.Trim(s, " ")
	s2 := strings.Trim(s1, "\"")
	s3 := strings.Trim(s2, "'")
	s4 := strings.Trim(s3, "\r\n")
	s5 := strings.Trim(s4, "\n")
	return s5
}

func SplitString(s string) []string {
	seq := "\r\n"
	if !strings.Contains(seq, "\r\n") {
		seq = "\n"
	}
	return strings.Split(s, seq)
}
