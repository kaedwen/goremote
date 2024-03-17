package http

import (
	"context"
	"net"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/kaedwen/goremote/pkg/common"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/health/grpc_health_v1"
)

func newLivenessHandler(ctx context.Context, lg common.Logger, gls net.Addr) runtime.HandlerFunc {
	if gls != nil {
		ghc, err := grpc.DialContext(ctx, gls.String(), grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			lg.Fatal("grpc dial error in liveness probe setup", zap.Error(err))
		}
		go func() {
			<-ctx.Done()
			_ = ghc.Close()
		}()

		return func(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
			stc := http.StatusOK

			res, err := grpc_health_v1.NewHealthClient(ghc).Check(r.Context(), &grpc_health_v1.HealthCheckRequest{})
			if err != nil {
				stc = http.StatusInternalServerError
				lg.Info("http get liveness, grpc call error", zap.Int("status", stc), zap.Error(err))
			} else if res.Status != grpc_health_v1.HealthCheckResponse_SERVING {
				stc = http.StatusInternalServerError
				lg.Info("http get liveness, grpc no health service", zap.Int("status", stc))
			} else {
				lg.Debug("http get liveness", zap.Int("status", stc))
			}

			w.WriteHeader(stc)
		}
	} else {
		return func(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
			stc := http.StatusOK
			w.WriteHeader(stc)
			lg.Debug("http get liveness", zap.Int("status", http.StatusOK))
		}
	}

}

func newReadinessHandler(lg common.Logger) runtime.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
		stc := http.StatusOK
		w.WriteHeader(stc)
		lg.Debug("http get readiness", zap.Int("status", stc))
	}
}
