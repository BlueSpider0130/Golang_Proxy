package main

import (
	"flag"
	"gitlab.com/devskiller-tasks/messaging-app-golang/restapi"
	"log"
	"net/http"
)

func main() {
	var port = flag.Int("port", 8080, "port")
	flag.Parse()

	server := restapi.NewServer(*port)
	server.BindEndpoints()
	if err := server.Run(); err != http.ErrServerClosed {
		panic(err)
	}
	log.Println("shutdown: completed")
}
