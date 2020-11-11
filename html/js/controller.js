
import * as Scene from "./app.js";

const DiscoveryMessageType      = 0
const StateMessageType          = 1
const MissionMessageType        = 2
const ReorganizationMessageType = 3
const RecalculatorMessageType   = 4

window.addEventListener("load", function(evt) {
    var ws;
    connectWS();

    function connectWS(evt) {
		console.log("OPEN MG");
        if (ws) {
            return false;
        }
        ws = new WebSocket("ws://localhost:8080/output");
        console.log(ws);
        ws.onopen = function(evt) {
            console.log("OPEN");
        }
        ws.onclose = function(evt) {
            console.log("CLOSE");
            ws = null;
            setTimeout(connectWS, 2000)
        }
        ws.onmessage = function(evt) {
            //console.log("RESPONSE: " + evt.data);
            let msg = null;
            try {
                msg = JSON.parse(evt.data);
            } catch (e) {
                console.log(evt.data);
                return;
            }
            //var drone = JSON.parse(evt.data);
            //console.log(msg);
            let drone = {}
            drone.id = msg.SenderId;
            drone.type = msg.SenderType;
            
            switch (msg.ContentType){
                case DiscoveryMessageType:
                    console.log("discovery mesg")
                    break;
                case StateMessageType:

                    if(msg.MissionBound!=null){
                        let max = msg.MissionBound["Max"];
                        let min = msg.MissionBound["Min"];
                        let normal = max[0]+max[1]+min[0]+min[1];
                        if(normal>0){
                            Scene.startScene(msg.MissionBound);
                            
                            Scene.addMissionPath(msg.StateMessage.Mission.SwarmGeometry[0],0xff0000,20, drone.id);
                            
                        }
                    }

                    //console.log(msg);
                    drone.mission = msg.StateMessage.Mission.Geometry;
                    drone.senderType = msg.SenderType;
                    drone.batteryLevel = msg.StateMessage.Battery;
                    let pos = {};
                    pos.x = msg.StateMessage.Position.X;
                    pos.y = msg.StateMessage.Position.Y;
                    pos.z = msg.StateMessage.Position.Z;
                    
                    
                    drone.position = pos;
                    
                    break;
                case MissionMessageType:
                    console.log("mission mesg")
                    break;
                case ReorganizationMessageType:
                    console.log("reorg mesg")
                    break;
                case RecalculatorMessageType:
                    console.log("recalc mesg")
                    break;
            }
        
            Scene.updateAgent(drone);
            

        }
        ws.onerror = function(evt) {
            console.log("ERROR: " + evt.data);
        }
        return false;
    };
    
});