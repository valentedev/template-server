package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/valentedev/template-server/internal/models"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {

	// Return 404 in case "/" does not refind a handler
	if r.URL.Path != "/" {
		//http.NotFound(w, r)
		app.notFound(w)
		return
	}

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

	// Will convert the parameter "id" which is a String in to Int
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
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

func (app *application) bookCreate(w http.ResponseWriter, r *http.Request) {
	// if r.Method != "POST" {
	// 	w.Header().Set("Allow", "POST")
	// 	// w.WriteHeader(405)
	// 	// w.Write([]byte("Method not allowed"))
	// 	http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	// 	return
	// }
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", "POST")
		//http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		app.clientError(w, http.StatusMethodNotAllowed)
		return
	}

	title := "O snail"
	content := "O snail\nClimb Mount Fuji,\nBut slowly, slowly!\n\n– Kobayashi Issa"
	//expires := "4"

	id, err := app.books.Insert(title, content)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("book/view?id=%d", id), http.StatusSeeOther)
	//w.Write([]byte("Create a new book..."))
}
