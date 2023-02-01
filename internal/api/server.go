package api

import (
	"fmt"
	"net/http"

	"go.uber.org/zap"

	"github.com/Quantum12k/healthcheck-service/internal/app_cache"
)

type (
	Config struct {
		Port string `yaml:"port"`
	}

	Server struct {
		log      *zap.SugaredLogger
		appCache *app_cache.Cache
	}
)

func NewServer(cfg Config, logger *zap.SugaredLogger, cache *app_cache.Cache) (*Server, error) {
	srv := &Server{
		log:      logger,
		appCache: cache,
	}

	go func() {
		if err := http.ListenAndServe(fmt.Sprintf(":%s", cfg.Port), srv.getRouter()); err != nil {
			logger.Errorf("listen and serve at port '%s': %v", cfg.Port, err)
		}
	}()

	return srv, nil
}
