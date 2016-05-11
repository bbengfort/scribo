package scribo

import (
	"encoding/json"
	"fmt"
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
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

// NodeList handles the listing of todo items
func NodeList(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json;charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(nodes); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// NodeCreate handles the creation of a todo item from post
func NodeCreate(w http.ResponseWriter, r *http.Request) {
	var node Node
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := r.Body.Close(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := json.Unmarshal(body, &node); err != nil {
		w.Header().Set("Content-Type", "application/json;charset=UTF-8")
		w.WriteHeader(422) // unprocessable entity
		if err := json.NewEncoder(w).Encode(err); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	t := RepoCreateNode(node)
	w.Header().Set("Content-Type", "application/json;charset=UTF-8")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(t); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// NodeDetail handles the return of a single item by index.
func NodeDetail(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	nodeID, err := strconv.ParseUint(vars["ID"], 0, 64)

	if err != nil {
		w.Header().Set("Content-Type", "application/json;charset=UTF-8")
		w.WriteHeader(http.StatusBadRequest)
		if err := json.NewEncoder(w).Encode(err); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	node := RepoFindNode(nodeID)

	w.Header().Set("Content-Type", "application/json;charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(node); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// NodeUpdate handles the return of a single item by index.
func NodeUpdate(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	nodeID := vars["ID"]
	fmt.Fprintln(w, "Node update:", nodeID)
}

// NodeDelete handles the return of a single item by index.
func NodeDelete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	nodeID := vars["ID"]
	fmt.Fprintln(w, "Node delete:", nodeID)
}

// PingList handles the listing of todo items
func PingList(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json;charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(pings); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// PingCreate handles the creation of a todo item from post
func PingCreate(w http.ResponseWriter, r *http.Request) {
	var ping Ping
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := r.Body.Close(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := json.Unmarshal(body, &ping); err != nil {
		w.Header().Set("Content-Type", "application/json;charset=UTF-8")
		w.WriteHeader(422) // unprocessable entity
		if err := json.NewEncoder(w).Encode(err); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	t := RepoCreatePing(ping)
	w.Header().Set("Content-Type", "application/json;charset=UTF-8")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(t); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// PingDetail handles the return of a single item by index.
func PingDetail(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	pingID := vars["ID"]
	fmt.Fprintln(w, "Ping show:", pingID)
}

// PingUpdate handles the return of a single item by index.
func PingUpdate(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	pingID := vars["ID"]
	fmt.Fprintln(w, "Ping update:", pingID)
}

// PingDelete handles the return of a single item by index.
func PingDelete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	pingID := vars["ID"]
	fmt.Fprintln(w, "Ping delete:", pingID)
}
