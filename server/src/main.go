package main

import (
	"fmt"
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

	app := router.Methods("GET", "POST", "OPTIONS").Subrouter()
	app.HandleFunc("/usage", consumptionHandler).Methods("GET", "POST")
	app.HandleFunc("/usage/resume/{uci}", consumptionResumeHandler).Methods("GET")
	app.Use(tokenMiddleware)

	headers := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"})
	origins := handlers.AllowedOrigins([]string{"*"}) // <- deberiamos aceptar solo al host de nuestra pagina web
	methods := handlers.AllowedMethods([]string{"GET", "POST", "OPTIONS"})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

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

func consumptionHandler(rw http.ResponseWriter, r *http.Request) {
	consumption := EnergyUsage{}
	err := deserializeBody(r.Body, consumption)
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		return
	}
	registerEnergyUsage(consumption)
}

func consumptionResumeHandler(rw http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	userCI := UserCI(params["uci"])

	resume := getEnergyUsageResume(UserCI(userCI))
	fmt.Fprint(rw, resume)
}
