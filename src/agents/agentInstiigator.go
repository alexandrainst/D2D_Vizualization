package main

import (
	"flag"
	"fmt"
	"time"
)

func main() {
	//fmt.Println("Usage: go agentInstigator #NoOfDrones")
	noOfAgents := flag.Int("number", 1, "The number of agents to be initiated")
	isDummy := flag.Bool("isDummy", true, "To run the agent with dummy data - for testing")
	debug := flag.Bool("debug", true, "Have the agent connect to a ws to visualize data")
	wsServerAddr := flag.String("wsAddr", "localhost:8080", "WS server for visualization")
	controllerAddr := flag.String("controllerAddr", "localhost:8044", "WS controller server")

	flag.Parse()
	fmt.Println("Starting", *noOfAgents, "agents")

	agents := make([]*Drone, 0)

	for i := 0; i < *noOfAgents; i++ {
		agent := startDrone(isDummy, debug, wsServerAddr, controllerAddr)
		agents = append(agents, agent)
	}
	for {
		fmt.Println("Status:")

		for i := range agents {
			fmt.Println(agents[i])
		}
		time.Sleep(5 * time.Second)
	}
}
