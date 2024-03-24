package impl

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/bendahl/uinput"
	"github.com/kaedwen/goremote/pkg/api/v1/gen"
	"github.com/kaedwen/goremote/pkg/common"
	"github.com/kaedwen/goremote/pkg/task"
	"github.com/kaedwen/goremote/pkg/utils"
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

	err = keyboard.KeyPress(int(req.Code))
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
			Message: utils.Ptr(fmt.Sprintf("task with id '%s' not found", req.Id)),
		}, nil
	}

	te := task.NewTaskExecutor()
	td, err := te.Execute(ctx, t, req.Args...)
	if err != nil {
		s.lg.Error("failed to execute task", zap.String("result", string(td)), zap.Error(err))
		return &gen.ExecTaskResponse{
			Success: false,
			Message: utils.Ptr(fmt.Sprintf("failed to execute task - %v", err)),
		}, nil
	}

	var result []string
	if err := json.Unmarshal(td, &result); err != nil {
		s.lg.Info("failed to parse result", zap.Error(err))
	}

	return &gen.ExecTaskResponse{
		Success: true,
		Results: result,
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

func (s *serveRemoteGRPC) MouseClick(ctx context.Context, req *gen.MouseClickRequest) (*gen.MouseClickResponse, error) {
	mouse, err := uinput.CreateMouse("/dev/uinput", []byte("goremote-mouse"))
	if err != nil {
		return &gen.MouseClickResponse{
			Success: false,
		}, nil
	}
	defer mouse.Close()

	switch req.Code {
	case 1:
		err = mouse.LeftClick()
	case 2:
		err = mouse.MiddleClick()
	case 3:
		err = mouse.RightClick()
	default:
		err = fmt.Errorf("code not found - %d", req.Code)
	}

	if err != nil {
		s.lg.Error("failed to click mouse", zap.Error(err))
		return &gen.MouseClickResponse{
			Success: false,
		}, nil
	}

	return &gen.MouseClickResponse{
		Success: true,
	}, nil
}

func (s *serveRemoteGRPC) MouseMove(ctx context.Context, req *gen.MouseMoveRequest) (*gen.MouseMoveResponse, error) {
	mouse, err := uinput.CreateMouse("/dev/uinput", []byte("goremote-mouse"))
	if err != nil {
		return &gen.MouseMoveResponse{
			Success: false,
		}, nil
	}
	defer mouse.Close()

	switch req.Direction {
	case 1:
		err = mouse.MoveUp(int32(req.Delta))
	case 2:
		err = mouse.MoveDown(int32(req.Delta))
	case 3:
		err = mouse.MoveLeft(int32(req.Delta))
	case 4:
		err = mouse.MoveRight(int32(req.Delta))
	}

	if err != nil {
		s.lg.Error("failed to move mouse", zap.Error(err))
		return &gen.MouseMoveResponse{
			Success: false,
		}, nil
	}

	return &gen.MouseMoveResponse{
		Success: true,
	}, nil
}

func (s *serveRemoteGRPC) MousePosition(ctx context.Context, req *gen.MousePositionRequest) (*gen.MousePositionResponse, error) {
	mouse, err := uinput.CreateMouse("/dev/uinput", []byte("goremote-mouse"))
	if err != nil {
		return &gen.MousePositionResponse{
			Success: false,
		}, nil
	}
	defer mouse.Close()

	// hack! first a huge amount up and left, ending up on 0x0 (also used by ydotool)
	// than move relative to the final position
	err = mouse.Move(-10000, -10000)
	if err != nil {
		s.lg.Error("failed to move mouse", zap.Error(err))
		return &gen.MousePositionResponse{
			Success: false,
		}, nil
	}

	err = mouse.Move(int32(req.X), int32(req.Y))
	if err != nil {
		s.lg.Error("failed to move mouse", zap.Error(err))
		return &gen.MousePositionResponse{
			Success: false,
		}, nil
	}

	return &gen.MousePositionResponse{
		Success: true,
	}, nil
}
