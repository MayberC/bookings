package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/MayberC/bookings/pkg/config"
	"github.com/MayberC/bookings/pkg/handlers"
	"github.com/MayberC/bookings/pkg/render"
	"github.com/alexedwards/scs/v2"
)

const portNumber = "8080"

var app config.AppConfig
var session *scs.SessionManager

// main is the main application function
func main() {

	// change this to true when in production
	app.InProduction = false

	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction

	app.Session = session

	ip := "127.0.0.1"
	address := fmt.Sprintf("%v:%v", ip, portNumber)

	tc, err := render.CreateTemplateCache()
	if err != nil {
		log.Fatal("cannot create template cache")
	}

	app.TemplateCache = tc
	app.UseCache = false

	repo := handlers.NewRepo(&app)
	handlers.NewHandlers(repo)

	render.NewTemplate(&app)

	log.Printf("Server open at: http://%v\n", address)

	srv := &http.Server{
		Addr:    address,
		Handler: routes(&app),
	}

	err = srv.ListenAndServe()

	log.Fatal(err)

}
