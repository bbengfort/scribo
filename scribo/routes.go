package scribo

import "net/http"

// Route allows easy definition of our API
type Route struct {
	Name    string
	Method  string
	Pattern string
	Handler http.HandlerFunc
}

// Routes defines the complete route set for the API
type Routes []Route

var routes = Routes{
	Route{
		"Index", "GET", "/", Index,
	},
	Route{
		"NodeList", "GET", "/nodes", NodeList,
	},
	Route{
		"NodeCreate", "POST", "/nodes", NodeCreate,
	},
	Route{
		"NodeDetail", "GET", "/nodes/{ID}", NodeDetail,
	},
	Route{
		"NodeUpdate", "PUT", "/nodes/{ID}", NodeUpdate,
	},
	Route{
		"NodeDelete", "DELETE", "/nodes/{ID}", NodeDelete,
	},
	Route{
		"PingList", "GET", "/pings", PingList,
	},
	Route{
		"PingCreate", "POST", "/pings", PingCreate,
	},
	Route{
		"PingDetail", "GET", "/pings/{ID}", PingDetail,
	},
	Route{
		"PingUpdate", "PUT", "/pings/{ID}", PingUpdate,
	},
	Route{
		"PingDelete", "DELETE", "/pings/{ID}", PingDelete,
	},
}
