package main

import (
	"log"
	"gopkg.in/gographics/imagick.v2/imagick"
	"go-usbmuxd/photobooth"
	"os"
	"bufio"
	"io/ioutil"
	// "os/exec"
	// "fmt"
)

var imageInstance PhotoBooth.Image
var senderInstance PhotoBooth.Sender

func main() {
	log.Println("Upload")

	f, _ := os.Open("here/IMG_8323.JPG")
	reader := bufio.NewReader(f)
	sample, _ := ioutil.ReadAll(reader)
	log.Println("Init")
	imageInstance.Init()
	senderInstance.Init()
	log.Println("Pre Gen")
	imageInstance.SetImage(sample)
	n := imageInstance.GeneratePrintable()
	log.Println("Post Gen")

	imageInstance.Print()
	senderInstance.Upload(n, "guillaume@vasypaulette.com")
	log.Println("Uploaded")

	


	imagick.Initialize()

	// Schedule cleanup
	defer imagick.Terminate()
	var err error
	mw := imagick.NewMagickWand()

	err = mw.ReadImageBlob(n)
	if err != nil {
		panic(err)
	}
	
	mw.WriteImage("generate.jpg")

	/*output, err := exec.Command(mw.GetImagesBlob(), "|", "lpr", "-P", "Dai_Nippon_Printing_DS-RX1", "generate.jpg").CombinedOutput()
	if err != nil {
	os.Stderr.WriteString(err.Error())
	}
	fmt.Println(string(output))*/
}

