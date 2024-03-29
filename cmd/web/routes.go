package main

import (
	"net/http"

	"github.com/valentedev/template-server/ui"

	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
)

func (app *application) routes() http.Handler {

	router := httprouter.New()

	router.NotFound = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.notFound(w)
	})

	// Without embed
	//fs := http.FileServer(http.Dir("./ui/static/"))
	// router.Handler(http.MethodGet, "/static/*filepath", http.StripPrefix("/static", fs))

	fs := http.FileServer(http.FS(ui.Files))
	router.Handler(http.MethodGet, "/static/*filepath", fs)

	router.HandlerFunc(http.MethodGet, "/ping", ping)

	// Unprotected Routes
	dynamic := alice.New(app.sessionManager.LoadAndSave, noSurf, app.authenticate)
	router.Handler(http.MethodGet, "/", dynamic.ThenFunc(app.home))
	router.Handler(http.MethodGet, "/book/view/:id", dynamic.ThenFunc(app.bookView))
	router.Handler(http.MethodGet, "/user/signup", dynamic.ThenFunc(app.userSignup))
	router.Handler(http.MethodPost, "/user/signup", dynamic.ThenFunc(app.userSignupPost))
	router.Handler(http.MethodGet, "/user/login", dynamic.ThenFunc(app.userLogin))
	router.Handler(http.MethodPost, "/user/login", dynamic.ThenFunc(app.userLoginPost))

	// Protected Routes
	protected := dynamic.Append(app.requireAuthentication)
	router.Handler(http.MethodGet, "/book/create", protected.ThenFunc(app.bookCreateForm))
	router.Handler(http.MethodPost, "/book/create", protected.ThenFunc(app.bookCreate))
	router.Handler(http.MethodPost, "/user/logout", protected.ThenFunc(app.userLogoutPost))

	stdMiddleWare := alice.New(app.recoverPanic, app.logRequest, secureHeaders)

	return stdMiddleWare.Then(router)
}
