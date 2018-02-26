package escpos

// ServerOption is a server option.
type ServerOption func(*Server) error

// WithLog is a server option to set a logging func.
func WithLog(f func(string, ...interface{})) ServerOption {
	return func(s *Server) error {
		s.logger = f
		return nil
	}
}
