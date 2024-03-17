package task

import (
	"bytes"
	"context"
	"os/exec"

	"github.com/kaedwen/goremote/pkg/common"
)

type TaskExecutor struct {
}

func NewTaskExecutor() *TaskExecutor {
	return &TaskExecutor{}
}

func (e *TaskExecutor) Execute(ctx context.Context, t *common.TaskDefinition) ([]byte, error) {
	if t.Script != nil {
		return e.executeScript(ctx, t)
	}

	return e.executeCommand(ctx, t)
}

func (e *TaskExecutor) executeCommand(ctx context.Context, t *common.TaskDefinition) ([]byte, error) {
	cmd := exec.CommandContext(ctx, t.Command, t.Arguments...)
	return cmd.CombinedOutput()
}

func (e *TaskExecutor) executeScript(ctx context.Context, t *common.TaskDefinition) ([]byte, error) {
	cmd := exec.CommandContext(ctx, "sh")
	cmd.Stdin = bytes.NewReader([]byte(*t.Script))
	return cmd.CombinedOutput()
}
