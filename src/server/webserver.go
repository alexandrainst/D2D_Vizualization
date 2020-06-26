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
	defer c.Close()
	for {
		/* mt, message, err := c.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		} */
		//log.Println("rdy")
		select {
		case agent := <-agentsInfo:
			log.Println("updated received")
			log.Println(agent)
			//drone, err := json.Marshal(update)
			err = c.WriteMessage(1, agent)
			if err != nil {
				log.Println("write:", err)
				break
			}
		default:
			//log.Println("no message from agents")
		}

	}
}

func input(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	log.Println(agentsInfo)
	defer c.Close()
	for {
		_, message, err := c.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}
		//log.Printf("recv: %s", message)
		select {
		case agentsInfo <- message:
			log.Println("update received and forwarded")
		default:
			//log.Println("updated receieved and ignored")
		}
		//log.Println("on channel")
		/* err = c.WriteMessage(mt, message)
		if err != nil {
			log.Println("write:", err)
			break
		} */
	}
}

func startWebServer() {
	flag.Parse()
	log.SetFlags(0)
	log.Println("Starting...")
	http.HandleFunc("/output", output)
	http.HandleFunc("/input", input)
	agentsInfo = make(chan []byte)

	fs := http.FileServer(http.Dir("../../html"))
	http.Handle("/", fs)
	log.Fatal(http.ListenAndServe(*addr, nil))

	/* for {
		log.Println(<-updates)
	} */

}
