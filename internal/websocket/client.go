package websocket

import (
	chat "GoChitChat/internal/chat/models"
	md "GoChitChat/internal/websocket/models"
	"GoChitChat/pkg"
	"encoding/json"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

// Constants for configuring the WebSocket client.
const (
	writeWait      = 10 * time.Second    // Time allowed to write a message to the peer.
	pongWait       = 60 * time.Second    // Time allowed to read the next pong message from the peer.
	pingPeriod     = (pongWait * 9) / 10 // Time period for sending ping messages to the peer.
	maxMessageSize = 4096                // Maximum message size allowed from the peer.
)

// Client is a middleman between the .websocket connection and the hub.
type Client struct {
	Conn *websocket.Conn // The .websocket connection.
	Send chan []byte     // Buffered channel of outbound messages.
	ID   string          // The unique ID for the client.
	Hub  *Hub            // The hub the client is connected to.
}

func MarshalData(data md.SocketMessage) chat.NewChatMessage {

	messageData, err := json.Marshal(data.DATA)
	if err != nil {
		log.Println("error marshalling data:", err)
	}

	var messageJson chat.NewChatMessage
	err = json.Unmarshal(messageData, &messageJson)
	if err != nil {
		log.Println("error unmarshalling new messages:", err)

	}
	log.Println(messageJson)
	return messageJson

}

// WritePump pumps messages from the hub to the .websocket connection.
func (client *Client) WritePump() {
	ticker := time.NewTicker(pingPeriod) // Ticker to control pings to peer.
	defer func() {
		ticker.Stop()       // Stop the ticker when we're done.
		client.Conn.Close() // Close the connection when we're done.
	}()

	// Process messages and pings.
	for {
		select {
		case message, ok := <-client.Send:
			client.Conn.SetWriteDeadline(time.Now().Add(writeWait)) // Set the write deadline.
			if !ok {
				// The hub closed the channel, so send a close message to the peer.
				client.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			// Create a new writer for the current message.
			w, err := client.Conn.NextWriter(websocket.TextMessage)
			if err != nil {
				log.Println(err)
				return
			}
			w.Write(message)

			// Add any queued chat messages to the current .websocket message.
			n := len(client.Send)
			for i := 0; i < n; i++ {
				w.Write(<-client.Send)
			}

			// Close the writer when we're done to flush the message to the network.
			if err := w.Close(); err != nil {
				log.Println(err)
				return
			}
		case <-ticker.C:
			// Every tick, send a ping message to the peer.
			client.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := client.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				log.Println(err)
				return
			}
		}
	}
}

// ReadPump pumps messages from the .websocket connection to the hub.
func (client *Client) ReadPump() {
	defer func() {
		client.Hub.Unregister <- client // Unregister the client from the hub when we're done.
		client.Conn.Close()             // Close the connection when we're done.
	}()

	client.Conn.SetReadLimit(maxMessageSize)              // Set the read limit.
	client.Conn.SetReadDeadline(time.Now().Add(pongWait)) // Set the read deadline.
	client.Conn.SetPongHandler(func(string) error {
		client.Conn.SetReadDeadline(time.Now().Add(pongWait)) // Reset the read deadline when we receive a pong.
		return nil
	})

	// Process incoming messages.
	for {
		_, message, err := client.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}

		// Unmarshal the incoming message into a ChatEvent.
		var newChatMessage chat.NewChatMessage
		var socketMessage md.SocketMessage

		err = json.Unmarshal(message, &socketMessage)
		if err != nil {
			log.Println("error unmarshalling message:", err)
			break
		}

		// Process the chat event.
		switch socketMessage.Event {
		case "new_message":

			organizedData := MarshalData(socketMessage)

			// Create a new chat message and broadcast it to the hub.
			newMessage := chat.NewChatMessage{
				Event:   "new_message",
				ID:      pkg.GenerateUniqueId(),
				TempId:  organizedData.ID,
				Message: organizedData.Message, // Include the original message content
			}
			log.Println(&newMessage)
			messageBytes, _ := json.Marshal(&newMessage)
			client.Hub.Broadcast <- messageBytes // Broadcast the new message to all clients
		default:
			log.Println("Unknown event:", newChatMessage.Event)
		}
	}
}
