package scribo_test

import (
	"crypto/sha256"
	"net/http"
	"net/http/httptest"
	"time"

	. "github.com/bbengfort/scribo/scribo"
	"github.com/tent/hawk-go"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Auth", func() {

	Describe("creating HAWK shared secrets", func() {

		It("should create two different keys at different times", func() {
			node := &Node{
				Name:    "apollo",
				Address: "108.51.64.223",
				DNS:     "bryant.bengfort.com",
			}

			node.UpdateKey()
			key1 := node.Key

			time.Sleep(1)
			node.UpdateKey()
			Ω(key1).ShouldNot(Equal(node.Key))
		})

		It("should create a base64 encoded sha256 hash as key", func() {
			node := &Node{
				Name:    "apollo",
				Address: "108.51.64.223",
				DNS:     "bryant.bengfort.com",
			}

			node.UpdateKey()
			Ω(node.Key).Should(HaveLen(44))
		})

	})

	Describe("authentication HAWK credentials", func() {

		It("should return a 401 error when no credentials are provided", func() {
			request, err := http.NewRequest("GET", "http://localhost:8080/nodes", nil)
			response := httptest.NewRecorder()
			Ω(err).ShouldNot(HaveOccurred())

			handler := Authenticate(app, staticHandler(200, []byte("worked!")))
			handler.ServeHTTP(response, request)

			Ω(response.Code).Should(Equal(http.StatusUnauthorized))
		})

		It("should return a 403 error if bad credentials are provided", func() {
			app := createTestApp()

			request, err := http.NewRequest("GET", "http://localhost:8080/nodes", nil)
			response := httptest.NewRecorder()
			Ω(err).ShouldNot(HaveOccurred())

			// Create credentials for a node not in the database
			creds := new(hawk.Credentials)
			creds.ID = "badguy"
			creds.Key = "unlockthedoor"
			creds.Hash = sha256.New

			// Add bad authorization header to the request
			request.Header.Add("Content-Type", "application/json")
			auth := hawk.NewRequestAuth(request, creds, time.Duration(0))
			request.Header.Add("Authorization", auth.RequestHeader())

			handler := Authenticate(app, staticHandler(200, []byte("worked!")))
			handler.ServeHTTP(response, request)

			Ω(response.Code).Should(Equal(http.StatusForbidden))
		})

		It("should return a 200 if good credentials are provided", func() {
			app := createTestApp()
			_, err := app.DB.Exec("INSERT INTO nodes (name, key) VALUES ($1, $2)", "spiderman", "tinglingspideysense")
			Ω(err).ShouldNot(HaveOccurred())

			request, err := http.NewRequest("GET", "http://localhost:8080/nodes", nil)
			response := httptest.NewRecorder()
			Ω(err).ShouldNot(HaveOccurred())

			// Create credentials for a node not in the database
			creds := new(hawk.Credentials)
			creds.ID = "spiderman"
			creds.Key = "tinglingspideysense"
			creds.Hash = sha256.New

			// Add bad authorization header to the request
			request.Header.Add("Content-Type", "application/json")
			auth := hawk.NewRequestAuth(request, creds, time.Duration(0))
			request.Header.Add("Authorization", auth.RequestHeader())

			handler := Authenticate(app, staticHandler(200, []byte("worked!")))
			handler.ServeHTTP(response, request)

			Ω(response.Code).Should(Equal(http.StatusOK))
		})

	})

})
