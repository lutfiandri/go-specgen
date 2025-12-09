package specgen

type Route struct {
	Tags        []string
	Summary     string
	Description string
	Path        string
	Method      string
	Request     any
	Responses   []RouteResponse
}

type RouteResponse struct {
	StatusCode int
	Response   any
}
