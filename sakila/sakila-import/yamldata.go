package main

type Data map[string]*Table

type Table struct {
	Config TableConfig `yaml:"config,omitempty"`
	Rows   []Row       `yaml:"rows,omitempty"`
}

type TableConfig struct {
	TableName string   `yaml:"table_name,omitempty"`
	Depends   []string `yaml:"depends,omitempty"`
}

type Row map[string]any

type RowConfig struct {
	RefID string   `yaml:"refid,omitempty"`
	Tags  []string `yaml:"tags,omitempty"`
}
