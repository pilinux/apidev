// Package handler receives input from controllers, performs specific tasks or CRUD
// (Create, Read, Update, Delete, etc.) operations and return the results to the controllers
package handler

import (
	"net/http"
	"strings"
	"time"

	gdatabase "github.com/pilinux/gorest/database"
	gmodel "github.com/pilinux/gorest/database/model"
	log "github.com/sirupsen/logrus"

	"apidev/database/model"
)

// GetUserProfile handles jobs for controller.GetUserProfile
func GetUserProfile(userIDAuth uint64) (httpResponse gmodel.HTTPResponse, httpStatusCode int) {
	db := gdatabase.GetDB()
	user := model.User{}

	// does the user have an existing profile
	if err := db.Where("id_auth = ?", userIDAuth).First(&user).Error; err != nil {
		httpResponse.Message = "user profile not found"
		httpStatusCode = http.StatusNotFound
		return
	}

	// return user profile
	httpResponse.Message = user
	httpStatusCode = http.StatusOK
	return
}

// CreateUserProfile handles jobs for controller.CreateUserProfile
func CreateUserProfile(userIDAuth uint64, user model.User) (httpResponse gmodel.HTTPResponse, httpStatusCode int) {
	db := gdatabase.GetDB()
	userFinal := model.User{}

	// remove all leading and trailing white spaces
	user.NickName = strings.TrimSpace(user.NickName)
	if user.NickName == "" {
		httpResponse.Message = "user nickname is required"
		httpStatusCode = http.StatusBadRequest
		return
	}

	// does the user have an existing profile
	if err := db.Where("id_auth = ?", userIDAuth).First(&userFinal).Error; err == nil {
		httpResponse.Message = "user profile found, no need to create a new one"
		httpStatusCode = http.StatusForbidden
		return
	}

	// security: user must not be able to manipulate all fields
	userFinal.NickName = user.NickName
	userFinal.IDAuth = userIDAuth

	// save in DB
	tx := db.Begin()
	if err := tx.Create(&userFinal).Error; err != nil {
		tx.Rollback()
		log.WithError(err).Error("error code: 1111")
		httpResponse.Message = "internal server error"
		httpStatusCode = http.StatusInternalServerError
		return
	}
	tx.Commit()

	httpResponse.Message = userFinal
	httpStatusCode = http.StatusCreated
	return
}

// UpdateUserProfile handles jobs for controller.UpdateUserProfile
func UpdateUserProfile(userIDAuth uint64, user model.User) (httpResponse gmodel.HTTPResponse, httpStatusCode int) {
	db := gdatabase.GetDB()
	userFinal := model.User{}

	// remove all leading and trailing white spaces
	user.NickName = strings.TrimSpace(user.NickName)
	if user.NickName == "" {
		httpResponse.Message = "user nickname is required"
		httpStatusCode = http.StatusBadRequest
		return
	}

	// does the user have an existing profile
	if err := db.Where("id_auth = ?", userIDAuth).First(&userFinal).Error; err != nil {
		httpResponse.Message = "no user profile found"
		httpStatusCode = http.StatusNotFound
		return
	}

	// if no new info is received, abort
	if user.NickName == userFinal.NickName {
		httpResponse.Message = "no new info to update"
		httpStatusCode = http.StatusBadRequest
		return
	}

	// security: user must not be able to manipulate all fields
	userFinal.UpdatedAt = time.Now()
	userFinal.NickName = user.NickName

	// update in DB
	tx := db.Begin()
	if err := tx.Save(&userFinal).Error; err != nil {
		tx.Rollback()
		log.WithError(err).Error("error code: 1121")
		httpResponse.Message = "internal server error"
		httpStatusCode = http.StatusInternalServerError
		return
	}
	tx.Commit()

	httpResponse.Message = userFinal
	httpStatusCode = http.StatusOK
	return
}
