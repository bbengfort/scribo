package scribo

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

const (
	// GET refers to the HTTP GET method
	GET = "GET"

	// POST refers to the HTTP POST method
	POST = "POST"

	// PUT refers to the HTTP PUT method
	PUT = "PUT"

	// DELETE refers to the HTTP DELETE method
	DELETE = "DELETE"

	// CTKEY is the content type key to add to the header.
	CTKEY = "Content-Type"

	// CTJSON is the content type for responses
	CTJSON = "application/json;charset=UTF-8"
)

// Resource defines a REST endpoint with the standard methods. The standard
// methods should return an int (the status code) as well as any other data
// that is required, e.g. the interface second argument.
type Resource interface {
	Get(app *App, request *http.Request) (int, interface{}, error)
	Post(app *App, request *http.Request) (int, interface{}, error)
	Put(app *App, request *http.Request) (int, interface{}, error)
	Delete(app *App, request *http.Request) (int, interface{}, error)
}

// CreateResourceRoute returns a Route object, constructing the Handler and
// Method from the Resource definition. This allows you to quickly add
// Resources to the routes object for inclusion with the web app.
func CreateResourceRoute(resource Resource, name string, pattern string) Route {

	var handler HandlerFunc
	handler = func(app *App) http.HandlerFunc {
		return func(w http.ResponseWriter, request *http.Request) {
			var data interface{}
			var code int
			var err error

			switch request.Method {
			case GET:
				code, data, err = resource.Get(app, request)
			case POST:
				code, data, err = resource.Post(app, request)
			case PUT:
				code, data, err = resource.Put(app, request)
			case DELETE:
				code, data, err = resource.Delete(app, request)
			default:
				app.Abort(w, http.StatusNotImplemented)
			}

			// Handle errors from the resource
			if err != nil {
				if code == 0 {
					code = http.StatusInternalServerError
				}

				app.Error(w, err, code)
				return
			}

			// Write the content type and the status code
			w.Header().Set(CTKEY, CTJSON)
			w.WriteHeader(code)

			// Write the response data as JSON to the stream.
			if err := json.NewEncoder(w).Encode(data); err != nil {
				app.Error(w, err, http.StatusInternalServerError)
				return
			}
		}
	}

	return Route{name, []string{GET, POST, PUT, DELETE}, pattern, handler, true}
}

type (
	// GetNotSupported allows the creation of Resources with no Get method.
	GetNotSupported struct{}

	// PostNotSupported allows the creation of Resources with no Post method.
	PostNotSupported struct{}

	// PutNotSupported allows the creation of Resources with no Put method.
	PutNotSupported struct{}

	// DeleteNotSupported allows the creation of Resources with no Delete method.
	DeleteNotSupported struct{}
)

// Get returns method not allowed (405) on GetNotSupported types.
func (r GetNotSupported) Get(app *App, request *http.Request) (int, interface{}, error) {
	return notSupported(GET)
}

// Post returns method not allowed (405) on PostNotSupported types.
func (r PostNotSupported) Post(app *App, request *http.Request) (int, interface{}, error) {
	return notSupported(POST)
}

// Put returns method not allowed (405) on PutNotSupported types.
func (r PutNotSupported) Put(app *App, request *http.Request) (int, interface{}, error) {
	return notSupported(PUT)
}

// Delete returns method not allowed (405) on DeleteNotSupported types.
func (r DeleteNotSupported) Delete(app *App, request *http.Request) (int, interface{}, error) {
	return notSupported(DELETE)
}

// Helper function to return a not reported status code and reason.
func notSupported(method string) (int, interface{}, error) {
	response := make(map[string]string)
	code := http.StatusMethodNotAllowed

	response["code"] = strconv.Itoa(code)
	response["reason"] = http.StatusText(code)
	response["message"] = fmt.Sprintf("This resource does not support HTTP %s.", method)
	return code, response, nil
}
