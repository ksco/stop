package main

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"github.com/joho/godotenv"
	"net/url"
	"os"
	"strconv"
)

func main() {
	_ = godotenv.Load()

	token := os.Getenv("TOKEN")
	memoryStr := os.Getenv("MEMORY")
	memory, err := strconv.ParseFloat(memoryStr, 64)
	if err != nil {
		memory = 70
	}
	cpuStr := os.Getenv("CPU")
	cpu, err := strconv.ParseFloat(cpuStr, 64)
	if err != nil {
		cpu = 80
	}
	diskStr := os.Getenv("DISK")
	disk, err := strconv.ParseFloat(diskStr, 64)
	if err != nil {
		disk = 80
	}

	host := os.Getenv("HOST")
	port := os.Getenv("PORT")
	u := url.URL{
		Scheme: "ws",
		User:   nil,
		Host:   host + ":" + port,
		Path:   "/ws",
	}
	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		return
	}
	defer func() { _ = c.Close() }()

	err = c.WriteMessage(websocket.TextMessage, []byte(token))
	if err != nil {
		return
	}

	for {
		_, message, err := c.ReadMessage()
		if err != nil {
			println(err.Error())
			return
		}

		frame := &Frame{}
		err = json.Unmarshal(message, frame)
		if err != nil {
			continue
		}

		if frame.MemoryUsagePercent() > memory || frame.DiskUsagePercent() > disk || frame.CPUUsagePercent > cpu {
			_ = Push(frame)
		}
	}
}
