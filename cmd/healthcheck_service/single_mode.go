package main

import (
	"context"

	"github.com/Quantum12k/healthcheck-service/internal/healthcheck"
)

func (a *App) single(ctx context.Context) error {
	results := healthcheck.HandleURLs(ctx, a.Cfg.URLs)

	for _, res := range results {
		a.Log.Info(res.String())
	}

	return nil
}
