package main

import (
	"fmt"
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
    rows:
      - tag_id: 559
        tenant_id: 987
        tag_name: "javascript"
        _metadata:
          !metadata
          filedata: "best"
        tagfilename:
          !copyfile
          id: tag_image
          setValue: true
          src: "images/tags/javascript.png"
          dest: "tenant/{valueref:tenant_id:tenants:tenant_id:name}/images/tags/{value:tag_id}.png"
`),
		},
	})

	_, loadOptions, resolveOptions := copyfile.NewOptions(
		copyfile.WithCallback(func(ctx debefix.ValueResolveContext, fieldname string, fileData copyfile.FileData) error {
			p := copyfile.Parse(fileData.Dest)
			rmap := map[string]string{}
			for _, fld := range p.Fields() {
				rmap[fld] = fld
			}
			replaceValues, err := ctx.ResolvedData().ExtractValues(ctx.Row(), rmap)
			if err != nil {
				return err
			}

			dest, err := copyfile.Replace(fileData.Dest, replaceValues)
			if err != nil {
				return err
			}
			fmt.Printf("$$ [%s] COPY FILE [%s] FROM '%s' to '%s' [%s]\n", fieldname, fileData.ID, fileData.Src, dest, fileData.Dest)
			return nil
		}),
		copyfile.WithSetValueCallback(func(ctx debefix.ValueCallbackResolveContext, fileData copyfile.FileData) (resolvedValue any, addField bool, err error) {
			switch ctx.Table().ID {
			case "tags":
				return fmt.Sprintf("file-%s.png", ctx.Row().Fields["tag_name"]), true, nil
			default:
				return nil, false, nil
			}
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
