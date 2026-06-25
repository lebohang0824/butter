package formatter

import (
	"regexp"
	"strings"
)

var keywordValueRe = regexp.MustCompile(`^\s*(app|product|description|version|feature|param|type|required|default|validate|length)\s+\S`)

var validateSpaceRe = regexp.MustCompile(`^(\s*validate\s+")([><=!]+)\s+(\d+(?:\.\d+)?)"`)

var blockKeywordRe = regexp.MustCompile(`^\s*(actions|params)\s*$`)

var commentRe = regexp.MustCompile(`^\s*#`)

func Format(content []byte) ([]byte, error) {
	lines := strings.Split(string(content), "\n")
	lines = pass1(lines)
	lines = pass2(lines)
	lines = normalizeValidateSpaces(lines)
	return []byte(strings.Join(lines, "\n")), nil
}

func normalizeValidateSpaces(lines []string) []string {
	result := make([]string, len(lines))
	for i, line := range lines {
		result[i] = validateSpaceRe.ReplaceAllString(line, `${1}${2}${3}"`)
	}
	return result
}

func isEmpty(s string) bool {
	return strings.TrimSpace(s) == ""
}

func startsWithFeature(s string) bool {
	return strings.HasPrefix(strings.TrimSpace(s), "feature ")
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
				if startsWithFeature(prev) {
					rightBelowFeature = true
				}
				break
			}

			if !rightBelowFeature {
				if len(result) > 0 && !isEmpty(result[len(result)-1]) {
					result = append(result, "")
				}
			}
		} else if startsWithFeature(line) {
			if len(result) > 0 && !isEmpty(result[len(result)-1]) {
				result = append(result, "")
			}
		}
		result = append(result, line)
	}
	return result
}
