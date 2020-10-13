package PhotoBooth

import (
	"gopkg.in/gographics/imagick.v2/imagick"
	"log"
	"bytes"
	"image"
	"image/draw"
	"image/jpeg"
	"os"
	"os/exec"
	"strings"
	"strconv"
)

type (
	Image struct {
		imageBytes []byte
		printableImage []byte
		backgroundImage image.Image
	}
)

 
func (i *Image) Init() {
	log.Println("init gabarit")
	filebg, err := os.Open("/home/pi/go/src/go-usbmuxd/gabarit.jpg")
	if (err != nil) {
		log.Println(err)
	}
	loaded_image, _, _ := image.Decode(filebg)
	i.backgroundImage = loaded_image
}

func (i *Image) SetImage(blob []byte) {
	i.imageBytes = blob
}

func (i *Image) Prepare(blob []byte) ([]byte){
	imagick.Initialize()

	i.imageBytes = blob

	// Schedule cleanup
	defer imagick.Terminate()
	var err error

	mw := imagick.NewMagickWand()
	// mw.SetColorspace(imagick.COLORSPACE_GRAY)
	err = mw.ReadImageBlob(blob)
	if err != nil {
		log.Println("prepare ReadImageBlob")
		log.Println(err)
	}

	// Get original logo size
	width := mw.GetImageWidth()
	height := mw.GetImageHeight()

	// Calculate half the size
	hWidth := uint(width / 2)
	hHeight := uint(height / 2)


	target := imagick.NewPixelWand()
	target.SetColor("white")

	// mw.SetColorspace(imagick.COLORSPACE_GRAY)
	

	log.Println("--->")
	// log.Println(imagick.COLORSPACE_GRAY)
	log.Println("<---")
	/*err = mw.RotateImage(target, -90.0)
	if err != nil {
		panic(err)
	}*/
	// Resize the image using the Lanczos filter
	// The blur factor is a float, where > 1 is blurry, < 1 is sharp
	err = mw.ResizeImage(hWidth, hHeight, imagick.FILTER_LANCZOS, 1)
	if err != nil {
		log.Println("prepare.ResizeImage")
		log.Println(err)
	}
	
	// mw.ModulateImage(100, 0, 100) // change to black n white only for Hexagone
	
	// Set the compression quality to 95 (high quality = low compression)
	
	err = mw.SetImageCompressionQuality(70)
	if err != nil {
		log.Println("prepare.SetImageCompressionQuality")
		log.Println(err)
	}

	return mw.GetImagesBlob()
}

