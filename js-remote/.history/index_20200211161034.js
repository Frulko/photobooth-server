
var exampleSocket = new WebSocket("ws://10.0.0.2:8080/echo");

exampleSocket.onopen = function (event) {
  exampleSocket.send(""); 
};