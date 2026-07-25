package cmd

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/briheet/nozarashi/examples/mix_docker_nix/backend/internal/api"
	"github.com/briheet/nozarashi/examples/mix_docker_nix/backend/internal/config"
	"github.com/briheet/nozarashi/examples/mix_docker_nix/backend/internal/db"
	"github.com/briheet/nozarashi/examples/mix_docker_nix/backend/internal/logger"
	"github.com/briheet/nozarashi/examples/mix_docker_nix/backend/internal/redis"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

func ApiCmd(ctx context.Context) *cobra.Command {

	var configPaths []string

	apiCmd := &cobra.Command{
		Use:   "api",
		Short: "GPT chat and inference ingestion API",
		RunE: func(cmd *cobra.Command, args []string) error {

			// Load config
			cfg, err := config.LoadConfig(ctx, configPaths...)
			if err != nil {
				return err
			}

			log, err := logger.NewLogger("api")
			if err != nil {
				return err
			}
			defer func() {
				if err := log.Sync(); err != nil {
					log.Error("sync logger", zap.Error(err))
				}
			}()

			dbClient, err := db.NewClient(ctx, cfg)
			if err != nil {
				return err
			}
			defer dbClient.Close()

			redisClient, err := redis.NewClient(ctx, cfg)
			if err != nil {
				return err
			}
			defer func() {
				if err := redisClient.Close(); err != nil {
					log.Error("close Redis client", zap.Error(err))
				}
			}()

			handler := api.NewAPI(cfg, log, dbClient, redisClient)
			srv := handler.Server(cfg.API.Port)

			apiErr := make(chan error, 1)
			go func() {
				apiErr <- srv.ListenAndServe()
			}()

			log.Info("started api", zap.Int("port", cfg.API.Port))

			select {
			case <-ctx.Done():
				shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
				defer cancel()

				if err := srv.Shutdown(shutdownCtx); err != nil && !errors.Is(err, http.ErrServerClosed) {
					return errors.Join(err, srv.Close())
				}
				return nil
			case err := <-apiErr:
				if errors.Is(err, http.ErrServerClosed) {
					return nil
				}
				return err
			}
		},
	}

	apiCmd.PersistentFlags().StringArrayVarP(&configPaths, "configPath", "c", nil, "configuration file path; later files override earlier files")

	return apiCmd
}
