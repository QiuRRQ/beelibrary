package controllers

import (
	"city/models"
	d "city/models"
	u "city/utils"
	"encoding/json"
	"fmt"
	"net/http"
)

var UserLoginController = func(w http.ResponseWriter, r *http.Request) {
	db := d.GetDB()
	user := &models.Users{}

	err := json.NewDecoder(r.Body).Decode(user)
	if err != nil {
		fmt.Print(err)
	}
	data := user.Login(user.Email, user.Password, db)

	db.Close()
	if data.UserData == nil {
		resp := u.Message(false, "login failed")
		u.Respond(w, resp)
	} else {
		resp := u.Message(true, "success")

		resp["data"] = data
		u.Respond(w, resp)
	}

}
