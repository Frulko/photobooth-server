
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

document.querySelectorAll('.button').forEach((el) => {
  el.addEventListener('click', function(ev) {
    const action = ev.target.getAttribute('data-action');

    ws.send(action);
  })
});

ws.addEventListener('open', function() {
  isConnected = true;
  updateStatus();
})

ws.addEventListener('close', function() {
  isConnected = false;
  updateStatus();
})

ws.addEventListener('message', function(data) {

  if (data.length > 100) {
    const out = document.querySelector("#out").innerHTML;
    document.querySelector("#out").innerHTML = out + "/n" + data.toString; 
  } else {
    console.log('>> data', data);
    var urlCreator = window.URL || window.webkitURL;
    var imageUrl = urlCreator.createObjectURL(this.response);
    document.querySelector("#outimage").src = imageUrl;
  }
  /* var imageElem = document.querySelector('#outimage');
// Just use the toString() method from your buffer instance
// to get date as base64 type
imageElem.src = 'data:image/jpeg;base64,' + data.toString('base64'); */
})