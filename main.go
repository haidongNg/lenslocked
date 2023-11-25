package main

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/haidongNg/lenslocked/controllers"
	"github.com/haidongNg/lenslocked/templates"
	"github.com/haidongNg/lenslocked/views"
)

func main() {
	r := chi.NewRouter()

	r.Get("/", controllers.StaticHandler(views.Must(views.ParseFS(templates.FS, "home-page.html", "tailwind.html"))))
	r.Get("/contact", controllers.StaticHandler(views.Must(views.ParseFS(templates.FS, "contact-page.html", "tailwind.html"))))
	r.Get("/faq", controllers.FAQ(views.Must(views.ParseFS(templates.FS, "faq-page.html", "tailwind.html"))))

	userC := controllers.Users{}
	userC.Templates.New = views.Must(views.ParseFS(templates.FS, "signup-page.html", "tailwind.html"))
	r.Get("/signup", userC.New)
	r.Post("/users", userC.Create)
	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Page not found", http.StatusNotFound)
	})
	fmt.Printf("Starting the server on :3000...")
	http.ListenAndServe(":3000", r)
}
