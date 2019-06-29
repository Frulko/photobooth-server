package main

import (
	"flag"
	"html/template"
	"log"
	"net/http"
    "github.com/gorilla/websocket"
    "go-usbmuxd/photobooth"
    "strconv"
)

var addr = flag.String("addr", "0.0.0.0:8080", "http service address")
var upgrader = websocket.Upgrader{} // use default options
var wsInstance PhotoBooth.WebSocket
var imageInstance PhotoBooth.Image

func echo(w http.ResponseWriter, r *http.Request) {
    c, err := upgrader.Upgrade(w, r, nil)
    
    
    wsInstance.OnConnection(c);
    

    defer func() {
        log.Println("hard disconnect")
    }()

	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer c.Close()
	for {
        _, message, err := c.ReadMessage()
        if err != nil {
            log.Println("Connection closed")
            wsInstance.OnDisconnection(c);
            return
        }
        command := string(message)
        //log.Println("Connection", command)
        wsInstance.OnMessage(command, c);
	}
}

type HomeData struct {
	TempCPU string
	PrinterCounter string
}

func home(w http.ResponseWriter, r *http.Request) {

    data := HomeData{
        PrinterCounter: strconv.Itoa(imageInstance.GetCounter() - 29),
        TempCPU: imageInstance.GetTemp(),
    }

	homeTemplate.Execute(w, data)
}

func main() {

    log.Println("Websocket server started")
    wsInstance.Init();
	flag.Parse()
	log.SetFlags(0)
	http.HandleFunc("/echo", echo)
	http.HandleFunc("/", home)
	log.Fatal(http.ListenAndServe(*addr, nil))
}

var homeTemplate = template.Must(template.New("").Parse(`
<!DOCTYPE html>
<html>
<head>
<meta charset="utf-8">
</head>
<body>
<ul>
<li>Compteur impression : {{.PrinterCounter}}</li>
<li>Température RaspberryPi : {{.TempCPU}}</li>
<ul>
</body>
</html>
`))
