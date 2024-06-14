package main

import (
	"fmt"
	"testing/fstest"

	"github.com/rrgmc/debefix"
	"github.com/rrgmc/debefix-samples/file/copyfile"
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
        _tagimage:
          !copyfile
          src: "images/tags/javascript.png"
          dest: "company/{companyID}/images/tags/{tagID}.png"
`),
		},
	})

	copyFilePlugin := &copyfile.CopyFile{}

	data, err := debefix.Load(provider, debefix.WithLoadValueParser(copyFilePlugin))
	if err != nil {
		return err
	}

	resolvedData, err := debefix.Resolve(data, func(ctx debefix.ResolveContext, fields map[string]any) error {
		return nil
	}, debefix.WithRowResolvedCallback(copyFilePlugin))
	if err != nil {
		return err
	}

	fmt.Println(resolvedData.Tables["tags"].Rows[0].Fields["_file"])
	return nil
}
