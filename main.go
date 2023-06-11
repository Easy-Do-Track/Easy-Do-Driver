package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"log"
	"net"
	"net/http"
	"strconv"
)

type values struct {
	Name    string            `json:"name"`
	Address map[string]string `json:"address"`
	Setting map[string]string `json:"setting"`
}

func main() {
	conf, err := ConfigFromFile("config.toml")
	if err != nil {
		log.Fatal(err)
	}

	r := mux.NewRouter()
	corsMw := mux.CORSMethodMiddleware(r)

	s := NewStreamer()

	stat := &Status{}

	t, err := NewTracker(conf.Tracker.Address, stat)

	fmt.Println("Starting tracker listener at", conf.Tracker.Address)

	if err != nil {
		log.Fatal(err)
	}

	go func() {
		for {
			data := <-t.data

			// 머리
			data[2] = multiplyQuaternion(inverseQuaternion(data[0]), data[2])

			// 왼쪽 팔
			data[3] = multiplyQuaternion(inverseQuaternion(data[1]), data[3])
			data[1] = multiplyQuaternion(inverseQuaternion(data[0]), data[1])

			// 오른쪽 팔
			data[7] = multiplyQuaternion(inverseQuaternion(data[6]), data[7])
			data[6] = multiplyQuaternion(inverseQuaternion(data[0]), data[6])

			// 왼쪽 다리
			data[5] = multiplyQuaternion(inverseQuaternion(data[4]), data[5])

			// 오른쪽 다리
			data[9] = multiplyQuaternion(inverseQuaternion(data[8]), data[9])

			result := make(map[string]Euler)
			for k, v := range conf.Tracker.Mappings {
				id, err := strconv.Atoi(k)
				if err != nil {
					log.Fatal(err)
				}

				temp := quaternionToEuler(data[id], "XYZ")

				if v.Multi.X != 0 {
					temp.X *= v.Multi.X
				}

				if v.Multi.Y != 0 {
					temp.Y *= v.Multi.Y
				}

				if v.Multi.Z != 0 {
					temp.Z *= v.Multi.Z
				}

				var e Euler

				switch v.Rotation[0] {
				case 'X':
					e.X = temp.X
				case 'Y':
					e.X = temp.Y
				case 'Z':
					e.X = temp.Z
				}

				switch v.Rotation[1] {
				case 'X':
					e.Y = temp.X
				case 'Y':
					e.Y = temp.Y
				case 'Z':
					e.Y = temp.Z
				}

				switch v.Rotation[2] {
				case 'X':
					e.Z = temp.X
				case 'Y':
					e.Z = temp.Y
				case 'Z':
					e.Z = temp.Z
				}

				e.X += v.Offset.X
				e.Y += v.Offset.Y
				e.Z += v.Offset.Z

				result[v.Name] = e
			}

			j, err := json.Marshal(result)
			if err != nil {
				log.Println(err)
			}

			//fmt.Println(string(j))

			s.Broadcast(j)
		}
	}()

	r.Handle("/stream", s)

	r.Handle("/status", stat)

	r.Use(corsMw)

	headersOk := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"})
	originsOk := handlers.AllowedOrigins([]string{"*"})
	methodsOk := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "OPTIONS"})

	ips, err := net.InterfaceAddrs()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Your IPs:")
	for _, i := range ips {
		fmt.Println(i, i.Network())
	}

	fmt.Println("Starting server at", conf.Server.Address)
	if err := http.ListenAndServe(conf.Server.Address,
		handlers.CORS(headersOk, originsOk, methodsOk)(r)); err != nil {
		log.Fatal(err)
	}
}
