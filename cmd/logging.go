package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/dagu-dev/dagu/internal/dag"
	"github.com/dagu-dev/dagu/internal/util"
)

// openLogFileForDAG opens a log file for the DAG.
func openLogFileForDAG(prefix string, logDir string, dg *dag.DAG, reqID string) (*os.File, error) {
	name := util.ValidFilename(dg.Name)
	if dg.LogDir != "" {
		logDir = filepath.Join(dg.LogDir, name)
	}
	// Check if the log directory exists
	if _, err := os.Stat(logDir); os.IsNotExist(err) {
		// Create the log directory
		if err := os.MkdirAll(logDir, 0755); err != nil {
			return nil, err
		}
	}
	file := filepath.Join(logDir, fmt.Sprintf("%s%s.%s.%s.log",
		prefix,
		name,
		time.Now().Format("20060102.15:04:05.000"),
		util.TruncString(reqID, 8),
	))
	// Open or create the log file
	return os.OpenFile(
		file, os.O_CREATE|os.O_WRONLY|os.O_APPEND|os.O_SYNC, 0644,
	)
}
