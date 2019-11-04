package models

import (
	"crypt-coin-payment/service"
	u "crypt-coin-payment/utils"
	"github.com/jinzhu/gorm"
)

type Order struct {
	gorm.Model
	PartnerOrderId string `json:"partner_order_id"`
	Amount uint `json:"amount"`
	ReceiveAddress string `json:"receive_address"`
	PaymentMethodId uint `json:"payment_method_id"`
	ApplicationId uint `json:"application_id"`
	Status uint `json:"status"`
	ExpiredTime uint `json:"expired_time"`
}

/*
 This struct function validate the required parameters sent through the http request body

returns message and true if the requirement is met
*/
func (order *Order) Validate() (map[string] interface{}, bool) {

	if order.Amount <= 0 {
		return u.Message(false, "Amount should greater than 0"), false
	}
	if order.ApplicationId <= 0 {
		return u.Message(false, "Application Id not found"), false
	}
	//All the required parameters are present
	return u.Message(true, "success"), true
}

func (order *Order) Create() (map[string] interface{}) {

	if resp, ok := order.Validate(); !ok {
		return resp
	}

	appPubKey := GetApplicationPubicKey(order.ApplicationId)
	if appPubKey == nil {
		return u.Message(false, "App public key not found")
	}

	keyService := service.KeyFactory(1)
	receiveAddress, err := keyService.GenerateAddressFromAccount(appPubKey.PublicKey, uint32(appPubKey.NumOfAddressGenerated))
	if err != nil {
		u.Message(false, "Error when generate address")
	}

	order.ReceiveAddress = receiveAddress
	order.Status = 0
	appPubKey.NumOfAddressGenerated = appPubKey.NumOfAddressGenerated + 1

	GetDB().Create(order)
	GetDB().Save(appPubKey)
	resp := u.Message(true, "success")
	resp["order"] = order
	return resp
}
