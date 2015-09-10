package main

import (
	"errors"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/codegangsta/negroni"
	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/context"
)

const publicKey = "keys/app.rsa.pub"

var (
	verifyKey []byte
)

func init() {
	var err error
	verifyKey, err = ioutil.ReadFile(publicKey)
	if err != nil {
		log.Fatal("Error reading Private key")
		return
	}
}

func JWTMiddleware() negroni.Handler {
	return negroni.HandlerFunc(func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		if h := r.Header.Get("Authorization"); h != "" {
			token, err := jwt.ParseFromRequest(r, func(token *jwt.Token) (interface{}, error) {
				return verifyKey, nil
			})

			switch err.(type) {
			case nil:
				if !token.Valid {
					NotAllowed(w, r)
					return
				}
				context.Set(r, "token", token)
				next(w, r)
			case *jwt.ValidationError:
				vErr := err.(*jwt.ValidationError)
				switch vErr.Errors {
				case jwt.ValidationErrorExpired:
					BR(w, r, errors.New("Token Expired"), http.StatusUnauthorized)
					return
				default:
					BR(w, r, errors.New("Bad Token"), http.StatusUnauthorized)
					log.Println(vErr.Error())
					return
				}
			default:
				ISR(w, r, err)
				return
			}
		} else {
			BR(w, r, errors.New("Missing Token"), http.StatusUnauthorized)
		}
	})
}
