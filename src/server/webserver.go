package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var addr = flag.String("addr", "localhost:8080", "http service address")

var upgrader = websocket.Upgrader{} // use default options
var agentsInfo chan []byte

func output(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	log.Println(agentsInfo)
	writeToClient(c, 1, []byte("From Server: connection established"))

	go func() {
		log.Println("sub started from website")
		defer c.Close()
		for {
			select {
			case agent := <-agentsInfo:
				//drone, err := json.Marshal(update)
				writeToClient(c, 1, agent)
			default:
				//log.Println("no message from agents")
			}

		}
	}()
}

func input(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	log.Println(agentsInfo)

	go func() {
		log.Println("sub started from agent")
		defer c.Close()
		for {

			//agentsInfo <- []byte("hej")
			_, message, err := c.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				break
			}
			agentsInfo <- []byte(message)
			/* log.Printf("recv: %s", message)
			msg := "hi"
			select {
			case agentsInfo <- []byte(msg):
				fmt.Println("sent message", msg)
			default:
				fmt.Println("no message sent")
			} */
			/* select {
			//case agentsInfo <- []byte(message):
			case agentsInfo <- []byte("message"):
				log.Println("update received and forwarded")
			default:
				//log.Println("updated receieved and ignored")
			} */
			//log.Println("on channel")
		}
	}()
}

func writeToClient(c *websocket.Conn, messageType int, message []byte) bool {
	err := c.WriteMessage(messageType, message)
	if err != nil {
		log.Println("write:", err)
		return false
	}
	return true
}

func startWebServer() {
	flag.Parse()
	log.SetFlags(0)
	log.Println("Starting...")
	agentsInfo = make(chan []byte)
	http.HandleFunc("/output", output)
	http.HandleFunc("/input", input)
	//go test()
	fs := http.FileServer(http.Dir("../../html"))
	http.Handle("/", fs)
	log.Fatal(http.ListenAndServe(*addr, nil))

	/* for {
		log.Println(<-updates)
	} */

}
