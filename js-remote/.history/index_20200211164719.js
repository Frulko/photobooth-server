
var ws = new WebSocket("ws://10.0.0.2:8080/echo");

let isConnected = false;
const TXT = {
  CONNECTED: 'connected',
  DISCONNECTED: 'disconnected'
}


function updateStatus() {
  document.querySelector('.status').innerHTML = isConnected ? TXT.CONNECTED: TXT.DISCONNECTED;
}

updateStatus();


ws.addEventListener('open', function() {
  isConnected = true;
  updateStatus();
})

ws.addEventListener('close', function() {
  isConnected = false;
  updateStatus();
})