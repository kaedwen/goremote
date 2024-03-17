package impl

import (
	"context"
	"fmt"

	"github.com/bendahl/uinput"
	"github.com/kaedwen/goremote/pkg/api/v1/gen"
	"github.com/kaedwen/goremote/pkg/common"
	"github.com/kaedwen/goremote/pkg/task"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

// RegisterGRPC will be used by service descriptor
// creates and registers GRPC server instance
func RegisterGRPC(ctx context.Context, lg common.Logger, cfg *common.Config, reg grpc.ServiceRegistrar) {
	gen.RegisterRemoteServer(reg, &serveRemoteGRPC{lg: lg, cfg: cfg})
}

// serveRemoteGRPC implements the expected handler methods
// private implementation based on service-specific generated model, entry point to service logic
type serveRemoteGRPC struct {
	gen.UnimplementedRemoteServer
	lg  common.Logger
	cfg *common.Config
}

func (s *serveRemoteGRPC) PressKey(ctx context.Context, req *gen.KeyRequest) (*gen.KeyResponse, error) {
	keyboard, err := uinput.CreateKeyboard("/dev/uinput", []byte("goremote-keyboard"))
	if err != nil {
		return &gen.KeyResponse{
			Success: false,
		}, nil
	}
	defer keyboard.Close()

	err = keyboard.KeyPress(uinput.KeyPause)
	if err != nil {
		return &gen.KeyResponse{
			Success: false,
		}, nil
	}

	return &gen.KeyResponse{
		Success: true,
	}, nil
}

func (s *serveRemoteGRPC) ExecTask(ctx context.Context, req *gen.ExecTaskRequest) (*gen.ExecTaskResponse, error) {
	t := s.cfg.Tasks.Find(req.Id)
	if t == nil {
		return &gen.ExecTaskResponse{
			Success: false,
			Result:  fmt.Sprintf("task with id '%s' not found", req.Id),
		}, nil
	}

	te := task.NewTaskExecutor()
	td, err := te.Execute(ctx, t)
	if err != nil {
		s.lg.Error("failed to execute task", zap.String("result", string(td)), zap.Error(err))
		return &gen.ExecTaskResponse{
			Success: false,
			Result:  fmt.Sprintf("failed to execute task - %v", err),
		}, nil
	}

	return &gen.ExecTaskResponse{
		Success: true,
		Result:  string(td),
	}, nil
}

func (s *serveRemoteGRPC) ListTask(ctx context.Context, req *emptypb.Empty) (*gen.ListTaskResponse, error) {
	tasks := []*gen.Task{}
	for _, t := range s.cfg.Tasks {
		tasks = append(tasks, &gen.Task{
			Id:   t.Id,
			Name: &t.Name,
		})
	}

	return &gen.ListTaskResponse{
		Tasks: tasks,
	}, nil
}
