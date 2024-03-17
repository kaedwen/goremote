package main

import (
	"context"
	"net"
	"sync"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/kaedwen/goremote/pkg/common"
	"github.com/kaedwen/goremote/pkg/server/grpc"
	"github.com/kaedwen/goremote/pkg/server/http"

	v1 "github.com/kaedwen/goremote/pkg/api/v1/impl"
	"go.uber.org/zap"
)

func main() {
	cfg := &common.Config{}
	cfg.MustParse()

	lg := common.NewLogger(cfg)

	ctx, end := context.WithCancel(context.Background())
	go common.SigWatch(end, cfg.Common.GracePeriod, lg)

	rg := func(s *grpc.GrpcServer) {
		v1.RegisterGRPC(ctx, lg, cfg, s)
	}

	rh := func(m *runtime.ServeMux, a net.Addr) {
		v1.RegisterHTTP(ctx, lg, m, a)
	}

	var wtg sync.WaitGroup
	gls := grpc.Serve(ctx, &wtg, &cfg.GRPC, lg, rg)
	hls := http.Serve(ctx, &wtg, &cfg.HTTP, lg, gls, rh)

	lg.Info("listening", zap.Stringer("http", hls), zap.Stringer("grpc", gls))

	wtg.Wait()

	lg.Info("shutdown completed")
	lg.Sync()
}
