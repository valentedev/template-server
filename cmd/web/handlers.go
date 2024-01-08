package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/julienschmidt/httprouter"
	"github.com/valentedev/template-server/internal/models"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {

	// Httprouter matches the "/" path exactly so we dont need below checking
	// Return 404 in case "/" does not refind a handler
	// if r.URL.Path != "/" {
	// 	//http.NotFound(w, r)
	// 	app.notFound(w)
	// 	return
	// }

	// panic("oops! something went wrong")

	books, err := app.books.Latest()
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	data := app.newTemplateData(r)
	data.Books = books

	// app.render(w, r, http.StatusOK, "home.html", templateData{
	// 	Books: books,
	// })
	app.render(w, r, http.StatusOK, "home.html", data)

	// // String slice with the paths to our templates
	// files := []string{
	// 	"./ui/html/base.html",
	// 	"./ui/html/components/nav.html",
	// 	"./ui/html/pages/home.html",
	// }

	// ts, err := template.ParseFiles(files...)
	// if err != nil {
	// 	log.Println(err.Error())
	// 	//http.Error(w, "Interal Server Error", http.StatusInternalServerError)
	// 	app.serverError(w, r, err)
	// }

	// data := templateData{
	// 	Books: books,
	// }

	// // Instead of ts.Execute we use ts.ExecuteTemplate to render the "base" Template
	// err = ts.ExecuteTemplate(w, "base", data)
	// if err != nil {
	// 	//log.Println(err.Error())
	// 	//http.Error(w, "Interal Server Error", http.StatusInternalServerError)
	// 	app.serverError(w, r, err)
	// }

}

func (app *application) bookView(w http.ResponseWriter, r *http.Request) {

	params := httprouter.ParamsFromContext(r.Context())

	// Will convert the parameter "id" which is a String in to Int
	//id, err := strconv.Atoi(r.URL.Query().Get("id"))
	id, err := strconv.Atoi(params.ByName("id"))
	if err != nil || id < 1 {
		//http.NotFound(w, r)
		app.notFound(w)
		return
	}

	book, err := app.books.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, r, err)
		}
		return
	}

	data := app.newTemplateData(r)
	data.Book = book

	// app.render(w, r, http.StatusOK, "view.html", templateData{
	// 	Book: book,
	// })
	app.render(w, r, http.StatusOK, "view.html", data)

	//fmt.Fprintf(w, "%+v", book)

	// w.Write([]byte("Display a specific book..."))
}

func (app *application) bookCreateForm(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	data.Form = bookCreateForm{}
	app.render(w, r, http.StatusOK, "create.html", data)
}

type bookCreateForm struct {
	Title       string
	Author      string
	FieldErrors map[string]string
}

func (app *application) bookCreate(w http.ResponseWriter, r *http.Request) {
	// if r.Method != "POST" {
	// 	w.Header().Set("Allow", "POST")
	// 	// w.WriteHeader(405)
	// 	// w.Write([]byte("Method not allowed"))
	// 	http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	// 	return
	// }
	// if r.Method != http.MethodPost {
	// 	w.Header().Set("Allow", "POST")
	// 	//http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	// 	app.clientError(w, http.StatusMethodNotAllowed)
	// 	return
	// }

	// title := "O snail"
	// content := "O snail\nClimb Mount Fuji,\nBut slowly, slowly!\n\n– Kobayashi Issa"
	//expires := "4"

	r.Body = http.MaxBytesReader(w, r.Body, 4096)

	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	// title := r.PostForm.Get("title")
	// author := r.PostForm.Get("author")

	// fieldErrors := make(map[string]string)

	form := bookCreateForm{
		Title:       r.PostForm.Get("title"),
		Author:      r.PostForm.Get("author"),
		FieldErrors: map[string]string{},
	}

	if strings.TrimSpace(form.Title) == "" {
		form.FieldErrors["title"] = "This field cannot be blank"
	} else if utf8.RuneCountInString(form.Title) > 100 {
		form.FieldErrors["title"] = "This field cannot be more than 100 characters long"
	}

	if strings.TrimSpace(form.Author) == "" {
		form.FieldErrors["author"] = "This field cannot be blank"
	} else if utf8.RuneCountInString(form.Title) > 100 {
		form.FieldErrors["author"] = "This field cannot be more than 100 characters long"
	}

	if len(form.FieldErrors) > 0 {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, r, http.StatusUnprocessableEntity, "create.html", data)
		return
	}

	id, err := app.books.Insert(form.Title, form.Author)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/book/view/%d", id), http.StatusSeeOther)
	//w.Write([]byte("Create a new book..."))
}
