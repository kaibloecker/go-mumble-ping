package main

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"time"
)

type Pong struct {
	Version        string `json:"server_version"`
	Ident          uint64 `json:"last_update"`
	ConnectedUsers uint32 `json:"connected_users"`
	MaxUsers       uint32 `json:"max_users"`
	Bandwidth      uint32 `json:"bandwidth"`
}

var (
	MUMBLE_HOST = "localhost"
	MUMBLE_PORT = "64738"
	PORT        = "8080"
)

func main() {
	if val, ok := os.LookupEnv("MUMBLE_HOST"); ok {
		MUMBLE_HOST = val
	}
	if val, ok := os.LookupEnv("MUMBLE_PORT"); ok {
		MUMBLE_PORT = val
	}
	if val, ok := os.LookupEnv("PORT"); ok {
		PORT = val
	}
	address := fmt.Sprintf(":%s", PORT)

	http.HandleFunc("/", getMumbleData)

	if err := http.ListenAndServe(address, nil); err != nil {
		log.Fatal(err)
	}
}

func getMumbleData(w http.ResponseWriter, req *http.Request) {
	server, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%s", MUMBLE_HOST, MUMBLE_PORT))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	conn, err := net.DialUDP("udp", nil, server)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer conn.Close()
	conn.SetReadDeadline(time.Now().Add(150 * time.Millisecond))

	identifier := uint64(time.Now().Unix())
	ping := make([]byte, 12)
	binary.BigEndian.PutUint64(ping[4:], identifier)
	if _, err = conn.Write(ping); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	received := make([]byte, 24)
	if _, err = conn.Read(received); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	pong := Pong{
		Version:        fmt.Sprintf("%d.%d.%d", received[1], received[2], received[3]),
		Ident:          binary.BigEndian.Uint64(received[4:12]),
		ConnectedUsers: binary.BigEndian.Uint32(received[12:16]),
		MaxUsers:       binary.BigEndian.Uint32(received[16:20]),
		Bandwidth:      binary.BigEndian.Uint32(received[20:24]),
	}

	if pong.Ident != identifier {
		http.Error(w, "received scrambled data", http.StatusInternalServerError)
		return
	}

	pongJSON, err := json.Marshal(pong)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "%v", string(pongJSON))
}
