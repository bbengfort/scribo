package scribo

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// Index handles the root route by rendering a small web page that uses the
// API to display information about the ping status.
func Index(app *App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := app.Templates.ExecuteTemplate(w, "index", nil)
		if err != nil {
			app.Error(w, err, http.StatusInternalServerError)
			return
		}
	}
}

type (
	// NodeCollection is a RESTful resource for listing and creating nodes.
	NodeCollection struct {
		PutNotSupported
		DeleteNotSupported
	}

	// NodeDetail is a RESTful resource for updating and deleting nodes.
	NodeDetail struct {
		PostNotSupported
	}

	// PingCollection is a RESTful resource for listing and creating pings.
	PingCollection struct {
		PutNotSupported
		DeleteNotSupported
	}

	// PingDetail is a RESTful resource for updating and deleting pings.
	PingDetail struct {
		PostNotSupported
	}
)

// Get returns the listing of nodes
func (r NodeCollection) Get(app *App, request *http.Request) (int, interface{}, error) {
	nodes, err := FetchNodes(app.DB, 10)

	if err != nil {
		return http.StatusInternalServerError, nil, err
	}

	return http.StatusOK, nodes, nil
}

// Post handles the creation of a node from JSON in the request body
func (r NodeCollection) Post(app *App, request *http.Request) (int, interface{}, error) {
	var node Node

	// Read the data from the request stream (limit the size to 1 MB)
	body, err := ioutil.ReadAll(io.LimitReader(request.Body, 1048576))

	// Todo return a 413 (entity too large) if it's the limit that's reached.
	if err != nil {
		return http.StatusInternalServerError, nil, err
	}

	// Attempt to close the body of the request for reading
	if err := request.Body.Close(); err != nil {
		return http.StatusInternalServerError, nil, err
	}

	// Unmarshal the Post into a Node struct
	if err := json.Unmarshal(body, &node); err != nil {
		// If JSON parsing fails send back a 422 "unprocessable entity"
		response := make(map[string]string)
		response["code"] = "422"
		response["reason"] = "Could not parse JSON into a Node object."
		response["error"] = err.Error()
		return 422, response, nil
	}

	// Create the node in the database
	// TODO: implement node.Create
	// node.Create()
	return http.StatusCreated, node, nil
}

// Get returns a single node from the database.
func (r NodeDetail) Get(app *App, request *http.Request) (int, interface{}, error) {
	// Parse the variables from the URL route.
	vars := mux.Vars(request)
	nodeID, err := strconv.ParseUint(vars["ID"], 0, 64)

	if err != nil {
		return http.StatusInternalServerError, nil, err
	}

	// Query the database for the node by the ID.
	node, err := GetNode(app.DB, nodeID)
	if err != nil {
		return http.StatusNotFound, nil, err
	}

	return http.StatusOK, node, nil
}

// Put updates a node in the database
func (r NodeDetail) Put(app *App, request *http.Request) (int, interface{}, error) {
	// Parse the variables from the URL route.
	vars := mux.Vars(request)
	nodeID, err := strconv.ParseUint(vars["ID"], 0, 64)

	if err != nil {
		return http.StatusInternalServerError, nil, err
	}

	// Query the database for the node by the ID.
	node, err := GetNode(app.DB, nodeID)
	if err != nil {
		return http.StatusNotFound, nil, err
	}

	// Now perform the update ...
	// TODO: Add the update handling code
	// node.Save()

	// Return the updated node
	return http.StatusOK, node, nil
}

// Delete a node from the database
func (r NodeDetail) Delete(app *App, request *http.Request) (int, interface{}, error) {
	// Parse the variables from the URL route.
	vars := mux.Vars(request)
	nodeID, err := strconv.ParseUint(vars["ID"], 0, 64)

	if err != nil {
		return http.StatusInternalServerError, nil, err
	}

	// Query the database for the node by the ID.
	node, err := GetNode(app.DB, nodeID)
	if err != nil {
		return http.StatusNotFound, nil, err
	}

	// Delete the Node from the database
	// TODO: implement node.Delete
	// node.Delete()

	// Return the updated node
	// return http.StatusNoContent, nil, nil
	return http.StatusNoContent, node, nil
}

// Get returns the listing of pings
func (r PingCollection) Get(app *App, request *http.Request) (int, interface{}, error) {
	pings, err := FetchPings(app.DB, 10)

	if err != nil {
		return http.StatusInternalServerError, nil, err
	}

	return http.StatusOK, pings, nil
}

// Post handles the creation of a ping from JSON in the request body
func (r PingCollection) Post(app *App, request *http.Request) (int, interface{}, error) {
	var ping Ping

	// Read the data from the request stream (limit the size to 1 MB)
	body, err := ioutil.ReadAll(io.LimitReader(request.Body, 1048576))

	// Todo return a 413 (entity too large) if it's the limit that's reached.
	if err != nil {
		return http.StatusInternalServerError, nil, err
	}

	// Attempt to close the body of the request for reading
	if err := request.Body.Close(); err != nil {
		return http.StatusInternalServerError, nil, err
	}

	// Unmarshal the Post into a Node struct
	if err := json.Unmarshal(body, &ping); err != nil {
		// If JSON parsing fails send back a 422 "unprocessable entity"
		response := make(map[string]string)
		response["code"] = "422"
		response["reason"] = "Could not parse JSON into a Ping object."
		response["error"] = err.Error()
		return 422, response, nil
	}

	// Create the node in the database
	// TODO: implement ping.Create
	// ping.Create()

	return http.StatusCreated, ping, nil
}

// Get returns a single ping from the database.
func (r PingDetail) Get(app *App, request *http.Request) (int, interface{}, error) {
	// Parse the variables from the URL route.
	vars := mux.Vars(request)
	pingID, err := strconv.ParseUint(vars["ID"], 0, 64)

	if err != nil {
		return http.StatusInternalServerError, nil, err
	}

	// Query the database for the ping by the ID.
	ping, err := GetPing(app.DB, pingID)
	if err != nil {
		return http.StatusNotFound, nil, err
	}

	return http.StatusOK, ping, nil
}

// Put updates a ping in the database
func (r PingDetail) Put(app *App, request *http.Request) (int, interface{}, error) {
	// Parse the variables from the URL route.
	vars := mux.Vars(request)
	pingID, err := strconv.ParseUint(vars["ID"], 0, 64)

	if err != nil {
		return http.StatusInternalServerError, nil, err
	}

	// Query the database for the ping by the ID.
	ping, err := GetPing(app.DB, pingID)
	if err != nil {
		return http.StatusNotFound, nil, err
	}

	// Now perform the update ...
	// TODO: Add the update handling code
	// ping.Save()

	// Return the updated node
	return http.StatusOK, ping, nil
}

// Delete a ping from the database
func (r PingDetail) Delete(app *App, request *http.Request) (int, interface{}, error) {
	// Parse the variables from the URL route.
	vars := mux.Vars(request)
	pingID, err := strconv.ParseUint(vars["ID"], 0, 64)

	if err != nil {
		return http.StatusInternalServerError, nil, err
	}

	// Query the database for the ping by the ID.
	ping, err := GetPing(app.DB, pingID)
	if err != nil {
		return http.StatusNotFound, nil, err
	}

	// Delete the Ping from the database
	// TODO: implement ping.Delete
	// ping.Delete()

	// Return the updated node
	return http.StatusNoContent, ping, nil
	// return http.StatusNoContent, nil, nil
}
