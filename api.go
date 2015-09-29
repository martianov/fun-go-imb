package main

import (
	"encoding/json"
	"net/http"
    "fmt"

    "gopkg.in/mgo.v2/bson"
)

func Me(w http.ResponseWriter, r *http.Request) {
    tokenData := GetToken(w, r)
    db := GetDB(w, r)

    user, errM := FindUserById(db, bson.ObjectIdHex(tokenData.ID))
    if errM != nil {
        fmt.Println("Failed to find user: %v", errM)
        ISR(w, r, errM.Reason)
        return
    }

    b, _ := json.Marshal(user)
    parse := &Response{}
    json.Unmarshal(b, parse)
    ServeJSON(w, r, parse, http.StatusOK)
}