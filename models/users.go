package models

import (
	"github.com/dgrijalva/jwt-go"
	u "crypt-coin-payment/utils"
	"strings"
	"github.com/jinzhu/gorm"
	"os"
	"golang.org/x/crypto/bcrypt"
	"fmt"
)

/*
JWT claims struct
*/
type Token struct {
	UserId uint
	jwt.StandardClaims
}

//a struct to rep user account
type User struct {
	gorm.Model
	Name string `json:"name"`
	Phone string `json:"phone"`
	Email string `gorm:"type:varchar(100);unique_index"`
	Password string `json:"password"`
	Token string `json:token`
	RoleId uint `json:role_id`
	RegisterKey string
	Role Role `gorm:"foreignkey:UserRole"`
}

//Validate incoming user details...
func (user *User) Validate() (map[string] interface{}, bool) {

	if !strings.Contains(user.Email, "@") {
		return u.Message(false, "Email address is required"), false
	}

	if len(user.Password) < 6 {
		return u.Message(false, "Password is required"), false
	}

	//Email must be unique
	temp := &User{}

	//check for errors and duplicate emails
	err := GetDB().Table("accounts").Where("email = ?", user.Email).First(temp).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return u.Message(false, "Connection error. Please retry"), false
	}
	if temp.Email != "" {
		return u.Message(false, "Email address already in use by another user."), false
	}

	return u.Message(false, "Requirement passed"), true
}

func (user *User) Create() (map[string] interface{}) {

	if resp, ok := user.Validate(); !ok {
		return resp
	}

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	user.Password = string(hashedPassword)

	registerUser := GetRegisterUserByEmail(user.Email)
	fmt.Print(user.RegisterKey)

	if len(user.RegisterKey) != 64 || user.RegisterKey != registerUser.RegisterKey {
		return u.Message(false, "This user has not accepted to register new account")
	}
	user.RoleId = 2
	user.Name = registerUser.Name
	user.Phone = registerUser.Phone
	GetDB().Create(user)

	if user.ID <= 0 {
		return u.Message(false, "Failed to create account, connection error.")
	}

	//Create new JWT token for the newly registered account
	tk := &Token{UserId: user.ID}
	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), tk)
	tokenString, _ := token.SignedString([]byte(os.Getenv("token_password")))
	user.Token = tokenString

	user.Password = "" //delete password

	response := u.Message(true, "Account has been created")
	response["account"] = user
	return response
}

func Login(email, password string) (map[string]interface{}) {

	user := &User{}
	err := GetDB().Table("users").Where("email = ?", email).First(user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return u.Message(false, "Email address not found")
		}
		return u.Message(false, "Connection error. Please retry")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword { //Password does not match!
		return u.Message(false, "Invalid login credentials. Please try again")
	}
	//Worked! Logged In
	user.Password = ""

	//Create JWT token
	tk := &Token{UserId: user.ID}
	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), tk)
	tokenString, _ := token.SignedString([]byte(os.Getenv("token_password")))
	user.Token = tokenString //Store the token in the response

	resp := u.Message(true, "Logged In")
	resp["account"] = user
	return resp
}

func GetUser(u uint) *User {

	acc := &User{}
	GetDB().Table("users").Where("id = ?", u).First(acc)
	if acc.Email == "" { //User not found!
		return nil
	}

	acc.Password = ""
	return acc
}
