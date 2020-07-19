package main

import (
	"log"
	"math/rand"
	"os"
	"strconv"
	"time"
)

/* type location struct {
	X int
	Y int
	Z int
} */

//Drone is a struct for handling data
type Drone struct {
	ID       int
	Position Vector
}

var toWS chan Drone
var fromController chan map[string]interface{}
var interrupt chan os.Signal
var ownPath = make([][]int, 0)
var swarmPath = make([][]int, 0)
var deltaMovement = float64(5)

func runAsDummy() *Drone {
	width := 2000
	height := 2000

	id := rand.Intn(667)
	x := rand.Intn(width) - 1000
	y := rand.Intn(height) - 1000
	z := rand.Intn(50)
	var agent Drone = Drone{id, Vector{float64(x), float64(y), float64(z)}}
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
	var agent Drone = Drone{id, Vector{0, 0, 0}}

	go func() {
		mission := <-fromController
		log.Println("Starting mission received:", mission)

		for true {
			log.Println("heppa")
			select {
			case newMission := <-fromController:
				log.Println("New mission recieved")
				log.Println(newMission)
			default:

				//normal work
				//first we the vector from current position to the current waypoint
				relevantWP := findRelevantWaypoint(mission["path"].([]interface{}))
				//first we get the direction
				direction := relevantWP.Sub(agent.Position)
				//check if relevantWP is the same as current pos - if so, mark as visited
				if direction.Length() == 0 {
					//TODO: mark as visited and get new WP - if no more wps unvisited, mark mission as done
				}
				//now we normalize
				normalizedDirection := direction.Normalize()

				//next we scale by delta
				newPos := normalizedDirection.MultiplyByScalar(deltaMovement)

				agent.Position = agent.Position.Add(newPos)

				//non-blocking send to channel
				select {
				case toWS <- agent:

				default:

				}

				time.Sleep(100 * time.Millisecond)

			}
		}
	}()

	return &agent
}

func findRelevantWaypoint(path []interface{}) Vector {
	var wp Vector

	for _, v := range path {
		wps := v.(map[string]interface{})
		if wps["visited"] == false {
			var coord []float64
			for _, c := range wps["coord"].([]interface{}) {
				point, _ := strconv.ParseFloat(c.(string), 64)
				coord = append(coord, point)
			}
			wp = Vector{coord[0], coord[1], coord[2]}
		}
	}
	return wp
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
