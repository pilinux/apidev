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

// GetNotes handles jobs for controller.GetNotes
func GetNotes(userIDAuth uint64) (httpResponse gmodel.HTTPResponse, httpStatusCode int) {
	db := gdatabase.GetDB()
	user := model.User{}
	notes := []model.Note{}

	// does the user have an existing profile
	if err := db.Where("id_auth = ?", userIDAuth).First(&user).Error; err != nil {
		httpResponse.Message = "no user profile found"
		httpStatusCode = http.StatusForbidden
		return
	}

	// find all notes written by this user
	if err := db.Where("id_user = ?", user.UserID).Find(&notes).Error; err != nil {
		log.WithError(err).Error("error code: 1201")
		httpResponse.Message = "internal server error"
		httpStatusCode = http.StatusInternalServerError
		return
	}

	if len(notes) == 0 {
		httpResponse.Message = "no note found"
		httpStatusCode = http.StatusNotFound
		return
	}

	httpResponse.Message = notes
	httpStatusCode = http.StatusOK
	return
}

// GetNote handles jobs for controller.GetNote
func GetNote(userIDAuth uint64, id string) (httpResponse gmodel.HTTPResponse, httpStatusCode int) {
	db := gdatabase.GetDB()
	user := model.User{}
	note := model.Note{}

	// does the user have an existing profile
	if err := db.Where("id_auth = ?", userIDAuth).First(&user).Error; err != nil {
		httpResponse.Message = "no user profile found"
		httpStatusCode = http.StatusForbidden
		return
	}

	// show the note if it is written by the user
	if err := db.Where("note_id = ?", id).Where("id_user = ?", user.UserID).First(&note).Error; err != nil {
		httpResponse.Message = "note not found"
		httpStatusCode = http.StatusNotFound
		return
	}

	httpResponse.Message = note
	httpStatusCode = http.StatusOK
	return
}

// CreateNote handles jobs for controller.CreateNote
func CreateNote(userIDAuth uint64, note model.Note) (httpResponse gmodel.HTTPResponse, httpStatusCode int) {
	db := gdatabase.GetDB()
	user := model.User{}
	noteFinal := model.Note{}

	// does the user have an existing profile
	if err := db.Where("id_auth = ?", userIDAuth).First(&user).Error; err != nil {
		httpResponse.Message = "no user profile found"
		httpStatusCode = http.StatusForbidden
		return
	}

	// remove all leading and trailing white spaces
	note.Title = strings.TrimSpace(note.Title)
	if note.Title == "" {
		httpResponse.Message = "title is required"
		httpStatusCode = http.StatusBadRequest
		return
	}

	// security: user must not be able to manipulate all fields
	noteFinal.Title = note.Title
	noteFinal.Body = note.Body
	noteFinal.IDUser = user.UserID

	// save in DB
	tx := db.Begin()
	if err := tx.Create(&noteFinal).Error; err != nil {
		tx.Rollback()
		log.WithError(err).Error("error code: 1211")
		httpResponse.Message = "internal server error"
		httpStatusCode = http.StatusInternalServerError
		return
	}
	tx.Commit()

	httpResponse.Message = noteFinal
	httpStatusCode = http.StatusCreated
	return
}

// UpdateNote handles jobs for controller.UpdateNote
func UpdateNote(userIDAuth uint64, id string, note model.Note) (httpResponse gmodel.HTTPResponse, httpStatusCode int) {
	db := gdatabase.GetDB()
	user := model.User{}
	noteFinal := model.Note{}

	// does the user have an existing profile
	if err := db.Where("id_auth = ?", userIDAuth).First(&user).Error; err != nil {
		httpResponse.Message = "no user profile found"
		httpStatusCode = http.StatusForbidden
		return
	}

	// does the note exist + does the user have right to modify this note
	if err := db.Where("note_id = ?", id).Where("id_user = ?", user.UserID).First(&noteFinal).Error; err != nil {
		httpResponse.Message = "user may not have access to perform this task"
		httpStatusCode = http.StatusForbidden
		return
	}

	// remove all leading and trailing white spaces
	note.Title = strings.TrimSpace(note.Title)
	if note.Title == "" {
		httpResponse.Message = "title is required"
		httpStatusCode = http.StatusBadRequest
		return
	}

	// if no new info is received, abort
	if note.Title == noteFinal.Title && note.Body == noteFinal.Body {
		httpResponse.Message = "no new info to update"
		httpStatusCode = http.StatusBadRequest
		return
	}

	// security: user must not be able to manipulate all fields
	noteFinal.UpdatedAt = time.Now()
	noteFinal.Title = note.Title
	noteFinal.Body = note.Body

	// update in DB
	tx := db.Begin()
	if err := tx.Save(&noteFinal).Error; err != nil {
		tx.Rollback()
		log.WithError(err).Error("error code: 1221")
		httpResponse.Message = "internal server error"
		httpStatusCode = http.StatusInternalServerError
		return
	}
	tx.Commit()

	httpResponse.Message = noteFinal
	httpStatusCode = http.StatusOK
	return
}

// DeleteNote handles jobs for controller.DeleteNote
func DeleteNote(userIDAuth uint64, id string) (httpResponse gmodel.HTTPResponse, httpStatusCode int) {
	db := gdatabase.GetDB()
	user := model.User{}
	note := model.Note{}

	// does the user have an existing profile
	if err := db.Where("id_auth = ?", userIDAuth).First(&user).Error; err != nil {
		httpResponse.Message = "no user profile found"
		httpStatusCode = http.StatusForbidden
		return
	}

	// does the note exist + does the user have right to delete this note
	if err := db.Where("note_id = ?", id).Where("id_user = ?", user.UserID).First(&note).Error; err != nil {
		httpResponse.Message = "user may not have access to perform this task"
		httpStatusCode = http.StatusForbidden
		return
	}

	// delete from DB
	tx := db.Begin()
	if err := tx.Delete(&note).Error; err != nil {
		tx.Rollback()
		log.WithError(err).Error("error code: 1231")
		httpResponse.Message = "internal server error"
		httpStatusCode = http.StatusInternalServerError
		return
	}
	tx.Commit()

	httpResponse.Message = "note ID# " + id + " deleted!"
	httpStatusCode = http.StatusOK
	return
}
