package api

type Endpoints []Endpoint
type Endpoint struct {
	Method             string
	Path               string
	EndpointHandler    Handler
	MiddlewareHandlers []MiddleWare
}

func NewEnpoint(method, path string, endpointHandler Handler, mid ...MiddleWare) Endpoint {
	return Endpoint{
		Method:             method,
		Path:               path,
		EndpointHandler:    endpointHandler,
		MiddlewareHandlers: mid,
	}
}

func (e Endpoints) Use(mid ...MiddleWare) {
	for i := 0; i < len(e); i++ {
		e[i].MiddlewareHandlers = append(e[i].MiddlewareHandlers, mid...)
	}
}
