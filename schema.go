package xatago

type Schema struct {
	Tables []SchemaTable
}

type SchemaTable struct {
	Name    string        `json:"name"`
	Columns []TableColumn `json:"columns"`
}

type TableColumn struct {
	Name         string    `json:"name"`
	Type         string    `json:"type"`
	Link         TableLink `json:"link,omitempty"`
	Unique       bool      `json:"unique"`
	NotNull      bool      `json:"notNull"`
	DefaultValue string    `json:"defaultValue"`
}

type TableLink struct {
	Table string `json:"table"`
}
