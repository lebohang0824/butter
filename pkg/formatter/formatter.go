package formatter

import (
	"regexp"
	"strings"
)

var keywordValueRe = regexp.MustCompile(`^\s*(app|description|version|feature|endpoint|method)\s+\S`)

var blockKeywordRe = regexp.MustCompile(`^\s*(actions|params|responses|returns|rules)\s*$`)

var commentRe = regexp.MustCompile(`^\s*(#|//)`)

var rootBlockRe = regexp.MustCompile(`^\s*(app|feature|endpoint)\s+\S`)

func Format(content []byte) ([]byte, error) {
	lines := strings.Split(string(content), "\n")
	lines = normalizeIndentation(lines)
	lines = pass1(lines)
	lines = pass2(lines)
	return []byte(strings.Join(lines, "\n")), nil
}

func normalizeIndentation(lines []string) []string {
	result := make([]string, len(lines))
	for i, line := range lines {
		indent := 0
		for indent < len(line) && (line[indent] == ' ' || line[indent] == '\t') {
			indent++
		}
		prefix := strings.Repeat(" ", strings.Count(line[:indent], " ")+2*strings.Count(line[:indent], "\t"))
		result[i] = prefix + line[indent:]
	}
	return result
}

func isEmpty(s string) bool {
	return strings.TrimSpace(s) == ""
}

func startsWithFeature(s string) bool {
	return strings.HasPrefix(strings.TrimSpace(s), "feature ")
}

func startsWithEndpoint(s string) bool {
	return strings.HasPrefix(strings.TrimSpace(s), "endpoint ")
}

func startsWithApp(s string) bool {
	return strings.HasPrefix(strings.TrimSpace(s), "app ")
}

func pass1(lines []string) []string {
	result := make([]string, 0, len(lines))
	for i := 0; i < len(lines); i++ {
		result = append(result, lines[i])
		if keywordValueRe.MatchString(lines[i]) {
			for i+1 < len(lines) && isEmpty(lines[i+1]) {
				i++
			}
		}
	}
	return result
}

func pass2(lines []string) []string {
	result := make([]string, 0, len(lines))
	for _, line := range lines {
		if blockKeywordRe.MatchString(line) {
			rightBelowFeature := false
			for j := len(result) - 1; j >= 0; j-- {
				prev := result[j]
				if isEmpty(prev) || commentRe.MatchString(prev) {
					continue
				}
				if startsWithFeature(prev) || startsWithEndpoint(prev) {
					rightBelowFeature = true
				}
				break
			}

			if !rightBelowFeature {
				if len(result) > 0 && !isEmpty(result[len(result)-1]) {
					result = append(result, "")
				}
			}
		} else if startsWithFeature(line) || startsWithEndpoint(line) {
			insertIdx := len(result)
			for insertIdx > 0 && (isEmpty(result[insertIdx-1]) || commentRe.MatchString(result[insertIdx-1])) {
				insertIdx--
			}
			if insertIdx > 0 && (insertIdx >= len(result) || !isEmpty(result[insertIdx])) {
				result = append(result[:insertIdx], append([]string{""}, result[insertIdx:]...)...)
			}
		} else if startsWithApp(line) {
			for len(result) > 0 && isEmpty(result[len(result)-1]) {
				result = result[:len(result)-1]
			}
		}
		result = append(result, line)
	}
	return result
}
