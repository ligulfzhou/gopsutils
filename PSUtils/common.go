package PSUtils

import (
	"strconv"
	"strings"
)

/*
I think It`s ok to just trim space, ', ", \n
not as dumb as TrimQuotes function.
*/
func StripString(s string) string {
	s1 := strings.Trim(s, " ")
	s2 := strings.Trim(s1, "\"")
	s3 := strings.Trim(s2, "'")
	s4 := strings.Trim(s3, "\r\n")
	s5 := strings.Trim(s4, "\n")
	return s5
}

func splitStringWithSeq(s, seq string) []string {
	var ret []string
	ss := strings.Split(s, seq)
	for _, s := range ss {
		sTrim := StripString(s)
		if sTrim != "" {
			ret = append(ret, sTrim)
		}
	}
	return ret
}

func SplitString(s string) []string {
	var ret []string

	ss := splitStringWithSeq(s, " ")
	for _, i := range ss {
		if strings.Contains(i, "\t") {
			ret = append(ret, splitStringWithSeq(i, "\t")...)
		} else if strings.Contains(i, "\n") {
			ret = append(ret, splitStringWithSeq(i, "\n")...)
		} else {
			ret = append(ret, i)
		}
	}
	return ret
}

func SplitStringToLines(s string) []string {
	// seems \r\n is impossible
	seq := "\r\n"
	if !strings.Contains(seq, "\r\n") {
		seq = "\n"
	}
	return strings.Split(s, seq)
}

func SplitStringWithDeeperLines(s string) []string {
	var names []string
	lines := SplitStringToLines(s)
	for _, line := range lines {
		sp := SplitString(line)
		names = append(names, sp...)
	}
	return names
}

func TrimQuotes(s string) string {
	if len(s) >= 2 {
		if s[0] == '"' && s[len(s)-1] == '"' {
			return s[1 : len(s)-1]
		}
	}
	return s
}

func ParseStringToMap(s, seq string) map[string]string {
	var ret map[string]string

	lines := SplitStringToLines(s)
	for _, line := range lines {
		kv := strings.Split(line, seq)
		if len(kv) != 2 {
			continue
		}
		ret[StripString(kv[0])] = StripString(kv[1])
	}
	return ret
}

func GetValueFromMapString(s, seq, key string) string {
	lines := SplitStringToLines(s)
	for _, line := range lines {
		kv := strings.Split(line, seq)
		if len(kv) != 2 {
			continue
		}

		if StripString(kv[0]) == key {
			return StripString(kv[1])
		}
	}
	return ""
}

func strToInt(s string, i int) int {
	ret, err := strconv.Atoi(s)
	if err != nil {
		return i
	}
	return ret
}

func strToInt64(s string, i int64) int64 {
	ret, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return i
	}
	return ret
}

func ListContain(ss []string, str string) bool {
	for _, s := range ss {
		if s == str {
			return true
		}
	}

	return false
}
