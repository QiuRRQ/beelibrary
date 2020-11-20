package models

import (
	"fmt"
	"log"
	u "city/utils"
	"time"
)

type City struct {
	Code string  `json:"code"`
	Name string `json:"name"`
	Updated_at string `json:"updated_at"` //The user that this contact belongs to
}

func GetCity(id string, dbname string) (*City) {
	
	log.Println(id)
	city := &City{}
	err := GetDB()
	if err != nil {
		return nil
	}
	return city
}

func GetCities(dbname string) ([]*City) {

	city := make([]*City, 0)
	err := GetDB()
	if err != nil {
		fmt.Println(err)
		return nil
	}

	return city
}

func (city *City) Validate() (map[string] interface{}, bool){
	if city.Code == "" {
		return u.Message(false, "Code tidak boleh kosong!"), false
	}

	if city.Name == "" {
		return u.Message(false, "Nama tidak boleh kosong!"), false
	}

	return u.Message(true,"success"), true
}

func (city *City) Create(dbname string) (map[string] interface{})  {

	if resp, ok := city.Validate(); !ok {
		return resp
	}
	loc, _ := time.LoadLocation("Asia/Jakarta")

	//set timezone,
	now := time.Now().In(loc)

	city.Updated_at = now.Format("2006-01-02	15:04:05")
	err := GetDB().Table("city").Create(&city).Error
	if err != nil {
		fmt.Println(err)
		return nil
	}
	resp := u.Message(true,"success")
	resp["city"] = city

	return resp
}