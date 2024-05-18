package server

// Server for server configuration
type Server struct {
}

// Run starts the server and blocks
func (s *Server) Run() {

}

// New to create a new server with configuration options
func New(opts ...func(s *Server)) *Server {
	s := &Server{}

	for _, opt := range opts {
		opt(s)
	}

	return s
}
