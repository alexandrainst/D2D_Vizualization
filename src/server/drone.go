package main

import (
	"math/rand"
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

func createDrone(id, x, y, z int) Drone {
	var agent Drone = Drone{id, location{x, y, z}}
	return agent
}

func createNumberOfDummies(number int) chan Drone {
	width := 2000
	height := 2000

	var drones []Drone
	for i := 0; i < number; i++ {
		x := rand.Intn(width) - 1000
		y := rand.Intn(height) - 1000
		z := rand.Intn(50)
		var agent Drone = createDrone(i, x, y, z)
		drones = append(drones, agent)
	}
	updates := make(chan Drone)
	go func() {
		for true {
			for _, drone := range drones {
				//fmt.Printf("id: %v position: %v\n", drone.id, drone.position)
				//fmt.Println("")
				updates <- drone
				time.Sleep(1000 * time.Millisecond)
			}

		}
	}()

	return updates
}
