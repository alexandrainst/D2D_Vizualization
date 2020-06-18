
import * as Scene from "./app.js";

function testDrones(number){
	console.log("Starting "+number+" drones");
	for(let x =0; x<number; x++){
		let droneData = {
			"id":x
		};
		Scene.addDrone(droneData);
	}

	setInterval(function(){ updateDrone(number); },100);
}

function updateDrone(number){
	for(let i =0; i<number; i++){
		let x =  Math.round(Math.random());
		let y = Math.round(Math.random());
		let z = Math.round(Math.random());
		Scene.updateDrone(i,x,y,z);
	}
}

export { testDrones};

window.addEventListener("load", function(evt) {
    var output = document.getElementById("output");
    var input = document.getElementById("input");
    var ws;
    var print = function(message) {
        var d = document.createElement("div");
        d.textContent = message;
        output.appendChild(d);
    };

    //document.getElementById("open").onclick = 
    connectWS();

    function connectWS(evt) {
		console.log("OPEN MG");
        if (ws) {
            return false;
        }
        ws = new WebSocket("ws://localhost:8080/echo");
        console.log(ws);
        ws.onopen = function(evt) {
            print("OPEN");
        }
        ws.onclose = function(evt) {
            print("CLOSE");
            ws = null;
        }
        ws.onmessage = function(evt) {
            //print("RESPONSE: " + evt.data);
            console.log(evt.data)
            var drone = JSON.parse(evt.data);
            console.log(drone)
        }
        ws.onerror = function(evt) {
            print("ERROR: " + evt.data);
        }
        return false;
    };
    document.getElementById("send").onclick = function(evt) {
        if (!ws) {
            return false;
        }
        print("SEND: " + input.value);
        ws.send(input.value);
        return false;
    };
    document.getElementById("close").onclick = function(evt) {
        if (!ws) {
            return false;
        }
        ws.close();
        return false;
    };
});