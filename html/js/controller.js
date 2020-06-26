
import * as Scene from "./app.js";

/* function testDrones(number){
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

export { testDrones}; */

window.addEventListener("load", function(evt) {
    var ws;

    //document.getElementById("open").onclick = 
    connectWS();
    /* var tmp = '{"ID":1,"Position":{"X":-0,"Y":-0,"Z":18}}';
    Scene.updateDrone(JSON.parse(tmp));
    var tmp = '{"ID":1,"Position":{"X":10,"Y":10,"Z":50}}';
    Scene.updateDrone(JSON.parse(tmp));
    var tmp = '{"ID":1,"Position":{"X":-20,"Y":-50,"Z":100}}';
    Scene.updateDrone(JSON.parse(tmp));
    var tmp = '{"ID":1,"Position":{"X":124,"Y":0,"Z":-200}}';
    Scene.updateDrone(JSON.parse(tmp)); */

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
        }
        ws.onmessage = function(evt) {
            //print("RESPONSE: " + evt.data);
            console.log(evt.data)
            var drone = JSON.parse(evt.data);
            
            Scene.updateDrone(drone);

        }
        ws.onerror = function(evt) {
            print("ERROR: " + evt.data);
        }
        return false;
    };
    
});