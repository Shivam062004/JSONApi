package main

import "log"

func main() {
	store, err := createPostGresStore()
	if err != nil {
		log.Fatal(err)
	}
	err2 := store.setUp()
	if err2 != nil {
		log.Fatal(err2)
	}
	server := makeServer(":8081", store)
	server.run()
}
