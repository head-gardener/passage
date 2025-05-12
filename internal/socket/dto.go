package socket

import (
	"encoding/json"
	"time"
)

type Command struct {
	Action  string `json:"action"`
	Payload any    `json:"payload,omitempty"`
}

type Response struct {
	Status  string          `json:"status"`
	Message string          `json:"message"`
	Data    json.RawMessage `json:"data,omitempty"`
}

type StatusResponse struct {
	Peers []StatusResponsePeer `json:"peers"`
}

type StatusResponsePeer struct {
	Addr      string
	Connected bool
	LastSeen  time.Time
}
