package main

import "github.com/vlasad/redislike/server"

func main() {
	s := server.New()
	s.Start(8080)
}
