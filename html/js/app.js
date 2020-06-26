

import * as THREE from './three/build/three.module.js';

import { OrbitControls } from './three/examples/jsm/controls/OrbitControls.js';
import Stats from './three/examples/jsm/libs/stats.module.js';
import { GUI } from './three/examples/jsm/libs/dat.gui.module.js';

import { DragControls } from './three/examples/jsm/controls/DragControls.js';
import { TransformControls } from './three/examples/jsm/controls/TransformControls.js'

import * as UnitController from "./controller.js";

String.prototype.format = function () {

		var str = this;

		for ( var i = 0; i < arguments.length; i ++ ) {

			str = str.replace( '{' + i + '}', arguments[ i ] );

		}
		return str;

	};

var container, stats;
var camera, scene, renderer, droneGroup;
var splineHelperObjects = [];
var splinePointsLength = 4;
var positions = [];
var point = new THREE.Vector3();

var drones = {};
var visiblePaths = {};

var geometry = new THREE.BoxBufferGeometry( 20, 20, 20 );
var transformControl;



var raycaster = new THREE.Raycaster();
var intersects;

var mouse = new THREE.Vector3();

document.addEventListener('mousemove', onDocumentMouseMove, false);
document.addEventListener('mousedown', onDocumentMouseDown, false);


function onDocumentMouseMove(event) {
	event.preventDefault();

	/* mouse.x = (event.clientX / window.innerWidth) * 2 - 1;
	mouse.y = -(event.clientY / window.innerHeight) * 2 + 1; */
	mouse.x = (event.clientX / document.getElementById("container").offsetWidth) * 2 - 1;
	mouse.y = -(event.clientY / document.getElementById("container").offsetHeight) * 2 + 1;
	
}
function onDocumentMouseDown(event) {
	event.preventDefault();
  	console.info('mouseDown');

	let ooi = getOOI(intersects);
	if(ooi==null){
		return;
	}
	
	let id = ooi.object.agentId;
	let clickedDrone = drones[id];
	if(clickedDrone.pathVisible){
		//remove path
		clickedDrone.pathVisible = false;
		let path = visiblePaths[id];
		scene.remove(path);
	}else{
		//draw path
		let material = clickedDrone.mesh.material;
		let points = [];
		for(let i in clickedDrone.path){
			let point = new THREE.Vector3(clickedDrone.path[i].x,clickedDrone.path[i].y,clickedDrone.path[i].z);
			points.push(point);
		}
		points.push(clickedDrone.mesh.position);
		var geometry = new THREE.BufferGeometry().setFromPoints( points );
		var line = new THREE.Line( geometry, material );
		visiblePaths[id] = line;

		scene.add( line );

		clickedDrone.pathVisible=true;
	}
  
}

/*var params = {
	uniform: true,
	tension: 0.5,
	centripetal: true,
	chordal: true,
	addPoint: addPoint,
	removePoint: removePoint,
	exportSpline: exportSpline
};*/

function getOOI(intersects){
	if(intersects.length>0){
		for(let i=0; i<intersects.length; i++){
			let obj = intersects[i];
			if(obj.object.name=="OOI"){
				//console.log(obj);
				return obj	
			}
		}
	}
	return null;
}


init();
animate();


