package db

import (
	"context"
	"errors"

	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"
)

type DemoRepo struct {
	db      orm.DB
	filters map[string][]Filter
	sort    map[string][]SortField
	join    map[string][]string
}

// NewDemoRepo returns new repository
func NewDemoRepo(db orm.DB) DemoRepo {
	return DemoRepo{
		db:      db,
		filters: map[string][]Filter{},
		sort: map[string][]SortField{
			Tables.Message.Name: {{Column: Columns.Message.ID, Direction: SortDesc}},
		},
		join: map[string][]string{
			Tables.Message.Name: {TableColumns},
		},
	}
}

// WithTransaction is a function that wraps DemoRepo with pg.Tx transaction.
func (dr DemoRepo) WithTransaction(tx *pg.Tx) DemoRepo {
	dr.db = tx
	return dr
}

// WithEnabledOnly is a function that adds "statusId"=1 as base filter.
func (dr DemoRepo) WithEnabledOnly() DemoRepo {
	f := make(map[string][]Filter, len(dr.filters))
	for i := range dr.filters {
		f[i] = make([]Filter, len(dr.filters[i]))
		copy(f[i], dr.filters[i])
		f[i] = append(f[i], StatusEnabledFilter)
	}
	dr.filters = f

	return dr
}

/*** Message ***/

// FullMessage returns full joins with all columns
func (dr DemoRepo) FullMessage() OpFunc {
	return WithColumns(dr.join[Tables.Message.Name]...)
}

// DefaultMessageSort returns default sort.
func (dr DemoRepo) DefaultMessageSort() OpFunc {
	return WithSort(dr.sort[Tables.Message.Name]...)
}

// MessageByID is a function that returns Message by ID(s) or nil.
func (dr DemoRepo) MessageByID(ctx context.Context, id int, ops ...OpFunc) (*Message, error) {
	return dr.OneMessage(ctx, &MessageSearch{ID: &id}, ops...)
}

// OneMessage is a function that returns one Message by filters. It could return pg.ErrMultiRows.
func (dr DemoRepo) OneMessage(ctx context.Context, search *MessageSearch, ops ...OpFunc) (*Message, error) {
	obj := &Message{}
	err := buildQuery(ctx, dr.db, obj, search, dr.filters[Tables.Message.Name], PagerTwo, ops...).Select()

	if errors.Is(err, pg.ErrMultiRows) {
		return nil, err
	} else if errors.Is(err, pg.ErrNoRows) {
		return nil, nil
	}

	return obj, err
}

// MessagesByFilters returns Message list.
func (dr DemoRepo) MessagesByFilters(ctx context.Context, search *MessageSearch, pager Pager, ops ...OpFunc) (messages []Message, err error) {
	err = buildQuery(ctx, dr.db, &messages, search, dr.filters[Tables.Message.Name], pager, ops...).Select()
	return
}

// CountMessages returns count
func (dr DemoRepo) CountMessages(ctx context.Context, search *MessageSearch, ops ...OpFunc) (int, error) {
	return buildQuery(ctx, dr.db, &Message{}, search, dr.filters[Tables.Message.Name], PagerOne, ops...).Count()
}

// AddMessage adds Message to DB.
func (dr DemoRepo) AddMessage(ctx context.Context, message *Message, ops ...OpFunc) (*Message, error) {
	q := dr.db.ModelContext(ctx, message)
	applyOps(q, ops...)
	_, err := q.Insert()

	return message, err
}

// UpdateMessage updates Message in DB.
func (dr DemoRepo) UpdateMessage(ctx context.Context, message *Message, ops ...OpFunc) (bool, error) {
	q := dr.db.ModelContext(ctx, message).WherePK()
	applyOps(q, ops...)
	res, err := q.Update()
	if err != nil {
		return false, err
	}

	return res.RowsAffected() > 0, err
}

// DeleteMessage deletes Message from DB.
func (dr DemoRepo) DeleteMessage(ctx context.Context, id int) (deleted bool, err error) {
	message := &Message{ID: id}

	res, err := dr.db.ModelContext(ctx, message).WherePK().Delete()
	if err != nil {
		return false, err
	}

	return res.RowsAffected() > 0, err
}
