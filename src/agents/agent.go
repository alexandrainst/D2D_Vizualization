package main

import (
	"log"
	"math/rand"
	"os"
	"time"
)

type location struct {
	X int
	Y int
	Z int
}

//Drone is a struct for handling data
type Drone struct {
	ID       int
	Position location
}

var toWS chan Drone
var fromController chan map[string]interface{}
var interrupt chan os.Signal

func runAsDummy() *Drone {
	width := 2000
	height := 2000

	id := rand.Intn(667)
	x := rand.Intn(width) - 1000
	y := rand.Intn(height) - 1000
	z := rand.Intn(50)
	var agent Drone = Drone{id, location{x, y, z}}
	go func() {

		for true {
			//fmt.Printf("id: %v position: %v\n", drone.id, drone.position)
			key := rand.Intn(3)
			switch key {
			case 0:
				agent.Position.X = agent.Position.X + 10
			case 1:
				agent.Position.Y = agent.Position.Y + 10
			case 2:
				agent.Position.Z = agent.Position.Z + 10
			}

			toWS <- agent
			time.Sleep(100 * time.Millisecond)
		}
	}()
	return &agent
}

func runAsActual() *Drone {
	id := rand.Intn(667)
	var agent Drone = Drone{id, location{0, 0, 0}}

	go func() {
		mission := <-fromController
		for true {
			select {
			case newMission := <-fromController:
				log.Println("New mission recieved")
				log.Println(newMission)
			default:
				//normal work
				log.Println(mission)
			}
		}
	}()

	return &agent
}

func startDrone(isDummy *bool, debug *bool, serverAddr *string, controllerAddr *string) *Drone {

	toWS = make(chan Drone)
	fromController = make(chan map[string]interface{})
	interrupt = make(chan os.Signal, 1)
	//signal.Notify(interrupt, os.Interrupt)
	if *debug {
		wsConn := connectToWsServer(serverAddr, "input")
		go sendToWsServer(wsConn, serverAddr)
	}
	var agent *Drone

	if *isDummy {
		log.Println("Start as dummy")
		agent = runAsDummy()
	} else {
		log.Println("Start as actual")
		agent = runAsActual()
	}
	controllerConn := connectToWsServer(controllerAddr, "agents")
	go controllerComm(agent, controllerConn, controllerAddr)
	return agent
}

/* func main() {

	isDummy := flag.Bool("isDummy", true, "To run the agent with dummy data - for testing")
	debug := flag.Bool("debug", true, "Have the agent connect to a ws to visualize data")
	addr := flag.String("addr", "localhost:8080", "http service address")

	flag.Parse()
	startDrone(isDummy, debug, addr)
} */
