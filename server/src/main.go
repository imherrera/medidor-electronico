package main

import (
	"fmt"
	"log"
	"net/http"

	"encoding/json"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/gorilla/mux"
)

// Numero de medidor
type MeterID string

type EnergyConsumption struct {
	OwnerCI        UserCI    `json:"user_cid"`
	MeterID        MeterID   `json:"meter_id"`
	Date           time.Time `json:"date"`
	ActualWattHour string    `json:"actual_watt_hour"`
	ExpectWattHour string    `json:"expect_watt_hour"`
}

// Numero de cedula
type UserCI string

type User struct {
	UserCI   UserCI `json:"user_cid"`
	Username string `json:"username"`
	Password string `json:"password"`
}

var rdb *redis.Client = redis.NewClient(&redis.Options{
	Addr:     "localhost:1300",
	Password: "", // no password set
	DB:       0,  // use default DB
})

func getUsageResume(ucid UserCI) ([]byte, error) {
	mock := EnergyConsumption{
		OwnerCI:        ucid,
		MeterID:        "41250050123",
		Date:           time.Now(),
		ActualWattHour: "1934W",
		ExpectWattHour: "2000W",
	}
	return json.Marshal(mock)
}

func putUsage(data EnergyConsumption) {
	println("Putooooo")
}

func main() {
	router := mux.NewRouter()

	router.HandleFunc("/", func(rw http.ResponseWriter, r *http.Request) {
		fmt.Fprint(rw, "Hey there!")
	})

	router.HandleFunc("/login", func(rw http.ResponseWriter, r *http.Request) {

	})

	/**
	* Tomamos el nro de medidor y registramos o devolvemos el uso de acuerdo al metodo de request
	**/
	router.HandleFunc("/usage/{mid}", energyConsumption).Methods("GET", "POST")

	/**
	* Tomamos el nro de cedula como parametro y devolvemos un resumen de uso
	**/
	router.HandleFunc("/usage/resume/{uci}", getEnergyConsumptionResume).Methods("GET")

	// Init http server
	err := http.ListenAndServe(":8080", router)
	if err != nil {
		log.Fatalln("Server Error: ", err)
	}
}

func energyConsumption(rw http.ResponseWriter, r *http.Request) {
	meterID := mux.Vars(r)["mid"]
	println("Meter id", meterID)
}

func getEnergyConsumptionResume(rw http.ResponseWriter, r *http.Request) {
	userCI := mux.Vars(r)["uci"]

	json, err := getUsageResume(UserCI(userCI))

	if err != nil {
		log.Fatalln("Error retrieving resume:", err)
	}

	fmt.Fprint(rw, string(json))
}
