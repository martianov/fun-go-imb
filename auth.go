package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
)

const (
	privateKey = "./keys/app.rsa"
)

var (
	signKey []byte
)

func init() {
	var err error
	signKey, err = ioutil.ReadFile(privateKey)
	if err != nil {
		log.Fatal("Error reading Private Key")
	}
}

type User struct {
	ID string
	Email string
	Password string
}

var (
	root = User{ID: "imroot", Email: "e@e.e", Password: "eee"}
)

func Login(w http.ResponseWriter, r *http.Request) {
	type UserData struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	decoder := json.NewDecoder(r.Body)
	var userData UserData
	err := decoder.Decode(&userData)
	if err != nil {
		log.Println(err)
		return
	}

	if userData.Email == "" || userData.Password == "" {
		BR(w, r, errors.New("Missing credentials"), http.StatusUnauthorized)
		return
	}

	if userData.Email != root.Email || userData.Password != root.Password {
		BR(w, r, errors.New("Bad credentials"), http.StatusUnauthorized)
		return
	} else {
		SetToken(w, r, &root)
	}
}

func SetToken(w http.ResponseWriter, r *http.Request, user *User) {
	t := jwt.New(jwt.GetSigningMethod("RS256"))
	t.Claims["ID"] = user.ID
	t.Claims["iat"] = time.Now().Unix()
	t.Claims["exp"] = time.Now().Add(time.Minute * 60 * 24 * 14).Unix()
	tokenString, err := t.SignedString(signKey)
	if err != nil {
		ISR(w, r, err)
		return
	}
	ServeJSON(w, r, &Response{"token": tokenString}, http.StatusOK)
	return
}