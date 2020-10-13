package PhotoBooth

import (
	"log"
	"encoding/base64"
	"github.com/gorilla/websocket"
)


var cameraInstance Camera



type (
	WebSocket struct {
		IsConnected bool
		email string
	}
)



func (ws *WebSocket) Init() {
	cameraInstance.Init();
	ws.email = ""
}

func (ws *WebSocket) ReConnect() {
	cameraInstance.ReConnect();
}


func (ws *WebSocket) OnConnection(socket *websocket.Conn) {
	log.Println("New OnConnection:")
	

	/*f, _ := os.Open("icon.png")
	reader := bufio.NewReader(f)
	content, _ := ioutil.ReadAll(reader)
	err := socket.WriteMessage(2, content)
	if err != nil {
		log.Println("write:", err)
		return
	}*/
}


func (ws *WebSocket) OnDisconnection(socket *websocket.Conn) {
	log.Println("New OnDisconnection:")
	

}

func (ws *WebSocket) OnMessage(command string, socket *websocket.Conn) {
	

	

	log.Println("[", command, "]")
		
		
	/*f, _ := os.Open("icon.png") // For read access.
		

	reader := bufio.NewReader(f)
	content, _ := ioutil.ReadAll(reader)

	log.Println("bytes", len(content))
	f.Close()*/

	

	/*err = socket.WriteMessage(2, content)
	if err != nil {
		log.Println("write:", err)
		return
	}*/


	/*command := string(message)

	log.Println("type", mt)*/
	//log.Println("recv", command)

	switch command {
		case "/photo/capture":
			cameraInstance.Reset()
			log.Println("-- Trigger Capture and Send")
			picture, hasError := cameraInstance.TakePicture();

			if (hasError) {
				ws.sendMessage("photo:error", socket)
				return
			}

			//picture := cameraInstance.GeneratePrintable()

			log.Println("-- Send")
			encoded := base64.StdEncoding.EncodeToString([]byte(picture))
			ws.sendMessage(encoded, socket)
			// err := socket.WriteMessage(2, picture)
			// if err != nil {
			// 	log.Println("write:", err)
			// 	return
			// }

			//cameraInstance.TakePicture();
			//device.SendData(picture[0:], 101)
		
		case "/photo/prepare": 
			cameraInstance.GeneratePrintable()
			ws.sendMessage("photo:ready", socket)
			

		case "/photo/print":
			log.Println("-- Print the photo")
			cameraInstance.Print()
		
		case "/photo/send":
			log.Println("-- Send the photo")
			statusCode := cameraInstance.Send(ws.email)
			if (statusCode == 200) {
				ws.sendMessage("photo:send", socket)
			} else {
				ws.sendMessage("photo:senderror", socket)
			}
			
		case "/photo/can/print":
			count := imageInstance.GetCounter()

			log.Println("-- Can print", count)
			if (count >= 700) {
				ws.sendMessage("photo:canPrint:0", socket)
			} else {
				ws.sendMessage("photo:canPrint:1", socket)
			}

		case "/photo/cancel":
			log.Println("-- Cancel Photo and restart")
			ws.email = ""
			cameraInstance.Reset()
			/*f, _ := os.Open("icon.png")
			reader := bufio.NewReader(f)
			_, _ := ioutil.ReadAll(reader)*/
			//device.SendData(content[0:], 101)
		default:
			log.Println("-- Unknown command set mail")
			//log.Println("Cancel Photo and restart")

			if (command == "reset@reset.fr") {
				cameraInstance.Reset()
				imageInstance.SetCounter("0");
				ws.sendMessage("photo:reset", socket)
			}

			ws.email = command
			/*f, _ := os.Open("IMG_6956.JPG_resized.jpg")
			reader := bufio.NewReader(f)
			_, _ := ioutil.ReadAll(reader)*/
			//device.SendData(content[0:], 101)
	}
}


func (ws *WebSocket) sendMessage(message string, socket *websocket.Conn) {
	err := socket.WriteMessage(1, []byte(message))
	if err != nil {
		log.Println("write:", err)
		return
	}
}