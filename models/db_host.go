package models

import (
	"log"
)

type Dbhost struct {
	Dbhost string
}

func GetDbhost(dbname string) (*Dbhost) {
	dbhost := &Dbhost{}
	conn := GetDB()
	err := conn.Table("vwregistrasi").Where("dbname = ?", dbname).First(&dbhost).Error
	if err != nil {
		log.Println(err)
		return nil
	}
	defer conn.Close()
	return dbhost
}