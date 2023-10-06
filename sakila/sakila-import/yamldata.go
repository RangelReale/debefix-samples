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
