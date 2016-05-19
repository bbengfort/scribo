package scribo

import (
	"log"
	"net/http"
	"net/http/httputil"
	"time"
)

// :remote-addr - :remote-user [:date[clf]] ":method :url HTTP/:http-version" :status :res[content-length]
const common = "%s - %s [%s] \"%s %s HTTP/%s\" %s %s"

// :method :url :status :response-time ms - :res[content-length]
const dev = "%s %s %d %s - %d"

type loggingResponseWriter interface {
	http.ResponseWriter
	http.Flusher
	Status() int
	Size() int
}

// responseLogger is a wrapper of http.ResponseWriter that keeps track of its
// HTTP status code and body size for reporting to the console.
type responseLogger struct {
	w      http.ResponseWriter
	status int
	size   int
}

func (l *responseLogger) Header() http.Header {
	return l.w.Header()
}

func (l *responseLogger) Write(b []byte) (int, error) {
	if l.status == 0 {
		l.status = http.StatusOK
	}

	size, err := l.w.Write(b)
	l.size += size
	return size, err
}

func (l *responseLogger) WriteHeader(s int) {
	l.w.WriteHeader(s)
	l.status = s
}

func (l *responseLogger) Status() int {
	return l.status
}

func (l *responseLogger) Size() int {
	return l.size
}

func (l *responseLogger) Flush() {
	f, ok := l.w.(http.Flusher)
	if ok {
		f.Flush()
	}
}

// Logger is a decorator for http handlers to record requests in dev format.
func Logger(app *App, inner http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		lw := &responseLogger{w: w}
		inner.ServeHTTP(lw, r)

		log.Printf(dev,
			r.Method, r.RequestURI, lw.Status(), time.Since(start), lw.Size(),
		)
	})
}

// Debugger is a decorator for http handlers to print out the incoming request
func Debugger(app *App, inner http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Print the request
		data, err := httputil.DumpRequest(r, true)
		if err != nil {
			log.Fatalf(err.Error())
		} else {
			log.Printf("%s\n", data)
		}

		// Now serve the request forward
		inner.ServeHTTP(w, r)
	})
}
