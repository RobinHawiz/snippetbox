package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"text/template"

	"github.com/robinhawiz/snippetbox/internal/models"
)

func (a *application) home (w http.ResponseWriter, r *http.Request){
	if r.URL.Path != "/"{
		a.notFound(w)
		return
	}
	files := []string{
		"./ui/html/base.tmpl",
		"./ui/html/pages/home.tmpl",
		"./ui/html/partials/nav.tmpl",
	}
	if tmpl, err := template.ParseFiles(files...); err != nil {
		a.serverError(w,r,err)
		return
	}else{
		if err = tmpl.ExecuteTemplate(w, "base", nil); err != nil{
			a.serverError(w,r,err)
			return
		}
	}
	if snippets, err := a.snippets.Latest(); err != nil{
		a.serverError(w,r,err)
		return
	}else{
		for _, snippet := range snippets {
			fmt.Fprintf(w, "%+v", snippet)
		}
	}
	
}

func (a *application) snippetViewHandler (w http.ResponseWriter, r *http.Request){
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
	files := []string{
		"./ui/html/base.tmpl",
		"./ui/html/partials/nav.tmpl",
		"./ui/html/pages/view.tmpl",
		
	}
	ts, err := template.ParseFiles(files...)
		if err != nil{
			a.serverError(w,r,err)
		}
	if err = ts.ExecuteTemplate(w, "base", snippet); err != nil{
		a.serverError(w,r,err)
		return
	}
}

func (a *application) snippetCreateHandler (w http.ResponseWriter, r *http.Request){
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