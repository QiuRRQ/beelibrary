package models

import (
	u "city/utils"
	"fmt"

	"github.com/jinzhu/gorm"
)

//ini detail peminjamannya
type Borrowd struct {
	Borrow_id int     `json:"borrow_id"`
	Book_id   int     `json:"book_id"`
	Qty       int     `json:"qty"`
	Price     float64 `json:"price"`
	Subtotal  float64 `json:"subtotal"`
}

func (borrowd *Borrowd) Borrowed(mydb *gorm.DB) map[string]interface{} {

	err := GetDB().Table("borrowd").Create(&borrowd).Error
	if err != nil {
		fmt.Println(err)
		return nil
	}
	resp := u.Message(true, "success")
	resp["data"] = borrowd

	return resp
}
