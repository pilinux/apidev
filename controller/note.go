package controller

import (
	"net/http"
	"reflect"
	"strings"

	"github.com/gin-gonic/gin"
	grenderer "github.com/pilinux/gorest/lib/renderer"

	"apidev/database/model"
	"apidev/handler"
)

// GetNotes - GET /notes
// no note is in public mode
// only an authorized user can access his notes
func GetNotes(c *gin.Context) {
	userIDAuth := c.GetUint64("authID")

	resp, statusCode := handler.GetNotes(userIDAuth)

	if reflect.TypeOf(resp.Message).Kind() == reflect.String {
		grenderer.Render(c, resp, statusCode)
		return
	}

	grenderer.Render(c, resp, statusCode)
}

// GetNote - GET /notes/:id
// fetch a note by its ID
// no note is in public mode
// only an authorized user can access his notes
func GetNote(c *gin.Context) {
	userIDAuth := c.GetUint64("authID")
	id := strings.TrimSpace(c.Params.ByName("id"))

	resp, statusCode := handler.GetNote(userIDAuth, id)

	if reflect.TypeOf(resp.Message).Kind() == reflect.String {
		grenderer.Render(c, resp, statusCode)
		return
	}

	grenderer.Render(c, resp.Message, statusCode)
}

// CreateNote - POST /notes
// only an authorized user can create a new note
// =================================
//
//	{
//	   "Title": "title_of_the_note",
//	   "Body": "body_of_the_note"
//	}
//
// =================================
func CreateNote(c *gin.Context) {
	userIDAuth := c.GetUint64("authID")
	note := model.Note{}

	// bind JSON
	if err := c.ShouldBindJSON(&note); err != nil {
		grenderer.Render(c, gin.H{"message": err.Error()}, http.StatusBadRequest)
		return
	}

	resp, statusCode := handler.CreateNote(userIDAuth, note)

	if reflect.TypeOf(resp.Message).Kind() == reflect.String {
		grenderer.Render(c, resp, statusCode)
		return
	}

	grenderer.Render(c, resp.Message, statusCode)
}

// UpdateNote - PUT /notes/:id
// only an authorized user can update his existing notes
// =====================================
//
//	{
//	   "Title": "new_title_of_the_note",
//	   "Body": "new_body_of_the_note"
//	}
//
// =====================================
func UpdateNote(c *gin.Context) {
	userIDAuth := c.GetUint64("authID")
	id := strings.TrimSpace(c.Params.ByName("id"))
	note := model.Note{}

	// bind JSON
	if err := c.ShouldBindJSON(&note); err != nil {
		grenderer.Render(c, gin.H{"message": err.Error()}, http.StatusBadRequest)
		return
	}

	resp, statusCode := handler.UpdateNote(userIDAuth, id, note)

	if reflect.TypeOf(resp.Message).Kind() == reflect.String {
		grenderer.Render(c, resp, statusCode)
		return
	}

	grenderer.Render(c, resp.Message, statusCode)
}

// DeleteNote - DELETE /notes/:id
// only an authorized user can delete his existing notes
// this example performs soft delete operation
func DeleteNote(c *gin.Context) {
	userIDAuth := c.GetUint64("authID")
	id := strings.TrimSpace(c.Params.ByName("id"))

	resp, statusCode := handler.DeleteNote(userIDAuth, id)

	if reflect.TypeOf(resp.Message).Kind() == reflect.String {
		grenderer.Render(c, resp, statusCode)
		return
	}

	grenderer.Render(c, resp, statusCode)
}
