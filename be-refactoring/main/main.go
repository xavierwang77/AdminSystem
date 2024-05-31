package main

import "admin-system-backend/impl"

func main() {
	// StartHttpServer()
	err := impl.StartHttpServer()
	if err != nil {
		panic(err)
	}
}
