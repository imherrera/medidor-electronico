package main

import (
	"encoding/json"
	"time"

	"github.com/go-redis/redis/v8"
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

var data = map[string]string{}

func getUsageResume(ucid UserCI) string {
	mock := EnergyConsumption{
		OwnerCI:        ucid,
		MeterID:        "41250050123",
		Date:           time.Now(),
		ActualWattHour: "1934W",
		ExpectWattHour: "2000W",
	}
	resume, _ := json.Marshal(mock)
	return string(resume)
}

func putUsage(data EnergyConsumption) {
	println("Putooooo")
}
