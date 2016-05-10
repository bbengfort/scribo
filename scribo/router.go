package scribo

import (
	"net/http"

	"github.com/gorilla/mux"
)

// NewRouter is the controller for our microservice. It creates a router
// from the API routes defined in the routes object then adds static asset
// handling and other types of template handling as well.
func NewRouter() *mux.Router {

	// Instantiate the router
	router := mux.NewRouter().StrictSlash(true)

	// Add the API routes
	for _, route := range routes {

		var handler http.Handler
		handler = route.Handler
		handler = logger(handler, route.Name)

		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(handler)
	}

	// Add a static file server pointing at the assets directory
	static := http.StripPrefix("/assets/", http.FileServer(http.Dir("./assets/")))
	router.PathPrefix("/assets/").Handler(logger(static, "Assets"))

	return router
}
