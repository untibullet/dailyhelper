package main

import (
	"flag"
)

const (
	host = "api.telegram.org"
)

func main() {
	// client := telegram.NewClient(host, mustToken())
}

func mustToken() string {
	token := flag.String("token", "", "need telegram bot token")
	return *token
}
