package PhotoBooth

import (
	"log"
	"bytes"
	"github.com/micahwedemeyer/gphoto2go"
)


type (
	Camera struct {
		IsConnected bool
	}
)


var camera *gphoto2go.Camera
var imageInstance Image
var senderInstance Sender

func (c *Camera) Init() {
	imageInstance.Init()
	senderInstance.Init()

	camera = new(gphoto2go.Camera)
	err := camera.Init()

	if err != 0 {
		log.Println("CAMERA Cannot Init")
		return
	}

	cameraInstance.GetConfig("captureTarget", "1")
	cameraInstance.SetConfig("captureTarget", "1")

	c.IsConnected = true
}

func (c *Camera) Close() {
	/*if (c.IsConnected)  {
		log.Fatal("CAMERA Cannot Close")
		c.Init()
		return
	}*/

	camera.Exit()
}

func (c *Camera) ReConnect() {
	c.Close()
	c.Init()
}

func (c *Camera) TakePicture() (resizedImage []byte, hasError bool) {
	hasError = false
	if (!c.IsConnected)  {
		log.Println("CAMERA Not Ready")
		c.ReConnect()
		hasError = true
		return
	}

	
	file, err := camera.TriggerCaptureToFile()
	if err != 0 {
		log.Println("CAMERA Cannot TakePicture")
		//c.ReConnect()
		hasError = true
		return
	} 

	log.Println(file)

	cameraFileReader := camera.FileReader(file.Folder, file.Name)
	
	buf := new(bytes.Buffer)
	buf.ReadFrom(cameraFileReader)
	

	/*img, _ := jpeg.Decode(cameraFileReader)

	m := resize.Resize(1400, 0, img, resize.Lanczos3)

	buf := new(bytes.Buffer)
	_ = jpeg.Encode(buf, m, nil)
	send_s3 := buf.Bytes()
	cameraFileReader.Close()*/
	//return buf.Bytes();
	resizedImage = imageInstance.Prepare(buf.Bytes())
	cameraFileReader.Close()
	return
}

func (c *Camera) GeneratePrintable() ([]byte){
	log.Println("Camera.GeneratePrintable")
	return imageInstance.GeneratePrintable()
}

func (c *Camera) Print() {
	imageInstance.Print()
}

func (c *Camera) Send(mail string) (int) {
	return senderInstance.Upload(imageInstance.printableImage, mail)
}

func (c *Camera) Reset() {
	imageInstance.Reset()
}

func (c *Camera) GetConfig(key string, value string) {
	conf, err := camera.GetConfig("capturetarget", value)
	log.Println("confi", conf)
	log.Println("err", err)
}

func (c *Camera) SetConfig(key string, value string) {
	err := camera.SetConfig("capturetarget", value)
	//log.Println("confi", conf)
	log.Println("SetConfind", key, value, err)
}