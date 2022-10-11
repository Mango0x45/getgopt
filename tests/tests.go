package main

import (
	"fmt"
	"os"

	"git.thomasvoss.com/getgopt"
)

func main() {
	for opt := byte(0); getgopt.Getopt(len(os.Args), os.Args, ":a:x", &opt); {
		switch opt {
		case 'a':
			fmt.Printf("Valid flag 'a' with arg '%s'\n", getgopt.Optarg)
		case 'x':
			fmt.Println("Valid flag 'x'")
		case ':':
			fmt.Printf("Valid flag '%c' with no arg\n", getgopt.Optopt)
		case '?':
			fmt.Printf("Invalid flag '%c'\n", getgopt.Optopt)
		}
	}
}
