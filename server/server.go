package server

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type ConfigOption func(s *Server)

// Server for server configuration
type Server struct {
	addr              string
	handler           http.Handler
	tlsEnabled        bool
	certFile          string
	keyFile           string
	readTimeout       time.Duration
	readHeaderTimeout time.Duration
	writeTimeout      time.Duration
	idelTimeout       time.Duration
	serverTimeout     time.Duration
}

// Run starts the server and blocks
func (s *Server) Run(ctx context.Context) {
	srv := &http.Server{
		Addr:              s.addr,
		Handler:           s.handler,
		ReadTimeout:       s.readTimeout,
		WriteTimeout:      s.writeTimeout,
		IdleTimeout:       s.idelTimeout,
		ReadHeaderTimeout: s.readHeaderTimeout,
	}

	go func() {
		var err error
		if s.tlsEnabled {
			err = srv.ListenAndServeTLS(s.certFile, s.keyFile)
		} else {
			err = srv.ListenAndServe()
		}

		if err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	// kill (no param) default send syscanll.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall. SIGKILL but can"t be catch, so don't need add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	slog.WarnContext(ctx, "Shutdown Server ...")

	ctx, cancel := context.WithTimeout(ctx, s.serverTimeout)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		slog.ErrorContext(ctx, "Server Shutdown:", err)
	}

	<-ctx.Done()
	slog.WarnContext(ctx, "Server exiting")
}

// WithServerTimeout for server timeout
func WithServerTimeout(timeout time.Duration) ConfigOption {
	return func(srv *Server) {
		srv.serverTimeout = timeout
	}
}

// WithWriteTimeout for read timeout
func WithWriteTimeout(timeout time.Duration) ConfigOption {
	return func(srv *Server) {
		srv.writeTimeout = timeout
	}
}

// WithIdelTimeout for read timeout
func WithIdelTimeout(timeout time.Duration) ConfigOption {
	return func(srv *Server) {
		srv.idelTimeout = timeout
	}
}

// WithReadHeaderTimeout for read timeout
func WithReadHeaderTimeout(timeout time.Duration) ConfigOption {
	return func(srv *Server) {
		srv.readHeaderTimeout = timeout
	}
}

// WithReadTimeout for read timeout
func WithReadTimeout(timeout time.Duration) ConfigOption {
	return func(srv *Server) {
		srv.readTimeout = timeout
	}
}

// WithTlsEnabled to enable TLS
func WithTlsEnabled(enable bool) ConfigOption {
	return func(srv *Server) {
		srv.tlsEnabled = enable
	}
}

// WithCertFile for updating the certfile path
func WithCertFile(certFile string) ConfigOption {
	return func(srv *Server) {
		srv.certFile = certFile
	}
}

// WithKeyFile for updating the certfile path
func WithKeyFile(keyFile string) ConfigOption {
	return func(srv *Server) {
		srv.keyFile = keyFile
	}
}

// New to create a new server with configuration options
func New(addr string, handler http.Handler, opts ...ConfigOption) *Server {
	srv := &Server{
		addr:              addr,
		handler:           handler,
		serverTimeout:     3 * time.Second,
		readTimeout:       4 * time.Second,
		writeTimeout:      3 * time.Second,
		readHeaderTimeout: 2 * time.Second,
		idelTimeout:       30 * time.Second,
	}

	for _, opt := range opts {
		opt(srv)
	}

	return srv
}
