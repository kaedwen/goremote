package task

import (
	"bytes"
	"context"
	"html/template"
	"os/exec"

	"github.com/kaedwen/goremote/pkg/common"
)

type TaskExecutor struct {
}

func NewTaskExecutor() *TaskExecutor {
	return &TaskExecutor{}
}

func (e *TaskExecutor) Execute(ctx context.Context, t *common.TaskDefinition, args ...string) ([]byte, error) {
	if t.Script != nil {
		return e.executeScript(ctx, t, args...)
	}

	return e.executeCommand(ctx, t, args...)
}

func (e *TaskExecutor) executeCommand(ctx context.Context, t *common.TaskDefinition, args ...string) ([]byte, error) {
	td, err := doTemplate(data(args...), t.Arguments...)
	if err != nil {
		return nil, err
	}

	cmd := exec.CommandContext(ctx, t.Command, td...)
	return cmd.CombinedOutput()
}

func (e *TaskExecutor) executeScript(ctx context.Context, t *common.TaskDefinition, args ...string) ([]byte, error) {
	td, err := doTemplate(data(args...), *t.Script)
	if err != nil {
		return nil, err
	}

	cmd := exec.CommandContext(ctx, "sh")
	cmd.Stdin = bytes.NewReader([]byte(td[0]))
	return cmd.CombinedOutput()
}

func doTemplate(d any, t ...string) ([]string, error) {
	tt := make([]string, 0, len(t))
	for _, ti := range t {
		tmpl, err := template.New("").Parse(ti)
		if err != nil {
			return nil, err
		}

		var buf bytes.Buffer
		if err = tmpl.Execute(&buf, d); err != nil {
			return nil, err
		}

		tt = append(tt, buf.String())
	}

	return tt, nil
}

func data(args ...string) map[string]any {
	return map[string]any{
		"args": args,
	}
}
