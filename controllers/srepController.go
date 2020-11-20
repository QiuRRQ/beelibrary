package controllers

import (
	"city/models"
	u "city/utils"
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

var GetSrep = func(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		//The passed path parameter is not an integer
		u.Respond(w, u.Message(false, "There was an error in your request"))
		return
	}

	data := models.GetSrep(id,r.Context().Value("dbname").(string))
	resp := u.Message(true, "success")
	resp["data"] = data
	u.Respond(w, resp)
}

var GetSreps = func(w http.ResponseWriter, r *http.Request) {
	data := models.GetSreps(r.Context().Value("dbname").(string))
	resp := u.Message(true,"success")
	resp["data"] = data
	u.Respond(w, resp)
}

var CreateSrep = func(w http.ResponseWriter, r *http.Request) {
	srep := &models.Srep{}

	err := json.NewDecoder(r.Body).Decode(srep)
	if err != nil {
		u.Respond(w, u.Message(false,"error while decoding request body"))
		return
	}

	resp := srep.Create(r.Context().Value("dbname").(string))
	u.Respond(w, resp)
}