package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
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

type Vector3 struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
	Z float64 `json:"z"`
}

type Quaternion struct {
	W float32
	X float32
	Y float32
	Z float32
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

		fmt.Println(data)
	}
}

func (t Tracker) DataChannel() <-chan map[string]Vector3 {
	return t.data
}
