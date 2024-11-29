package main

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/davecgh/go-spew/spew"
	"github.com/rrgmc/debefix-value/v2/copyfile"
	"github.com/rrgmc/debefix/v2"
)

func main() {
	err := run()
	if err != nil {
		panic(err)
	}
}

func run() error {
	ctx := context.Background()
	data := debefix.NewData()

	tableTenants := debefix.TableName("tenants")
	tableTags := debefix.TableName("tags")

	tenantIID := data.AddWithID(tableTenants,
		debefix.MapValues{
			"tenant_id": 987,
			"name":      "Joomla",
		})

	data.AddDependencies(tableTags, tableTenants)

	data.AddValues(tableTags, debefix.MapValues{
		"tag_id":    559,
		"tenant_id": tenantIID.ValueForField("tenant_id"),
		"tag_name":  "javascript",
		"_copyfile": copyfile.New(nil,
			copyfile.Filename("images/tags/javascript.png"),
			copyfile.FilenameFormatTemplate("tenant/{{.tenant_name}}/images/tags/{{.tag_id}}.png", map[string]any{
				"tenant_name": tenantIID.ValueForField("name"),
				"tag_id":      debefix.ValueFieldValue("tag_id"),
			})),
	})

	// _, loadOptions, resolveOptions := copyfile.NewOptions(
	// 	copyfile.WithSourcePath("/tmp/source"),
	// 	copyfile.WithDestinationPath("/tmp/destination"),
	// 	copyfile.WithGetPathsCallback(func(ctx debefix.ValueResolveContext, fieldname string, fileData copyfile.FileData) (source string, destination string, err error) {
	// 		return copyfile.DefaultGetPathsCallback(ctx, fieldname, fileData)
	// 	}),
	// 	copyfile.WithGetValueCallback(func(ctx debefix.ValueCallbackResolveContext, fileData copyfile.FileData) (value any, addField bool, err error) {
	// 		return copyfile.DefaultGetValueCallback(ctx, fileData)
	// 		// switch ctx.Table().ID {
	// 		// case "tags":
	// 		// 	return fmt.Sprintf("file-%s.png", ctx.Row().Fields["tag_name"]), true, nil
	// 		// default:
	// 		// 	return copyfile.DefaultGetValueCallback(ctx, fileData)
	// 		// }
	// 	}),
	// 	copyfile.WithCopyFileCallback(func(sourcePath, sourceFilename string, destinationPath, destinationFilename string) error {
	// 		fmt.Printf("Copying file %s to %s\n", filepath.Join(sourcePath, sourceFilename),
	// 			filepath.Join(destinationPath, destinationFilename))
	// 		return nil
	// 	}),
	// )

	resolvedData, err := debefix.Resolve(ctx, data,
		func(ctx context.Context, resolveInfo debefix.ResolveInfo, values debefix.ValuesMutable) error {
			return nil
		},
		debefix.WithResolveOptionProcess(copyfile.NewProcess(
			copyfile.WithProcessFilenameProvider(func(ctx context.Context, fileField copyfile.FileField,
				item copyfile.Value, tableID debefix.TableID, filename string) (string, error) {
				switch fileField {
				case copyfile.FileFieldSource:
					return filepath.Join("/tmp/source", filename), nil
				case copyfile.FileFieldDestination:
					return filepath.Join("/tmp/destination", filename), nil
				default:
					return "", fmt.Errorf("unknown FileField %v", fileField)
				}
			}),
			copyfile.WithProcessResolveCallback(func(ctx context.Context, resolvedData *debefix.ResolvedData,
				tableID debefix.TableID, fieldName string, values debefix.ValuesMutable, item copyfile.Value,
				reader copyfile.FileReader, writer copyfile.FileWriter) error {
				fmt.Printf("Copying file %s to %s\n", reader, writer)
				return nil
			}),
		)),
	)
	if err != nil {
		return err
	}

	spew.Dump(resolvedData)

	return nil
}
