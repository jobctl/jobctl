// Copyright (C) 2024 The Dagu Authors
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program. If not, see <https://www.gnu.org/licenses/>.

package client

import (
	"context"
	"errors"
	"fmt"
	"math"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"github.com/dagu-org/dagu/internal/dag"
	"github.com/dagu-org/dagu/internal/dag/scheduler"
	"github.com/dagu-org/dagu/internal/frontend/gen/restapi/operations/dags"
	"github.com/dagu-org/dagu/internal/logger"
	"github.com/dagu-org/dagu/internal/persistence"
	"github.com/dagu-org/dagu/internal/persistence/history"
	"github.com/dagu-org/dagu/internal/persistence/model"
	"github.com/dagu-org/dagu/internal/sock"
)

var _ Client = (*client)(nil)

// New creates a new Client instance.
// The Client is used to interact with the DAG.
func New(
	dataStore persistence.ClientFactory,
	executable string,
	workDir string,
	lg logger.Logger,
) Client {
	return &client{
		dataStore:  dataStore,
		executable: executable,
		workDir:    workDir,
		logger:     lg,
	}
}

type client struct {
	dataStore  persistence.ClientFactory
	executable string
	workDir    string
	logger     logger.Logger
}

var (
	dagTemplate = []byte(`steps:
  - name: step1
    command: echo hello
`)
)

var (
	errCreateDAGFile = errors.New("failed to create DAG file")
	errGetStatus     = errors.New("failed to get status")
	errDAGIsRunning  = errors.New("the DAG is running")
)

func (e *client) GetDAGSpec(_ context.Context, id string) (string, error) {
	dagStore := e.dataStore.DAGStore()
	return dagStore.GetSpec(id)
}

func (e *client) CreateDAG(_ context.Context, name string) (string, error) {
	dagStore := e.dataStore.DAGStore()
	id, err := dagStore.Create(name, dagTemplate)
	if err != nil {
		return "", fmt.Errorf("%w: %s", errCreateDAGFile, err)
	}
	return id, nil
}

func (e *client) GrepDAGs(_ context.Context, pattern string) (
	[]*persistence.GrepResult, []string, error,
) {
	dagStore := e.dataStore.DAGStore()
	return dagStore.Grep(pattern)
}

func (e *client) Rename(ctx context.Context, oldID, newID string) error {
	dagStore := e.dataStore.DAGStore()
	oldDAG, err := dagStore.Find(oldID)
	if err != nil {
		return err
	}
	if err := dagStore.Rename(oldID, newID); err != nil {
		return err
	}
	newDAG, err := dagStore.Find(newID)
	if err != nil {
		return err
	}
	historyStore := e.dataStore.HistoryStore()
	return historyStore.RenameDAG(ctx, oldDAG.Location, newDAG.Location)
}

func (e *client) Stop(_ context.Context, dAG *dag.DAG) error {
	// TODO: fix this not to connect to the DAG directly
	client := sock.NewClient(dAG.SockAddr())
	_, err := client.Request("POST", "/stop")
	return err
}

func (e *client) StartAsync(ctx context.Context, dAG *dag.DAG, opts StartOptions) {
	go func() {
		if err := e.Start(ctx, dAG, opts); err != nil {
			e.logger.Error("Workflow start operation failed", "error", err)
		}
	}()
}

func (e *client) Start(_ context.Context, dAG *dag.DAG, opts StartOptions) error {
	args := []string{"start"}
	if opts.Params != "" {
		args = append(args, "-p")
		args = append(args, fmt.Sprintf(`"%s"`, escapeArg(opts.Params)))
	}
	if opts.Quiet {
		args = append(args, "-q")
	}
	args = append(args, dAG.Location)
	// nolint:gosec
	cmd := exec.Command(e.executable, args...)
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true, Pgid: 0}
	cmd.Dir = e.workDir
	cmd.Env = os.Environ()
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Start()
	if err != nil {
		return err
	}
	return cmd.Wait()
}

func (e *client) Restart(_ context.Context, dAG *dag.DAG, opts RestartOptions) error {
	args := []string{"restart"}
	if opts.Quiet {
		args = append(args, "-q")
	}
	args = append(args, dAG.Location)
	// nolint:gosec
	cmd := exec.Command(e.executable, args...)
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true, Pgid: 0}
	cmd.Dir = e.workDir
	cmd.Env = os.Environ()
	err := cmd.Start()
	if err != nil {
		return err
	}
	return cmd.Wait()
}

