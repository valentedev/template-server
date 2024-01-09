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

	// Book Routes
	router.Handler(http.MethodGet, "/", dynamic.ThenFunc(app.home))
	router.Handler(http.MethodGet, "/book/view/:id", dynamic.ThenFunc(app.bookView))
	router.Handler(http.MethodGet, "/book/create", dynamic.ThenFunc(app.bookCreateForm))
	router.Handler(http.MethodPost, "/book/create", dynamic.ThenFunc(app.bookCreate))

	// User Routes
	router.Handler(http.MethodGet, "/user/signup", dynamic.ThenFunc(app.userSignup))
	router.Handler(http.MethodPost, "/user/signup", dynamic.ThenFunc(app.userSignupPost))
	router.Handler(http.MethodGet, "/user/login", dynamic.ThenFunc(app.userLogin))
	router.Handler(http.MethodPost, "/user/login", dynamic.ThenFunc(app.userLoginPost))
	router.Handler(http.MethodPost, "/user/logout", dynamic.ThenFunc(app.userLogoutPost))

	//return app.recoverPanic(app.logRequest(secureHeaders(mux)))

	stdMiddleWare := alice.New(app.recoverPanic, app.logRequest, secureHeaders)

	return stdMiddleWare.Then(router)
}
