package main

import "log"

func main() {
	//updates := make(chan struct{})
	//create some dummy drones for testing
	//updates = createNumberOfDummies(5)

	//start the webserver
	log.Println("go go")
	startWebServer()

	//time.Sleep(1000 * time.Millisecond)

}
