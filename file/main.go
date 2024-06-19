package main

import (
	"fmt"
	"path/filepath"
	"testing/fstest"

	"github.com/davecgh/go-spew/spew"
	"github.com/rrgmc/debefix"
	"github.com/rrgmc/debefix-copyfile"
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
  tenants:
    rows:
      - tenant_id: 987
        name: "Joomla"
  tags:
    config:
      depends: ["tenants"]
      default_values:
        tagfilename:
          !copyfile
          id: tag_image
          value: "{value:tag_id}.png"
          source: "images/tags/javascript.png"
          destination: "tenant/{valueref:tenant_id:tenants:tenant_id:name}/images/tags/{value:tag_id}.png"
    rows:
      - tag_id: 559
        tenant_id: 987
        tag_name: "javascript"
        _metadata:
          !metadata
          filedata: "best"
`),
		},
	})

	_, loadOptions, resolveOptions := copyfile.NewOptions(
		copyfile.WithSourcePath("/tmp/source"),
		copyfile.WithDestinationPath("/tmp/destination"),
		copyfile.WithGetPathsCallback(func(ctx debefix.ValueResolveContext, fieldname string, fileData copyfile.FileData) (source string, destination string, err error) {
			return copyfile.DefaultGetPathsCallback(ctx, fieldname, fileData)
		}),
		copyfile.WithGetValueCallback(func(ctx debefix.ValueCallbackResolveContext, fileData copyfile.FileData) (value any, addField bool, err error) {
			return copyfile.DefaultGetValueCallback(ctx, fileData)
			// switch ctx.Table().ID {
			// case "tags":
			// 	return fmt.Sprintf("file-%s.png", ctx.Row().Fields["tag_name"]), true, nil
			// default:
			// 	return copyfile.DefaultGetValueCallback(ctx, fileData)
			// }
		}),
		copyfile.WithCopyFileCallback(func(sourcePath, sourceFilename string, destinationPath, destinationFilename string) error {
			fmt.Printf("Copying file %s to %s\n", filepath.Join(sourcePath, sourceFilename),
				filepath.Join(destinationPath, destinationFilename))
			return nil
		}),
	)

	data, err := debefix.Load(provider, loadOptions...)
	if err != nil {
		return err
	}

	resolvedData, err := debefix.Resolve(data, func(ctx debefix.ResolveContext, fields map[string]any) error {
		return nil
	}, resolveOptions...)
	if err != nil {
		return err
	}

	spew.Dump(resolvedData)

	return nil
}
