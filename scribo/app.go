package scribo

import (
	"html/template"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

// App defines the complete structure of a web application including a
// router for multiplexing, storages, assets, templates, cookies, and more.
// This should be the primary interface for working with Scribo.
type App struct {
	Router      *mux.Router
	StaticDir   string
	TemplateDir string
	Templates   map[string]*template.Template
}

// CreateApp allows you to easily instantiate an App instance.
// TODO: Allow the app to be available to other resources
func CreateApp() *App {
	app := new(App)
	app.Router = NewRouter()

	return app
}

// Run the web application via the associated router.
func (app *App) Run(addr string) {
	log.Printf("Starting server at http://%s (use CTRL+C to quit)", addr)
	log.Fatal(http.ListenAndServe(addr, app.Router))
}
