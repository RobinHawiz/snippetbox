package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"unicode/utf8"

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
	data := a.newTemplateData(r)
	a.render(w,r,http.StatusOK,"create.tmpl",data)
}

func (a *application) snippetCreatePostHandler(w http.ResponseWriter, r *http.Request){

	//Add any data in POST request bodies to the r.PostForm map.
	err := r.ParseForm()
	if err != nil{
		a.clientError(w, http.StatusBadRequest)
		return
	}

	//Get input values from the r.PostForm map.
	title := r.PostForm.Get("title")
	content := r.PostForm.Get("content")
	//We are expecting our expires value to be a number, and want to represent it in our Go code as an integer.
	expires, err := strconv.Atoi((r.PostForm.Get("expires")))
	if err != nil{
		a.clientError(w, http.StatusBadRequest)
		return
	}
	//Initialize a new map to hold any validation errors for the form fields.
	fieldErrors := make(map[string]string)

	//Check that the title value is not blank and is not more than 100 characters long.
	if strings.TrimSpace(title) == "" {
		fieldErrors["title"] = "This field cannot be blank"
	} else if utf8.RuneCountInString(title) > 100{
		fieldErrors["title"] = "This field cannot be more than 100 characters long"
	}

	//Check that the Content value isn't blank.
	if strings.TrimSpace(content) == "" {
		fieldErrors["content"] = "This field cannot be blank"
	}

	//Check the expires value matches one of the permitted values (1, 7 or 365).
	if expires != 1 && expires != 7 && expires != 365{
		fieldErrors["expires"] = "This field must equal 1, 7 or 365"
	}

	//If there are any errors, dump them in a plain text HTTP response and return from the handler.
	if len(fieldErrors) > 0 {
		fmt.Fprint(w, fieldErrors)
		return
	}

	//Insert snippet with the form values into the db.
	id, err := a.snippets.Insert(title, content, expires)
	if err != nil {
		a.serverError(w,r,err)
		return
	}

	http.Redirect(w,r,fmt.Sprintf("/snippet/view/%d", id), http.StatusSeeOther)
}