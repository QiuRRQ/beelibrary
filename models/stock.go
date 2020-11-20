package models

type Stock struct {
	Book_id    int    `json:"book_id"`
	Qty        int    `json:"qty"`
	Created_at string `json:"created_at"`
	Updated_at string `json:"updated_at"`
}
