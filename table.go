package xatago

// Table represents a database table
type Table[T Record] interface {
	Select(columns ...string) *Query[T]
	Filter(key string, filterKey FilterKey, value any) *Query[T]

	GetMany() ([]*T, error)
	GetFirst() (*T, error)

	Create(item *T) (*T, error)
	Update(item *T) (*T, error)
	Delete(item *T) error
	DeleteById(id string) error
}

// TableImpl implements all the Table functions for a given generic type
type TableImpl[T Record] struct {
	client    *Client
	tableName string
}

// NewTableImpl creates an implementation for the given table
func NewTableImpl[T Record](baseClient *Client, tableName string) *TableImpl[T] {
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

func (ti *TableImpl[T]) GetMany() ([]*T, error) {
	q := NewQuery[T](ti.client, ti.tableName)
	return q.GetMany()
}

func (ti *TableImpl[T]) GetFirst() (*T, error) {
	q := NewQuery[T](ti.client, ti.tableName)
	return q.GetFirst()
}

func (ti *TableImpl[T]) Create(item *T) (*T, error) {
	t := new(T)
	return t, ti.client.create(ti.tableName, item, t)
}

func (ti *TableImpl[T]) Update(item *T) (*T, error) {
	t := new(T)
	return t, ti.client.update(ti.tableName, (*item).ID(), item, t)
}

func (ti *TableImpl[T]) Delete(item *T) error {
	return ti.DeleteById((*item).ID())
}

func (ti *TableImpl[T]) DeleteById(id string) error {
	return ti.client.delete(ti.tableName, id)
}
