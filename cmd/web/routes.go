package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
)

func (app *application) routes() http.Handler {

	// mux := http.NewServeMux()
	router := httprouter.New()

	router.NotFound = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.notFound(w)
	})

	fs := http.FileServer(http.Dir("./ui/static/"))
	//mux.Handle("/static/", http.StripPrefix("/static", fs))
	router.Handler(http.MethodGet, "/static/*filepath", http.StripPrefix("/static", fs))

	router.HandlerFunc(http.MethodGet, "/", app.home)
	router.HandlerFunc(http.MethodGet, "/book/view/:id", app.bookView)
	router.HandlerFunc(http.MethodGet, "/book/create", app.bookCreateForm)
	router.HandlerFunc(http.MethodPost, "/book/create", app.bookCreate)

	//return app.recoverPanic(app.logRequest(secureHeaders(mux)))

	stdMiddleWare := alice.New(app.recoverPanic, app.logRequest, secureHeaders)

	return stdMiddleWare.Then(router)
}
