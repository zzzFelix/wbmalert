package main

import "regexp"

func removeByRegexp(input string, pattern string) string {
	if pattern == "" {
		return input
	}
	regexpPattern := regexp.MustCompile(pattern)
	return regexpPattern.ReplaceAllString(input, "")
}
