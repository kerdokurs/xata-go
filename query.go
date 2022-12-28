package xatago

type Query[T any] struct {
	client    *Client
	tableName string
	columns   []string
	limit     int
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

type apiQuery struct {
	Page    apiPage  `json:"page"`
	Columns []string `json:"columns"`
}

func (q *Query[T]) Select(columns []string) *Query[T] {
	q.columns = columns
	return q
}

func (q *Query[T]) Filter() *Query[T] {
	return q
}

func (q *Query[T]) GetFirst() (*T, error) {
	query := apiQuery{
		Page: apiPage{
			Size: 1,
		},
		Columns: q.columns,
	}

	ts, err := q.doQuery(&query)
	if err != nil {
		return nil, err
	}

	return ts[0], nil
}

func (q *Query[T]) GetMany() ([]*T, error) {
	query := apiQuery{
		Page: apiPage{
			Size: q.limit,
		},
		Columns: q.columns,
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
