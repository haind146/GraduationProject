package models

import (
	u "crypt-coin-payment/utils"
	"github.com/jinzhu/gorm"
)

type MasterPublicKey struct {
	gorm.Model
	PublicKey string `json:"public_key"`
	GenerateNumber uint `json:"generate_number"`
	UserId uint `json:"user_id"` //The user that this key belongs to
}


func (masterPublicKey *MasterPublicKey) Create() (map[string] interface{}) {
	isMasterPublicKeyImported := GetMasterPublicKeyByUser(masterPublicKey.UserId)
	if (isMasterPublicKeyImported != nil) {
		return u.Message(false, "Master public key already imported")
	}

    isMasterPubKey := u.ValidateMasterPublicKey(masterPublicKey.PublicKey)
    if !isMasterPubKey {
    	return u.Message(false, "This master public key is not valid")
	}
    masterPublicKey.GenerateNumber = 0

	result := GetDB().Create(masterPublicKey)
	if result.Error != nil {
		return u.Message(false, "Please try another key")
	}

	resp := u.Message(true, "success")
	resp["public_key"] = masterPublicKey
	return resp
}

func GetMasterPublicKeyByUser(userId uint) (*MasterPublicKey) {
	masterPublicKey := &MasterPublicKey{}
	err := GetDB().Table("master_public_keys").Where("user_id = ?", userId).Last(masterPublicKey).Error
	if err != nil {
		return nil
	}
	return masterPublicKey
}
