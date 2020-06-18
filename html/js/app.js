

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
var camera, scene, renderer;
var splineHelperObjects = [];
var splinePointsLength = 4;
var positions = [];
var point = new THREE.Vector3();

var drones = [];

var geometry = new THREE.BoxBufferGeometry( 20, 20, 20 );
var transformControl;

var ARC_SEGMENTS = 200;

var splines = {};

/*var params = {
	uniform: true,
	tension: 0.5,
	centripetal: true,
	chordal: true,
	addPoint: addPoint,
	removePoint: removePoint,
	exportSpline: exportSpline
};*/

init();
animate();
initDrones();

function initDrones(){
	//UnitController.testDrones(5);
}


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
/*
	controls.addEventListener( 'start', function () {
		cancelHideTransform();
	} );

	controls.addEventListener( 'end', function () {
		delayHideTransform();
	} );*/

	/*transformControl = new TransformControls( camera, renderer.domElement );
	transformControl.addEventListener( 'change', render );
	transformControl.addEventListener( 'dragging-changed', function ( event ) {
		controls.enabled = ! event.value;
	} );

	scene.add( transformControl );

	// Hiding transform situation is a little in a mess :()
	transformControl.addEventListener( 'change', function () {
		cancelHideTransform();
	} );

	transformControl.addEventListener( 'mouseDown', function () {
		cancelHideTransform();
	} );

	transformControl.addEventListener( 'mouseUp', function () {
		delayHideTransform();
	} );

	transformControl.addEventListener( 'objectChange', function () {
		updateSplineOutline();
	} );
*/
	/*var dragcontrols = new DragControls( splineHelperObjects, camera, renderer.domElement ); //
	dragcontrols.enabled = false;
	dragcontrols.addEventListener( 'hoveron', function ( event ) {
		transformControl.attach( event.object );
		cancelHideTransform();
	} );

	dragcontrols.addEventListener( 'hoveroff', function () {
		delayHideTransform();
	} );

	var hiding;

	function delayHideTransform() {

		cancelHideTransform();
		hideTransform();

	}

	function hideTransform() {

		hiding = setTimeout( function () {

			transformControl.detach( transformControl.object );

		}, 2500 );

	}

	function cancelHideTransform() {

		if ( hiding ) clearTimeout( hiding );

	}
*/
	/*******
	 * Curves
	 *********/

	/*for ( var i = 0; i < splinePointsLength; i ++ ) {

		addSplineObject( positions[ i ] );

	}

	positions = [];

	for ( var i = 0; i < splinePointsLength; i ++ ) {

		positions.push( splineHelperObjects[ i ].position );

	}

	var geometry = new THREE.BufferGeometry();
	geometry.setAttribute( 'position', new THREE.BufferAttribute( new Float32Array( ARC_SEGMENTS * 3 ), 3 ) );

	var curve = new THREE.CatmullRomCurve3( positions );
	curve.curveType = 'catmullrom';
	curve.mesh = new THREE.Line( geometry.clone(), new THREE.LineBasicMaterial( {
		color: 0xff0000,
		opacity: 0.35
	} ) );
	curve.mesh.castShadow = true;
	splines.uniform = curve;

	curve = new THREE.CatmullRomCurve3( positions );
	curve.curveType = 'centripetal';
	curve.mesh = new THREE.Line( geometry.clone(), new THREE.LineBasicMaterial( {
		color: 0x00ff00,
		opacity: 0.35
	} ) );
	curve.mesh.castShadow = true;
	splines.centripetal = curve;

	curve = new THREE.CatmullRomCurve3( positions );
	curve.curveType = 'chordal';
	curve.mesh = new THREE.Line( geometry.clone(), new THREE.LineBasicMaterial( {
		color: 0x0000ff,
		opacity: 0.35
	} ) );
	curve.mesh.castShadow = true;
	splines.chordal = curve;

	for ( var k in splines ) {

		var spline = splines[ k ];
		scene.add( spline.mesh );

	}

	load( [ new THREE.Vector3( 289.76843686945404, 452.51481137238443, 56.10018915737797 ),
		new THREE.Vector3( - 53.56300074753207, 171.49711742836848, - 14.495472686253045 ),
		new THREE.Vector3( - 91.40118730204415, 176.4306956436485, - 6.958271935582161 ),
		new THREE.Vector3( - 383.785318791128, 491.1365363371675, 47.869296953772746 ) ] );*/

}





