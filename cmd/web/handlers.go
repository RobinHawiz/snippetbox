package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
	"github.com/robinhawiz/snippetbox/internal/models"
	"github.com/robinhawiz/snippetbox/internal/validator"
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
	//in this instance params becomes a slice with a key (id) and a corresponding value (taken from the URL).
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
//At the end of the field I've included struct tags which tells the decoder how to map HTML form values into the different struct fields.
type snippetCreateForm struct {
	Title 				string `form:"title"`
	Content 			string `form:"content"`
	Expires 			int	`form:"expires"`
	validator.Validator `form:"-"`
}


func (a *application) snippetCreatePostHandler(w http.ResponseWriter, r *http.Request){

	var form snippetCreateForm

	err := a.decodePostForm(r, &form)
	if err != nil{
		a.clientError(w, http.StatusBadRequest)
		return
	}

	//CheckField() will add the provided key and error message to the FieldErrors map if the check does not evaluate to true.
	form.CheckField(validator.NotBlank(form.Title), "title", "This field cannot be blank")
	form.CheckField(validator.MaxChars(form.Title, 100), "title", "This field cannot be more than 100 characters long")
	form.CheckField(validator.NotBlank(form.Content), "content", "This field cannot be blank")
	form.CheckField(validator.PermittedValue(form.Expires, 1, 7, 365), "Expires", "This field must equal 1, 7 or 365")

	//Check if any of the above checks failed. If they did, re-display the create.tmpl template, passing in the snippetCreateForm instance as dynamic data in the Form field.
	if !form.Valid() {
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

	//We add a string key ("flash") and a corresponding value to the session data.
	a.sessionManager.Put(r.Context(), "flash", "Snippet successfully created!")

	http.Redirect(w,r,fmt.Sprintf("/snippet/view/%d", id), http.StatusSeeOther)
}

type userSignupForm struct {
	Name 				string `form:"name"`
	Email 				string `form:"email"`
	Password			string `form:"password"`
	validator.Validator `form:"-"`
}

func (a *application) userSignup(w http.ResponseWriter, r *http.Request) {
	data := a.newTemplateData(r)
	//Set any initial display values for the form.
	data.Form = userSignupForm{}
	a.render(w,r,http.StatusOK,"signup.tmpl",data)
}

func (a *application) userSignupPost(w http.ResponseWriter, r *http.Request) {
	var form userSignupForm

	//Parse the form data into the userSignupForm struct.
	err := a.decodePostForm(r, &form)
	if err != nil {
		a.clientError(w, http.StatusBadRequest)
		return
	}

	//CheckField() will add the provided key and error message to the FieldErrors map if the check does not evaluate to true.
	form.CheckField(validator.NotBlank(form.Name), "name", "This field cannot be blank")
	form.CheckField(validator.NotBlank(form.Email), "email", "This field cannot be blank")
	form.CheckField(validator.Matches(form.Email, validator.EmailRX), "email", "This field must be a valid email address")
	form.CheckField(validator.NotBlank(form.Password), "password", "This field cannot be blank")
	form.CheckField(validator.MinChars(form.Password, 8), "password", "This field must be at least 8 characters long")

	//Check if any of the above checks failed. If they did, re-display the signup.tmpl template, passing in the userSignupForm instance as dynamic data in the Form field.
	if !form.Valid() {
		data := a.newTemplateData(r)
		data.Form = form
		a.render(w,r,http.StatusUnprocessableEntity, "signup.tmpl", data)
		return
	}

	//Tries to create a new user record in the database. If the email already exists then add an error message to the form and re-display it.
	err = a.users.Insert(form.Name, form.Email, form.Password)
	if err != nil {
		if errors.Is(err, models.ErrDuplicateEmail) {
			form.AddFieldError("email", "Email is already in use")

			data := a.newTemplateData(r)
			data.Form = form
			a.render(w,r,http.StatusUnprocessableEntity, "signup.tmpl", data)
		}else {
			a.serverError(w, r, err)
		}

		return
	}

	//If creating a user was successfull, add a confirmation flash message to the session confirming that their signup worked.
	a.sessionManager.Put(r.Context(), "flash", "Your signup was successfull. Please log in.")

	//And redirect the user to the login page.
	http.Redirect(w, r, "/user/login", http.StatusSeeOther)
}

type userLoginForm struct {
	Email 				string `form:"email"`
	Password			string `form:"password"`
	validator.Validator `form:"-"`
}

func (a *application) userLogin(w http.ResponseWriter, r *http.Request) {
	data := a.newTemplateData(r)
	data.Form = userLoginForm{}
	a.render(w,r,http.StatusUnprocessableEntity, "login.tmpl", data)
}

func (a *application) userLoginPost(w http.ResponseWriter, r *http.Request) {
	//Decode the form data into the userLoginForm struct
	var form userLoginForm

	err := a.decodePostForm(r, &form)
	if err != nil {
		a.clientError(w, http.StatusBadRequest)
		return
	}

	//CheckField() will add the provided key and error message to the FieldErrors map if the check does not evaluate to true.
	form.CheckField(validator.NotBlank(form.Email), "email", "This field cannot be blank")
	form.CheckField(validator.Matches(form.Email, validator.EmailRX), "email", "This field must be a valid email address")
	form.CheckField(validator.NotBlank(form.Password), "password", "This field cannot be blank")

	//Check if any of the above checks failed. If they did, re-display the login.tmpl template, passing in the userLoginForm instance as dynamic data in the Form field.
	if !form.Valid() {
		data := a.newTemplateData(r)
		data.Form = form
		a.render(w,r,http.StatusUnprocessableEntity, "login.tmpl", data)
		return
	}

	id, err := a.users.Authenticate(form.Email, form.Password)
	if err != nil {
		if errors.Is(err, models.ErrInvalidCredentials) {
			form.AddNonFieldError("Email or password is incorrect")

			data := a.newTemplateData(r)
			data.Form = form
			a.render(w,r,http.StatusUnprocessableEntity, "login.tmpl", data)
		} else {
			a.serverError(w, r, err)
		}
		return
	}

	//Changes the current session ID.
	//It's considered good practice to generate a new session ID when the authentication state or privilige levels change for the user (e.g. login and logout operations).
	err = a.sessionManager.RenewToken(r.Context())
	if err != nil {
		a.serverError(w, r, err)
		return
	}

	//Add the ID of the current user to the session, so that they are now "logged in".
	a.sessionManager.Put(r.Context(), "authenticatedUserID", id)

	//Redirect the user to the create snippet page.
	http.Redirect(w, r, "/snippet/create", http.StatusSeeOther)
}

func (a *application) userLogoutPost(w http.ResponseWriter, r *http.Request) {
	//Changes the current session ID.
	err := a.sessionManager.RenewToken(r.Context())
	if err != nil {
		a.serverError(w, r, err)
		return
	}

	//Remove the authenticatedUserID from the session data so that user is "logged out".
	a.sessionManager.Remove(r.Context(), "authenticatedUserID")

	//Add a flash message to the session to confirm to the user that they've been logged out.
	a.sessionManager.Put(r.Context(), "flash", "You've been logged out successfully!")

	//Redirect the user to the application home page.
	http.Redirect(w, r, "/", http.StatusSeeOther)
}