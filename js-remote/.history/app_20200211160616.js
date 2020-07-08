var WebSocketClient = require('websocket').client;
 
var client = new WebSocketClient();
 
client.on('connectFailed', function(error) {
    console.log('Connect Error: ' + error.toString());
});
 

client.connect('ws://10.0.0.2:8080/echo', 'echo-protocol');