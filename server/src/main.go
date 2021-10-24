package main

import (
	"fmt"
	"html"
	"log"
	"net/http"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/gorilla/mux"
)

type EnergyConsumption struct {
	MeterID        string    `json:"meter_id"`
	Date           time.Time `json:"date"`
	ActualWattHour string    `json:"actual_watt_hour"`
	ExpectWattHour string    `json:"expect_watt_hour"`
}

// Global reference to redis database client
var rdb *redis.Client

func main() {
	router := mux.NewRouter()
	rdb = redis.NewClient(&redis.Options{
		Addr:     "localhost:1300",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
	})

	router.HandleFunc("/hi", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hi")
	})

	// Init http server
	err := http.ListenAndServe(":8080", router)
	if err != nil {
		log.Fatalln("Server Error: ", err)
	}
}

func putEnergyConsumption(w http.ResponseWriter, r *http.Request) {
	meterId := mux.Vars(r)["mid"]

}

func getEnergyConsumption(w http.ResponseWriter, r *http.Request) {

}
