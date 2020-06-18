package main

import "log"

var updates chan Drone

func main() {
	//updates := make(chan Drone)
	//create some dummy drones for testing
	updates = createNumberOfDummies(5)

	//start the webserver
	log.Println("go go")
	startWebServer(updates)

	//time.Sleep(1000 * time.Millisecond)

}
