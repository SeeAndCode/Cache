package server

import "Cache/internal/app/cache_server/cache"

const (
	SvrHTTP = "http"
	SvrCSP  = "csp"
)

type Server interface {
	cache.Cache
	Listen(addr string)
}

func New(typ string, c cache.Cache) Server {
	var svr Server
	switch typ {
	case SvrHTTP:
		svr = newHTTPServer(c)
	case SvrCSP:
		svr = newCSPServer(c)
	default:
		panic("unknown server type")
	}
	return svr
}