func (e *client) Retry(_ context.Context, dAG *dag.DAG, requestID string) error {
	args := []string{"retry"}
	args = append(args, fmt.Sprintf("--req=%s", requestID))
	args = append(args, dAG.Location)
	// nolint:gosec
	cmd := exec.Command(e.executable, args...)
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true, Pgid: 0}
	cmd.Dir = e.workDir
	cmd.Env = os.Environ()
	err := cmd.Start()
	if err != nil {
		return err
	}
	return cmd.Wait()
}

func (*client) GetCurrentStatus(_ context.Context, dAG *dag.DAG) (*model.Status, error) {
	client := sock.NewClient(dAG.SockAddr())
	ret, err := client.Request("GET", "/status")
	if err != nil {
		if errors.Is(err, sock.ErrTimeout) {
			return nil, err
		}
		return model.NewStatusDefault(dAG), nil
	}
	return model.StatusFromJSON(ret)
}

func (e *client) GetStatusByRequestID(ctx context.Context, dAG *dag.DAG, requestID string) (
	*model.Status, error,
) {
	ret, err := e.dataStore.HistoryStore().GetStatusByRequestID(
		ctx, dAG.Location, requestID,
	)
	if err != nil {
		return nil, err
	}
	status, _ := e.GetCurrentStatus(ctx, dAG)
	if status != nil && status.RequestID != requestID {
		// if the request id is not matched then correct the status
		ret.Status.CorrectRunningStatus()
	}
	return ret.Status, err
}

func (*client) currentStatus(dAG *dag.DAG) (*model.Status, error) {
	client := sock.NewClient(dAG.SockAddr())
	ret, err := client.Request("GET", "/status")
	if err != nil {
		return nil, fmt.Errorf("%w: %s", errGetStatus, err)
	}
	return model.StatusFromJSON(ret)
}

func (e *client) GetLatestStatus(ctx context.Context, dAG *dag.DAG) (*model.Status, error) {
	currStatus, _ := e.currentStatus(dAG)
	if currStatus != nil {
		return currStatus, nil
	}
	status, err := e.dataStore.HistoryStore().GetLatestStatus(ctx, dAG.Location)
	if errors.Is(err, history.ErrNoStatusDataToday) ||
		errors.Is(err, history.ErrNoStatusData) {
		return model.NewStatusDefault(dAG), nil
	}
	if err != nil {
		return model.NewStatusDefault(dAG), err
	}
	status.CorrectRunningStatus()
	return status, nil
}

func (e *client) ListRecentHistory(ctx context.Context, dAG *dag.DAG, n int) []*model.History {
	return e.dataStore.HistoryStore().ListRecentStatuses(ctx, dAG.Location, n)
}

func (e *client) UpdateStatus(ctx context.Context, dAG *dag.DAG, status *model.Status) error {
	client := sock.NewClient(dAG.SockAddr())
	res, err := client.Request("GET", "/status")
	if err != nil {
		if errors.Is(err, sock.ErrTimeout) {
			return err
		}
	} else {
		unmarshalled, _ := model.StatusFromJSON(res)
		if unmarshalled != nil && unmarshalled.RequestID == status.RequestID &&
			unmarshalled.Status == scheduler.StatusRunning {
			return errDAGIsRunning
		}
	}
	return e.dataStore.HistoryStore().UpdateStatus(
		ctx, dAG.Location, status.RequestID, status,
	)
}

func (e *client) UpdateDAGSpec(_ context.Context, id string, spec string) error {
	dagStore := e.dataStore.DAGStore()
	return dagStore.UpdateSpec(id, []byte(spec))
}

func (e *client) DeleteDAG(ctx context.Context, name, loc string) error {
	err := e.dataStore.HistoryStore().DeleteAllStatuses(ctx, loc)
	if err != nil {
		return err
	}
	dagStore := e.dataStore.DAGStore()
	return dagStore.Delete(name)
}

func (e *client) ListDAGStatusObsolete(ctx context.Context) (
	statuses []*DAGStatus, errs []string, err error,
) {
	dagStore := e.dataStore.DAGStore()
	dagList, errs, err := dagStore.List()

	var ret []*DAGStatus
	for _, d := range dagList {
		status, err := e.readStatus(ctx, d)
		if err != nil {
			errs = append(errs, err.Error())
		}
		ret = append(ret, status)
	}

	return ret, errs, err
}

