package controllers

import (
	"crypt-coin-payment/src/models"
	u "crypt-coin-payment/src/utils"
	"encoding/json"
	"net/http"
)

var CreateRegisterUser = func(w http.ResponseWriter, r *http.Request) {
	registerUser := &models.RegisterUser{}
	err := json.NewDecoder(r.Body).Decode(registerUser)
	if err != nil {
		u.Respond(w, u.Message(false, "Error while decoding request body"))
		return
	}
	resp := registerUser.Create()
	u.Respond(w, resp)
}

var AcceptRegisterUser = func(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value("user") . (uint)
	user := models.GetUser(userId)
	if(user.RoleId != 1) {
		u.Respond(w, u.Message(false, "User permission denied"))
		return
	}
	registerUser := &models.RegisterUser{}
	err := json.NewDecoder(r.Body).Decode(registerUser)
	if err != nil {
		u.Respond(w, u.Message(false, "Error while decoding request body"))
		return
	}
	registerUser = models.GetRegisterUserByEmail(registerUser.Email)
	registerUser.CreateRegisterKey()
	resp := registerUser.Save()
	u.Respond(w, resp)
}
