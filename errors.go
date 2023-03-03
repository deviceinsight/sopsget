package main

import (
	"fmt"
	"os"
)

func check(e error) {
	if e != nil {
		exitGracefully(e)
	}
}

func exitGracefully(err error) {
	_, err = fmt.Fprintf(os.Stderr, "error: %v\n", err)
	os.Exit(1)
}
