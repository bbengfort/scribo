package scribo

import (
	"crypto/sha256"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/tent/hawk-go"
)

func getCredentials(c *hawk.Credentials) error {
	if c.ID == "benjamin" {
		c.Key = "werxhqb98rpaxn39848xrunpaw3489ruxnpa98w4rxn"
		c.Hash = sha256.New
		return nil
	}

	err := new(hawk.CredentialError)
	err.Type = hawk.UnknownID
	err.Credentials = c

	return err
}

// Authenticate is decorator that implements Hawk authorization.
func Authenticate(inner http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		auth, err := hawk.NewAuthFromRequest(r, getCredentials, nil)

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
