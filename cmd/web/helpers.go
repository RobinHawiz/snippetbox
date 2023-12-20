package main

import (
	"bytes"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/go-playground/form/v4"
	"github.com/justinas/nosurf"
)

func (a *application) serverError(w http.ResponseWriter, r *http.Request, err error){
	a.logger.Error(err.Error(), slog.String("method", r.Method), slog.String("uri", r.URL.RequestURI()))
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func (a *application) clientError(w http.ResponseWriter, status int){
	http.Error(w, http.StatusText(status), status)
}

func (a *application) notFound(w http.ResponseWriter){
	a.clientError(w, http.StatusNotFound)
}

func (a *application) render(w http.ResponseWriter, r *http.Request, status int, page string, data templateData){
	ts, ok := a.templateCache[page]
	if !ok{
		err := fmt.Errorf("the template %s does not exist", page)
		a.serverError(w,r,err)
		return
	}
	buf := new(bytes.Buffer)
	if err := ts.ExecuteTemplate(buf, "base", data); err != nil{
		a.serverError(w,r,err)
		return
	}
	w.WriteHeader(status)
	buf.WriteTo(w)
}

func (a *application) newTemplateData(r *http.Request) templateData{
	return templateData{
		CurrentYear: time.Now().Year(),
		//If there's a string key "flash" with a value, we'll retrieve it and delete it from the session data.
		//It there's no matching key in the session data then this will return an empty string.
		Flash: a.sessionManager.PopString(r.Context(), "flash"),
		//Adds the authentication status to the template data.
		IsAuthenticated: a.isAuthenticated(r),
		CSRFToken: nosurf.Token(r),
	}
}

func (a *application) decodePostForm(r *http.Request, dst any) error{
	//Add any data in POST request bodies to the r.PostForm map.
	err := r.ParseForm()
	if err != nil{
		return err
	}

	//This will fill the target destination (dst) with the relevant values from the HTML form.
	err = a.formDecoder.Decode(dst, r.PostForm)
	if err != nil {
		
		//If the target destination (dst) is invalid, the Decode() method will return an error with the type *form.InvalidDecoderError.
		//We use errors.As() to check if we got that exact error.
		var invalidDecodeError *form.InvalidDecoderError
		if errors.As(err, &invalidDecodeError){
			panic(err)
		}

		//For all other errors, we return as normal.
		return err
	}
	return nil
}

func (a *application) isAuthenticated(r *http.Request) bool {
	return a.sessionManager.Exists(r.Context(), "authenticatedUserID")
}