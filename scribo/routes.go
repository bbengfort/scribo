package scribo

import "net/http"

// HandlerFunc is wrapper for creating handler functions with app callbacks.
type HandlerFunc func(app *App) http.HandlerFunc

// Route allows easy definition of our API
type Route struct {
	Name    string
	Method  string
	Pattern string
	Handler HandlerFunc
}

// Routes defines the complete route set for the API
type Routes []Route

var routes = Routes{
	Route{
		"Index", "GET", "/", Index,
	},
}
