package controllers

import (
	"crypt-coin-payment/models"
	"encoding/json"
	"net/http"
	//"u"
	u "crypt-coin-payment/utils"
)


var CreateOrder = func(w http.ResponseWriter, r *http.Request) {
	order := &models.Order{}
	userId := r.Context().Value("user") . (uint)
	err := json.NewDecoder(r.Body).Decode(order)
	if err != nil {
		u.Respond(w, u.Message(false, "Error while decoding request body"))
		return
	}
	application := models.GetApplication(order.ApplicationId)
	if application == nil || application.UserId != userId {
		u.Respond(w, u.Message(false, "Application Id not found"))
		return
	}

	resp := order.Create()
	u.Respond(w, resp)
}