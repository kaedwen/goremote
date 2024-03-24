package grpc

import (
	"context"
	"net"
	"sync"
	"time"

	"github.com/kaedwen/goremote/pkg/common"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/reflection"
)

const (
	MaxConcurrentStreams      uint32 = 100
	MaxRecvMessageSize               = 1 * 1024 * 1024
	MaxSendMessageSize               = 1 * 1024 * 1024
	ConnReadBufferSize               = 32 * 1024
	ConnWriteBufferSize              = 32 * 1024
	GrpcMaxConnectionIdle            = 90 * time.Second
	GrpcMaxConnectionAge             = 0 * time.Second
	GrpcMaxConnectionAgeGrace        = 0 * time.Second
)

type GrpcServer struct {
	*grpc.Server
	lg common.Logger
}

type RegisterGrpc func(*GrpcServer)

func Serve(ctx context.Context, wtg *sync.WaitGroup, cfg *common.ConfigGRPC, lg common.Logger, r RegisterGrpc) net.Addr {
	srv := newServer(cfg, lg)
	r(srv)

	shz := health.NewServer()
	grpc_health_v1.RegisterHealthServer(srv, shz)

	lac := make(chan net.Addr, 1) // listener address notification channel

	wtg.Add(2)
	go srv.shutdownIfDone(ctx, wtg)
	go srv.listenAndServe(lac, wtg, cfg.Address())

	return <-lac
}

func newServer(cfg *common.ConfigGRPC, lg common.Logger) *GrpcServer {
	opts := []grpc.ServerOption{
		grpc.MaxConcurrentStreams(MaxConcurrentStreams),
		grpc.KeepaliveParams(keepalive.ServerParameters{
			MaxConnectionIdle:     GrpcMaxConnectionIdle,
			MaxConnectionAge:      GrpcMaxConnectionAge,
			MaxConnectionAgeGrace: GrpcMaxConnectionAgeGrace,
			Time:                  30 * time.Second,
			Timeout:               20 * time.Second,
		}),
		grpc.KeepaliveEnforcementPolicy(keepalive.EnforcementPolicy{
			MinTime:             0,
			PermitWithoutStream: true,
		}),
		grpc.MaxRecvMsgSize(MaxRecvMessageSize),
		grpc.MaxSendMsgSize(MaxSendMessageSize),
		grpc.ReadBufferSize(ConnReadBufferSize),
		grpc.WriteBufferSize(ConnWriteBufferSize),
		grpc.UnaryInterceptor(func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
			lg.Info("Received request", zap.Any("request", req), zap.String("method", info.FullMethod))
			return handler(ctx, req)
		}),
	}

	s := &GrpcServer{grpc.NewServer(opts...), lg}

	if cfg.Reflection {
		reflection.Register(s)
	}

	return s
}

func (s *GrpcServer) listenAndServe(lac chan<- net.Addr, wtg *sync.WaitGroup, addr string) {
	defer wtg.Done()

	ln, err := net.Listen("tcp", addr)
	if err != nil {
		s.lg.Fatal("listener failed", zap.Error(err))
	}

	// make address public
	lac <- ln.Addr()
	close(lac)

	switch err := s.Serve(ln); err {
	case nil:
	default:
		s.lg.Error("server stopped", zap.Error(err))
	}
}

func (s *GrpcServer) shutdownIfDone(ctx context.Context, wtg *sync.WaitGroup) {
	defer wtg.Done()

	<-ctx.Done()
	s.GracefulStop()
}
