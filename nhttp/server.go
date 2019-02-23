package nhttp

import (
	"context"
	"crypto/tls"
	"errors"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"sync/atomic"
	"syscall"
	"time"

	netutils "github.com/gokit/npkg/nnet"
	"golang.org/x/crypto/acme/autocert"
)

const (
	shutdownDuration = time.Second * 30
)

var (
	// ErrUnhealthy is returned when a server is considered unhealthy.
	ErrUnhealthy = errors.New("Service is unhealthy")
)

// HealthPinger exposes what we expect to provide us a health check for a server.
type HealthPinger interface {
	Ping() error
}

type healthPinger struct {
	healthy uint32
}

// Ping implements the HealthPinger interface.
func (h *healthPinger) Ping() error {
	if atomic.LoadUint32(&h.healthy) == 1 {
		return ErrUnhealthy
	}
	return nil
}

func (h *healthPinger) setUnhealthy() {
	atomic.StoreUint32(&h.healthy, 1)
}
func (h *healthPinger) setHealthy() {
	atomic.StoreUint32(&h.healthy, 0)
}

// Server implements a http server wrapper.
type Server struct {
	http2           bool
	shutdownTimeout time.Duration
	handler         http.Handler
	health          *healthPinger
	server          *http.Server
	listener        net.Listener
	tlsConfig       *tls.Config
	waiter          sync.WaitGroup
	man             *autocert.Manager
	closer          chan struct{}
}

// NewServer returns a new server which uses http instead of https.
func NewServer(handler http.Handler, shutdown ...time.Duration) *Server {
	var shutdownDur = shutdownDuration
	if len(shutdown) != 0 {
		shutdownDur = shutdown[0]
	}

	var health healthPinger
	var server Server
	server.http2 = false
	server.health = &health
	server.handler = handler
	server.shutdownTimeout = shutdownDur
	return &server
}

// NewServerWithTLS returns a new server which uses the provided tlsconfig for https connections.
func NewServerWithTLS(http2 bool, tconfig *tls.Config, handler http.Handler, shutdown ...time.Duration) *Server {
	var shutdownDur = shutdownDuration
	if len(shutdown) != 0 {
		shutdownDur = shutdown[0]
	}

	var health healthPinger
	var server Server
	server.http2 = http2
	server.health = &health
	server.handler = handler
	server.tlsConfig = tconfig
	server.shutdownTimeout = shutdownDur
	return &server
}

// NewServerWithCertMan returns a new server which uses the provided autocert certificate
// manager to provide http certificate.
func NewServerWithCertMan(http2 bool, man *autocert.Manager, handler http.Handler, shutdown ...time.Duration) *Server {
	var shutdownDur = shutdownDuration
	if len(shutdown) != 0 {
		shutdownDur = shutdown[0]
	}

	var health healthPinger
	var server Server
	server.man = man
	server.http2 = http2
	server.health = &health
	server.handler = handler
	server.shutdownTimeout = shutdownDur
	return &server
}

// Listen creates new http listen for giving addr and returns any error
// that occurs in attempt to starting the server.
func (s *Server) Listen(ctx context.Context, addr string) error {
	s.closer = make(chan struct{})

	var tlsConfig = s.tlsConfig
	if tlsConfig == nil && s.man == nil {
		tlsConfig = &tls.Config{
			GetCertificate: s.man.GetCertificate,
		}
	}

	if s.http2 {
		tlsConfig.NextProtos = append(tlsConfig.NextProtos, "h2")
	}

	listener, err := netutils.MakeListener("tcp", addr, tlsConfig)
	if err != nil {
		return err
	}

	tlsListener, ok := listener.(*net.TCPListener)
	if !ok {
		return errors.New("not tcp listener")
	}

	var server = &http.Server{
		Addr:           addr,
		Handler:        s.handler,
		ReadTimeout:    30 * time.Second,
		WriteTimeout:   30 * time.Second,
		MaxHeaderBytes: 1 << 20,
		TLSConfig:      tlsConfig,
	}

	s.health.setHealthy()

	var errs = make(chan error, 1)
	s.waiter.Add(1)
	go func() {
		defer s.waiter.Done()
		if err := server.Serve(netutils.NewKeepAliveListener(tlsListener)); err != nil {
			s.health.setUnhealthy()
			errs <- err
		}
	}()

	var signals = make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGTERM, os.Interrupt, syscall.SIGSTOP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	s.waiter.Add(1)
	go func() {
		defer s.waiter.Done()
		select {
		case <-ctx.Done():
			// server was called to close.
			s.gracefulShutdown(server)
		case <-s.closer:
			// server was closed intentionally.
			s.gracefulShutdown(server)
		case <-signals:
			// server received signal to close entirely.
			s.gracefulShutdown(server)
		}
	}()

	return <-errs
}

func (s *Server) gracefulShutdown(server *http.Server) {
	s.health.setUnhealthy()
	time.Sleep(20 * time.Second)
	var ctx, cancel = context.WithTimeout(context.Background(), s.shutdownTimeout)
	server.Shutdown(ctx)
	cancel()
}

// Close closes giving server
func (s *Server) Close() {
	select {
	case <-s.closer:
		return
	default:
		close(s.closer)
	}
}

// Wait blocks till server is closed.
func (s *Server) Wait(after ...func()) {
	s.waiter.Wait()
	for _, cb := range after {
		cb()
	}
}

// TLSManager returns the autocert.Manager associated with the giving server
// for its tls certificates.
func (s *Server) TLSManager() *autocert.Manager {
	return s.man
}

// Health returns the HealthPinger for giving server.
func (s *Server) Health() HealthPinger {
	return s.health
}
