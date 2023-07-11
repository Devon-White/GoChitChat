const ws = new WebSocket(`wss://${document.location.host}/ws`);
let clientId;

const messageInput = document.getElementById('message-input');
const sendButton = document.getElementById('send-button');
const chatbox = document.getElementById('chatbox');

ws.onopen = function() {
    console.log('WebSocket connection opened');
};

ws.onmessage = function(event) {
    let msg;
    try {
        msg = JSON.parse(event.data);
        console.log(msg)
    } catch(e) {
        console.log("Received non-JSON message:", event.data);
        return;
    }

    switch (msg.event) {
        case "clientId":
            clientId = msg.clientId;
            messageInput.disabled = false;
            sendButton.disabled = false;
            break;
        case "new_message":
            let messageElement = document.getElementById(msg.temp_id);
            if (messageElement) {
                messageElement.style.color = "black";
                messageElement.id = msg.id;
            } else if (!document.getElementById(msg.id)) {
                console.log("Can't find Element ID")
                messageElement = document.createElement("p");
                messageElement.style.color = "black";
                messageElement.id = msg.id;
                messageElement.textContent = msg.message;
                chatbox.appendChild(messageElement);
            }
            break;
        default:
            console.log("Unknown message type:", msg.event);
    }
};

ws.onclose = function(event) {
    console.log('WebSocket connection closed:', event);
};

ws.onerror = function(error) {
    console.log('WebSocket error:', error);
};

async function sendMessage() {
    let currentTempId = clientId + ":" + Date.now().toString();

    let messageElement = document.createElement("p");
    messageElement.style.color = "blue";
    messageElement.id = currentTempId;
    messageElement.textContent = messageInput.value;
    chatbox.appendChild(messageElement);

    let msg = {
        "event": "new_message",
        "data": {
            "id": currentTempId,
            "clientid": clientId,
            "message": messageInput.value,
        }
    };

    await new Promise(resolve => setTimeout(resolve, 500));
    ws.send(JSON.stringify(msg));

    messageInput.value = '';
    // Add a short delay to allow the DOM to update
}

function showMembers() {
    alert('Showing active members...');
}