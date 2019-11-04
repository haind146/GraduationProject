package controllers

import (
	"crypt-coin-payment/models"
	"encoding/json"
	"net/http"
	//"u"
	u "crypt-coin-payment/utils"
)



var CreateApplication = func(w http.ResponseWriter, r *http.Request) {
	application := &models.Application{}
	userId := r.Context().Value("user") . (uint)
	err := json.NewDecoder(r.Body).Decode(application)
	if err != nil {
		u.Respond(w, u.Message(false, "Error while decoding request body"))
		return
	}
	application.UserId = userId
	resp := application.Create()
	u.Respond(w, resp)
}
