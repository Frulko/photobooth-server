
var exampleSocket = new WebSocket("ws://10.0.0.2:8080/echo");

exampleSocket.onopen = function (event) {
  exampleSocket.send("Voici un texte que le serveur attend de recevoir dès que possible !"); 
};