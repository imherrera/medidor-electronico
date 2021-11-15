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
	MeterID  MeterID   `json:"meter_id"`
	Date     time.Time `json:"date"`
	WattHour float64   `json:"watt_hour"`
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
			MeterID:  "41250050123",
			Date:     time.Now().Add(time.Hour * time.Duration(1)),
			WattHour: 200.0,
		},
		{
			MeterID:  "41250050123",
			Date:     time.Now().Add(time.Hour * time.Duration(2)),
			WattHour: 300.0,
		},
		{
			MeterID:  "41250050123",
			Date:     time.Now().Add(time.Hour * time.Duration(3)),
			WattHour: 290,
		},
		{
			MeterID:  "41250050123",
			Date:     time.Now().Add(time.Hour * time.Duration(4)),
			WattHour: 200.0,
		},
		{
			MeterID:  "41250050123",
			Date:     time.Now().Add(time.Hour * time.Duration(5)),
			WattHour: 300.0,
		},
		{
			MeterID:  "41250050123",
			Date:     time.Now(),
			WattHour: 290,
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