func (e *client) getPageCount(total int64, limit int64) int {
	pageCount := int(math.Ceil(float64(total) / float64(limit)))
	if pageCount == 0 {
		pageCount = 1
	}

	return pageCount
}

func (e *client) ListDAGStatus(ctx context.Context, params dags.ListDagsParams) ([]*DAGStatus, *DagListPaginationSummaryResult, error) {
	var (
		dagListPaginationResult *persistence.DagListPaginationResult
		err                     error
		dagStore                = e.dataStore.DAGStore()
		dagStatusList           = make([]*DAGStatus, 0)
		currentStatus           *DAGStatus
	)

	if dagListPaginationResult, err = dagStore.ListPagination(persistence.DAGListPaginationArgs{
		Page:  int(params.Page),
		Limit: int(params.Limit),
		Name:  params.SearchName,
		Tag:   params.SearchTag,
	}); err != nil {
		return dagStatusList, &DagListPaginationSummaryResult{PageCount: 1}, err
	}

	for _, currentDag := range dagListPaginationResult.DagList {
		if currentStatus, err = e.readStatus(ctx, currentDag); err != nil {
			dagListPaginationResult.ErrorList = append(dagListPaginationResult.ErrorList, err.Error())
		}
		dagStatusList = append(dagStatusList, currentStatus)
	}

	return dagStatusList, &DagListPaginationSummaryResult{
		PageCount: e.getPageCount(int64(dagListPaginationResult.Count), params.Limit),
		ErrorList: dagListPaginationResult.ErrorList,
	}, nil
}

func (e *client) ListHistoryByDate(ctx context.Context, date string) ([]*model.History, error) {
	d, err := time.Parse("2006-01-02", date)
	if err != nil {
		return nil, err
	}

	historyStore := e.dataStore.HistoryStore()
	return historyStore.ListStatusesByDate(ctx, d)
}

func (e *client) getDAG(name string) (*dag.DAG, error) {
	dagStore := e.dataStore.DAGStore()
	dagDetail, err := dagStore.GetDetails(name)
	return e.emptyDAGIfNil(dagDetail, name), err
}

func (e *client) GetLatestDAGStatus(ctx context.Context, id string) (*DAGStatus, error) {
	dg, err := e.getDAG(id)
	if dg == nil {
		// TODO: fix not to use location
		dg = &dag.DAG{Name: id, Location: id}
	}
	if err == nil {
		// check the dag is correct in terms of graph
		_, err = scheduler.NewExecutionGraph(e.logger, dg.Steps...)
	}
	latestStatus, _ := e.GetLatestStatus(ctx, dg)
	return newDAGStatus(
		dg, latestStatus, e.IsSuspended(ctx, id), err,
	), err
}

func (e *client) ToggleSuspend(_ context.Context, id string, suspend bool) error {
	flagStore := e.dataStore.FlagStore()
	return flagStore.ToggleSuspend(id, suspend)
}

func (e *client) IsSuspended(_ context.Context, id string) bool {
	flagStore := e.dataStore.FlagStore()
	return flagStore.IsSuspended(id)
}

func (e *client) ListTags(_ context.Context) ([]string, []string, error) {
	return e.dataStore.DAGStore().TagList()
}

func (e *client) readStatus(ctx context.Context, dAG *dag.DAG) (*DAGStatus, error) {
	latestStatus, err := e.GetLatestStatus(ctx, dAG)
	id := strings.TrimSuffix(
		filepath.Base(dAG.Location),
		filepath.Ext(dAG.Location),
	)

	return newDAGStatus(
		dAG, latestStatus, e.IsSuspended(ctx, id), err,
	), err
}

func (*client) emptyDAGIfNil(dAG *dag.DAG, dagLocation string) *dag.DAG {
	if dAG != nil {
		return dAG
	}
	return &dag.DAG{Location: dagLocation}
}

func escapeArg(input string) string {
	escaped := strings.Builder{}

	for _, char := range input {
		if char == '\r' {
			_, _ = escaped.WriteString("\\r")
		} else if char == '\n' {
			_, _ = escaped.WriteString("\\n")
		} else {
			_, _ = escaped.WriteRune(char)
		}
	}

	return escaped.String()
}