function init() {

	container = document.getElementById( 'container' );

	scene = new THREE.Scene();
	scene.background = new THREE.Color( 0xf0f0f0 );

	camera = new THREE.PerspectiveCamera( 70, window.innerWidth / window.innerHeight, 1, 10000 );
	camera.position.set( 0, 250, 1000 );
	scene.add( camera );

	scene.add( new THREE.AmbientLight( 0xf0f0f0 ) );
	var light = new THREE.SpotLight( 0xffffff, 1.5 );
	light.position.set( 0, 1500, 200 );
	light.angle = Math.PI * 0.2;
	light.castShadow = true;
	light.shadow.camera.near = 200;
	light.shadow.camera.far = 2000;
	light.shadow.bias = - 0.000222;
	light.shadow.mapSize.width = 1024;
	light.shadow.mapSize.height = 1024;
	scene.add( light );

	var planeGeometry = new THREE.PlaneBufferGeometry( 2000, 2000 );
	planeGeometry.rotateX( - Math.PI / 2 );
	var planeMaterial = new THREE.ShadowMaterial( { opacity: 0.2 } );

	var plane = new THREE.Mesh( planeGeometry, planeMaterial );
	plane.position.y = - 200;
	plane.receiveShadow = true;
	scene.add( plane );

	var helper = new THREE.GridHelper( 2000, 100 );
	helper.position.y = - 199;
	helper.material.opacity = 0.25;
	helper.material.transparent = true;
	scene.add( helper );

	//var axes = new AxesHelper( 1000 );
	//axes.position.set( - 500, - 500, - 500 );
	//scene.add( axes );

	renderer = new THREE.WebGLRenderer( { antialias: true } );
	renderer.setPixelRatio( window.devicePixelRatio );
	renderer.setSize( window.innerWidth, window.innerHeight );
	renderer.shadowMap.enabled = true;
	container.appendChild( renderer.domElement );

	stats = new Stats();
	container.appendChild( stats.dom );

	droneGroup = new THREE.Group();
	scene.add(droneGroup);


	/*var gui = new GUI();

	gui.add( params, 'uniform' );
	gui.add( params, 'tension', 0, 1 ).step( 0.01 ).onChange( function ( value ) {
		splines.uniform.tension = value;
		updateSplineOutline();
	} );

	gui.add( params, 'centripetal' );
	gui.add( params, 'chordal' );
	gui.add( params, 'addPoint' );
	gui.add( params, 'removePoint' );
	gui.add( params, 'exportSpline' );
	gui.open();*/

	// Controls
	var controls = new OrbitControls( camera, renderer.domElement );
	controls.damping = 0.2;
	controls.addEventListener( 'change', render );
	
}

function animate() {

	requestAnimationFrame( animate );
	render();
	stats.update();

}


function render() {

	raycaster.setFromCamera(mouse, camera);
	intersects = raycaster.intersectObjects(droneGroup.children, true);
	let obj = getOOI(intersects);
	//show info about any ooi:
	let infoHolder = document.getElementById("info");
	if(obj!==null){
		let pos = obj.object.position;
		let content = "Id: "+obj.object.agentId+"<br> Current Position: ("+pos.x+","+pos.y+","+pos.z+")";
		infoHolder.innerHTML=content;
		infoHolder.style.opacity = "1.0"; 
	}else{
		infoHolder.style.opacity = "0.0";
	}
	renderer.render( scene, camera );
}

function updateDrone(data){
	
	if(data.ID in drones){
		let drone = drones[data.ID];
		drone.path.push(JSON.parse(JSON.stringify(drone.mesh.position)));
		drone.mesh.position.x = data.Position.X;
		drone.mesh.position.y = data.Position.Y;
		drone.mesh.position.z = data.Position.Z;
	}else{
		addDrone(data);
	}
}


function addDrone(droneData){
	console.log("drone added");
	
	var width = 40;

	var droneId = droneData.ID;
	// console.log(droneData);
	var hColor = Math.floor(Math.random() * 361);
	var color = new THREE.Color("hsl("+hColor+", 100%, 50%)");
	
	var geometry = new THREE.TetrahedronGeometry(width);
	var material = new THREE.MeshLambertMaterial({ color:color, transparent: true });
	var mesh = new THREE.Mesh( geometry, material ) ;
	material.opacity = 0.6;
	
	
	mesh.position.x = droneData.Position.X;
	mesh.position.y = droneData.Position.Y;
	mesh.position.z = droneData.Position.Z;
	var drone = {"data":droneData,"mesh":mesh,"path":[],"pathVisible":false};
	drones[droneId] = drone;
	mesh.name = "OOI";
	mesh.agentId = droneId

	//scene.add( mesh );
	droneGroup.add( mesh );
}

export {updateDrone}