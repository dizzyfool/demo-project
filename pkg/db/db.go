//lint:file-ignore U1000 ignore unused code, it's generated
//nolint:structcheck,unused
package db

import (
	"context"
	"fmt"
	"reflect"

	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"
)

type Config struct {
	Host     string `env:"POSTGRES_HOST,default=localhost"`
	Port     uint16 `env:"POSTGRES_PORT,default=5432"`
	Database string `env:"POSTGRES_DB,required"`
	User     string `env:"POSTGRES_USER,required"`
	Pass     string `env:"POSTGRES_PASSWORD"`
	ShowSQL  bool   `env:"POSTGRES_VERBOSE,default=true"`
}

func (conf Config) ToPG() *pg.Options {
	return &pg.Options{
		Addr:     fmt.Sprintf("%s:%d", conf.Host, conf.Port),
		User:     conf.User,
		Password: conf.Pass,
		Database: conf.Database,
	}
}

// DB stores db connection
type DB struct {
	*pg.DB
}

// New is a function that returns DB as wrapper on postgres connection.
func New(config Config, hook QueryHook) (DB, orm.DB) {
	db := pg.Connect(config.ToPG())

	if config.ShowSQL {
		db.AddQueryHook(hook)
	}

	return DB{
		DB: db,
	}, db
}

// Version is a function that returns Postgres version.
func (db *DB) Version() (string, error) {
	var v string
	if _, err := db.QueryOne(pg.Scan(&v), "select version()"); err != nil {
		return "", err
	}

	return v, nil
}

// runInTransaction runs chain of functions in transaction until first error
func (db *DB) runInTransaction(ctx context.Context, fns ...func(*pg.Tx) error) error {
	return db.RunInTransaction(ctx, func(tx *pg.Tx) error {
		for _, fn := range fns {
			if err := fn(tx); err != nil {
				return err
			}
		}
		return nil
	})
}

// buildQuery applies all functions to orm query.
func buildQuery(ctx context.Context, db orm.DB, model interface{}, search Searcher, filters []Filter, pager Pager, ops ...OpFunc) *orm.Query {
	q := db.ModelContext(ctx, model)
	for _, filter := range filters {
		filter.Apply(q)
	}

	if reflect.ValueOf(search).IsValid() && !reflect.ValueOf(search).IsNil() { // is it good?
		search.Apply(q)
	}

	q = pager.Apply(q)
	applyOps(q, ops...)

	return q
}
