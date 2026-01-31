package main

import "regexp"

func headings(s string) string {
	
	re := regexp.MustCompile("^# (.*)")
	replace := "<h1>$1</h1>"

	ret := re.ReplaceAllString(s, replace)

	return ret
}
