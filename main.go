package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"log"
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

	t, err := NewTracker(conf.Tracker.Address)
	fmt.Println("Starting tracker listener at", conf.Tracker.Address)

	if err != nil {
		log.Fatal(err)
	}

	go func() {
		for {
			data := <-t.data

			result := make(map[string]Euler)
			for k, v := range conf.Tracker.Mappings {
				id, err := strconv.Atoi(k)
				if err != nil {
					log.Fatal(err)
				}

				result[v] = quaternionToEuler(data[id], "XYZ")
			}

			j, err := json.Marshal(result)
			if err != nil {
				log.Println(err)
			}

			fmt.Println(string(j))

			s.Broadcast(j)
		}
	}()

	r.Handle("/stream", s)

	r.HandleFunc("/profile", func(res http.ResponseWriter, req *http.Request) {
		var result values
		if err := json.NewDecoder(req.Body).Decode(&result); err != nil {
			log.Println(err)
			res.Write([]byte(`"result": "error"`))
			return
		}
		fmt.Println(result)
		res.Write([]byte(`{"result":"ok"}`))
	})
	r.Use(corsMw)

	headersOk := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"})
	originsOk := handlers.AllowedOrigins([]string{"*"})
	methodsOk := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "OPTIONS"})

	fmt.Println("Starting server at", conf.Server.Address)
	if err := http.ListenAndServe(conf.Server.Address,
		handlers.CORS(headersOk, originsOk, methodsOk)(r)); err != nil {
		log.Fatal(err)
	}
}
