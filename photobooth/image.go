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
	filebg, _ := os.Open("./gabarit.jpg")
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

	err = mw.ReadImageBlob(blob)
	if err != nil {
		panic(err)
	}

	// Get original logo size
	width := mw.GetImageWidth()
	height := mw.GetImageHeight()

	// Calculate half the size
	hWidth := uint(width / 2)
	hHeight := uint(height / 2)


	target := imagick.NewPixelWand()
	target.SetColor("white")

	/*err = mw.RotateImage(target, -90.0)
	if err != nil {
		panic(err)
	}*/
	// Resize the image using the Lanczos filter
	// The blur factor is a float, where > 1 is blurry, < 1 is sharp
	err = mw.ResizeImage(hWidth, hHeight, imagick.FILTER_LANCZOS, 1)
	if err != nil {
		panic(err)
	}
	

	// Set the compression quality to 95 (high quality = low compression)
	err = mw.SetImageCompressionQuality(70)
	if err != nil {
		panic(err)
	}

	return mw.GetImagesBlob()
}

func (i *Image) GeneratePrintable() ([]byte){
	log.Println("generate", len(i.imageBytes))
	imagick.Initialize()

	// Schedule cleanup
	defer imagick.Terminate()
	var err error

	mw := imagick.NewMagickWand()


	err = mw.ReadImageBlob(i.imageBytes)
	if err != nil {
		panic(err)
	}

	hWidth := uint(2889)
	hHeight := uint(1926)
	

	target := imagick.NewPixelWand()
	target.SetColor("white")

	// Resize the image using the Lanczos filter
	// The blur factor is a float, where > 1 is blurry, < 1 is sharp
	err = mw.ResizeImage(hWidth, hHeight, imagick.FILTER_LANCZOS, 1)
	if err != nil {
		panic(err)
	}

	err = mw.RotateImage(target, 90.0)
	if err != nil {
		panic(err)
	}

	// Set the compression quality to 95 (high quality = low compression)2889 * 800 
	err = mw.SetImageCompressionQuality(95)
	if err != nil {
		panic(err)
	}



	canvas := image.NewRGBA(image.Rect(0, 0, 2720, 4080))
	
	foreground_image, _, err := image.Decode(bytes.NewReader(mw.GetImagesBlob()))

	draw.Draw(canvas, canvas.Bounds(), i.backgroundImage, image.Point{0, 0}, draw.Src)
	draw.Draw(canvas, canvas.Bounds(), foreground_image, image.Point{-397, -359}, draw.Src)

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

	i.SetCounter(strconv.Itoa(newCount))

	imagick.Initialize()

	// Schedule cleanup
	defer imagick.Terminate()
	var err error
	mw := imagick.NewMagickWand()

	err = mw.ReadImageBlob(i.printableImage)
	if err != nil {
		panic(err)
	}
	
	mw.WriteImage("generate.jpg")
	log.Println("Print IS OK")
	/*output, err := exec.Command("lpr", "-P", "Dai_Nippon_Printing_DS-RX1", "generate.jpg").CombinedOutput()
	if err != nil {
	os.Stderr.WriteString(err.Error())
	}
	fmt.Println(string(output))*/
}

func (i *Image) GetCounter() (int){
	log.Println("Counter")

	output, err := exec.Command("sudo", "./dnpds40", "-n").CombinedOutput()
	if err != nil {
		log.Println("error")
		return 0
	}

	m := strings.Split(string(output), "\n")
	if (len(m) > 1 && len(m) > 14) {
		g := strings.Split(m[14], ":")

		val := strings.TrimLeft(g[2], " ")
		intVal, _ := strconv.Atoi(val)
		return intVal
	} else {
		return 0
	}
}	

func (i *Image) SetCounter(number string){
	log.Println("SetCounter", number)

	exec.Command("sudo", "./dnpds40", "-p", number).CombinedOutput()
}	


func (i *Image) Reset() {
	i.imageBytes = make([]byte, 0)
	i.printableImage = make([]byte, 0)
	log.Println("Reset", len(i.imageBytes))
}