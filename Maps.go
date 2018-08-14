package main

import (
	"golang.org/x/tour/wc"
	"strings"
)

func WordCount(s string) map[string]int {
	m := make(map[string]int)
	
	// Split https://stackoverflow.com/questions/16551354/how-to-split-a-string-and-assign-it-to-variables-in-golang
	for _,c := range strings.Split(s, " ") {
		if m[c] > 0 {
			m[c] = m[c] + 1
		} else {
			m[c] = 1
		}
	}
	
	return map[string]int(m)
}

func main() {
	wc.Test(WordCount)
}

