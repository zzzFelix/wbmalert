package main

import "regexp"

func removeByRegex(input string, regex string) string {
	if regex == "" {
		return input
	}
	regexpPattern := regexp.MustCompile(regex)
	return regexpPattern.ReplaceAllString(input, "")
}
