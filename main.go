package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
	"time"
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
			data := <-t.DataChannel()
			log.Println(data)
		}
	}()

	go func() {
		for {
			time.Sleep(500 * time.Millisecond)
			data, err := os.ReadFile("tracker.json")
			if err != nil {
				log.Fatal(err)
			}
			s.Broadcast(data)
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
