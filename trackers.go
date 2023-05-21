package main

import (
	"fmt"
	"log"
	"net"
)

const (
	KeyHead       = "head"
	KeyChest      = "chest"
	KeyElbowLeft  = "elbow_left"
	KeyElbowRight = "elbow_right"
	KeyWrestLeft  = "wrest_left"
	KeyWrestRight = "wrest_right"
	KeyKneeLeft   = "knee_left"
	KeyKneeRight  = "knee_right"
	KeyAnkleLeft  = "ankle_left"
	KeyAnkleRight = "ankle_right"
)

type Vector3 struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
	Z float64 `json:"z"`
}

type Tracker struct {
	conn net.PacketConn
	data chan map[string]Vector3
}

func NewTracker(addr string) (Tracker, error) {
	conn, err := net.ListenPacket("udp", addr)

	if err != nil {
		return Tracker{}, err
	}

	t := Tracker{conn: conn}

	go t.readJSON()

	return t, nil
}

func (t Tracker) readJSON() {
	for {
		//data := make(map[string]Vector3)
		resp := make([]byte, 512)
		n, _, err := t.conn.ReadFrom(resp)
		if err != nil {
			log.Println(err)
		}
		fmt.Println(string(resp[:n]))
	}
}

func (t Tracker) DataChannel() <-chan map[string]Vector3 {
	return t.data
}
