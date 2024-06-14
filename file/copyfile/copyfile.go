package copyfile

import (
	"fmt"

	"github.com/goccy/go-yaml"
	"github.com/goccy/go-yaml/ast"
	"github.com/rrgmc/debefix"
)

type CopyFile struct {
	debefix.ValueImpl
}

var (
	_ debefix.ValueParser         = (*CopyFile)(nil)
	_ debefix.RowResolvedCallback = (*CopyFile)(nil)
)

func (c *CopyFile) ParseValue(tag *ast.TagNode) (bool, any, error) {
	if tag.Start.Value != "!copyfile" {
		return false, nil, nil
	}

	var fileData FileData
	err := yaml.NodeToValue(tag.Value, &fileData)
	if err != nil {
		return false, nil, err
	}

	return true, &copyFileValue{fileData: fileData}, nil
}

func (c *CopyFile) RowResolved(ctx debefix.ValueResolveContext) {
	md := getMetadata(ctx.Row().Metadata)
	for fieldName, file := range md.Fields {
		fmt.Printf("$$ [%s] COPY FILE FROM '%s' to '%s'\n", fieldName, file.Src, file.Dest)
	}

	// spew.Dump(ctx.Row().Metadata)
}

type copyFileValue struct {
	debefix.ValueImpl
	fileData FileData
}

var (
	_ debefix.ValueCallback = (*copyFileValue)(nil)
)

func (c *copyFileValue) GetValueCallback(ctx debefix.ValueCallbackResolveContext) (resolvedValue any, addField bool, err error) {
	setMetadata(ctx, c.fileData)
	// don't add a data field
	return nil, false, nil
}

const (
	metadataName = "__copyfile__"
)

func getMetadata(metadata map[string]any) *FileDataList {
	if md, ok := metadata[metadataName]; ok {
		if mdfl, ok := md.(*FileDataList); ok {
			return mdfl
		}
	}
	return &FileDataList{
		Fields: map[string]FileData{},
	}
}

func setMetadata(ctx debefix.ValueCallbackResolveContext, fileData FileData) {
	md := getMetadata(ctx.Metadata())
	md.Fields[ctx.FieldName()] = fileData
	ctx.SetMetadata(metadataName, md)
}
