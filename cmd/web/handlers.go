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
	//Set any initial display values for the form.
	data.Form = snippetCreateForm{
		Expires: 365,
	}
	a.render(w,r,http.StatusOK,"create.tmpl",data)
}

//Represents the form data and validation erros for the form fields.
//We export the struct fields in order to be read by the html/template package when rendering the template.
type snippetCreateForm struct {
	Title 		string
	Content 	string
	Expires 	int
	FieldErrors map[string]string
}

func (a *application) snippetCreatePostHandler(w http.ResponseWriter, r *http.Request){

	//Add any data in POST request bodies to the r.PostForm map.
	err := r.ParseForm()
	if err != nil{
		a.clientError(w, http.StatusBadRequest)
		return
	}

	//We are expecting our expires value to be a number, and want to represent it in our Go code as an integer.
	expires, err := strconv.Atoi((r.PostForm.Get("expires")))
	if err != nil{
		a.clientError(w, http.StatusBadRequest)
		return
	}

	//Create an instance of the snippetCreateForm struct containing the values from the form and an empty map for any validation errors.
	form := snippetCreateForm{
		Title: r.PostForm.Get("title"),
		Content: r.PostForm.Get("content"),
		Expires: expires,
		FieldErrors: map[string]string{},
	}

	//Check that the title value is not blank and is not more than 100 characters long.
	if strings.TrimSpace(form.Title) == "" {
		form.FieldErrors["title"] = "This field cannot be blank"
	} else if utf8.RuneCountInString(form.Title) > 100{
		form.FieldErrors["title"] = "This field cannot be more than 100 characters long"
	}

	//Check that the Content value isn't blank.
	if strings.TrimSpace(form.Content) == "" {
		form.FieldErrors["content"] = "This field cannot be blank"
	}

	//Check the expires value matches one of the permitted values (1, 7 or 365).
	if expires != 1 && expires != 7 && expires != 365{
		form.FieldErrors["expires"] = "This field must equal 1, 7 or 365"
	}

	//If there are any errors, then re-display the create.tmpl template, passing in the snippetCreateForm instance as dynamic data in the Form field.
	if len(form.FieldErrors) > 0 {
		data := a.newTemplateData(r)
		data.Form = form
		a.render(w,r,http.StatusUnprocessableEntity, "create.tmpl", data)
		return
	}

	//Insert snippet with the form values into the db.
	id, err := a.snippets.Insert(form.Title, form.Content, form.Expires)
	if err != nil {
		a.serverError(w,r,err)
		return
	}

	http.Redirect(w,r,fmt.Sprintf("/snippet/view/%d", id), http.StatusSeeOther)
}