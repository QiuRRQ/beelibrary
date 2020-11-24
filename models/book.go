package models

import (
	u "city/utils"
	"fmt"
	"log"

	"github.com/jinzhu/gorm"
)

type Books struct {
	Id         int     `json:"id"`
	Name       string  `json:"name"`
	Author     string  `json:"author"`
	Isbn       string  `json:"isbn"`
	Isbn13     string  `json:"isbn13"`
	Genre      string  `json:"genre"`
	Language   string  `json:"language"`
	Date_pub   string  `json:"date_pub"`
	Pages      string  `json:"pages"`
	Sinopsis   string  `json:"sinopsis"`
	Price      float64 `json:"price"`
	Created_at string  `json:"created_at"`
	Updated_at string  `json:"updated_at"`
	Fineamt    float64 `json:"fineamt"`
}

type BookDetail struct {
	DataBook      *Books `json:"data_book"`
	ThisBookStock *Stock `json:"stock_book"`
}

//buku terbaru berdasarkan tanggal publish
func GetNewestBook(pages string, perpage string, mydb *gorm.DB) []*Books {

	book := make([]*Books, 0)
	err := mydb.Table("book").Joins("join stock on book.id = stock.book_id").Order("date_pub ASC").Limit(perpage).Offset(pages).Find(&book).Error

	if err != nil {
		log.Fatal(err)
		return book
	}
	return book
}

//buku terpopuler
func GetPopularBook(pages string, perpage string, mydb *gorm.DB) []*Books {

	book := make([]*Books, 0)
	err := mydb.Table("borrowd").Joins("book on borrowd.book_id = book.id").Order("date_pub ASC").Limit(perpage).Offset(pages).Find(&book).Error

	if err != nil {
		log.Fatal(err)
		return book
	}
	return book
}

//detail book
func GetBooksByID(id int, mydb *gorm.DB) *BookDetail {
	book := &Books{}
	stock := &Stock{}
	err := mydb.Table("book").Joins("join stock on book.id = stock.book_id").Where("book.id = ?", id).First(book).Error
	if err != nil {
		return nil
	}

	err = mydb.Table("stock").Joins("join book on book.id = stock.book_id").Where("book.id = ?", id).First(stock).Error
	if err != nil {
		return nil
	}

	res := &BookDetail{
		DataBook:      book,
		ThisBookStock: stock,
	}

	return res
}

func (book *Books) CreatBook(mydb *gorm.DB) (map[string]interface{}, *Books) {

	err := mydb.Table("book").Create(&book).Error
	if err != nil {
		fmt.Println(err)
		return nil, nil
	}

	resp := u.Message(true, "success")
	resp["data"] = book

	return resp, book
}

func FindByID(id int, mydb *gorm.DB) *Books {
	book := &Books{}
	err := mydb.Table("book").Where("id = ?", id).First(book).Error
	if err != nil {
		return nil
	}

	return book
}

func (book *Books) UpdateBook(id int, mydb *gorm.DB) (map[string]interface{}, *Books) {

	fmt.Println(id)
	fmt.Println(book)

	err := mydb.Table("book").Where("id = ?", id).Update(book).Error
	if err != nil {
		fmt.Println(err)
		return nil, nil
	}

	resp := u.Message(true, "success")
	resp["data"] = book

	return resp, book
}
