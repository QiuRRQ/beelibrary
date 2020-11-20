package controllers

import (
	d "city/models"
	u "city/utils"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

var GetNewestBook = func(w http.ResponseWriter, r *http.Request) {

	db := d.GetDB()
	query := r.URL.Query()
	pages := query.Get("pages")
	perpages := query.Get("perpages")
	data := d.GetNewestBook(
		pages, perpages, db)

	db.Close()
	resp := u.Message(true, "success")

	resp["data"] = data
	u.Respond(w, resp)
}

var GetBookByID = func(w http.ResponseWriter, r *http.Request) {
	db := d.GetDB()
	params := mux.Vars(r)
	id := (params["id"])

	book_id, err := strconv.Atoi(id)
	if err != nil {
		fmt.Println(err)
	}
	data := d.GetBooksByID(book_id, db)

	db.Close()

	resp := u.Message(true, "success")
	resp["data"] = data
	u.Respond(w, resp)
}
