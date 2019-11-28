package models

import (
	"crypt-coin-payment/keychain"
	u "crypt-coin-payment/utils"
	"fmt"
	"github.com/jinzhu/gorm"
)

type Application struct {
	gorm.Model
	Name string `json:"name"`
	UserId uint `json:"user_id"`
}


type ApplicationPublicKey struct {
	gorm.Model
	PublicKey string `json:"public_key"`
	ApplicationId uint `json:"application_id"`
	Index uint `json:"index"`
	NumOfAddressGenerated uint `json: num_address_generated`
}

/*
 This struct function validate the required parameters sent through the http request body

returns message and true if the requirement is met
*/
func (application *Application) Validate() (map[string] interface{}, bool) {

	if application.Name == "" {
		return u.Message(false, "Application name should be on the payload"), false
	}
	//All the required parameters are present
	return u.Message(true, "success"), true
}

func (application *Application) Create() (map[string] interface{}) {

	if resp, ok := application.Validate(); !ok {
		return resp
	}
	applicationPublicKey := &ApplicationPublicKey{}

	masterPublicKey := GetMasterPublicKeyByUser(application.UserId)
	if masterPublicKey == nil {

	}
	keyService := keychain.KeyFactory(1)
	genNumber := GetUserApplicationCount(application.UserId)
	appPublicKey, err := keyService.GenerateAccountFromMasterPubKey(masterPublicKey.PublicKey, uint32(genNumber))
	if err != nil {
		return u.Message(false, "Master public key not found")
	}

	GetDB().Create(application)

	applicationPublicKey.ApplicationId = application.ID
	applicationPublicKey.Index = GetUserApplicationCount(application.UserId) - 1
	applicationPublicKey.PublicKey = appPublicKey
	applicationPublicKey.NumOfAddressGenerated = 0

	err = GetDB().Create(applicationPublicKey).Error
	if err != nil {
		fmt.Println(err)
	}

	resp := u.Message(true, "success")
	resp["application"] = application
	resp["extend_public_key"] = applicationPublicKey
	return resp
}

func GetApplication(id uint) (*Application) {

	application := &Application{}
	err := GetDB().Table("applications").Where("id = ?", id).First(application).Error
	if err != nil {
		return nil
	}
	return application
}

func GetUserApplicationCount(userId uint) uint{
	var count uint
	db.Table("applications").Where("user_id = ?", userId).Count(&count)
	return count
}

func GetApplicationPubicKey(applicationId uint) (*ApplicationPublicKey)  {
	appPubKey := &ApplicationPublicKey{}
	err := GetDB().Table("application_public_keys").Where("application_id = ?", applicationId).First(appPubKey).Error
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return appPubKey
}

