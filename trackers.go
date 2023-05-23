package main

import (
	"bytes"
	"encoding/binary"
	"log"
	"net"
	"unsafe"
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

type Euler struct {
	X float32 `json:"x"`
	Y float32 `json:"y"`
	Z float32 `json:"z"`
}

type Quaternion struct {
	W float32
	X float32
	Y float32
	Z float32
}

type Tracker struct {
	conn net.PacketConn
	data chan [16]Quaternion
}

func NewTracker(addr string) (Tracker, error) {
	conn, err := net.ListenPacket("udp", addr)

	if err != nil {
		return Tracker{}, err
	}

	t := Tracker{conn: conn}
	t.data = make(chan [16]Quaternion, 1)

	go t.readPump()

	return t, nil
}

func (t Tracker) readPump() {
	for {
		var data [16]Quaternion

		resp := make([]byte, unsafe.Sizeof(data))

		_, _, err := t.conn.ReadFrom(resp)
		if err != nil {
			log.Println(err)
		}

		if err = binary.Read(bytes.NewReader(resp), binary.LittleEndian, &data); err != nil {
			log.Println(err)
		}

		t.data <- data
	}
}

func (t Tracker) DataChannel() <-chan [16]Quaternion {
	return t.data
}
