package main

import (
	"flag"
	"notify/app"
)

func main() {
	bind := flag.String("bind", ":8000", "Bind address")
	bindctl := flag.String("bindctl", "127.0.0.1:8001", "Control server bind address")
	secret := flag.String("secret", "secret", "Secret key")
	flag.Parse()
	server := app.NewServer(*bind, *bindctl, *secret)
	server.Start()
}
