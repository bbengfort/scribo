package scribo

import (
	"html/template"
	"log"
	"net/http"
	"os"
	"path"

	"github.com/gorilla/mux"
)

// App defines the complete structure of a web application including a
// router for multiplexing, storages, assets, templates, cookies, and more.
// This should be the primary interface for working with Scribo.
type App struct {
	StaticDir   string
	TemplateDir string
	Templates   *template.Template
	Router      *mux.Router
}

// CreateApp allows you to easily instantiate an App instance.
// TODO: Allow the app to be available to other resources
func CreateApp() *App {
	// Instantiate the app
	app := new(App)

	// Set the static and template directories
	root, _ := os.Getwd()
	app.StaticDir = path.Join(root, "assets")
	app.TemplateDir = path.Join(root, "templates")

	// Load the templates from the template directory
	tmplGlob := path.Join(app.TemplateDir, "*")
	app.Templates = template.Must(template.ParseGlob(tmplGlob))

	app.Router = mux.NewRouter().StrictSlash(true)

	// Add all routes
	for _, route := range routes {
		app.AddRoute(route)
	}

	// Add a static file server pointing at the assets directory
	app.addStaticHandler(app.StaticDir)

	return app
}

// Run the web application via the associated router.
func (app *App) Run(addr string) {
	log.Printf("Starting server at http://%s (use CTRL+C to quit)", addr)
	log.Fatal(http.ListenAndServe(addr, app.Router))
}

// AddRoute allows you to add a handler for a specific route to the router.
func (app *App) AddRoute(route Route) {
	var handler http.Handler
	handler = route.Handler(app)
	handler = logger(handler, route.Name)

	app.Router.
		Methods(route.Method).
		Path(route.Pattern).
		Name(route.Name).
		Handler(handler)
}

// Add the static file server handler to the app
func (app *App) addStaticHandler(staticDir string) {
	static := http.StripPrefix("/assets/", http.FileServer(http.Dir(staticDir)))
	app.Router.PathPrefix("/assets/").Handler(logger(static, "Assets"))
}
