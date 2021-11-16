package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

/**
 * Middleware de utilizad que añade el header json
**/
func jsonMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		rw.Header().Add("Content-Type", "application/json")
		next.ServeHTTP(rw, r)
	})
}

/**
 * Middleware que interceptara los pedidos a este servidor
 * y validara los tokens de autorizacion obtenidos en el login de usuario,
 * en caso de invalidez el pedido sera rechazado
**/
func webappTokenMiddleware(next http.Handler) http.Handler {
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

var secret string

/**
 * Autorizamos solo a nuestro arduino a enviar datos
**/
func deviceTokenMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		if token != secret {
			rw.WriteHeader(http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(rw, r)
	})
}

/**
 * Punto de entrada del programa
**/
func main() {
	router := mux.NewRouter()
	router.Use(jsonMiddleware)

	/**
	 * Rutas para la autentificacion de cliente
	**/
	auth := router.Methods("POST").Subrouter()
	auth.HandleFunc("/login", loginHandler).Methods("POST")

	/**
	 * Rutas para conexion del cliente web
	**/
	app := router.Methods("GET", "POST", "OPTIONS").Subrouter()
	app.HandleFunc("/usage/resume/{uci}", consumptionResumeHandler).Methods("GET")
	app.Use(webappTokenMiddleware)

	/**
	 * Rutas para la recoleccion de metricas de uso
	**/
	arduino := router.Methods("POST").Subrouter()
	arduino.HandleFunc("/report/watt-hour", consumptionReportHandler).Methods("POST")
	arduino.Use(deviceTokenMiddleware)

	/**
	 * Setup global de headers y CORS
	**/
	headers := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"})
	origins := handlers.AllowedOrigins([]string{"*"}) // <- deberiamos aceptar solo al host de nuestra pagina web
	methods := handlers.AllowedMethods([]string{"GET", "POST", "OPTIONS"})

	/**
	 * Utilizamos puerto definido por las variables de entorno
	**/
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	secret = os.Getenv("ARDUINO_SECRET")

	/**
	 * Iniciamos el servidor HTTP con un logger para debug
	**/
	err := http.ListenAndServe(":"+port, handlers.CombinedLoggingHandler(os.Stderr, handlers.CORS(headers, origins, methods)(router)))
	if err != nil {
		log.Fatalln("Server Error: ", err)
	}
}

func loginHandler(rw http.ResponseWriter, r *http.Request) {
	credentials := UserLoginCredential{}
	err := deserializeBody(r.Body, &credentials)
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	token := login(credentials)
	if len(token) <= 0 {
		rw.WriteHeader(http.StatusUnauthorized)
		return
	}
	fmt.Fprint(rw, token)
}

func consumptionReportHandler(rw http.ResponseWriter, r *http.Request) {
	consumption := EnergyUsage{}
	bytes, _ := ioutil.ReadAll(r.Body)
	json.Unmarshal(bytes, &consumption)
	// TODO: eliminar linea cuando arduino pueda enviar fecha
	consumption.Date = time.Now()
	registerEnergyUsage(consumption)
}

func consumptionResumeHandler(rw http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	userCI := UserCI(params["uci"])

	resume := getEnergyUsageResume(UserCI(userCI))
	fmt.Fprint(rw, resume)
}
