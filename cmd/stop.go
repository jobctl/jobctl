package cmd

import (
	"log"
	"os"

	"github.com/dagu-dev/dagu/internal/config"
	"github.com/dagu-dev/dagu/internal/dag"
	"github.com/dagu-dev/dagu/internal/logger"
	"github.com/spf13/cobra"
)

func stopCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "stop <DAG file>",
		Short: "Stop the running DAG",
		Long:  `dagu stop <DAG file>`,
		Args:  cobra.ExactArgs(1),
		Run: func(_ *cobra.Command, args []string) {
			cfg, err := config.Load()
			if err != nil {
				log.Fatalf("Failed to load config: %v", err)
			}
			logger := logger.NewLogger(logger.NewLoggerArgs{
				Config: cfg,
			})

			loadedDAG, err := dag.Load(cfg.BaseConfig, args[0], "")
			if err != nil {
				logger.Error("Failed to load DAG", "error", err)
				os.Exit(1)
			}

			logger.Info("Stopping the DAG", "dag", loadedDAG.Name)

			eng := newEngine(cfg)

			if err := eng.Stop(loadedDAG); err != nil {
				logger.Error("Failed to stop the DAG", "error", err)
				os.Exit(1)
			}
		},
	}
}
