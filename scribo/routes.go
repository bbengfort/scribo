package scribo

import "net/http"

// HandlerFunc is wrapper for creating handler functions with app callbacks.
type HandlerFunc func(app *App) http.HandlerFunc

// Route allows easy definition of our API
type Route struct {
	Name      string
	Methods   []string
	Pattern   string
	Handler   HandlerFunc
	Authorize bool
}

// Routes defines the complete route set for various non-API handlers.
type Routes []Route

var routes = Routes{
	Route{
		"Index", []string{GET}, "/", Index, false,
	},
	CreateResourceRoute(NodeCollection{}, "NodeCollection", "/nodes"),
	CreateResourceRoute(NodeDetail{}, "NodeDetail", "/nodes/{ID}"),
	CreateResourceRoute(PingCollection{}, "PingCollection", "/pings"),
	CreateResourceRoute(PingDetail{}, "PingDetail", "/pings/{ID}"),
}
