package controllers

import (
	"city/models"
	d "city/models"
	u "city/utils"
	"encoding/json"
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

var CreateBook = func(w http.ResponseWriter, r *http.Request) {
	Input := &models.BookDetail{}
	db := d.GetDB()
	err := json.NewDecoder(r.Body).Decode(Input)
	if err != nil {
		fmt.Print(err)
	}

	data, insertedData := Input.DataBook.CreatBook(db)

	Input.ThisBookStock.Book_id = insertedData.Id

	Input.ThisBookStock.CreatStock(db)

	data["data"] = insertedData

	db.Close()

	u.Respond(w, data)
}

var UpdateBook = func(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := (params["id"])

	Input := &models.BookDetail{}
	db := d.GetDB()
	err := json.NewDecoder(r.Body).Decode(Input)
	if err != nil {
		fmt.Print(err)
	}

	myID, errr := strconv.Atoi(id)
	if errr != nil {
		fmt.Println(err)
	}
	data, insertedData := Input.DataBook.UpdateBook(myID, db)

	Input.ThisBookStock.Book_id = insertedData.Id

	Input.ThisBookStock.UpdateStock(myID, db)

	data["data"] = insertedData

	db.Close()

	u.Respond(w, data)
}

var GetPopularBook = func(w http.ResponseWriter, r *http.Request) {

	db := d.GetDB()
	query := r.URL.Query()
	pages := query.Get("pages")
	perpages := query.Get("perpages")
	data := d.GetPopularBook(
		pages, perpages, db)

	books := d.GetBooks(db)

	var popularBooks []d.Books
	for _, e := range data {
		for _, j := range books {
			if e.Id == j.Id {
				popularBooks = append(popularBooks, *j)
			}
		}
	}

	db.Close()
	resp := u.Message(true, "success")

	resp["data"] = popularBooks
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
