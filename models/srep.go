package models

import (
	u "city/utils"
	"fmt"
	"log"

)

type Srep struct {
	Id       int    `json:"id"`
	Code     string `json:"code"`
	Name     string `json:"name"`
	CityCode *string `json:"city_code"`
	Address	 string `json:"address"`
	Zipcode  string `json:"zipcode"`
	Phone    string `json:"phone"`
	Mobile   string `json:"mobile"`
	Email    string `json:"email"`
	Active   bool	`json:"active"`
}

func (srep *Srep) Validate() (map[string] interface{}, bool){
	if srep.Code == "" {
		return u.Message(false, "Code tidak boleh kosong!"), false
	}

	if srep.Name == "" {
		return u.Message(false, "Nama tidak boleh kosong!"), false
	}

	return u.Message(true,"success"), true
}

func GetSrep(id int, dbname string) (*Srep)  {

	log.Println(id)
	srep := &Srep{}
	err := GetDB().Table("srep").Where("id = ?", id).First(srep).Error
	if err != nil {
		return nil
	}
	return srep
}

func GetSreps(dbname string) ([]*Srep)  {

	srep := make([]*Srep,0)
	err := GetDB().Table("srep").Find(&srep).Error
	if err != nil {
		return nil
	}
	return srep
}

func (srep *Srep) Create(dbname string) (map[string] interface{}){

	if resp, ok := srep.Validate(); !ok {
		return resp
	}
	err := GetDB().Table("srep").Create(&srep).Error
	if err != nil {
		fmt.Println(err)
		return nil
	}
	resp := u.Message(true,"success")
	resp["srep"] = srep

	return resp
}