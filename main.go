package main

import (
	"encoding/binary"
	"fmt"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

type Pong struct {
	Version        string `json:"version"`
	ConnectedUsers uint32 `json:"connected_users"`
	MaxUsers       uint32 `json:"max_users"`
	Bandwidth      uint32 `json:"bandwidth"`
}

func main() {
	router := gin.Default()
	router.GET("/", getMumbleData)

	router.Run("localhost:8080")
}

func getMumbleData(c *gin.Context) {
	server, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%v:%d", os.Getenv("HOST"), 64738))
	if err != nil {
		c.IndentedJSON(http.StatusServiceUnavailable, gin.H{"message": err.Error()})
		return
	}

	conn, err := net.DialUDP("udp", nil, server)
	if err != nil {
		c.IndentedJSON(http.StatusServiceUnavailable, gin.H{"message": err.Error()})
		return
	}
	defer conn.Close()
	conn.SetReadDeadline(time.Now().Add(150 * time.Millisecond))

	ping := []byte{0, 0, 0, 0, 23, 23, 23, 23, 23, 23, 23, 23}
	_, err = conn.Write(ping)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	received := make([]byte, 24)
	_, err = conn.Read(received)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	answer := Pong{
		Version:        fmt.Sprintf("%d.%d.%d", received[1], received[2], received[3]),
		ConnectedUsers: binary.BigEndian.Uint32(received[12:16]),
		MaxUsers:       binary.BigEndian.Uint32(received[16:20]),
		Bandwidth:      binary.BigEndian.Uint32(received[20:24]),
	}

	c.IndentedJSON(http.StatusOK, answer)
}
