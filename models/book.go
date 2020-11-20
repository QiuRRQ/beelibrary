package models

import (
	"log"

	"github.com/jinzhu/gorm"
)

type Books struct {
	Id         int     `json:"id"`
	Name       string  `json:"name"`
	Author     string  `json:"author"`
	Isbn       string  `json:"isbn"`
	Isbn13     string  `json:"genre"`
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

//detail book
func GetBooksByID(id int, mydb *gorm.DB) *BookDetail {
	book := &Books{}
	stock := &Stock{}
	err := GetDB().Table("book").Joins("join stock on book.id = stock.book_id").Where("book.id = ?", id).First(book).Error
	if err != nil {
		return nil
	}

	err = GetDB().Table("stock").Joins("join book on book.id = stock.book_id").Where("book.id = ?", id).First(stock).Error
	if err != nil {
		return nil
	}

	res := &BookDetail{
		DataBook:      book,
		ThisBookStock: stock,
	}

	return res
}
