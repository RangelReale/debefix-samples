package main

import (
	"bytes"
	"database/sql"
	"fmt"
	"time"

	"github.com/goccy/go-yaml"
)

type TaggedString struct {
	Tag   string
	Value string
}

func (t TaggedString) MarshalYAML() ([]byte, error) {
	v, err := yaml.Marshal(t.Value)
	if err != nil {
		return nil, err
	}
	return []byte(fmt.Sprintf("%s %s", t.Tag, string(v))), nil
}

type TaggedValue struct {
	Tag   string
	Value any
}

func (t TaggedValue) MarshalYAML() ([]byte, error) {
	var out bytes.Buffer
	// _, _ = fmt.Fprintf(&out, "\n%s\n", t.Tag)
	_, _ = fmt.Fprintf(&out, "%s\n", t.Tag)
	v, err := yaml.ValueToNode(t.Value, yaml.Flow(false))
	if err != nil {
		return nil, err
	}
	_, _ = fmt.Fprintf(&out, "%s", v)
	return out.Bytes(), nil

	// node := ast.Tag(token.New("", "", &token.Position{}))
	// node.Start.Value = t.Tag
	// var err error
	// node.Value, err = yaml.ValueToNode(t.Value, yaml.Flow(false))
	// if err != nil {
	// 	return nil, err
	// }
	// return []byte(node.String()), nil
}

func rowToMap(row *sql.Rows) (map[string]any, error) {
	cols, err := row.Columns()
	if err != nil {
		return nil, err
	}

	// Create a slice of interface{}'s to represent each column,
	// and a second slice to contain pointers to each item in the columns slice.
	columns := make([]interface{}, len(cols))
	columnPointers := make([]interface{}, len(cols))
	for i, _ := range columns {
		columnPointers[i] = &columns[i]
	}

	// Scan the result into the column pointers...
	if err := row.Scan(columnPointers...); err != nil {
		return nil, err
	}

	// Create our map, and retrieve the value for each column from the pointers slice,
	// storing it in the map with the name of the column as the key.
	m := make(map[string]interface{})
	for i, colName := range cols {
		val := columnPointers[i].(*interface{})
		colVal := *val
		switch v := colVal.(type) {
		case time.Time:
			colVal = TaggedString{
				Tag:   "!!timestamp",
				Value: v.Format(time.RFC3339Nano),
			}
		}
		m[colName] = colVal
	}

	return m, nil
}
