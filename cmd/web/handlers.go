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
	if id, err := strconv.Atoi(r.URL.Query().Get("id")); err != nil || id < 1 {
		a.notFound(w)
		return
	}else{
		if snippet, err := a.snippets.Get(id); err != nil{
			if errors.Is(err, models.ErrNoRecord){
				a.notFound(w)
			}else{
				a.serverError(w,r,err)
			}
			return
		}else{
		// fmt.Fprintf(w, "Your snippet:\nID:%d\nTitle:\n%s\nContent:\n%s\nCreated:\n%s\nExpires:\n%s\n", snippet.ID, snippet.Title, snippet.Content, snippet.Created, snippet.Expires)
		fmt.Fprintf(w, "%+v", snippet)
		}
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