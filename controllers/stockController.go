package controllers

import (
	"city/models"
	d "city/models"
	u "city/utils"
	"encoding/json"
	"fmt"
	"net/http"
)

var CreateStock = func(w http.ResponseWriter, r *http.Request) {
	Input := &models.Stock{}

	db := d.GetDB()
	err := json.NewDecoder(r.Body).Decode(Input)
	if err != nil {
		fmt.Print(err)
	}

	data, insertedData := Input.CreatStock(db)

	data["data"] = insertedData

	db.Close()

	u.Respond(w, data)
}
