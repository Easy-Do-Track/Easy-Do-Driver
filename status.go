package main

import (
	"encoding/json"
	"net/http"
	"time"
)

type Status struct {
	Sensors    []int `json:"sensors"`
	data       [16]Quaternion
	IP         string    `json:"addr"`
	LastUpdate time.Time `json:"last_update"`
}

func (s *Status) Dump(data [16]Quaternion, ip string) {
	s.data = data
	s.IP = ip
	s.LastUpdate = time.Now()
}

func (s Status) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	s.Sensors = make([]int, 0)
	for i, d := range s.data {
		if d.X != 0 || d.Y != 0 || d.Z != 0 {
			s.Sensors = append(s.Sensors, i)
		}
	}

	err := json.NewEncoder(res).Encode(s)
	if err != nil {
		http.Error(res, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}
