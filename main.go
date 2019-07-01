package main

import (
	"flag"
	"fmt"
)

func main() {
	filename := flag.String("path", "", "path to csv file")
	hostname := flag.String("hostname", "http://localhost:8080", "lora-app-server full hostname (e.g., https://example.com:8080)")
	username := flag.String("username", "", "lora-app-server user username")
	password := flag.String("password", "", "lora-app-server user password")
	flag.Parse()
	if *filename == "" {
		fmt.Println("error: no path given")
		return
	}
	if *username == "" {
		fmt.Println("error: no username given")
		return
	}
	if *password == "" {
		fmt.Println("error: no password given")
		return
	}
	err := Load(*filename, *hostname, *username, *password)
	if err != nil {
		fmt.Printf("load error: %s\n", err)
	}
	return
}
