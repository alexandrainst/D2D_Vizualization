package main

import "log"

var optionsDimensions = [4]int{0, 1, 2, 3}
var optionsHardware = [2]string{"camera", "human detection"}

//key: index in dimension
//value: number of agents needed with this
var neededDimensions = map[int]int{
	0: 1,
	4: 4,
}

//key: dimension
//value: hardware index
var neededHardware = map[int]int{
	0: 1,
	3: 0,
}

//key: dimension
//value number created
var fullfilledActions = map[int]int{
	0: 0,
	3: 0,
}

func getAddress() string {
	return "URL"
}

func getPosition() Vector {
	return Vector{0, 0, 0}
}

func getPublicKey() string {
	return "Key"
}

func getBatteryLevel() int {
	return 100
}

func getDimensions() int {
	for dimension, needed := range neededDimensions {
		status := fullfilledActions[dimension]
		if status < needed {
			return dimension

		}
	}
	log.Println("Something is wrong! All dimension needs fullfilled")
	return -1
}

func getOnbardHardware(numberOfDimensions int) string {
	hardwareIndex := neededHardware[numberOfDimensions]
	return optionsHardware[hardwareIndex]
}

//GetMetadataForAgent returns all the meta needed
func GetMetadataForAgent() (string, Vector, string, int, int, string) {

	var dimensions = getDimensions()
	var hardware = getOnbardHardware(dimensions)
	fullfilledActions[dimensions]++

	return getAddress(), getPosition(), getPublicKey(), getBatteryLevel(), dimensions, hardware

}
