package main

import "log"

//go:generate bash build-proto.sh
//go:generate bash build-inject-tag.sh
//go:generate bash build-reverse-proxy.sh

func main() {
	log.Println("Hello")
}
