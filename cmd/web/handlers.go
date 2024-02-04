package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/coderjojo/goapi/pkg/models"
)

func (app *application) home(rw http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		app.notFound(rw)
		return
	}

	s, err := app.snippets.Latest()
	if err != nil {
		app.serverError(rw, err)
		return
	}

	app.render(rw, r, "home.page.tmpl", &templateData{
		Snippets: s,
	})

}

func (app *application) showSnippet(rw http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 1 {
		app.notFound(rw)
		return
	}
	app.infoLog.Printf("Requested details for ID %d", id)

	s, err := app.snippets.Get(id)
	if err == models.ErrNoRecord {
		app.notFound(rw)
		return
	} else if err != nil {
		app.serverError(rw, err)
		return
	}

	app.render(rw, r, "show.page.tmpl", &templateData{
		Snippet: s,
	})

}

func (app *application) createSnippet(rw http.ResponseWriter, r *http.Request) {

	if r.Method != "POST" {
		rw.Header().Set("Allow", "POST")
		app.clientError(rw, http.StatusMethodNotAllowed)
		return
	}

	// dummy data
	title := "O snail"
	content := "O snail\nClimb Mount Fuji,\nBut slowly, slowly!\n\n– Kobayash"
	expires := "7"

	id, err := app.snippets.Insert(title, content, expires)
	if err != nil {
		app.serverError(rw, err)
		return
	}

	http.Redirect(rw, r, fmt.Sprintf("/snippet?id=%d", id), http.StatusSeeOther)

}
