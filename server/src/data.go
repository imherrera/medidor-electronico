package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"time"

	"github.com/go-redis/redis/v8"
)

// Numero de medidor
type MeterID string

// Fecha de reporte
type ReportTime time.Time

/*// Implement Marshaler and Unmarshaler interface
func (j *ReportTime) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), "\"")
	t, err := time.Parse("2006-01-02", s)
	if err != nil {
		return err
	}
	*j = ReportTime(t)
	return nil
}

func (j ReportTime) MarshalJSON() ([]byte, error) {
	return json.Marshal(time.Time(j))
}*/

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
		/*{
			MeterID: "41250050123",
			Date:     time.Now(),
			WattHour: 290,
		},
		{
			MeterID:  "41250050123",
			Date:     ReportTime(time.Now().Add(time.Hour * time.Duration(1))),
			WattHour: 200.0,
		},
		{
			MeterID:  "41250050123",
			Date:     ReportTime(time.Now().Add(time.Hour * time.Duration(2))),
			WattHour: 300.0,
		},
		{
			MeterID:  "41250050123",
			Date:     ReportTime(time.Now().Add(time.Hour * time.Duration(3))),
			WattHour: 290,
		},
		{
			MeterID:  "41250050123",
			Date:     ReportTime(time.Now().Add(time.Hour * time.Duration(4))),
			WattHour: 200.0,
		},
		{
			MeterID:  "41250050123",
			Date:     ReportTime(time.Now().Add(time.Hour * time.Duration(5))),
			WattHour: 300.0,
		},*/
	},
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
	time.Now()
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
