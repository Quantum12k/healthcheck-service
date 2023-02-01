package main

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/Quantum12k/healthcheck-service/internal/healthcheck"
	"github.com/Quantum12k/healthcheck-service/internal/http_client"
	"github.com/Quantum12k/healthcheck-service/internal/postgresql"
)

const (
	notifyURL = "http://httpbin.org/post"
)

func (a *App) withDB(ctx context.Context) error {
	db, err := postgresql.New(a.Cfg.PostgreSQL)
	if err != nil {
		return fmt.Errorf("get postgreSQL instance: %v", err)
	}

	// забираем последние статусы по переданным url из БД
	lastChecks, err := db.GetLastHealthCheckEntries(a.Cfg.URLsToSlice())
	if err != nil {
		return fmt.Errorf("get last checks: %v", err)
	}

	for _, check := range lastChecks {
		a.cache.LastChecks.Add(check.URL, check.Result)
	}

	// получаем актуальные статусы
	results := healthcheck.HandleURLs(ctx, a.Cfg.URLs)

	if err := db.CreateHealthCheckEntries(results); err != nil {
		return fmt.Errorf("create healthcheck entry in DB: %v", err)
	}

	a.Log.Debugf("created: %v", results)

	// оповещаем при изменении статуса
	for _, res := range results {
		status, ok := a.cache.LastChecks.Get(res.URL)
		if !ok || (status != res.Result) {
			body := bytes.NewReader([]byte(res.String()))

			if _, err := http_client.DoRequest(ctx, http.MethodPost, notifyURL, body, time.Second); err != nil {
				return fmt.Errorf("do notification post request: %v", err)
			}
		}
	}

	return nil
}
