package models

// Client is a middleman between the .websocket connection and the hub.

type ChatEvent struct {
	Event    string `json:"event"`
	ID       string `json:"id"`
	ClientId string `json:"clientId"`
	Message  string `json:"message"`
}

type NewChatMessage struct {
	Event   string `json:"event"`
	ID      string `json:"id"`
	TempId  string `json:"temp_id"`
	Message string `json:"message"` // Add this field
}

type EditChatMessage struct {
	ID string `json:"id"`
}
