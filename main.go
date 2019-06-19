package main

import (
	"bufio"
	"bytes"
	"fmt"
	_ "fmt"
	"go-usbmuxd/USB"
	"go-usbmuxd/frames"
	"go-usbmuxd/transmission"
	"image/jpeg"
	_ "image/jpeg"
	"io"
	"io/ioutil"
	"log"
	"os"

	"github.com/micahwedemeyer/gphoto2go"
	"github.com/nfnt/resize"
	_ "github.com/nfnt/resize"
)

// some global vars
var connectHandle USB.ConnectedDevices
var port = 4986
var pluggedUSBDevices map[int]frames.USBDeviceAttachedDetachedFrame
var connectedUSB int // only stores the device id
var scanningInstance USB.Scan
var self USBDeviceDelegate

func main() {

	// inti section
	connectedUSB = -1
	pluggedUSBDevices = map[int]frames.USBDeviceAttachedDetachedFrame{}
	scanningInstance = USB.Scan{}
	self = USBDeviceDelegate{}

	// logger
	logFile, err := os.OpenFile("kusb_ios.log", os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	if err != nil {
		log.Println(err)
	}
	defer logFile.Close()
	mw := io.MultiWriter(os.Stdout, logFile)
	log.SetOutput(mw)

	// create a USB.Listen(USBDeviceDelegate) instance. Pass a delegate to resolve the attached and detached callbacks
	// then on device added save ot to array/ map and send connect to a port with proper tag
	listenConnection := USB.Listen(transmission.Tunnel(), self)
	defer listenConnection.Close()

	// connect to a random usb device, if Number == 0 then
	connectHandle = USB.ConnectedDevices{Delegate: self, Connection: transmission.Tunnel()}
	defer connectHandle.Connection.Close()

	// scan defer
	defer scanningInstance.Stop()

	// run loop
	select {}
}

// USBDeviceDelegate - USB Delegate Methods
type USBDeviceDelegate struct{}

// USBDeviceDidPlug - device plugged callback
func (usb USBDeviceDelegate) USBDeviceDidPlug(frame frames.USBDeviceAttachedDetachedFrame) {
	// usb has been plugged DO: startScanning
	log.Printf("[USB-INFO] : Device Plugged %s ID: %d\n", frame.Properties.SerialNumber, frame.DeviceID)
	pluggedUSBDevices[frame.DeviceID] = frame
	scanningInstance.Start(&connectHandle, frame, port)
}

// USBDeviceDidUnPlug - device unplugged callback
func (usb USBDeviceDelegate) USBDeviceDidUnPlug(frame frames.USBDeviceAttachedDetachedFrame) {
	// usb has been unplugged
	// stop scan
	log.Printf("[USB-INFO] : Device UnPlugged %s ID: %d\n", pluggedUSBDevices[frame.DeviceID].Properties.SerialNumber, frame.DeviceID)
	delete(pluggedUSBDevices, frame.DeviceID)
	scanningInstance.Stop()
}

// USBDidReceiveErrorWhilePluggingOrUnplugging - device plugging/unplugging callback
func (usb USBDeviceDelegate) USBDidReceiveErrorWhilePluggingOrUnplugging(err error, stringResponse string) {
	// plug or unplug error
	// stop scan
	if stringResponse != "" {
		//some unresolved message came
		//TODO - Implement some resolver to understand message received
	}
	log.Println("[USB-EM-1] : Some error encountered wile pluging and unpluging. ", err.Error())
	scanningInstance.Stop()
}

// USBDeviceDidSuccessfullyConnect - device successful connection callback
func (usb USBDeviceDelegate) USBDeviceDidSuccessfullyConnect(device USB.ConnectedDevices, deviceID int, toPort int) {
	// successfully connected to the port mentioned
	// stop the scan
	connectedUSB = deviceID

	f, _ := os.Open("icon.png") // For read access.
	// decode jpeg into image.Image
	/*img, err := jpeg.Decode(f)
	if err != nil {
		log.Fatal(err)
	}
	f.Close()

	m := resize.Resize(1400, 0, img, resize.Lanczos3)

	buf := new(bytes.Buffer)
	_ = jpeg.Encode(buf, m, nil)
	send_s3 := buf.Bytes()*/

	reader := bufio.NewReader(f)
	content, _ := ioutil.ReadAll(reader)

	// Encode as base64.
	//encoded := base64.StdEncoding.EncodeToString(content)

	/*file, err := os.Open("icon.png") // For read access.
	if err != nil {
		log.Fatal(err)
	}
	data := make([]byte, 100)
	count, err := file.Read(data)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("read %d bytes: %q\n", count, data[:count])*/

	// Print encoded data to console.
	// ... The base64 image can be used as a data URI in a browser.
	//fmt.Println("ENCODED: " + encoded)
	//fmt.Printf("read %d bytes: %q\n", count, data[:count])
	//b := []byte(encoded)
	//data := make([]byte, 20)

	device.SendData(content[0:], 101)
	scanningInstance.Stop()
}

// USBDeviceDidFailToConnect - device connection failure callback
func (usb USBDeviceDelegate) USBDeviceDidFailToConnect(device USB.ConnectedDevices, deviceID int, toPort int, err error) {
	// error while communication in the socket
	// start scan
	connectedUSB = -1
	pluggedDeviceID := getFirstPluggedDeviceId()
	if pluggedDeviceID != -1 {
		scanningInstance.Start(&connectHandle, pluggedUSBDevices[pluggedDeviceID], port)
	}

}

// USBDeviceDidReceiveData - data received callback
func (usb USBDeviceDelegate) USBDeviceDidReceiveData(device USB.ConnectedDevices, deviceID int, messageTAG uint32, data []byte) {
	// received data from the device
	command := string(data[16:])

	switch command {
	case "/photo/capture":

		fmt.Println("Trigger Capture and Send")

		camera := new(gphoto2go.Camera)
		camera.Init()
		file, _ := camera.TriggerCaptureToFile()
		cameraFileReader := camera.FileReader(file.Folder, file.Name)

		img, _ := jpeg.Decode(cameraFileReader)

		m := resize.Resize(1400, 0, img, resize.Lanczos3)

		buf := new(bytes.Buffer)
		_ = jpeg.Encode(buf, m, nil)
		send_s3 := buf.Bytes()
		cameraFileReader.Close()
		camera.Exit()

		device.SendData(send_s3[0:], 101)

	case "/photo/print":
		fmt.Println("Print the photo")
	case "/photo/cancel":
		fmt.Println("Cancel Photo and restart")
		f, _ := os.Open("icon.png")
		reader := bufio.NewReader(f)
		content, _ := ioutil.ReadAll(reader)
		device.SendData(content[0:], 101)
	default:
		fmt.Println("Unknown command")
		fmt.Println("Cancel Photo and restart")
		f, _ := os.Open("IMG_6954.JPG")
		reader := bufio.NewReader(f)
		content, _ := ioutil.ReadAll(reader)
		device.SendData(content[0:], 101)
	}

	if command == "/photo/capture" {
		fmt.Println("OK")
	}

	fmt.Println(command, messageTAG)
	//device.SendData(data[20:], 106)
}

// USBDeviceDidDisconnect - device disconnect callback
func (usb USBDeviceDelegate) USBDeviceDidDisconnect(devices USB.ConnectedDevices, deviceID int, toPort int) {
	// socket disconnect
	// start scan
	connectedUSB = -1
	pluggedDeviceID := getFirstPluggedDeviceId()
	if pluggedDeviceID != -1 {
		scanningInstance.Start(&connectHandle, pluggedUSBDevices[pluggedDeviceID], port)
	}
}

// MARK - helper functions here
// Needs restructuring, removal or other implementation
func getFirstPluggedDeviceId() int {
	var deviceID int = -1
	for deviceID, _ = range pluggedUSBDevices {
		break
	}
	return deviceID
}
