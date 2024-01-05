package main

import (
	"net/http"

	"github.com/justinas/alice"
)

func (app *application) routes() http.Handler {

	mux := http.NewServeMux()

	fs := http.FileServer(http.Dir("./ui/static/"))
	mux.Handle("/static/", http.StripPrefix("/static", fs))

	mux.HandleFunc("/", app.home)
	mux.HandleFunc("/book/view", app.bookView)
	mux.HandleFunc("/book/create", app.bookCreate)

	//return app.recoverPanic(app.logRequest(secureHeaders(mux)))

	stdMiddleWare := alice.New(app.recoverPanic, app.logRequest, secureHeaders)

	return stdMiddleWare.Then(mux)
}
