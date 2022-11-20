package joystick_adapter

import (
	"errors"
	"fmt"
	"github.com/splace/joysticks"
	"joystick-adapter/client"
	"joystick-adapter/client/commands"
	"log"
)

// this maps to my csl-controller
const LEFT_STICK = 1
const RIGHT_STICK = 2
const LEFT_SHOULDER_TRIGGER = 7

const SPEED_FACTOR = 0.25
const STEER_FACTOR = 0.25

type JoystickAdapter struct {
	device *joysticks.HID
	conn   *client.Connection
}

func NewJoystickAdapter(hidIndex int, options client.ConnectionOptions) (*JoystickAdapter, error) {
	log.Printf("connecting to joystick")
	device := joysticks.Connect(hidIndex)
	log.Printf("connecting to joystick done")
	if device == nil {
		return nil, errors.New(fmt.Sprintf("hid %v not found", hidIndex))
	}
	log.Printf("connecting to ws")
	conn, err := client.Connect(options)
	if err != nil {
		return nil, err
	}
	log.Printf("connecting to ws done")
	adapter := JoystickAdapter{
		device,
		conn,
	}
	p := &adapter
	return p, nil
}

func (ja *JoystickAdapter) Run() {
	log.Printf("X %v", ja.device.Buttons)
	log.Printf("X %v", ja.device.HatAxes)
	leftStick := ja.device.OnMove(LEFT_STICK)
	rightStick := ja.device.OnMove(RIGHT_STICK)
	leftShoulderTriggerPress := ja.device.OnClose(LEFT_SHOULDER_TRIGGER)
	leftShoulderTriggerLift := ja.device.OnOpen(LEFT_SHOULDER_TRIGGER)

	go ja.device.ParcelOutEvents()

	speed := 0.0
	steer := 0.0
	allowFullspeed := false

	for {
		select {
		case event := <-leftStick:
			coordsEvent := event.(joysticks.CoordsEvent)
			log.Printf("leftStick: %v %v\n", coordsEvent.X, coordsEvent.Y)
			if allowFullspeed {
				speed = float64(coordsEvent.Y) * 1000
			} else {
				speed = float64(coordsEvent.Y) * 1000 * SPEED_FACTOR
			}
			ja.conn.SendCmd(commands.SetDriveCommand{steer, speed})
		case event := <-rightStick:
			coordsEvent := event.(joysticks.CoordsEvent)
			log.Printf("rightStick: %v %v\n", coordsEvent.X, coordsEvent.Y)
			steer = float64(coordsEvent.X) * 1000 * STEER_FACTOR
			ja.conn.SendCmd(commands.SetDriveCommand{steer, speed})
		case <-leftShoulderTriggerPress:
			allowFullspeed = true
		case <-leftShoulderTriggerLift:
			allowFullspeed = false
		}
	}
}
