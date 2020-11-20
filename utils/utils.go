package utils

import (
	"fmt"
	"encoding/json"
	"net/http"
)

func Message(status bool, message string) (map[string]interface{}) {
	fmt.Println(1)
	return map[string]interface{} {"status" : status, "message" : message}
}

func Respond(w http.ResponseWriter, data map[string] interface{})  {
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}