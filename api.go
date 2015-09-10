package main

import (
	"encoding/json"
	"net/http"
	"fmt"
)

func Me(w http.ResponseWriter, r *http.Request) {
	tokenData := GetToken(w, r)

	fmt.Printf("Token Data: %v", tokenData)

	b, _ := json.Marshal(tokenData)
	parse := &Response{}
	json.Unmarshal(b, parse)
	ServeJSON(w, r, parse, http.StatusOK)
}