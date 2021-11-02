package main

import (
	"errors"
	"fmt"
	"github.com/White-AK111/snippetbox/pkg/models"
	"github.com/jackc/pgtype"
	"net/http"
	"strconv"
	"time"
)

// home page handler
func (a *app) home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		a.notFound(w)
		return
	}

	// get latest snippet
	s, err := a.snippets.LatestSnippets(1, 10)
	if err != nil {
		a.serverError(w, err)
		return
	}

	// fill template
	a.render(w, r, "home.page.tmpl", &templateData{
		Snippets: s,
	})
}

// showSnippet view snippet content
func (a *app) showSnippet(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 1 {
		a.notFound(w)
		return
	}

	s, err := a.snippets.GetSnippet(id, 1)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			a.notFound(w)
		} else {
			a.serverError(w, err)
		}
		return
	}

	// fill template
	a.render(w, r, "show.page.tmpl", &templateData{
		Snippet: s,
	})
}

// createSnippet create new snippet
func (a *app) createSnippet(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		a.clientError(w, http.StatusMethodNotAllowed)
		return
	}

	snippet := models.Snippet{
		UserId:  1,
		Title:   "AnySnippet",
		Content: "Test test test ...",
		Created: pgtype.Timestamptz{Time: time.Now()},
		Expires: pgtype.Timestamptz{Time: time.Now().Local().Add(time.Hour * time.Duration(240))},
		Changed: pgtype.Timestamptz{Time: time.Now()},
		Deleted: false,
	}

	id, err := a.snippets.InsertSnippet(&snippet)
	if err != nil {
		a.serverError(w, err)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/snippet?id=%d", id), http.StatusSeeOther)

	w.Write([]byte("Create snippet form"))
}
