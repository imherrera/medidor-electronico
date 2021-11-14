package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"time"
)

/**
 * Credencial para login de usuarios
**/
type UserLoginCredential struct {
	UserCI   UserCI `json:"uci"`
	Password string `json:"pass"`
}

/**
 * Credencial de acceso para un usuario
**/
type UserAccessCredential struct {
	Token      string    `json:"tkn"`
	ValidUntil time.Time `json:"vud"`
}

var tokens = map[UserCI]UserAccessCredential{}
var passwords = map[UserCI]string{
	UserCI("5444854"): "sss",
}

/**
 * Tomamos los credenciales del usuario C.I y contraseña
 *
 * Devolvemos token(vacio en caso de credencial incorrecta)
**/
func login(credential UserLoginCredential) (token string) {
	if passwords[credential.UserCI] == credential.Password {
		token := UserAccessCredential{
			Token:      randSeq(256),
			ValidUntil: time.Now().Add(time.Hour * time.Duration(24)), // <- validez de 24 horas
		}
		// guardamos el token generado para este C.I
		tokens[credential.UserCI] = token

		jwt, _ := json.Marshal(token)
		return string(jwt)
	}
	return ""
}

/**
 * Tomamos un token de sesion
 *
 * Devolvemos un true si el token es valido
**/
func tokenIsValid(uci UserCI, token string) bool {
	mtoken, exists := tokens[uci]
	// check que este token corresponda a este usuario
	if !exists {
		fmt.Println("token does not exist!")
		return false
	}
	// check que el token que se le asigno a @uci corresponda al token
	if token != mtoken.Token {
		fmt.Println("token does not match!")
		return false
	}
	// check que el token no haya expirado
	if time.Now().After(mtoken.ValidUntil) {
		fmt.Println("toke is expired!")
		return false
	}
	// el token existe, corresponde al usuario y no expiro aun
	return true
}

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randSeq(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
