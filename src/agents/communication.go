package main

import (
	"encoding/json"
	"log"
	"net/url"
	"time"

	"github.com/gorilla/websocket"
)

func controllerComm(agent *Drone, controllerConn *websocket.Conn, controllerAddr *string) {
	defer controllerConn.Close()

	registerAtController(agent, controllerConn)

	for {

		select {

		case inter := <-interrupt:
			log.Println("interrupt")
			interrupt <- inter

			// Cleanly close the connection by sending a close message and then
			// waiting (with timeout) for the server to close the connection.
			err := controllerConn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Println("write close:", err)
				return
			}

		default:
			_, message, err := controllerConn.ReadMessage()
			if err != nil {
				time.Sleep(500 * time.Millisecond)
				for i := 0; i < 5; i++ {
					log.Println("trying to reconnect")
					controllerConn = connectToWsServer(controllerAddr, "agents")
					if controllerConn != nil {
						registerAtController(agent, controllerConn)
						break
					}
					time.Sleep(1 * time.Second)
				}
			}
			log.Println(message)
		}
	}
}

func registerAtController(agent *Drone, controllerConn *websocket.Conn) {
	drone, err := json.Marshal(&agent)
	if err != nil {
		log.Println("marshal:", err)
		return
	}
	log.Println(string(drone))
	payload := map[string]interface{}{
		"type":  "register",
		"agent": &agent,
	}

	jsonPayload, _ := json.Marshal(payload)

	wsErr := controllerConn.WriteMessage(websocket.TextMessage, jsonPayload)
	if wsErr != nil {
		log.Println("registering at controller failed:", wsErr)
	}
}

func connectToWsServer(addr *string, path string) *websocket.Conn {
	u := url.URL{Scheme: "ws", Host: *addr, Path: path}
	log.Printf("connecting to %s", u.String())
	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)

	if err == nil {
		log.Println("connected to ", *addr)
	} else {
		log.Println("Not able to connect", err)
	}
	return c
}

func sendToWsServer(c *websocket.Conn, addr *string) {

	defer c.Close()

	for {

		select {
		case update := <-toWS:

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
					c = connectToWsServer(addr, "input")
					if c != nil {
						break
					}
					time.Sleep(1 * time.Second)
				}

			}

		case inter := <-interrupt:
			log.Println("interrupt")
			interrupt <- inter
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
