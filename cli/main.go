package main

import (
	"flag"
	"fmt"

	"github.com/iegomez/lsp"
	log "github.com/sirupsen/logrus"
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
	devices, err := lsp.Load(*filename)
	if err != nil {
		fmt.Printf("load error: %s\n", err)
	}
	token, err := lsp.Login(*username, *password, *hostname)
	if err != nil {
		fmt.Printf("login error: %s\n", err)
	}

	//Provision only logs errors.
	lsp.Provision(devices, *hostname, fmt.Sprintf("Bearer %s", token))

	log.Infoln("finished provisioning devices")

	return
}
