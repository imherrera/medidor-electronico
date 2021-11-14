package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func jsonMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		rw.Header().Add("Content-Type", "application/json")
		next.ServeHTTP(rw, r)
	})
}

/**
 * Middleware que interceptara los pedidos a este servidor
 * y validara los tokens de autorizacion, en caso de invalidez el pedido sera rechazado
**/
func tokenMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		userCI := UserCI(params["uci"])
		// Hacemos una pequeña sanitizacion del token
		token := r.Header.Get("Authorization")
		token = strings.Replace(token, "Bearer ", "", 1)
		token = strings.TrimSpace(token)
		// Rechazamos pedido en caso de invalidez
		if !tokenIsValid(userCI, token) {
			rw.WriteHeader(http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(rw, r) // Dejamos pasar al pedido
	})
}

func main() {
	router := mux.NewRouter()
	router.Use(jsonMiddleware)

	auth := router.Methods("POST").Subrouter()
	auth.HandleFunc("/login", loginHandler).Methods("POST")

	app := router.Methods("GET", "POST").Subrouter()
	app.HandleFunc("/usage/{mid}", consumptionHandler).Methods("GET", "POST")
	app.HandleFunc("/usage/resume/{uci}", consumptionResumeHandler).Methods("GET")
	app.Use(tokenMiddleware)

	headersOK := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type"})
	originsOK := handlers.AllowedOrigins([]string{"*"})
	methodsOK := handlers.AllowedMethods([]string{"GET", "POST", "OPTIONS", "DELETE", "PUT"})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	err := http.ListenAndServe(":"+port, handlers.CombinedLoggingHandler(os.Stderr, handlers.CORS(headersOK, originsOK, methodsOK)(router)))
	if err != nil {
		log.Fatalln("Server Error: ", err)
	}
}

func loginHandler(rw http.ResponseWriter, r *http.Request) {
	body, readErr := ioutil.ReadAll(r.Body)
	if readErr != nil {
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	credentials := UserLoginCredential{}
	parseErr := json.Unmarshal(body, &credentials)
	if parseErr != nil {
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	rw.Header().Set("Access-Control-Allow-Origin", "*")
	rw.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	rw.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type")

	token := login(credentials)
	if len(token) == 0 {
		rw.WriteHeader(http.StatusUnauthorized)
		return
	}

	fmt.Fprint(rw, token)
}

func consumptionHandler(rw http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	meterID := params["mid"]
	println("Meter id", meterID)
}

func consumptionResumeHandler(rw http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	userCI := UserCI(params["uci"])
	fmt.Println("Resume request for C.I: ", userCI)

	resume := getUsageResume(UserCI(userCI))
	fmt.Fprint(rw, resume)
}
