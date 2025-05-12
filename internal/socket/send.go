package socket

import (
	"encoding/json"
	"log"
	"net"
	"os"

	"github.com/jedib0t/go-pretty/v6/table"

	"github.com/head-gardener/passage/internal/config"
)

func Status(conf *config.Config) {
	conn, err := net.Dial("unix", conf.Socket.Path)
	if err != nil {
		log.Fatal("Connection error:", err)
	}
	defer conn.Close()

	// Send the command as JSON
	jsonCmd, err := json.Marshal(Command{
		"status",
		"",
	})

	if err != nil {
		log.Fatal("JSON encode error:", err)
	}

	if _, err := conn.Write(jsonCmd); err != nil {
		log.Fatal("Write error:", err)
	}

	buf := make([]byte, 4096)
	n, err := conn.Read(buf)
	if err != nil {
		log.Fatal("Read error:", err)
	}

	var response Response
	if err := json.Unmarshal(buf[:n], &response); err != nil {
		log.Fatal("JSON decode error:", err)
	}

	if response.Status != "success" {
		log.Fatal("Server error:", response.Message)
	}

	if response.Data == nil {
		log.Fatal("Server returned no data with message:", response.Message)
	}

	var data StatusResponse
	if err := json.Unmarshal(response.Data, &data); err != nil {
		log.Fatal("Couldn't unmarshall StatusResponse:", err, response.Data)
	}

	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"Remote", "Connected", "Last seen"})
	for _, p := range data.Peers {
		lastSeen := "never"
		if !p.LastSeen.IsZero() {
			lastSeen = p.LastSeen.String()
		}
		t.AppendRow(table.Row{p.Addr, p.Connected, lastSeen})
	}
	t.Render()
}
