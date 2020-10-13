package main

import (
	"flag"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"github.com/gorilla/websocket"
)

var addr = flag.String("addr", "0.0.0.0:8080", "http service address")
var upgrader = websocket.Upgrader{} // use default options

func echo(w http.ResponseWriter, r *http.Request) {
		upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	
    c, err := upgrader.Upgrade(w, r, nil)
  
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
            return
        }
        command := string(message)
        log.Println("Connection", command)
	}
}

type HomeData struct { 
	TempCPU string
	PrinterCounter string
}

func home(w http.ResponseWriter, r *http.Request) {

    data := HomeData{
        PrinterCounter: "12",
        TempCPU: "24",
    }

	homeTemplate.Execute(w, data)
}

func reset(w http.ResponseWriter, r *http.Request) {
    log.Println("Reset")

    w.Write([]byte("Reset"))
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


func websocketClientConnection() {
	u := url.URL{Scheme: "ws", Host: "54.36.98.177:8080", Path: "/"}
	log.Printf("connecting to %s", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer c.Close()

	done := make(chan struct{})

	go func() {
		defer close(done)
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				return
			}
			log.Printf("recv: %s", message)
		}
	}()
}

func main() {

	

	log.Println("Websocket server started")
	flag.Parse()
	log.SetFlags(0)
	http.HandleFunc("/echo", echo)
	http.HandleFunc("/reset", reset)
	http.HandleFunc("/", home)

	websocketClientConnection();
	
	log.Fatal(http.ListenAndServe(*addr, nil))

	
}