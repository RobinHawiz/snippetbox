package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
	"github.com/robinhawiz/snippetbox/internal/models"
)

func (a *application) home(w http.ResponseWriter, r *http.Request){
	snippets, err := a.snippets.Latest()
	if err != nil {
		a.serverError(w,r,err)
	}
	data := a.newTemplateData(r)
	data.Snippets = snippets
	a.render(w,r,http.StatusOK,"home.tmpl",data)
}

func (a *application) snippetViewHandler(w http.ResponseWriter, r *http.Request){
	//httprouter when parsing a request puts any named parameters in the request context.
	//in this instance params becomes a slice with a key (id) and a corresponding value.
	params := httprouter.ParamsFromContext(r.Context())

	//We use the ByName() method to get the value of the "id" named parameter from the slice and valitdate it.
	id, err := strconv.Atoi(params.ByName("id"))
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

	data := a.newTemplateData(r)
	data.Snippet = snippet

	a.render(w,r,http.StatusOK,"view.tmpl",data)
}

func (a *application) snippetCreateHandler(w http.ResponseWriter, r *http.Request){
	w.Write([]byte("Display the form for creating a new snippet..."))
}

func (a *application) snippetCreatePostHandler(w http.ResponseWriter, r *http.Request){

	id, err := a.snippets.Insert("Snails: 101", "Snails move surprisingly slowly. But while snails may not be the fastest creatures, their steady pace shows that sometimes perseverance is more important than speed.", 1)
	if err != nil {
		a.serverError(w,r,err)
		return
	}

	http.Redirect(w,r,fmt.Sprintf("/snippet/view/%d", id), http.StatusSeeOther)
}