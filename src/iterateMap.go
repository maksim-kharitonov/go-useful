package main

import (
	"golang.org/x/tour/wc"
	"strings"
)

func WordCount(s string) map[string]int {
	wc :=  make(map[string]int)
	for _,v := range strings.Fields(s) {
		cnt,ok := wc[v]
		if ok {
			wc[v] = cnt + 1
		} else {
			wc[v] = 1
		}
	}

	return wc
}

func main() {
	wc.Test(WordCount)
}
