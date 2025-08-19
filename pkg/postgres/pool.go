package postgres

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kotyari-bots_backend/pkg/utils"
	"time"
)

func GetPool(ctx context.Context, config Config) (*pgxpool.Pool, error) {
	poolCfg, err := pgxpool.ParseConfig(ToUrl(config))
	if err != nil {
		return nil, fmt.Errorf("parse config error: %w", err)
	}

	poolCfg.MaxConns = utils.GetValueOrDefault(config.MaxConns, defaultMaxConns)
	poolCfg.MinConns = utils.GetValueOrDefault(config.MinConns, defaultMinConns)
	poolCfg.MinIdleConns = utils.GetValueOrDefault(config.MinIdleConns, defautlMinIdleConns)
	poolCfg.HealthCheckPeriod = utils.GetValueOrDefault(config.HealthCheckPeriod, defaultHealthCheckPeriod)
	poolCfg.MaxConnIdleTime = utils.GetValueOrDefault(config.MaxConnIdleTime, defaultMaxConnIdleTime)
	poolCfg.MaxConnLifetime = utils.GetValueOrDefault(config.MaxConnLifetime, defaultMaxConnLifetime)

	pool, err := pgxpool.NewWithConfig(ctx, poolCfg)
	if err != nil {
		return nil, fmt.Errorf("unable to create connection pool: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err = TestPing(ctx, pool)
	if err != nil {
		return nil, fmt.Errorf("get pool failed: %w", err)
	}

	return pool, nil
}
