package main

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/gorilla/csrf"
	"github.com/haidongNg/lenslocked/controllers"
	"github.com/haidongNg/lenslocked/models"
	"github.com/haidongNg/lenslocked/templates"
	"github.com/haidongNg/lenslocked/views"
)

func main() {
	email := models.Email{
		From:      models.DefaultSender,
		To:        "abc@doamain.com",
		Subject:   "Test mail",
		Plaintext: "TEST MAILTRAP",
	}
	es := models.NewMailService(models.SMTPConfig{
		Host:     "sandbox.smtp.mailtrap.io",
		Port:     "25",
		Username: "",
		Password: "",
	})

	// test send mail
	err := es.Send(email)
	if err != nil {
		panic(err)
	}

	// Test send mail reset password
	err = es.ForgotPassword(email.To, "http://localhost:3000/signup")
	if err != nil {
		panic(err)
	}

	fmt.Println("Email Send")
	r := chi.NewRouter()

	r.Get("/", controllers.StaticHandler(views.Must(views.ParseFS(templates.FS, "home-page.html", "paper.html"))))
	r.Get("/contact", controllers.StaticHandler(views.Must(views.ParseFS(templates.FS, "contact-page.html", "paper.html"))))
	r.Get("/faq", controllers.FAQ(views.Must(views.ParseFS(templates.FS, "faq-page.html", "paper.html"))))

	// Set up a database connection
	cfg := models.DefaultPostgresConfig()
	db, err := models.Open(cfg)

	if err != nil {
		panic(err)
	}

	defer db.Close()

	// Setup our model services
	userService := models.UserService{
		DB: db,
	}

	sessionService := models.SessionService{
		DB: db,
	}

	// Set up middleware
	umw := controllers.UserMiddleware{
		SessionService: &sessionService,
	}

	csrfKey := "gFvi45R4fy5xNBlnEeZtQbfAVCYEIAUX"
	// csrf.Secure(true) fix after deyploy
	csrfMw := csrf.Protect([]byte(csrfKey), csrf.Secure(false))

	// Setup our controllers
	userC := controllers.Users{
		UserService:    &userService,
		SessionService: &sessionService,
	}
	userC.Templates.New = views.Must(views.ParseFS(templates.FS, "signup-page.html", "paper.html"))
	userC.Templates.SignIn = views.Must(views.ParseFS(templates.FS, "signin-page.html", "paper.html"))
	// Set up router and routes
	r.Get("/signup", userC.New)
	r.Post("/users", userC.Create)
	r.Get("/signin", userC.SignIn)
	r.Post("/signin", userC.ProcessSignIn)
	r.With(umw.RequireUser).Post("/signout", userC.ProcessSignOut)
	r.Route("/users/me", func(r chi.Router) {
		r.Use(umw.RequireUser)
		r.Get("/", userC.CurrentUser)
	})
	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Page not found", http.StatusNotFound)
	})

	fmt.Printf("Starting the server on :3000...")
	http.ListenAndServe("127.0.0.1:3000", csrfMw(umw.SetUser(r)))
}
