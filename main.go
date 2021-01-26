package main

import (
	"flag"
	"strings"

	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	flagMode := flag.String("mode", "server", "start in client or server mode")
	flag.Parse()

	if strings.ToLower(*flagMode) == "server" {
		StartServerMode()
	} else {
		StartClientMode()
	}
}
