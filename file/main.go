package main

import (
	"fmt"
	"testing/fstest"

	"github.com/davecgh/go-spew/spew"
	"github.com/goccy/go-yaml"
	"github.com/goccy/go-yaml/ast"
	"github.com/rrgmc/debefix"
)

func main() {
	err := run()
	if err != nil {
		panic(err)
	}
}

func run() error {
	provider := debefix.NewFSFileProvider(fstest.MapFS{
		"users.dbf.yaml": &fstest.MapFile{
			Data: []byte(`tables:
  tags:
    rows:
      - tag_id: 1
        tag_name: "javascript"
        _metadata:
          !metadata
          filedata: "best"
        _file:
          !file
          src: "images/tags/javascript.png"
          dest: "company/{companyID}/images/tags/{tagID}.png"
`),
		},
	})

	data, err := debefix.Load(provider, debefix.WithLoadValueParser(&parseFile{}))
	if err != nil {
		return err
	}

	resolvedData, err := debefix.Resolve(data, func(ctx debefix.ResolveContext, fields map[string]any) error {
		return nil
	}, debefix.WithRowResolvedCallback(debefix.RowResolvedCallbackFunc(func(ctx debefix.ValueResolveContext) {
		spew.Dump(ctx.Row().Metadata)
	})))
	if err != nil {
		return err
	}

	fmt.Println(resolvedData.Tables["tags"].Rows[0].Fields["_file"])
	return nil
}

type FileData struct {
	Src  string `yaml:"src"`
	Dest string `yaml:"dest"`
}

type parseFile struct {
}

func (t *parseFile) ParseValue(tag *ast.TagNode) (bool, any, error) {
	if tag.Start.Value != "!file" {
		return false, nil, nil
	}

	var fc FileData
	err := yaml.NodeToValue(tag.Value, &fc)
	if err != nil {
		return false, nil, err
	}

	return true, debefix.ValueCallbackFunc(func(ctx debefix.ValueCallbackResolveContext) (any, bool, error) {
		ctx.AddMetadata("src", 123)
		return fmt.Sprintf("v=%s -- %s", fc.Src, fc.Dest), true, nil
	}), nil
}
