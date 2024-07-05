package executor

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"sync"
	"syscall"

	"github.com/dagu-dev/dagu/internal/util"

	"github.com/dagu-dev/dagu/internal/dag"
)

type CommandExecutor struct {
	cmd  *exec.Cmd
	lock sync.Mutex
}

func NewCommandExecutor(ctx context.Context, step dag.Step) (Executor, error) {
	// nolint: gosec
	cmd := exec.CommandContext(ctx, step.Command, step.Args...)
	if len(step.Dir) > 0 && !util.FileExists(step.Dir) {
		return nil, fmt.Errorf("directory %q does not exist", step.Dir)
	}
	cmd.Dir = step.Dir
	cmd.Env = append(cmd.Env, os.Environ()...)
	cmd.Env = append(cmd.Env, step.Variables...)
	step.OutputVariables.Range(func(_, value any) bool {
		cmd.Env = append(cmd.Env, value.(string))
		return true
	})
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Setpgid: true,
		Pgid:    0,
	}

	return &CommandExecutor{
		cmd: cmd,
	}, nil
}

func (e *CommandExecutor) Run() error {
	e.lock.Lock()
	err := e.cmd.Start()
	e.lock.Unlock()
	if err != nil {
		return err
	}
	return e.cmd.Wait()
}

func (e *CommandExecutor) SetStdout(out io.Writer) {
	e.cmd.Stdout = out
}

func (e *CommandExecutor) SetStderr(out io.Writer) {
	e.cmd.Stderr = out
}

func (e *CommandExecutor) Kill(sig os.Signal) error {
	e.lock.Lock()
	defer e.lock.Unlock()
	if e.cmd == nil || e.cmd.Process == nil {
		return nil
	}
	return syscall.Kill(-e.cmd.Process.Pid, sig.(syscall.Signal))
}

func init() {
	Register("", NewCommandExecutor)
	Register("command", NewCommandExecutor)
}
