package main

import (
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

	_, loadOptions, resolveOptions := copyfile.NewOptions()

	data, err := debefix.Load(provider, loadOptions...)
	if err != nil {
		return err
	}

	_, err = debefix.Resolve(data, func(ctx debefix.ResolveContext, fields map[string]any) error {
		return nil
	}, resolveOptions...)
	if err != nil {
		return err
	}

	return nil
}
