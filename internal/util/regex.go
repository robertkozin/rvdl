package util

import "regexp"

func FirstSubmatch(pattern *regexp.Regexp, s string) string {
	if submatches := pattern.FindStringSubmatch(s); submatches != nil && len(submatches) >= 2 {
		return submatches[1]
	}
	return ""
}
