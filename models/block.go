package models

import (
	"github.com/jinzhu/gorm"
)

type Block struct {
	gorm.Model
	BlockHash string `json:"block_hash"`
	BlockNumber string `json:"block_number"`
	PaymentMethodId string `json:"payment_method_id"`
}

//func (block *Block) Create() (map[string] interface{}) {
//	//GetDB().Create(address)
//}

//func GetBlock(id uint) (*Address) {
//	//
//	//address := &Address{}
//	//err := GetDB().Table("addresses").Where("id = ?", id).First(application).Error
//	//if err != nil {
//	//	return nil
//	//}
//	//return address
//}
