package api

type Endpoints []endpoint
type endpoint struct {
	Verb               string
	Path               string
	EndpointHandler    Handler
	MiddlewareHandlers []MiddleWare
}

func NewEndpoint(method, path string, endpointHandler Handler, mid ...MiddleWare) endpoint {
	return endpoint{
		Verb:               method,
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
