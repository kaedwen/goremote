package http

import (
	"context"
	"net"
	"net/http"
	"sync"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/kaedwen/goremote/pkg/common"
	"go.uber.org/zap"
)

type RegisterHttp func(*runtime.ServeMux, net.Addr)

// Serve starts a HTTP server
// lifetime depends on ctx, wtg allows for lifetime synchronization
func Serve(ctx context.Context, wtg *sync.WaitGroup, cfg *common.ConfigHTTP, lg common.Logger, gls net.Addr, r RegisterHttp) net.Addr {
	srv := newServer(ctx, cfg, lg, gls)
	r(srv.mux, gls)

	lac := make(chan net.Addr, 1) // listener address notification channel

	wtg.Add(2)
	go srv.shutdownIfDone(ctx, wtg)
	go srv.listenAndServe(lac, wtg)

	return <-lac
}

type HttpServer struct {
	*http.Server
	lg  common.Logger
	mux *runtime.ServeMux
}

func newServer(ctx context.Context, cfg *common.ConfigHTTP, lg common.Logger, gls net.Addr) *HttpServer {
	mux := runtime.NewServeMux()

	if cfg.PathHealthz != nil {
		if err := mux.HandlePath(http.MethodGet, *cfg.PathHealthz, newLivenessHandler(ctx, lg, gls)); err != nil {
			lg.Fatal("failed to handle healthz", zap.Error(err))
		}
	}

	if cfg.PathReadyz != nil {
		if err := mux.HandlePath(http.MethodGet, *cfg.PathReadyz, newReadinessHandler(lg)); err != nil {
			lg.Fatal("failed to handle readyz", zap.Error(err))
		}
	}

	s := &HttpServer{
		&http.Server{
			Addr:              cfg.Address(),
			Handler:           mux,
			TLSConfig:         nil,
			ReadTimeout:       0,
			ReadHeaderTimeout: 0,
			WriteTimeout:      0,
			IdleTimeout:       0,
			TLSNextProto:      nil,
			ErrorLog:          nil,
			BaseContext:       nil,
			ConnContext:       nil,
		},
		lg, mux,
	}

	return s
}

func (s *HttpServer) listenAndServe(lac chan<- net.Addr, wtg *sync.WaitGroup) {
	defer wtg.Done()

	ln, err := net.Listen("tcp", s.Addr)
	if err != nil {
		s.lg.Fatal("listener failed", zap.Error(err))
	}

	// make address public
	lac <- ln.Addr()
	close(lac)

	switch err := s.Serve(ln); err {
	case http.ErrServerClosed:
	default:
		s.lg.Error("server stopped", zap.Error(err))
	}
}

func (s *HttpServer) shutdownIfDone(ctx context.Context, wtg *sync.WaitGroup) {
	defer wtg.Done()

	<-ctx.Done()

	_ = s.Shutdown(context.Background())
}
