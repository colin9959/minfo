package main

import "log"

func main() {
	server, err := newServer()
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("minfo listening on http://localhost%s", server.Addr)
	log.Fatal(server.ListenAndServe())
}
