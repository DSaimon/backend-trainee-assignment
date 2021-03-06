package app

import (
	"net/http"

	"github.com/go-chi/chi"
)

func (app *App) routes() http.Handler {
	r := chi.NewRouter()
	r.Route("/", func(r chi.Router) {
		r.Post("/users/add", app.addUser)
		r.Post("/chats/add", app.addChat)
		r.Post("/messages/add", app.addMessage)
		r.Post("/chats/get", app.getChats)
		r.Post("/messages/get", app.getMessages)
	})

	return r
}
