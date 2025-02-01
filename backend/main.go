package main

import (
	"financial-helper/server"
	"log"
)

func main() {
	gin_server, err := server.GetNewServer()
	if err != nil {
		log.Fatal("Could not get the server object: ", err)
	}

	err = gin_server.Router.Run(":3000")
	if err != nil {
		log.Fatal("Could not start the server: ", err)
	}

	server.TestHoldings()
}
