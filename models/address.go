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
	PendingSent float64 `json:"pending_sent"`
	PendingReceive float64 `json:"pending_receive"`
	MnemonicPath string `json:"mnemonic_path"`
}

//func (address *Address) Create() (map[string] interface{}) {
//	GetDB().Create(address)
//}

func GetAddress(addr string) *Address {
	address := &Address{}
	err := GetDB().Table("addresses").Where("address = ?", addr).First(address).Error
	if err != nil {
		return nil
	}
	return address
}
