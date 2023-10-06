package main

import "github.com/dungnguyen/clean-architecture/infrastructure"

func main() {
	infrastructure.NewHTTPServer().Start()
}
