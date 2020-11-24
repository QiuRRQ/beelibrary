package models

import (
	u "city/utils"
	"fmt"

	"github.com/jinzhu/gorm"
)

type Stock struct {
	Book_id    int    `json:"book_id"`
	Qty        int    `json:"qty"`
	Created_at string `json:"created_at"`
	Updated_at string `json:"updated_at"`
}

func (stock *Stock) CreatStock(mydb *gorm.DB) (map[string]interface{}, *Stock) {

	err := mydb.Table("stocks").Create(&stock).Error
	if err != nil {
		fmt.Println(err)
		return nil, nil
	}

	resp := u.Message(true, "success")
	resp["data"] = stock

	return resp, stock
}

func (stock *Stock) UpdateStock(id int, mydb *gorm.DB) (map[string]interface{}, *Stock) {

	err := mydb.Table("stocks").Where("book_id = ?", id).Update(stock).Error
	if err != nil {
		fmt.Println(err)
		return nil, nil
	}

	resp := u.Message(true, "success")
	resp["data"] = stock

	return resp, stock
}

func CheckStockByID(id int, mydb *gorm.DB) *Stock {
	stock := &Stock{}
	err := mydb.Table("stocks").Where("book_id = ?", id).First(stock).Error
	if err != nil {
		return nil
	}

	return stock
}
