package scribo_test

import (
	"net/http"

	. "github.com/bbengfort/scribo/scribo"

	. "github.com/onsi/ginkgo"
	// . "github.com/onsi/gomega"
)

var _ = Describe("App", func() {

})

func createTestApp() *App {
	app := new(App)

	app.DB = db

	return app
}

func staticHandler(code int, data []byte) http.Handler {
	var handler http.HandlerFunc
	handler = func(w http.ResponseWriter, request *http.Request) {
		w.WriteHeader(code)
		w.Write(data)
	}
	return handler
}
