package main

import (
	"encoding/json"
	"flag"
	"log"
	"math/rand"
	"net/url"
	"os"
	"time"

	"github.com/gorilla/websocket"
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

func runAsDummy(updates chan Drone) {
	width := 2000
	height := 2000

	id := rand.Intn(667)
	x := rand.Intn(width) - 1000
	y := rand.Intn(height) - 1000
	z := rand.Intn(50)
	var agent Drone = Drone{id, location{x, y, z}}

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

		updates <- agent
		time.Sleep(100 * time.Millisecond)
	}
}

func connectToWsServer(addr *string) *websocket.Conn {
	u := url.URL{Scheme: "ws", Host: *addr, Path: "/input"}
	log.Printf("connecting to %s", u.String())
	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)

	if err == nil {
		log.Println("connected")
	}
	return c
}
func sendToWsServer(updates chan Drone, interrupt chan os.Signal, c *websocket.Conn, addr *string) {

	defer c.Close()

	for {

		select {
		case update := <-updates:

			drone, err := json.Marshal(update)
			if err != nil {
				log.Println("marshal:", err)
				return
			}

			wsErr := c.WriteMessage(websocket.TextMessage, drone)
			if wsErr != nil {
				log.Println("write:", wsErr)
				time.Sleep(500 * time.Millisecond)
				for i := 0; i < 5; i++ {
					log.Println("trying to reconnect")
					c = connectToWsServer(addr)
					if c != nil {
						break
					}
					time.Sleep(500 * time.Millisecond)
				}

			}

		case <-interrupt:
			log.Println("interrupt")

			// Cleanly close the connection by sending a close message and then
			// waiting (with timeout) for the server to close the connection.
			err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Println("write close:", err)
				return
			}
		}
	}

}

func main() {

	isDummy := flag.Bool("isDummy", true, "To run the agent with dummy data - for testing")
	debug := flag.Bool("debug", true, "Have the agent connect to a ws to visualize data")
	addr := flag.String("addr", "localhost:8080", "http service address")

	flag.Parse()

	updates := make(chan Drone)
	interrupt := make(chan os.Signal, 1)
	//signal.Notify(interrupt, os.Interrupt)
	if *debug {
		url := connectToWsServer(addr)
		go sendToWsServer(updates, interrupt, url, addr)
	}

	if *isDummy {
		runAsDummy(updates)
	}
}
