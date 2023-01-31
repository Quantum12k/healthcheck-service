package postgresql

import (
	"fmt"

	"github.com/doug-martin/goqu/v9"

	"github.com/Quantum12k/healthcheck-service/internal/healthcheck"
)

func (p *PostgreSQL) CreateHealthCheckEntries(entries []healthcheck.CheckResult) error {
	query, _, _ := goqu.Insert("healthcheck").Rows(entries).ToSQL()

	if _, err := p.db.Exec(query); err != nil {
		return fmt.Errorf("exec query: %v", err)
	}

	return nil
}

func (p *PostgreSQL) GetLastHealthCheckEntries(urls []string) ([]healthcheck.CheckResult, error) {
	result := make([]healthcheck.CheckResult, 0, len(urls))

	query, _, _ := goqu.From("healthcheck").Select("url", "result").Where(
			goqu.Ex{"url":  urls},
			goqu.C("time").Eq(goqu.From(goqu.T("healthcheck").As("h2")).Select(goqu.MAX("h2.time")).Where(goqu.L(`h2.url="url"`))),
		).Order(goqu.C("time").Desc()).ToSQL()

	if err := p.db.Select(&result, query); err != nil {
		return nil, fmt.Errorf("get results: %v", err)
	}

	return result, nil
}
