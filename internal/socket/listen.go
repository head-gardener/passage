package socket

import (
	"encoding/json"
	"log"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/head-gardener/passage/internal/config"
	state "github.com/head-gardener/passage/internal/net"
)

func Listen(log *slog.Logger, st *state.State, conf *config.Config) {
	if err := os.RemoveAll(conf.Socket.Path); err != nil {
		log.Error("removing socket file", "err", err)
	}

	listener, err := net.Listen("unix", conf.Socket.Path)
	if err != nil {
		log.Error("creating socket", "err", err)
		os.Exit(1)
	}
	defer listener.Close()

	if err := os.Chmod(conf.Socket.Path, 0777); err != nil {
		log.Error("setting socket permissions", "err", err)
	}

	done := make(chan struct{}, 2)
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sigChan
		done <- struct{}{}
		listener.Close()
		os.Remove(conf.Socket.Path)
		os.Exit(0)
	}()

	errors := 0
	for {
		select {
		case <-done:
			break
		default:
		}

		conn, err := listener.Accept()
		if err != nil {
			log.Warn("socket accept error", "err", err)
			errors += 1
			if errors >= 3 {
				return
			}
			continue
		}

		go handleConnection(conn, st, conf)
	}
}

func handleConnection(conn net.Conn, st *state.State, conf *config.Config) {
	defer conn.Close()

	buf := make([]byte, 4096)

	n, err := conn.Read(buf)
	if err != nil {
		log.Println("Read error:", err)
		return
	}

	var cmd Command
	if err := json.Unmarshal(buf[:n], &cmd); err != nil {
		sendError(conn, "Invalid command format")
		return
	}

	response := processCommand(cmd, st, conf)

	jsonResponse, err := json.Marshal(response)
	if err != nil {
		sendError(conn, "Error creating response")
		return
	}

	if _, err := conn.Write(jsonResponse); err != nil {
		log.Println("Write error:", err)
	}
}

func processCommand(cmd Command, st *state.State, conf *config.Config) Response {
	switch cmd.Action {
	case "ping":
		return Response{
			Status:  "success",
			Message: "pong",
		}
	case "status":
		peers := make([]StatusResponsePeer, len(conf.Peers))
		for i := range conf.Peers {
			peers[i] = StatusResponsePeer{
				conf.Peers[i].Addr.String(),
				st.IsConnected(i),
				st.LastSeen(i),
			}
		}
		resp := StatusResponse{peers}
		data, err := json.Marshal(resp)
		if err != nil {
			return Response{
				Status:  "error",
				Message: "Couldn't marshall response",
			}
		}
		return Response{
			Status:  "success",
			Message: "status",
			Data:    data,
		}
	default:
		return Response{
			Status:  "error",
			Message: "unknown command",
		}
	}
}

func sendError(conn net.Conn, message string) {
	errorResponse := Response{
		Status:  "error",
		Message: message,
	}
	jsonResponse, _ := json.Marshal(errorResponse)
	conn.Write(jsonResponse)
}
