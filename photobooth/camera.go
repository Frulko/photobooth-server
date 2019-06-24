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

func (c *Camera) Init() {
	camera = new(gphoto2go.Camera)
	err := camera.Init()

	if err != 0 {
		log.Fatal("CAMERA Cannot Init")
		return
	}

	c.IsConnected = true
}

func (c *Camera) Close() {
	if (c.IsConnected)  {
		log.Fatal("CAMERA Cannot TakePicture")
		return
	}

	camera.Exit()
}

func (c *Camera) TakePicture() (resizedImage []byte) {
	if (!c.IsConnected)  {
		log.Fatal("CAMERA Not Ready")
	}

	
	file, err := camera.TriggerCaptureToFile()
	if err != 0 {
		log.Fatal("CAMERA Cannot TakePicture")
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

	resizedImage = imageInstance.Prepare(buf.Bytes())
	cameraFileReader.Close()
	return
}