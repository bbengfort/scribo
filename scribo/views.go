package scribo

import (
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

const (
	// StatusUnprocessableEntity is a missing status code constant in http
	StatusUnprocessableEntity = 422
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
		response["code"] = strconv.Itoa(StatusUnprocessableEntity)
		response["reason"] = "Could not parse JSON into a Node object."
		response["error"] = err.Error()
		return StatusUnprocessableEntity, response, nil
	}

	// Create the node in the database
	_, dberr := node.Save(app.DB)

	// Handle the creation conditions
	switch {
	case dberr != nil:
		return http.StatusConflict, nil, dberr
	default:
		return http.StatusCreated, node, nil
	}

}

// Get returns a single node from the database.
func (r NodeDetail) Get(app *App, request *http.Request) (int, interface{}, error) {
	// Parse the variables from the URL route.
	vars := mux.Vars(request)
	nodeID, err := strconv.ParseInt(vars["ID"], 0, 64)

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
	nodeID, err := strconv.ParseInt(vars["ID"], 0, 64)

	if err != nil {
		return http.StatusInternalServerError, nil, err
	}

	// Query the database for the node by the ID.
	node, err := GetNode(app.DB, nodeID)
	if err != nil {
		return http.StatusNotFound, nil, err
	}

	// Now perform the update ...
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

	// Create a temporary node to  unmarshall data to.
	var fields map[string]interface{}

	// Unmarshal the Put into a Node struct
	if err := json.Unmarshal(body, &fields); err != nil {
		// If JSON parsing fails send back a 422 "unprocessable entity"
		response := make(map[string]string)
		response["code"] = strconv.Itoa(StatusUnprocessableEntity)
		response["reason"] = "Could not parse JSON into a Node object."
		response["error"] = err.Error()
		return StatusUnprocessableEntity, response, nil
	}

	// Update the node with the fields that are updateable.
	if val, ok := fields["name"]; ok {
		node.Name = val.(string)
	}

	if val, ok := fields["address"]; ok {
		node.Address = val.(string)
	}

	if val, ok := fields["dns"]; ok {
		node.DNS = val.(string)
	}

	// Save the node updates in the database
	_, dberr := node.Save(app.DB)

	// Handle the creation conditions
	switch {
	case dberr != nil:
		return http.StatusConflict, nil, dberr
	default:
		return http.StatusOK, node, nil
	}

}

// Delete a node from the database
func (r NodeDetail) Delete(app *App, request *http.Request) (int, interface{}, error) {
	// Parse the variables from the URL route.
	vars := mux.Vars(request)
	nodeID, err := strconv.ParseInt(vars["ID"], 0, 64)

	if err != nil {
		return http.StatusInternalServerError, nil, err
	}

	// Query the database for the node by the ID.
	node, err := GetNode(app.DB, nodeID)
	if err != nil {
		return http.StatusNotFound, nil, err
	}

	// Delete the Node from the database
	deleted, err := node.Delete(app.DB)

	switch {
	case err != nil:
		return http.StatusInternalServerError, nil, err
	case !deleted:
		return http.StatusConflict, nil, errors.New("Unable to delete node!")
	default:
		return http.StatusNoContent, nil, nil
	}
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
		response["code"] = strconv.Itoa(StatusUnprocessableEntity)
		response["reason"] = "Could not parse JSON into a Ping object."
		response["error"] = err.Error()
		return StatusUnprocessableEntity, response, nil
	}

	// Create the ping in the database
	_, dberr := ping.Save(app.DB)

	// Handle the creation conditions
	switch {
	case dberr != nil:
		return http.StatusConflict, nil, dberr
	default:
		return http.StatusCreated, ping, nil
	}
}

// Get returns a single ping from the database.
func (r PingDetail) Get(app *App, request *http.Request) (int, interface{}, error) {
	// Parse the variables from the URL route.
	vars := mux.Vars(request)
	pingID, err := strconv.ParseInt(vars["ID"], 0, 64)

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
	pingID, err := strconv.ParseInt(vars["ID"], 0, 64)

	if err != nil {
		return http.StatusInternalServerError, nil, err
	}

	// Query the database for the ping by the ID.
	ping, err := GetPing(app.DB, pingID)
	if err != nil {
		return http.StatusNotFound, nil, err
	}

	// Now perform the update ...
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

	// Create a temporary node to  unmarshall data to.
	var fields map[string]interface{}

	// Unmarshal the Put into a Node struct
	if err := json.Unmarshal(body, &fields); err != nil {
		// If JSON parsing fails send back a 422 "unprocessable entity"
		response := make(map[string]string)
		response["code"] = strconv.Itoa(StatusUnprocessableEntity)
		response["reason"] = "Could not parse JSON into a Ping object."
		response["error"] = err.Error()
		return StatusUnprocessableEntity, response, nil
	}

	// Update the node with the fields that are updateable.
	if val, ok := fields["source"]; ok {
		ping.Source = val.(int64)
	}

	if val, ok := fields["target"]; ok {
		ping.Target = val.(int64)
	}

	if val, ok := fields["payload"]; ok {
		ping.Payload = val.(int)
	}

	if val, ok := fields["latency"]; ok {
		ping.Latency = val.(float64)
	}

	if val, ok := fields["timeout"]; ok {
		ping.Timeout = val.(bool)
	}

	// Save the node updates in the database
	_, dberr := ping.Save(app.DB)

	// Handle the creation conditions
	switch {
	case dberr != nil:
		return http.StatusConflict, nil, dberr
	default:
		return http.StatusOK, ping, nil
	}
}

// Delete a ping from the database
func (r PingDetail) Delete(app *App, request *http.Request) (int, interface{}, error) {
	// Parse the variables from the URL route.
	vars := mux.Vars(request)
	pingID, err := strconv.ParseInt(vars["ID"], 0, 64)

	if err != nil {
		return http.StatusInternalServerError, nil, err
	}

	// Query the database for the ping by the ID.
	ping, err := GetPing(app.DB, pingID)
	if err != nil {
		return http.StatusNotFound, nil, err
	}

	// Delete the Ping from the database
	deleted, err := ping.Delete(app.DB)

	switch {
	case err != nil:
		return http.StatusInternalServerError, nil, err
	case !deleted:
		return http.StatusConflict, nil, errors.New("Unable to delete ping!")
	default:
		return http.StatusNoContent, nil, nil
	}
}
