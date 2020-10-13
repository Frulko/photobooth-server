package main

import (
	"log"
	"os"
	"go-usbmuxd/photobooth"
	// "github.com/davidbyttow/govips/pkg/vips"
)

var imageInstance PhotoBooth.Image

func main() {
	log.Println("[Image] : test")


	
	file, err := os.Open("/home/pi/go/src/go-usbmuxd/here/IMG_8323.JPG")
	if (err != nil) {
		log.Println(err)
	}

	defer file.Close()

	fileinfo, err := file.Stat()
	if err != nil {
		log.Println(err)
		return
	}

	filesize := fileinfo.Size()
	buffer := make([]byte, filesize)

	
	log.Println(len(buffer))

	imageInstance.Init();
	
	imageInstance.Prepare(buffer);
	imageInstance.GeneratePrintable();
	log.Println("ready to Print");
}