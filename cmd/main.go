package main

import (
	joystick_adapter "joystick-adapter"
	"joystick-adapter/client"
	"log"
)

func main() {
	options := client.ConnectionOptions{
		Url:      "ws://127.0.0.1:8000/ws",
		Protocol: "ws",
		Origin:   "http://127.0.0.1:8000",
	}
	adapter, err := joystick_adapter.NewJoystickAdapter(1, options)
	if err != nil {
		log.Fatal("error during initialization", err)
	}
	adapter.Run()
}
