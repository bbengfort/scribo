// Package scribo provides the library code for a lightweight, fast RESTful microservice.
//
// This service is an uptime collection mechanism associated with the Mora project. Mora contains three pieces: Oro and Scio which ping each other to measure latency inside of the network, then report those pings to Scribo, which is simply a RESTful API designed to record experimental data.
//
// The package is implemented by three commands: scribo, scribo-migrate, and scribo-register. To run the application locally:
//
//     $ scribo-migrate --all
//     $ scribo -port 8080
//
// So long as you have environment variables configured correctly, the database should be created and the web application will run. See the README for more information on getting started.
package scribo

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path"
	"strconv"

	"github.com/gorilla/mux"
)

// Version specifies what version of Scribo we're on.
const Version = "1.0"

// App defines the complete structure of a web application including a
// router for multiplexing, storages, assets, templates, cookies, and more.
// This should be the primary interface for working with Scribo.
type App struct {
	StaticDir   string
	TemplateDir string
	Templates   *template.Template
	Router      *mux.Router
	DB          *sql.DB
}

// CreateApp allows you to easily instantiate an App instance.
func CreateApp() *App {
	// Instantiate the app
	app := new(App)

	// Connect to the database
	app.DB = ConnectDB()

	// Set the static and template directories
	// BUG(bbengfort): relative import for static/template directories needs to be configured rather than guessed.
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
	app.AddStatic(app.StaticDir)

	return app
}

// Run the web application via the associated router.
func (app *App) Run(port int) {
	addr := fmt.Sprintf(":%d", port)

	name, err := os.Hostname()
	if err != nil {
		name = "localhost"
	}

	log.Printf("Starting server at http://%s:%d (use CTRL+C to quit)", name, port)
	log.Fatal(http.ListenAndServe(addr, app.Router))
}

// AddRoute allows you to add a handler for a specific route to the router.
func (app *App) AddRoute(route Route) {
	var handler http.Handler
	handler = route.Handler(app)

	if route.Authorize {
		handler = Authenticate(app, handler)
	}

	handler = Logger(app, handler)
	// handler = Debugger(app, handler)

	app.Router.
		Methods(route.Methods...).
		Path(route.Pattern).
		Name(route.Name).
		Handler(handler)
}

// AddStatic creates a handler to serve static files.
func (app *App) AddStatic(staticDir string) {
	static := http.StripPrefix("/assets/", http.FileServer(http.Dir(staticDir)))
	app.Router.PathPrefix("/assets/").Handler(Logger(app, static))
}

// Abort is a handler to terminate the request with no error message
func (app *App) Abort(w http.ResponseWriter, statusCode int) {
	w.WriteHeader(statusCode)
}

// Error is a handler to terminate the request with an error message.
func (app *App) Error(w http.ResponseWriter, err error, statusCode int) {
	http.Error(w, err.Error(), statusCode)
}

// JSONAbort is a handler to terminate the request with a JSON response
func (app *App) JSONAbort(w http.ResponseWriter, statusCode int) {
	w.Header().Set(CTKEY, CTJSON)
	w.WriteHeader(statusCode)

	response := make(map[string]string)
	response["code"] = strconv.Itoa(statusCode)
	response["reason"] = http.StatusText(statusCode)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// JSONError is a handler to terminate the request with a JSON response
func (app *App) JSONError(w http.ResponseWriter, err error, statusCode int) {
	w.Header().Set(CTKEY, CTJSON)
	w.WriteHeader(statusCode)

	response := make(map[string]string)
	response["code"] = strconv.Itoa(statusCode)
	response["error"] = err.Error()

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
