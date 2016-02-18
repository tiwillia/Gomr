package main

import (
	"log"
	"regexp"
)

// Test if 'match' matches 'line' and, if it does, return the substring indicated by regex 'pull'
func MatchAndPull(line, regex, pull string) (substr string) {
	matchRegex, err := regexp.Compile(regex)
	if err != nil {
		log.Println("Unable to parse regex:", regex)
		return ""
	}

	if matchRegex.MatchString(line) {
		pullRegex := regexp.MustCompile(pull)
		substr := pullRegex.FindStringSubmatch(line)
		if len(substr) == 0 {
			return ""
		}
		return substr[1]
	} else {
		return ""
	}
}

// Test if line matches the provided regex string
func Match(line, regex string) bool {
	matchRegex, err := regexp.Compile(regex)
	if err != nil {
		log.Println("Unable to parse regex:", regex)
		return false
	}

	return matchRegex.MatchString(line)
}
