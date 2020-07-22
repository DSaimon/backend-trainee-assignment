package app

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/rdnply/backend-trainee-assignment/internal/user"
)

func (app *App) addUser(w http.ResponseWriter, r *http.Request) {
	var u user.User

	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		app.BadRequest(w, err, "incorrect json")
		return
	}

	fromDB, err := app.UserStorage.Find(u.Username)
	if err != nil {
		app.ServerError(w, err, "")
		return
	}
	if fromDB.ID != 0 {
		app.BadRequest(w, err, fmt.Sprintf("user %s is already exist", u.Username))
		return
	}

	if err := app.UserStorage.Add(&u); err != nil {
		app.ServerError(w, err, "")
		return
	}

	respondJSON(w, http.StatusCreated, map[string]int{"id": u.ID})
}

func (app *App) addChat(w http.ResponseWriter, r *http.Request) {
	type chatInfo struct {
		Name  string `json:"name"`
		Users []int  `json:"users"`
	}

	var c chatInfo
	if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
		app.BadRequest(w, err, "incorrect json")
		return
	}

	exists, id, err := app.UserStorage.AllExists(c.Users)
	if err != nil {
		app.ServerError(w, err, "")
		return
	}
	if !exists {
		app.NotFound(w, err, fmt.Sprintf("not found user with id: %v", id))
		return
	}

	fromDB, err := app.ChatStorage.Find(c.Name)
	if err != nil {
		app.ServerError(w, err, "")
		return
	}
	if fromDB.ID != 0 {
		app.BadRequest(w, err, fmt.Sprintf("chat %s is already exist", c.Name))
		return
	}

	newID, err := app.ChatStorage.Add(c.Name, c.Users)
	if err != nil {
		app.ServerError(w, err, "")
		return
	}

	respondJSON(w, http.StatusCreated, map[string]int{"id": newID})
}

func respondJSON(w http.ResponseWriter, successCode int, payload interface{}) {
	response, err := json.Marshal(payload)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(successCode)
	w.Write(response)
}
