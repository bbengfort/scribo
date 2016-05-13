package scribo

import (
	"crypto/sha256"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/tent/hawk-go"
)

// Helper function that looks up a Node's credentials by their name.
func getCredentials(app *App, c *hawk.Credentials) error {
	var key string

	// Lookup node by name (the ID specified in the request)
	query := "SELECT key FROM nodes WHERE name=$1"
	row := app.DB.QueryRow(query, c.ID)
	err := row.Scan(&key)

	if err != nil {

		if err == sql.ErrNoRows {
			err := new(hawk.CredentialError)
			err.Type = hawk.UnknownID
			err.Credentials = c

			return err
		}

		return err
	}

	// Otherwise we're good to go! Update the Credentials
	c.Key = key
	c.Hash = sha256.New
	return nil
}

// Authenticate is decorator that implements Hawk authorization.
func Authenticate(app *App, inner http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// Get the authentication from the request, using a closure.
		auth, err := hawk.NewAuthFromRequest(r, func(c *hawk.Credentials) error {
			return getCredentials(app, c)
		}, nil)

		// If the parsing didn't fail, check to see if auth is valid.
		if err == nil {
			err = auth.Valid()
		}

		if err != nil {
			var statusCode int

			if r.Header.Get("Authorization") == "" {
				statusCode = http.StatusUnauthorized
				w.Header().Set("WWW-Authenticate", "Hawk")
			} else {
				statusCode = http.StatusForbidden
			}

			w.Header().Set(CTKEY, CTJSON)
			w.WriteHeader(statusCode)

			response := make(map[string]string)
			response["code"] = strconv.Itoa(statusCode)
			response["error"] = err.Error()

			if err := json.NewEncoder(w).Encode(response); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			return
		}

		inner.ServeHTTP(w, r)
	})
}

// UpdateKey is the key creation mechanism for the node. It combines the Node
// name, address, and dns fields with the current time and a simple secret
// that is loaded from the environment. This method then sets the base64
// encoded sha256 hash of the generated string as the key on the node.
// Note: this method does not update the database!
func (node *Node) UpdateKey() {
	// Generate the plaintext version of the key.
	secret := os.Getenv("SCRIBO_SECRET")
	rawkey := fmt.Sprintf("%s:%s:%s:%s:%s", secret, node.Name, node.Address, node.DNS, time.Now())

	// Write the Hash
	hash := sha256.New()
	hash.Write([]byte(rawkey))

	// Set the base64 encoded key on the node.
	node.Key = base64.URLEncoding.EncodeToString(hash.Sum(nil))
}
