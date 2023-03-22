package db

import (
	"context"
	"go.uber.org/zap"
	"time"

	"mackey/pkg/logger"

	"github.com/go-pg/pg/v10"
)

type QueryHook struct {
	log *logger.Logger
}

func NewQueryHook(log *logger.Logger) QueryHook {
	return QueryHook{
		log: log.WithOptions(zap.AddCallerSkip(8)),
	}
}

func (q QueryHook) BeforeQuery(ctx context.Context, event *pg.QueryEvent) (context.Context, error) {
	event.Stash = make(map[interface{}]interface{})
	event.Stash["startedAt"] = time.Now()

	return ctx, nil
}

func (q QueryHook) AfterQuery(ctx context.Context, event *pg.QueryEvent) error {
	query, err := event.FormattedQuery()
	if err != nil {
		q.log.Warn(ctx, "formatted query error", zap.Error(err))
		return nil
	}

	var fields []zap.Field

	if event.Stash != nil {
		if v, ok := event.Stash["startedAt"]; ok {
			if startAt, ok := v.(time.Time); ok {
				fields = append(fields, zap.Duration("duration", time.Since(startAt)))
			}
		}
	}

	if event.Err == nil {
		fields = append(fields,
			zap.Int("affected", event.Result.RowsAffected()),
			zap.Int("returned", event.Result.RowsReturned()),
		)
	}

	q.log.Info(ctx, string(query), fields...)
	return nil
}
