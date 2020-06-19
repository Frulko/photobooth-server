package main

import (
	// "io/ioutil"
	"log"
	"time"

	"github.com/davidbyttow/govips/pkg/vips"
)

func main() {

	log.Println("[Image] : test")

	defer timeTrack(time.Now(), "process vips")
	t, err := vips.NewImageFromFile("/home/pi/go/src/go-usbmuxd/IMG_7921.JPG")
	
	log.Println(len(t.ToBuffer()))
	if err != nil {
		log.Fatal("ok")
	}
	/*t.Reduce(1200.0, 990.9)

	if err != nil {
		log.Fatal("ok")
	}

	scale := float64(t.Width()) / float64(t.Height())

	log.Println("Width: " + strconv.Itoa(t.Width()) + " -- Height: " + strconv.Itoa(t.Height()) + " // " + fmt.Sprintf("%f", scale))

	t.Resize(2)

	err = ioutil.WriteFile("./IMG_7921_r4.JPG", t.ToBuffer(), 0644)
	if err != nil {
		panic(err)
	}*/

	/*buff, _, err := vips.NewTransform().
		LoadFile("./IMG_7921.JPG").
		// ResizeStrategy(vips.ExtendBlack).
		Resize(1200, 1200).
		OutputFile("./IMG_7921_s5.JPG").
		Apply()
	
	log.Println(len(buff))

	err = ioutil.WriteFile("./IMG_7921_s6.JPG", buff, 0644)
	if err != nil {
		panic(err)
	}*/
}

func timeTrack(start time.Time, name string) {
	elapsed := time.Since(start)
	log.Printf("%s took %s", name, elapsed)
}
