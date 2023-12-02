package controllers

import (
	"fmt"
	"net/http"

	"github.com/haidongNg/lenslocked/models"
	"github.com/haidongNg/lenslocked/views"
)

type Users struct {
	Templates   UsersTemplates
	UserService *models.UserService
}

type UsersTemplates struct {
	New    views.Template
	SignIn views.Template
}

func (u Users) New(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Email    string
		Password string
	}

	data.Email = r.FormValue("email")
	data.Password = r.FormValue("password")
	u.Templates.New.Execute(w, data)
}

func (u Users) Create(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	email := r.FormValue("email")
	password := r.FormValue("password")
	user, err := u.UserService.Create(email, password)
	if err != nil {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "User create: %+v", user)
}

func (u Users) SignIn(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Email    string
		Password string
	}

	data.Email = r.FormValue("email")
	data.Password = r.FormValue("password")
	u.Templates.SignIn.Execute(w, data)
}

func (u Users) ProcessSignIn(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	email := r.FormValue("email")
	password := r.FormValue("password")

	user, err := u.UserService.Authenticate(email, password)

	if err != nil {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	cookie := http.Cookie{
		Name:     "email",
		Value:    user.Email,
		Path:     "/",
		HttpOnly: true,
	}
	http.SetCookie(w, &cookie)
	fmt.Fprintf(w, "User create: %+v", user)
}

func (u Users) CurrentUser(w http.ResponseWriter, r *http.Request) {
	email, err := r.Cookie("email")

	if err != nil {
		fmt.Fprint(w, "The email cookie could not be read.")
	}
	fmt.Fprintf(w, "Email cookie %s\n", email)
}
