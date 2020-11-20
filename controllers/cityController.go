package controllers

import (
	"fmt"
	"city/models"
	u "city/utils"
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
)

var GetCity = func(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)
	id := (params["id"])

	data := models.GetCity(id,r.Context().Value("dbname").(string))
	resp := u.Message(true, "success")
	resp["data"] = data
	u.Respond(w, resp)
}

var GetCities = func(w http.ResponseWriter, r *http.Request){
	fmt.Println("there")
	data := models.GetCities(r.Context().Value("dbname").(string))
	resp := u.Message(true,"success")
	resp["data"] = data
	u.Respond(w, resp)
}

var CreateCity = func(w http.ResponseWriter, r *http.Request) {
	city := &models.City{}

	err := json.NewDecoder(r.Body).Decode(city)
	if err != nil {
		u.Respond(w, u.Message(false,"error while decoding request body"))
		return
	}

	resp := city.Create(r.Context().Value("dbname").(string))
	u.Respond(w, resp)
}