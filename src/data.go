package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"time"
)

// Numero de medidor
type MeterID string

type EnergyUsage struct {
	MeterID  MeterID   `json:"meter_id"`
	Date     time.Time `json:"date"`
	WattHour float64   `json:"watt_hour"`
	AmpsHour float64   `json:"amps_hour"`
}

// Numero de cedula
type UserCI string

type User struct {
	UserCI   UserCI `json:"user_cid"`
	Username string `json:"username"`
	Password string `json:"password"`
}


/**
 * Mapa para relacionar los C.I -> Uso de Medidor
**/
var energyUsage = map[UserCI][]EnergyUsage{
	UserCI("5444854"): {},
}

/**
 * Mapa para relacionar Nro. Medidor -> Nro. C.I
**/
var energyMonitorOwners = map[MeterID]UserCI{
	MeterID("41250050123"): UserCI("5444854"),
}

/**
 * TODO: send only last 100 reports
**/
func getEnergyUsageResume(userCI UserCI) string {
	usage := energyUsage[userCI]
	resume, _ := json.Marshal(usage)
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
		fmt.Println("Read Error: ", readError)
		return readError
	}

	parseError := json.Unmarshal(bytes, &v)
	if parseError != nil {
		fmt.Println("Parse Error: ", parseError)
		return parseError
	}

	return
}
