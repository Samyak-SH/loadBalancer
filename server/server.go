package server

type Server struct {
	ServerURL string
	isAlive   bool
}

func NewServer(url string) *Server {
	return &Server{
		ServerURL: url,
		isAlive:   true,
	}
}

func (s *Server) IsAlive() bool {
	return s.isAlive
}

func (s *Server) GetServerURL() string {
	return s.ServerURL
}
