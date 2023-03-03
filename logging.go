package main

import "fmt"

func log(text string) {
	fmt.Printf(text + "\n")
}

func debug(text string) {
	if debugEnabled {
		log(text)
	}
}
