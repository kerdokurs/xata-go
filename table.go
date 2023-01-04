package xatago

// Table represents a database table
type Table[T any] interface {
	Select(columns ...string) *Query[T]
	Filter(key string, filterKey FilterKey, value any) *Query[T]

	Create(item *T) (string, error)
	Delete(id string) error
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

func (ti *TableImpl[T]) Filter(key string, filterKey FilterKey, value any) *Query[T] {
	q := NewQuery[T](ti.client, ti.tableName)
	return q.Filter(key, filterKey, value)
}

func (ti *TableImpl[T]) Create(item *T) (string, error) {
	return ti.client.create(ti.tableName, item)
}

func (ti *TableImpl[T]) Delete(id string) error {
	return ti.client.delete(ti.tableName, id)
}
