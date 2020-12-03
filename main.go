package main

import (
	"log"

	"alexandra.dk/D2D-WebVisualization/webserver"
	comm "github.com/alexandrainst/D2D-communication"
)

func main() {
	log.Println("Starting webserver")
	webserver.StartWebServer()
	log.Println("Starting P2P communication")
	go startVizwork()
	select {}

}

func startVizwork() {
	comm.InitD2DCommuncation(comm.VisualizationAgentType)
	vizChannel := comm.InitVisualizationMessages(true)
	log.Println("Wating for visualization message")
	go func() {
		for {

			msg := <-vizChannel.Messages
			if msg.StateMessage.Mission.SwarmGeometry != nil {
				msg.MissionBound = msg.StateMessage.Mission.SwarmGeometry.Bound()
			}

			select {
			case webserver.AgentsInfo <- *msg:
				//log.Println("sent message")
			default:

			}

		}
	}()
	for {
		goal := <-webserver.GoalInfo
		comm.SendVizGoal(goal)
	}

}
