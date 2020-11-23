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
	BorrowCard    Borrow        `json:"borrow_card"`
	MyBorrowdBook []BorrowdBook `json:"borrowd_book"`
}

type BorrowdBook struct {
	Borrow_id  int     `json:"borrow_id"`
	Book_id    int     `json:"book_id"`
	Qty        int     `json:"qty"`
	Price      float64 `json:"price"`
	Subtotal   float64 `json:"subtotal"`
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
	Created_at string  `json:"created_at"`
	Updated_at string  `json:"updated_at"`
	Fineamt    float64 `json:"fineamt"`
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
	// var resa []BorrowdBook
	var borrowdBook []BorrowdBook
	// var books []Books
	// var book Books
	err := mydb.Table("borrow").Where("id = ?", id).First(borrow).Error
	if err != nil {
		fmt.Println(2)
		return nil
	}

	mydb.Table("borrowd").Select("*").Joins("join book on book.id = borrowd.book_id").Where("borrowd.borrow_id = ?", id).Scan(&borrowdBook)
	// rows, errr := mydb.Table("borrowd").Select("*").Joins("join book on book.id = borrowd.book_id").Where("borrowd.borrow_id = ?", id).Rows()
	// defer rows.Close()
	// if errr != nil {
	// 	fmt.Println(3)
	// 	return nil
	// }

	// for rows.Next() {
	// 	mydb.ScanRows(rows, &borrowdBook)
	// 	resa = append(resa, borrowdBook)
	// }

	// rows, errr = mydb.Table("book").Joins("join borrowd on book.id = borrowd.book_id").Where("borrowd.borrow_id = ?", id).Rows()
	// defer rows.Close()
	// if errr != nil {
	// 	fmt.Println(3)
	// 	return nil
	// }

	// for rows.Next() {
	// 	mydb.ScanRows(rows, &book)
	// 	books = append(books, book)
	// }

	res := &DetailsBorrow{
		BorrowCard:    *borrow,
		MyBorrowdBook: borrowdBook,
	}

	return res
}

//get daftar pinjaman
func GetBorrowByUser(id int, mydb *gorm.DB) []DetailsBorrow {
	borrows := make([]*Borrow, 0)
	var borrowdBook []BorrowdBook
	err := mydb.Table("borrow").Where("usr_id = ?", id).Find(&borrows).Error
	if err != nil {
		fmt.Println(2)
		return nil
	}

	var res []DetailsBorrow

	mydb.Table("borrowd").Select("*").Joins("join book on book.id = borrowd.book_id").Joins("join borrow on borrow.id = borrowd.borrow_id").Where("borrow.usr_id = ?", id).Scan(&borrowdBook)

	for _, e := range borrows {
		for _, k := range borrowdBook {
			var thisBook []*BorrowdBook
			if e.Id == k.Borrow_id {
				thisBook = append(thisBook, &BorrowdBook{
					Borrow_id: k.Borrow_id,
					Name:      k.Name,
				})
				fmt.Println("append here")
			} else {
				fmt.Println(k.Name + k.Author)
				fmt.Println(thisBook)
				res = append(res, DetailsBorrow{
					BorrowCard: *e,
				})
			}
		}

	}

	return res
}
