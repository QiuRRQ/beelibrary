package models

import (
	u "city/utils"
	"fmt"

	"github.com/jinzhu/gorm"
)

//ini kartu pinjamnya
type Borrow struct {
	Id         int     `json:"id"`
	Start_date string  `json:"start_date"`
	End_date   string  `json:"end_date"`
	Usr_id     int     `json:"usr_id"`
	Status     string  `json:"status"`
	Created_at string  `json:"created_at"`
	Updated_at string  `json:"updated_at"`
	Total      float64 `json:"total"`
	Terlambat  bool    `json:"terlambat"`
}

type InputBorrowd struct {
	BorrowCard  Borrow    `json:"borrow_card"`
	BorrowdBook []Borrowd `json:"borrowd_book"`
}

type DetailsBorrow struct {
	BorrowCard  Borrow     `json:"borrow_card"`
	BorrowdBook []*Borrowd `json:"borrowd_book"`
}

func (borrow *Borrow) Borrowing(borrowCard Borrow, mydb *gorm.DB) (map[string]interface{}, *Borrow) {

	err := mydb.Table("borrow").Create(&borrow).Error
	if err != nil {
		fmt.Println(err)
		return nil, nil
	}

	resp := u.Message(true, "success")
	resp["data"] = borrow

	return resp, borrow
}

//detail book
func GetBorrowByID(id int, mydb *gorm.DB) *DetailsBorrow {
	borrow := &Borrow{}
	err := GetDB().Table("borrow").Where("id = ?", id).First(borrow).Error
	if err != nil {
		return nil
	}

	borrowd := make([]*Borrowd, 0)
	err = GetDB().Table("borrowd").Where("borrow_id = ?", id).Find(&borrowd).Error
	if err != nil {
		return nil
	}

	res := &DetailsBorrow{
		BorrowCard:  *borrow,
		BorrowdBook: borrowd,
	}

	return res
}
