const ws = new WebSocket("wss://" + location.host + "/ws/");

ws.onmessage = function (e) {
  console.log(e.data);
};

function send(msg) {
  ws.send(msg);
}
