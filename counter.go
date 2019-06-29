package main

import (
	"log"
	"go-usbmuxd/photobooth"
)

var cameraInstance PhotoBooth.Camera

func main() {
	log.Println("Counter")

	cameraInstance.Init()
	cameraInstance.TakePicture()
	cameraInstance.GetConfig("captureTarget", "1")
	
}

