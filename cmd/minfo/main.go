package main

import (
	"log"

	"minfo"
	"minfo/internal/app"
)

func main() {
	server, err := app.NewServer(minfo.EmbeddedWebUI())
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("minfo listening on http://localhost%s", server.Addr)
	log.Fatal(server.ListenAndServe())
}
