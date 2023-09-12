package main

import (
	"encoding/binary"
	"fmt"
	"net"
	"net/http"
	"os"
	"time"

	cache "github.com/chenyahui/gin-cache"
	"github.com/chenyahui/gin-cache/persist"
	"github.com/gin-gonic/gin"
)

type Pong struct {
	Version        string `json:"server_version"`
	Ident          uint64 `json:"last_update"`
	ConnectedUsers uint32 `json:"connected_users"`
	MaxUsers       uint32 `json:"max_users"`
	Bandwidth      uint32 `json:"bandwidth"`
}

func main() {
	router := gin.Default()
	memoryStore := persist.NewMemoryStore(1 * time.Minute)

	router.GET("/", cache.CacheByRequestURI(memoryStore, 15*time.Second), getMumbleData)

	port := "8080"
	value, ok := os.LookupEnv("PORT")
	if ok {
		port = value
	}
	address := fmt.Sprintf(":%s", port)

	if err := router.Run(address); err != nil {
		panic(err)
	}
}

func getMumbleData(c *gin.Context) {
	mumblePort := "64738"
	if v, ok := os.LookupEnv("MUMBLE_PORT"); ok {
		mumblePort = v
	}

	server, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%s", os.Getenv("MUMBLE_HOST"), mumblePort))
	if err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"message": err.Error()})
		return
	}

	conn, err := net.DialUDP("udp", nil, server)
	if err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"message": err.Error()})
		return
	}
	defer conn.Close()
	conn.SetReadDeadline(time.Now().Add(150 * time.Millisecond))

	identifier := uint64(time.Now().Unix())
	ping := make([]byte, 12)
	binary.BigEndian.PutUint64(ping[4:], identifier)
	_, err = conn.Write(ping)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	received := make([]byte, 24)
	_, err = conn.Read(received)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
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
		c.JSON(http.StatusInternalServerError, gin.H{"message": "received scrambled data"})
		return
	}

	c.JSON(http.StatusOK, pong)
}
