// Package controller contains all controllers
package controller

import (
	"net/http"
	"reflect"

	"github.com/gin-gonic/gin"
	grenderer "github.com/pilinux/gorest/lib/renderer"

	"apidev/database/model"
	"apidev/handler"
)

// GetUserProfile - GET /users
// - find the profile of a logged-in user
// - authID is fetched from the access token
func GetUserProfile(c *gin.Context) {
	userIDAuth := c.GetUint64("authID")

	resp, statusCode := handler.GetUserProfile(userIDAuth)

	if reflect.TypeOf(resp.Message).Kind() == reflect.String {
		grenderer.Render(c, resp, statusCode)
		return
	}

	grenderer.Render(c, resp, statusCode)
}

// CreateUserProfile - POST /users
// - only a registered user can create his personal profile
// - authID is fetched from the access token
// - if the user already has a profile, he is not allowed to create a new one
// ===============================
//
//	{
//	   "nickName": "your_nickname"
//	}
//
// ===============================
func CreateUserProfile(c *gin.Context) {
	userIDAuth := c.GetUint64("authID")
	user := model.User{}

	// bind JSON
	if err := c.ShouldBindJSON(&user); err != nil {
		grenderer.Render(c, gin.H{"message": err.Error()}, http.StatusBadRequest)
		return
	}

	resp, statusCode := handler.CreateUserProfile(userIDAuth, user)

	if reflect.TypeOf(resp.Message).Kind() == reflect.String {
		grenderer.Render(c, resp, statusCode)
		return
	}

	grenderer.Render(c, resp.Message, statusCode)
}

// UpdateUserProfile - PUT /users
// - only a registered user can update his existing personal profile
// - authID is fetched from the access token
// - if the user has no profile, he has to create a new one
// ===================================
//
//	{
//	   "nickName": "your_new_nickname"
//	}
//
// ===================================
func UpdateUserProfile(c *gin.Context) {
	userIDAuth := c.GetUint64("authID")
	user := model.User{}

	// bind JSON
	if err := c.ShouldBindJSON(&user); err != nil {
		grenderer.Render(c, gin.H{"message": err.Error()}, http.StatusBadRequest)
		return
	}

	resp, statusCode := handler.UpdateUserProfile(userIDAuth, user)

	if reflect.TypeOf(resp.Message).Kind() == reflect.String {
		grenderer.Render(c, resp, statusCode)
		return
	}

	grenderer.Render(c, resp.Message, statusCode)
}
