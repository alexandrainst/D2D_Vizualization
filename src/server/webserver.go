package main

import (
	"encoding/json"
	"flag"
	"html/template"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var addr = flag.String("addr", "localhost:8080", "http service address")

var upgrader = websocket.Upgrader{} // use default options

func echo(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer c.Close()
	for {
		/* mt, message, err := c.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		} */
		update := <-updates
		log.Println(update)
		drone, err := json.Marshal(update)
		err = c.WriteMessage(1, drone)
		if err != nil {
			log.Println("write:", err)
			break
		}
	}
}

func home(w http.ResponseWriter, r *http.Request) {

	t, err := template.ParseFiles("edit.html")
	log.Println(t)
	log.Println(err)
	// title := "D2D Visualization"
	//t.Execute(w, title)
	//homeTemplate.Execute(w, "ws://"+r.Host+"/echo")
}

func startWebServer(updates chan Drone) {
	flag.Parse()
	log.SetFlags(0)
	log.Println("Starting...")
	http.HandleFunc("/echo", echo)

	fs := http.FileServer(http.Dir("../../html"))
	http.Handle("/", fs)
	log.Fatal(http.ListenAndServe(*addr, nil))

	/* for {
		log.Println(<-updates)
	} */

}
