package models

import (
	//u "crypt-coin-payment/src/utils"
	//"fmt"
	"github.com/jinzhu/gorm"
)

type Permission struct {
	gorm.Model
	Name string `json:name`
	Description string `json:description`

}

type Role struct {
	gorm.Model
	Name string `json:"name"`
	Description string `json: Description`
	Permission []Permission
}

