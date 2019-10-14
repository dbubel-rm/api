package api

type Endpoints []endpoint
type endpoint struct {
	Method             string
	Path               string
	EndpointHandler    Handler
	MiddlewareHandlers []MiddleWare
}

func NewEnpoint(method, path string, endpointHandler Handler, mid ...MiddleWare) endpoint {
	return endpoint{
		Method:             method,
		Path:               path,
		EndpointHandler:    endpointHandler,
		MiddlewareHandlers: mid,
	}
}

func (ep Endpoints) Use(mid ...MiddleWare) {
	for i := 0; i < len(ep); i++ {
		ep[i].MiddlewareHandlers = append(ep[i].MiddlewareHandlers, mid...)
	}
}