/*			function addSplineObject( position ) {

	var material = new THREE.MeshLambertMaterial( { color: Math.random() * 0xffffff } );
	var object = new THREE.Mesh( geometry, material );

	if ( position ) {

		object.position.copy( position );

	} else {

		object.position.x = Math.random() * 1000 - 500;
		object.position.y = Math.random() * 600;
		object.position.z = Math.random() * 800 - 400;

	}

	object.castShadow = true;
	object.receiveShadow = true;
	scene.add( object );
	splineHelperObjects.push( object );
	return object;

}*/

/*			function addPoint() {

	splinePointsLength ++;

	positions.push( addSplineObject().position );

	updateSplineOutline();

}*/

/*		function removePoint() {

	if ( splinePointsLength <= 4 ) {

		return;

	}
	splinePointsLength --;
	positions.pop();
	scene.remove( splineHelperObjects.pop() );

	updateSplineOutline();

}*/

/*		function updateSplineOutline() {

	for ( var k in splines ) {

		var spline = splines[ k ];

		var splineMesh = spline.mesh;
		var position = splineMesh.geometry.attributes.position;

		for ( var i = 0; i < ARC_SEGMENTS; i ++ ) {

			var t = i / ( ARC_SEGMENTS - 1 );
			spline.getPoint( t, point );
			position.setXYZ( i, point.x, point.y, point.z );

		}

		position.needsUpdate = true;

	}

}*/

/*			function exportSpline() {

	var strplace = [];

	for ( var i = 0; i < splinePointsLength; i ++ ) {

		var p = splineHelperObjects[ i ].position;
		strplace.push( 'new THREE.Vector3({0}, {1}, {2})'.format( p.x, p.y, p.z ) );

	}

	console.log( strplace.join( ',\n' ) );
	var code = '[' + ( strplace.join( ',\n\t' ) ) + ']';
	prompt( 'copy and paste code', code );

}

function load( new_positions ) {

	while ( new_positions.length > positions.length ) {

		addPoint();

	}

	while ( new_positions.length < positions.length ) {

		removePoint();

	}

	for ( var i = 0; i < positions.length; i ++ ) {

		positions[ i ].copy( new_positions[ i ] );

	}

	updateSplineOutline();

}*/

function animate() {

	requestAnimationFrame( animate );
	render();
	stats.update();

}

function render() {
	//console.log(renderer.info.render);
	/*splines.uniform.mesh.visible = params.uniform;
	splines.centripetal.mesh.visible = params.centripetal;
	splines.chordal.mesh.visible = params.chordal;*/
	renderer.render( scene, camera );

}

function updateDrone(droneId, dx,dy,dz){
	
	for(let pos in drones){
		let drone = drones[pos];
		if(drone.data.ID==droneId){
			drone.mesh.position.x+=dx;
			drone.mesh.position.y+=dy;
			drone.mesh.position.z+=dz;
		}
	}
}


function addDrone(droneData){
	console.log("drone added");
	console.log(droneData);
	var width = 4;


	var droneId = droneData.ID;
	// console.log(droneData);
	var hColor = Math.floor(Math.random() * 361);
	var color = new THREE.Color("hsl("+hColor+", 100%, 50%)");
	
	var geometry = new THREE.TetrahedronGeometry(width);
	var material = new THREE.MeshLambertMaterial({ color:color, transparent: true });
	var mesh = new THREE.Mesh( geometry, material ) ;
	material.opacity = 0.6;
	
	//mesh.position.x = mesh.position.x+(droneId*(width+10));
	mesh.position.x = droneData.Position;
	var drone = {"data":droneData,"mesh":mesh};
	drones.push(drone);
	

	scene.add( mesh );
}

export {addDrone,updateDrone}