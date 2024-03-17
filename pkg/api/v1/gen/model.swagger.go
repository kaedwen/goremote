package gen

import (
	"context"
	_ "embed"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/kaedwen/goremote/pkg/common"
	"go.uber.org/zap"
)

//go:embed *.swagger.json
var SwaggerModelContent []byte

func RegisterSwaggerHandler(_ context.Context, lg common.Logger, mux *runtime.ServeMux, path string) error {
	return mux.HandlePath(http.MethodGet, path, func(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		cnt, err := w.Write(SwaggerModelContent)
		if err != nil {
			lg.Info("http", zap.String("method", r.Method), zap.String("path", r.URL.Path), zap.Error(err))
		} else {
			lg.Info("http", zap.String("method", r.Method), zap.String("path", r.URL.Path), zap.Int("status", http.StatusOK), zap.Int("bytes", cnt))
		}
	})
}
