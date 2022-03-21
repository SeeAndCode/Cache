package client

const (
	CliTypCSP  = "csp"
	CliTypHTTP = "http"
)

// Client to communicate with cache server
type Client interface {
	Set(key string, value []byte) error

	Get(key string) ([]byte, error)
}

// New return a Client will connect to cache server on srvAddr
func New(typ string, srvAddr string) Client {
	var cli Client
	switch typ {
	case CliTypCSP:
		cli = newCSPClient(srvAddr)
	default:
		panic("unknown client type, " + typ)
	}
	return cli
}
