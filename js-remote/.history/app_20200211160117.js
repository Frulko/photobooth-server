const WebSocket = require('ws');

const ws = new WebSocket('ws://10.0.0.2:8080');

ws.on('open', function open() {
  console.log('open');
  // ws.send('something');
});

ws.on('message', function incoming(data) {
  console.log(data);
});