package xatago

// Table represents a database table
type Table[T any] interface {
	Select(columns ...string) *Query[T]
}

type TableRecord interface {
	ID() string
}

// TableImpl implements all the Table functions for a given generic type
type TableImpl[T any] struct {
	client    *Client
	tableName string
}

// NewTableImpl creates an implementation for the given table
func NewTableImpl[T any](baseClient *Client, tableName string) *TableImpl[T] {
	return &TableImpl[T]{
		client:    baseClient,
		tableName: tableName,
	}
}

func (ti *TableImpl[T]) Select(columns ...string) *Query[T] {
	q := NewQuery[T](ti.client, ti.tableName)
	return q.Select(columns)
}
