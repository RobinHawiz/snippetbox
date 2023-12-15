package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/robinhawiz/snippetbox/internal/models"
)

func (a *application) home(w http.ResponseWriter, r *http.Request){
	if r.URL.Path != "/"{
		a.notFound(w)
		return
	}else{
		snippets, err := a.snippets.Latest()
		if err != nil {
			a.serverError(w,r,err)
		}
		a.render(w,r,http.StatusOK,"home.tmpl",templateData{
			Snippets: snippets,
		})
	}
}

func (a *application) snippetViewHandler(w http.ResponseWriter, r *http.Request){
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 1 {
		a.notFound(w)
		return
	}
	snippet, err := a.snippets.Get(id)
	if err != nil{
		if errors.Is(err, models.ErrNoRecord){
			a.notFound(w)
		}else{
			a.serverError(w,r,err)
		}
			return
		}
	a.render(w,r,http.StatusOK,"view.tmpl",templateData{
		Snippet: snippet,
	})
}

func (a *application) snippetCreateHandler(w http.ResponseWriter, r *http.Request){
	if(r.Method != http.MethodPost){
		w.Header().Set("Allow", http.MethodPost)
		a.clientError(w,http.StatusMethodNotAllowed)
		return
	}
	id, err := a.snippets.Insert("Snails: 101", "Snails move surprisingly slowly. But while snails may not be the fastest creatures, their steady pace shows that sometimes perseverance is more important than speed.", 1)
	if err != nil {
		a.serverError(w,r,err)
		return
	}

	http.Redirect(w,r,fmt.Sprintf("/snippet/view?id=%d", id), http.StatusSeeOther)
}