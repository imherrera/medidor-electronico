package main

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"time"

	"github.com/go-redis/redis/v8"
)

// Numero de medidor
type MeterID string

type EnergyUsage struct {
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

/**
 * Mapa para relacionar los C.I -> Uso de Medidor
**/
var energyUsage = map[UserCI][]EnergyUsage{
	UserCI("5444854"): {
		{
			MeterID:        "41250050123",
			Date:           time.Now(),
			ActualWattHour: "1934W",
			ExpectWattHour: "2000W",
		},
	},
}

/**
 * Mapa para relacionar Nro. Medidor -> Nro. C.I
**/
var energyMonitorOwners = map[MeterID]UserCI{}

func getEnergyUsageResume(userCI UserCI) string {
	resume, _ := json.Marshal(energyUsage[userCI])
	return string(resume)
}

func registerEnergyUsage(data EnergyUsage) {
	// Buscamos a que C.I corresponde este Nro. de medidor
	ownerCI := energyMonitorOwners[data.MeterID]
	// Añadimos el uso, al registro de este C.I
	energyUsage[ownerCI] = append(energyUsage[ownerCI], data)
}

func deserializeBody(body io.ReadCloser, v interface{}) (err error) {
	bytes, readError := ioutil.ReadAll(body)
	if readError != nil {
		return readError
	}

	parseError := json.Unmarshal(bytes, &v)
	if parseError != nil {
		return parseError
	}

	return
}
