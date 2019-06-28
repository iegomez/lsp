package main

import(
	"fmt"
	"flag"
)

func main() {
	filename := flag.String("path", "", "path to csv file")
	flag.Parse()
	if *filename == "" {
		fmt.Println("error: no path given")
		return
	}
	err := Load(*filename)
	if err != nil {
		fmt.Printf("load error: %s\n", err)
	}
	return
}
