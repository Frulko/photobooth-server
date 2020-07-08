const WebSocket = require('ws');

const ws = new WebSocket('ws://10.0.0.2:8080/echo');

ws.on('open', function open() {
  console.log('open');
  // ws.send('something');
});

ws.on('message', function incoming(data) {
  console.log(data);
});

ws.on('close', function close() {
  console.log('disconnected');
});