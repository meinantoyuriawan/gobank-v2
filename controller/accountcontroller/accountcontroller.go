package accountcontroller

import (
	"encoding/json"
	"net/http"

	"github.com/meinantoyuriawan/gobank-v2/auth"
	"github.com/meinantoyuriawan/gobank-v2/helper"
	"github.com/meinantoyuriawan/gobank-v2/models"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func Login(w http.ResponseWriter, r *http.Request) {
	userInput := models.User{}
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&userInput); err != nil {
		response := map[string]string{"message": err.Error()}
		helper.ResponseJSON(w, http.StatusBadRequest, response)
		return
	}

	defer r.Body.Close()

	// get user data
	user := models.User{}
	// process the password (hash and match it with db)
	if err := models.DB.Where("Username = ?", userInput.Username).First(&user).Error; err != nil {
		switch err {
		case gorm.ErrRecordNotFound:
			response := map[string]string{"message": err.Error()}
			helper.ResponseJSON(w, http.StatusBadRequest, response)
			return
		default:
			response := map[string]string{"message": err.Error()}
			helper.ResponseJSON(w, http.StatusInternalServerError, response)
			return
		}
	}

	// verify is password valid
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(userInput.Password)); err != nil {
		response := map[string]string{"message": "Username atau password salah"}
		helper.ResponseJSON(w, http.StatusUnauthorized, response)
		return
	}

	token, err := auth.GenerateJWT(w, r, user.Username)
	if err != nil {
		response := map[string]string{"message": err.Error()}
		helper.ResponseJSON(w, http.StatusInternalServerError, response)
		return
	}

	// set token to cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Path:     "/",
		Value:    token,
		HttpOnly: true,
	})
	response := map[string]string{"message": "logged in"}
	helper.ResponseJSON(w, http.StatusOK, response)
}

func Registeruser(w http.ResponseWriter, r *http.Request) {
	userInput := models.User{}
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&userInput); err != nil {
		response := map[string]string{"message": err.Error()}
		helper.ResponseJSON(w, http.StatusBadRequest, response)
		return
	}

	defer r.Body.Close()

	// hash password with bcrypt
	hashPassword, _ := bcrypt.GenerateFromPassword([]byte(userInput.Password), bcrypt.DefaultCost)
	userInput.Password = string(hashPassword)

	// insert into db
	if err := models.DB.Create(&userInput).Error; err != nil {
		response := map[string]string{"message": err.Error()}
		helper.ResponseJSON(w, http.StatusInternalServerError, response)
		return
	}

	response := map[string]string{"message": "success"}
	helper.ResponseJSON(w, http.StatusOK, response)
}

func Createaccount(w http.ResponseWriter, r *http.Request) {
	userInput := models.Account{}
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&userInput); err != nil {
		response := map[string]string{"message": err.Error()}
		helper.ResponseJSON(w, http.StatusBadRequest, response)
		return
	}

	defer r.Body.Close()

	hashPin, _ := bcrypt.GenerateFromPassword([]byte(userInput.Pin), bcrypt.DefaultCost)
	userInput.Pin = string(hashPin)

	// insert into db
	if err := models.DB.Create(&userInput).Error; err != nil {
		response := map[string]string{"message": err.Error()}
		helper.ResponseJSON(w, http.StatusInternalServerError, response)
		return
	}

	// create status db relation
	userStatus := models.Status{
		UserId:     userInput.UserId,
		AccNumber:  userInput.Number,
		PinAttempt: 2,
		IsBlocked:  0,
	}

	if err := models.DB.Create(&userStatus).Error; err != nil {
		response := map[string]string{"message": err.Error()}
		helper.ResponseJSON(w, http.StatusInternalServerError, response)
		return
	}

	response := map[string]string{"message": "success"}
	helper.ResponseJSON(w, http.StatusOK, response)
}

func Logout(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Path:     "/",
		Value:    "",
		HttpOnly: true,
		MaxAge:   -1,
	})
	response := map[string]string{"message": "logged out"}
	helper.ResponseJSON(w, http.StatusOK, response)
}
