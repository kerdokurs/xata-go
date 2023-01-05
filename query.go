package xatago

import "errors"

type Query[T any] struct {
	client    *Client
	tableName string
	columns   []string
	limit     int
	filter    *apiFilter
}

func NewQuery[T any](client *Client, tableName string) *Query[T] {
	return &Query[T]{
		client:    client,
		tableName: tableName,
		columns:   []string{"*"},
		limit:     0,
	}
}

type apiPage struct {
	Size int `json:"size"`
}

type apiFilter struct {
	All []map[string]any `json:"$all"`
}

type apiQuery struct {
	Page    apiPage    `json:"page,omitempty"`
	Columns []string   `json:"columns,omitempty"`
	Filter  *apiFilter `json:"filter,omitempty"`
}

func (q *Query[T]) Select(columns []string) *Query[T] {
	q.columns = columns
	return q
}

func (q *Query[T]) Filter(key string, filterKey FilterKey, value any) *Query[T] {
	if q.filter == nil {
		q.filter = new(apiFilter)
	}

	f := make(map[string]any)
	switch filterKey {
	case Is:
		f[key] = value
	default:
		{
			inner := make(map[string]any)
			inner[string(filterKey)] = value
			f[key] = inner
		}
	}

	q.filter.All = append(q.filter.All, f)

	return q
}

func (q *Query[T]) GetFirst() (*T, error) {
	query := apiQuery{
		Page: apiPage{
			Size: 1,
		},
		Columns: q.columns,
		Filter:  q.filter,
	}

	ts, err := q.doQuery(&query)
	if err != nil {
		return nil, err
	}

	if len(ts) == 0 {
		return nil, errors.New("record not found")
	}

	return ts[0], nil
}

func (q *Query[T]) GetMany() ([]*T, error) {
	query := apiQuery{
		Page: apiPage{
			Size: q.limit,
		},
		Columns: q.columns,
		Filter:  q.filter,
	}

	return q.doQuery(&query)
}

func (q *Query[T]) doQuery(query *apiQuery) ([]*T, error) {
	rawData, err := q.client.query(q.tableName, query)
	if err != nil {
		return nil, err
	}

	return MapToStructs[any, *T](rawData)
}
