package models

import (
	"github.com/jinzhu/gorm"
)

type Address struct {
	gorm.Model
	Address string `json:"address"`
	PaymentMethodId uint `json:"payment_method_id"`
	OrderId uint `json:"order_id"`
	Balance uint `json:"balance"`
	PendingSent uint `json:"pending_sent"`
	PendingReceive uint `json:"pending_receive"`
	MnemonicPath uint `json:"mnemonic_path"`
}

//func (address *Address) Create() (map[string] interface{}) {
//	GetDB().Create(address)
//}

func GetAddress(id uint) (*Address) {

	address := &Address{}
	err := GetDB().Table("addresses").Where("id = ?", id).First(address).Error
	if err != nil {
		return nil
	}
	return address
}
