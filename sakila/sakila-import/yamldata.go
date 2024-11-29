package main

type Data struct {
	Tables map[string]*Table `yaml:"tables"`
}

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
	RefID string `yaml:"refid,omitempty"`
}
