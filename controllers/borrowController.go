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

//input peminjaman
var BorrowC = func(w http.ResponseWriter, r *http.Request) {
	Inputborrow := &models.InputBorrowd{}

	db := d.GetDB()
	err := json.NewDecoder(r.Body).Decode(Inputborrow)
	if err != nil {
		fmt.Print(err)
	}

	data, insertedData := Inputborrow.BorrowCard.Borrowing(Inputborrow.BorrowCard, db)

	for _, e := range Inputborrow.BorrowdBook {
		e.Borrow_id = insertedData.Id
		data = e.Borrowed(db)

		stock := d.CheckStockByID(e.Book_id, db)
		stock.Qty = stock.Qty - e.Qty
		stock.UpdateStock(e.Book_id, db)
	}

	db.Close()

	u.Respond(w, data)
}

//get detail peminjaman
var BorrowDetail = func(w http.ResponseWriter, r *http.Request) {
	db := d.GetDB()
	params := mux.Vars(r)
	id := (params["id"])

	Borrow_id, err := strconv.Atoi(id)
	if err != nil {
		fmt.Println(err)
	}
	data := d.GetBorrowByID(Borrow_id, db)
	fmt.Println(data)
	db.Close()

	resp := u.Message(true, "success")
	resp["data"] = data
	u.Respond(w, resp)
}

//get daftar pinjaman
var Borrowing = func(w http.ResponseWriter, r *http.Request) {
	db := d.GetDB()
	query := r.URL.Query()
	id := query.Get("usr_id")

	Usr_id, err := strconv.Atoi(id)
	if err != nil {
		fmt.Println(err)
	}
	data := d.GetBorrowByUser(Usr_id, db)
	fmt.Println(data)
	db.Close()

	resp := u.Message(true, "success")
	resp["data"] = data
	u.Respond(w, resp)
}
