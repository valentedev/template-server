package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
)

func (app *application) routes() http.Handler {

	router := httprouter.New()

	router.NotFound = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.notFound(w)
	})

	fs := http.FileServer(http.Dir("./ui/static/"))
	router.Handler(http.MethodGet, "/static/*filepath", http.StripPrefix("/static", fs))

	dynamic := alice.New(app.sessionManager.LoadAndSave)

	router.Handler(http.MethodGet, "/", dynamic.ThenFunc(app.home))
	router.Handler(http.MethodGet, "/book/view/:id", dynamic.ThenFunc(app.bookView))
	router.Handler(http.MethodGet, "/book/create", dynamic.ThenFunc(app.bookCreateForm))
	router.Handler(http.MethodPost, "/book/create", dynamic.ThenFunc(app.bookCreate))

	//return app.recoverPanic(app.logRequest(secureHeaders(mux)))

	stdMiddleWare := alice.New(app.recoverPanic, app.logRequest, secureHeaders)

	return stdMiddleWare.Then(router)
}
