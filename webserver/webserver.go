package webserver

import (
	"encoding/json"
	"flag"
	"log"
	"net/http"

	comm "github.com/alexandrainst/D2D-communication"

	"github.com/gorilla/websocket"
)

var addr = flag.String("addr", "localhost:8080", "http service address")

var upgrader = websocket.Upgrader{} // use default options
var AgentsInfo = make(chan comm.VisualizationMessage, 128)

func output(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}

	if !writeToClient(c, 1, []byte("From Server: connection established")) {
		return
	}

	go func() {
		log.Println("sub started from website")
		defer c.Close()
		for {
			//log.Println("in loop")
			select {
			case vsMessage := <-AgentsInfo:
				agent, err := json.Marshal(vsMessage)
				if err != nil {
					log.Println("ERRR WS!")
					log.Println(err)
				}
				//drone, err := json.Marshal(update)
				if !writeToClient(c, 1, agent) {
					log.Println("faled")
					return
				}
			default:
				//log.Println("no message from agents")
			}
		}
	}()
}

func input(w http.ResponseWriter, r *http.Request) {
	// c, err := upgrader.Upgrade(w, r, nil)
	// if err != nil {
	// 	log.Print("upgrade:", err)
	// 	return
	// }

	// go func() {
	// 	log.Println("sub started from agent")
	// 	defer c.Close()
	// 	for {

	// 		//agentsInfo <- []byte("hej")
	// 		_, message, err := c.ReadMessage()
	// 		if err != nil {
	// 			log.Println("read:", err)
	// 			break
	// 		}
	// 		//AgentsInfo <- []byte(message)
	// 		/* log.Printf("recv: %s", message)
	// 		msg := "hi"
	// 		select {
	// 		case agentsInfo <- []byte(msg):
	// 			fmt.Println("sent message", msg)
	// 		default:
	// 			fmt.Println("no message sent")
	// 		} */
	// 		/* select {
	// 		//case agentsInfo <- []byte(message):
	// 		case agentsInfo <- []byte("message"):
	// 			log.Println("update received and forwarded")
	// 		default:
	// 			//log.Println("updated receieved and ignored")
	// 		} */
	// 		//log.Println("on channel")
	// 	}
	// }()
}

func writeToClient(c *websocket.Conn, messageType int, message []byte) bool {
	err := c.WriteMessage(messageType, message)
	if err != nil {
		log.Println("write:", err)
		return false
	}
	return true
}

func StartWebServer() {
	flag.Parse()
	log.SetFlags(0)
	log.Println("Starting...")

	http.HandleFunc("/output", output)
	//go test()
	dir := http.Dir("html")
	fs := http.FileServer(dir)
	http.Handle("/", fs)

	go func() {
		log.Fatal(http.ListenAndServe(*addr, nil))
	}()
}
