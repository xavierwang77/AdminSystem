package main

import "be/impl"

func main() {
	// StartHttpServer()
	err := impl.StartHttpServer()
	if err != nil {
		panic(err)
	}
}
