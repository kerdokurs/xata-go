package xatago

// Table represents a database table
type Table[T any] interface {
	GetAll() ([]*T, error)
	GetFirst() (*T, error)

	Create(t *T) (string, error)
	Delete(id string) error

	Filter(f Filter) ([]*T, error)
}

// TableImpl implements all the Table functions for a given generic type
type TableImpl[T any] struct {
	client    *Client
	tableName string
}

// Filter given fields by by FilterElement criterion
type Filter map[string]FilterElement

// FilterElement describes criterion to filter by
type FilterElement struct {
	Is string `json:"$is,omitempty"`
}

// Page describes info about a query page
type Page struct {
	Size   int    `json:"size,omitempty"`
	Offset int    `json:"offset,omitempty"`
	Before string `json:"before,omitempty"`
	After  string `json:"after,omitempty"`
	Cursor string `json:"cursor,omitempty"`
}

// Query describes a query
type Query struct {
	Columns []string `json:"columns,omitempty"`
	Page    Page     `json:"page,omitempty"`
	Filter  Filter   `json:"filter,omitempty"`
}

// GetAll retrieves all elements of type T in the table
func (pt *TableImpl[T]) GetAll() ([]*T, error) {
	query := Query{
		Page: Page{
			Size: 0,
		},
	}
	rawData, err := pt.client.query(pt.tableName, &query)
	if err != nil {
		return nil, err
	}
	data, err := Map[any, *T](rawData, MapToStruct[T])
	if err != nil {
		return nil, err
	}
	return data, nil
}

// GetFirst retrieves the first element in the table
func (pt *TableImpl[T]) GetFirst() (*T, error) {
	query := Query{
		Page: Page{
			Size: 1,
		},
	}
	rawData, err := pt.client.query(pt.tableName, &query)
	if err != nil {
		return nil, err
	}
	data, err := Map[any, *T](rawData, MapToStruct[T])
	if err != nil {
		return nil, err
	}
	if len(data) == 0 {
		return nil, nil
	}
	return data[0], nil
}

// Create creates an element of type T in the table
func (pt *TableImpl[T]) Create(t *T) (string, error) {
	return pt.client.create(pt.tableName, t)
}

// Delete deletes an object with the given id
func (pt *TableImpl[T]) Delete(id string) error {
	return pt.client.delete(pt.tableName, id)
}

// Filter filters objects by a given filter
func (pt *TableImpl[T]) Filter(f Filter) ([]*T, error) {
	query := Query{
		Filter: f,
	}
	rawData, err := pt.client.query(pt.tableName, &query)
	if err != nil {
		return nil, err
	}
	data, err := Map[any, *T](rawData, MapToStruct[T])
	if err != nil {
		return nil, err
	}
	return data, nil
}
