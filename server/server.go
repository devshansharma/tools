package server

type ConfigOption func(s *Server)

// Server for server configuration
type Server struct {
	port string
}

// Run starts the server and blocks
func (s *Server) Run() {

}

// WithPort provide port with colon, e.g., :8080
func WithPort(port string) ConfigOption {
	return func(s *Server) {
		s.port = port
	}
}

// New to create a new server with configuration options
func New(opts ...ConfigOption) *Server {
	s := &Server{}

	for _, opt := range opts {
		opt(s)
	}

	return s
}