func (i *Image) GeneratePrintable() ([]byte){
	log.Println("generate", len(i.imageBytes))
	imagick.Initialize()

	image_is_black_and_white := false
	// image_output_size := [2]int{
	// 	2720, // width
	// 	4080, // height
	// }
	image_output_size := [2]int{
		4080, // width
		2720, // height
	}
	image_rotation := float64(0.0) // in degrees closewise
	image_compression := uint(95)
	brightness_contrast := float64(0.0) // int
	// picture_size := [2]int{
	// 	2889, // width
	// 	1926, // height
	// }
	picture_size := [2]int{
		3076, // height
		2049, // width
	}
	// picture_position := [2]int{
	// 	397, // x position
	// 	359, // y position
	// }
	// picture_position := [2]int{
	// 	499, // x position
	// 	325, // y position
	// }
		picture_position := [2]int{
			379, // x position
			205, // y position
		}

	// Schedule cleanup
	defer imagick.Terminate()
	var err error

	mw := imagick.NewMagickWand()


	err = mw.ReadImageBlob(i.imageBytes)
	if err != nil {
		log.Println("generate.ReadImageBlob")
		log.Println(err)
	}

	hWidth := uint(picture_size[0])
	hHeight := uint(picture_size[1])
	

	target := imagick.NewPixelWand()
	// target.SetColor("white")
	

	// Resize the image using the Lanczos filter
	// The blur factor is a float, where > 1 is blurry, < 1 is sharp
	err = mw.ResizeImage(hWidth, hHeight, imagick.FILTER_LANCZOS, 1)
	if err != nil {
		log.Println("ResizeImage")
		log.Println(err)
	}

	if (image_is_black_and_white) {
		// mw.SetColorspace(imagick.COLORSPACE_GRAY)
		// mw.ModulateImage(100, 0, 100) // change to black n white only for Hexagone
	}


	mw.BrightnessContrastImage(0, brightness_contrast) // force contrast for printing

	err = mw.RotateImage(target, image_rotation)
	if err != nil {
		log.Println("RotateImage")
		log.Println(err)
	}

	// Set the compression quality to 95 (high quality = low compression)2889 * 800 
	err = mw.SetImageCompressionQuality(image_compression)
	if err != nil {
		log.Println("SetImageCompressionQuality")
		log.Println(err)
	}



	canvas := image.NewRGBA(image.Rect(0, 0, image_output_size[0], image_output_size[1]))
	
	foreground_image, _, err := image.Decode(bytes.NewReader(mw.GetImagesBlob()))

	draw.Draw(canvas, canvas.Bounds(), i.backgroundImage, image.Point{0, 0}, draw.Src)
	draw.Draw(canvas, canvas.Bounds(), foreground_image, image.Point{-picture_position[0], -picture_position[1]}, draw.Src)

	buf := new(bytes.Buffer)
	_ = jpeg.Encode(buf, canvas, nil)

	i.printableImage = buf.Bytes();

	return i.printableImage
}

func (i *Image) Print() {
	log.Println("Print", len(i.printableImage))

	if (len(i.printableImage) < 1) {
		log.Println("Need to prepare image", len(i.printableImage))
		return
	}

	currentCount := i.GetCounter()
	newCount := currentCount + 1

	log.Println(currentCount, "< - OLD && NEW - >", newCount)

	i.SetCounter(strconv.Itoa(currentCount))

	imagick.Initialize()

	// Schedule cleanup
	defer imagick.Terminate()
	var err error
	mw := imagick.NewMagickWand()

	err = mw.ReadImageBlob(i.printableImage)
	if err != nil {
		log.Println("print.ReadImageBlob")
		log.Println(err)
	}
	
	mw.WriteImage("generate.jpg")
	log.Println("Print IS OK")
	output, err := exec.Command("lpr", "-P", "Dai_Nippon_Printing_DS-RX1", "generate.jpg").CombinedOutput()
	if err != nil {
		log.Println(err)
	}
	log.Println(string(output))

	// output2, err2 := exec.Command("lpr", "-P", "Dai_Nippon_Printing_DS-RX1", "generate.jpg").CombinedOutput()
	// if err2 != nil {
	// 	log.Println(err2)
	// }
	// log.Println(string(output2))
}

func (i *Image) GetCounter() (int){

	output, err := exec.Command("/home/pi/go/src/go-usbmuxd/dnpds40", "-n").CombinedOutput()
	if err != nil {
		log.Println("error", err)
		return 1000
	}

	m := strings.Split(string(output), "\n")
	if (len(m) > 1 && len(m) > 14) {
		g := strings.Split(m[14], ":")

		val := strings.TrimLeft(g[2], " ")
		intVal, _ := strconv.Atoi(val)
		return intVal
	} else {
		return 1000
	}
}	
func (i *Image) GetTemp() (string){

	output, err := exec.Command("/opt/vc/bin/vcgencmd", "measure_temp").CombinedOutput()
	if err != nil {
		log.Println("error", err)
		return ""
	}

	return string(output)
}	

func (i *Image) SetCounter(number string){
	log.Println("SetCounter", number)

	exec.Command("/home/pi/go/src/go-usbmuxd/dnpds40", "-p", number).CombinedOutput()
}	


func (i *Image) Reset() {
	i.imageBytes = make([]byte, 0)
	i.printableImage = make([]byte, 0)
	log.Println("Reset", len(i.imageBytes))
}

