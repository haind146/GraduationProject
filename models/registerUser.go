package models

import (
	u "crypt-coin-payment/utils"
	"fmt"
	"github.com/jinzhu/gorm"
)

const UNACTIVE  = 0
const ACTIVATED  = 1

type RegisterUser struct {
	gorm.Model
	Name string `json:"name"`
	Email string `gorm:"type:varchar(100);unique_index"`
	Phone string `json:"phone"`
	Status uint `json:"status"`
	RegisterKey string `json:"registerKey"`
}

func (registerUser *RegisterUser) Validate() (map[string] interface{}, bool)  {
	if(registerUser.Name == "") {
		return u.Message(false, "User name should be on the payload"), false
	}
	if(!u.ValidateEmail(registerUser.Email)) {
		fmt.Print(registerUser.Email)
		return u.Message(false, "Email is not valid"), false
	}
	if(registerUser.Phone == "") {
		return u.Message(false, "Phone number should be on the payload"), false
	}

	return u.Message(true, "success"), true
}

func (registerUser *RegisterUser) Create() (map[string] interface{}) {
	if resp, ok := registerUser.Validate(); !ok {
		return resp
	}

	registerUser.Status = UNACTIVE
	createResult := GetDB().Create(registerUser)
	if(createResult.Error != nil) {
		return u.Message(false, "Email have already been used")
	}

	resp := u.Message(true, "success")
	resp["registerUser"] = registerUser
	return resp
}

func GetRegisterUser(id uint) (*RegisterUser) {
	registerUser := &RegisterUser{}
	err := GetDB().Table("register_users").Where("id = ?", id).First(registerUser).Error
	if err != nil {
		return nil
	}
	return registerUser
}

func GetRegisterUserByEmail(email string) (*RegisterUser)  {
	registerUser := &RegisterUser{}
	err := GetDB().Table("register_users").Where("email = ?", email).First(registerUser).Error
	if err != nil {
		return registerUser
	}
	return registerUser
}

func (registerUser *RegisterUser) CreateRegisterKey() {
	registerUser.RegisterKey = u.Generate64BytesRandom()
	registerUser.Status = ACTIVATED
}

func (registerUser *RegisterUser) Save() (map[string] interface{}) {
	saveResult := GetDB().Save(&registerUser)
	if(saveResult.Error != nil) {
		return u.Message(false, "Email is invalid")
	}
	resp := u.Message(true, "success")
	resp["registerUser"] = registerUser
	return resp
}