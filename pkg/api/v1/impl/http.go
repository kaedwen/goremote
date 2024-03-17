package impl

import (
	"context"
	"log"
	"net"
	"time"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/kaedwen/goremote/pkg/api/v1/gen"
	"github.com/kaedwen/goremote/pkg/common"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"
)

func RegisterHTTP(ctx context.Context, lg common.Logger, mux *runtime.ServeMux, gls net.Addr) {
	// provide generated swagger file
	if err := gen.RegisterSwaggerHandler(ctx, lg, mux, "/service/v1/swagger.json"); err != nil {
		log.Fatal("http swagger handler registration failed", zap.Error(err))
	}

	// provide generated GRPC gateway
	if err := gen.RegisterRemoteHandlerFromEndpoint(ctx, mux, gls.String(), []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithKeepaliveParams(keepalive.ClientParameters{
			Time:                30 * time.Second,
			Timeout:             15 * time.Second,
			PermitWithoutStream: true,
		}),
	}); err != nil {
		log.Panic("http gateway handler registration failed", zap.Error(err))
	}
}
