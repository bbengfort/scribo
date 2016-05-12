package scribo_test

import (
	"net/http"
	"net/http/httptest"

	. "github.com/bbengfort/scribo/scribo"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

type BookResource struct{}

// Get returns a single book object
func (r BookResource) Get(app *App, request *http.Request) (int, interface{}, error) {
	data := make(map[string]string)
	data["title"] = "Data Analytics with Hadoop"
	data["author"] = "Benjamin Bengfort and Jenny Kim"
	return 200, data, nil
}

// Post returns a single book object and a created status code
func (r BookResource) Post(app *App, request *http.Request) (int, interface{}, error) {
	data := make(map[string]string)
	data["title"] = "Data Analytics with Hadoop"
	data["author"] = "Benjamin Bengfort and Jenny Kim"
	return 201, data, nil
}

// Put returns a 202 with no content associated (pending)
func (r BookResource) Put(app *App, request *http.Request) (int, interface{}, error) {
	return 202, nil, nil
}

// Delete returns a 204 (no content) after delete
func (r BookResource) Delete(app *App, request *http.Request) (int, interface{}, error) {
	return 204, nil, nil
}

var _ = Describe("Resource", func() {

	Describe("Route Creation", func() {

		It("should return a correct Route", func() {
			var aroute Route
			route := CreateResourceRoute(BookResource{}, "BookResource", "/books")
			Ω(route).Should(BeAssignableToTypeOf(aroute))
			Ω(route.Name).Should(Equal("BookResource"))
			Ω(route.Pattern).Should(Equal("/books"))
		})

		It("should route to GET, POST, PUT, and DELETE methods", func() {
			route := CreateResourceRoute(BookResource{}, "BookResource", "/books")
			Ω(route.Methods).Should(ContainElement(GET))
			Ω(route.Methods).Should(ContainElement(POST))
			Ω(route.Methods).Should(ContainElement(PUT))
			Ω(route.Methods).Should(ContainElement(DELETE))
		})

		It("should handle GET requests properly", func() {
			var handle http.HandlerFunc
			app := &App{"", "", nil, nil}
			request, _ := http.NewRequest(GET, "/books", nil)
			route := CreateResourceRoute(BookResource{}, "BookResource", "/books")

			Ω(route.Handler(app)).Should(BeAssignableToTypeOf(handle))

			handle = route.Handler(app)
			writer := httptest.NewRecorder()
			handle(writer, request)
			headers := writer.Header()

			Ω(writer.Code).Should(Equal(200))
			Ω(headers).Should(HaveKey(CTKEY))
			Ω(headers[CTKEY][0]).Should(Equal(CTJSON))
			Ω(writer.Body).Should(MatchJSON("{\"title\": \"Data Analytics with Hadoop\", \"author\": \"Benjamin Bengfort and Jenny Kim\"}"))
		})

		It("should handle POST requests properly", func() {
			var handle http.HandlerFunc
			app := &App{"", "", nil, nil}
			request, _ := http.NewRequest(POST, "/books", nil)
			route := CreateResourceRoute(BookResource{}, "BookResource", "/books")

			Ω(route.Handler(app)).Should(BeAssignableToTypeOf(handle))

			handle = route.Handler(app)
			writer := httptest.NewRecorder()
			handle(writer, request)
			headers := writer.Header()

			Ω(writer.Code).Should(Equal(201))
			Ω(headers).Should(HaveKey(CTKEY))
			Ω(headers[CTKEY][0]).Should(Equal(CTJSON))
			Ω(writer.Body).Should(MatchJSON("{\"title\": \"Data Analytics with Hadoop\", \"author\": \"Benjamin Bengfort and Jenny Kim\"}"))
		})

		It("should handle PUT requests properly", func() {
			var handle http.HandlerFunc
			app := &App{"", "", nil, nil}
			request, _ := http.NewRequest(PUT, "/books", nil)
			route := CreateResourceRoute(BookResource{}, "BookResource", "/books")

			Ω(route.Handler(app)).Should(BeAssignableToTypeOf(handle))

			handle = route.Handler(app)
			writer := httptest.NewRecorder()
			handle(writer, request)
			headers := writer.Header()

			Ω(writer.Code).Should(Equal(202))
			Ω(headers).Should(HaveKey(CTKEY))
			Ω(headers[CTKEY][0]).Should(Equal(CTJSON))
			Ω(writer.Body).Should(MatchJSON("null"))
		})

		It("should handle DELETE requests properly", func() {
			var handle http.HandlerFunc
			app := &App{"", "", nil, nil}
			request, _ := http.NewRequest(DELETE, "/books", nil)
			route := CreateResourceRoute(BookResource{}, "BookResource", "/books")

			Ω(route.Handler(app)).Should(BeAssignableToTypeOf(handle))

			handle = route.Handler(app)
			writer := httptest.NewRecorder()
			handle(writer, request)
			headers := writer.Header()

			Ω(writer.Code).Should(Equal(204))
			Ω(headers).Should(HaveKey(CTKEY))
			Ω(headers[CTKEY][0]).Should(Equal(CTJSON))
			Ω(writer.Body).Should(MatchJSON("null"))
		})

	})

	Describe("Unresponsive Resource", func() {

		type Unresponsive struct {
			GetNotSupported
			PostNotSupported
			PutNotSupported
			DeleteNotSupported
		}

		type JSONObject map[string]string

		var (
			resource Unresponsive  // The unresponsive resource being tested
			app      *App          // The Web Application hook
			request  *http.Request // A mock HTTP request
			url      string        // A mock URL endpoint for the resource
		)

		BeforeEach(func() {
			// Set up the test suite
			resource = Unresponsive{}
			app = &App{"", "", nil, nil}
			url = "/test"
		})

		It("should not respond to GET requests", func() {
			request, _ = http.NewRequest(GET, url, nil)
			code, data, err := resource.Get(app, request)

			Ω(code).Should(Equal(http.StatusMethodNotAllowed))
			Ω(err).Should(BeNil())

			result, _ := data.(map[string]string)

			Ω(result).Should(HaveKey("code"))
			Ω(result).Should(HaveKey("reason"))
			Ω(result).Should(HaveKey("message"))
			Ω(result["message"]).Should(Equal("This resource does not support HTTP GET."))
		})

		It("should not respond to POST requests", func() {
			request, _ = http.NewRequest(POST, url, nil)
			code, data, err := resource.Post(app, request)

			Ω(code).Should(Equal(http.StatusMethodNotAllowed))
			Ω(err).Should(BeNil())

			result, _ := data.(map[string]string)

			Ω(result).Should(HaveKey("code"))
			Ω(result).Should(HaveKey("reason"))
			Ω(result).Should(HaveKey("message"))
			Ω(result["message"]).Should(Equal("This resource does not support HTTP POST."))
		})

		It("should not respond to PUT requests", func() {
			request, _ = http.NewRequest(PUT, url, nil)
			code, data, err := resource.Put(app, request)

			Ω(code).Should(Equal(http.StatusMethodNotAllowed))
			Ω(err).Should(BeNil())

			result, _ := data.(map[string]string)

			Ω(result).Should(HaveKey("code"))
			Ω(result).Should(HaveKey("reason"))
			Ω(result).Should(HaveKey("message"))
			Ω(result["message"]).Should(Equal("This resource does not support HTTP PUT."))
		})

		It("should not respond to DELETE requests", func() {
			request, _ = http.NewRequest(GET, url, nil)
			code, data, err := resource.Delete(app, request)

			Ω(code).Should(Equal(http.StatusMethodNotAllowed))
			Ω(err).Should(BeNil())

			result, _ := data.(map[string]string)

			Ω(result).Should(HaveKey("code"))
			Ω(result).Should(HaveKey("reason"))
			Ω(result).Should(HaveKey("message"))
			Ω(result["message"]).Should(Equal("This resource does not support HTTP DELETE."))
		})

	})

})
