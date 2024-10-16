//operations send by ws

let ops = {
  // remove element
  0: removeElement,
  1: addElement,
};

function handleOperation(op) {
  let [id, ...args] = op;
  let fn = ops[id];
  if (fn) {
    fn(...args);
  } else {
    console.error("Invalid operation: " + id);
  }
}

function parseMessage(message) {
  // op is the first character of the message
  let op = message[0];
  let data = message.slice(1);
  handleOperation([op, data]);
}

function removeElement(id) {
  document.getElementById(id).remove();
}

function addElement(text) {
  let parentID = text.split(",", 2);
  let parent = document.getElementById(parentID[0]);
  if (parent) {
    let innerText = text.slice(parentID[0].length + 1);
    parent.innerHTML += innerText;
  } else {
    console.error("Parent not found: " + parentID[0]);
  }
}

let ws;
function connectWebSocket() {
  ws = new WebSocket(window.location.href.replace("http", "ws") + "ws");

  ws.onopen = function () {
    console.log("Connected");
  };

  ws.onmessage = function (event) {
    parseMessage(event.data);
  };

  ws.onclose = function () {
    console.log("Disconnected");
  };
}

function sendMessage() {
  let message = document.getElementById("message").value;
  if (ws && ws.readyState === WebSocket.OPEN) {
    ws.send(message);
  }
}

// Connect WebSocket when the page loads
window.onload = connectWebSocket;
