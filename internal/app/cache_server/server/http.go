package server

import "Cache/internal/app/cache_server/cache"

type httpServer struct {
	cache.Cache
}

func newHTTPServer(c cache.Cache) *httpServer {
	return &httpServer{c}
}

func (s *httpServer) Listen(addr string) {

}
