package copyfile

import (
	"github.com/davecgh/go-spew/spew"
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

	return true, &copyFileValue{}, nil
}

func (c *CopyFile) RowResolved(ctx debefix.ValueResolveContext) {
	spew.Dump(ctx.Row().Metadata)
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

func getMetadata(ctx debefix.ValueCallbackResolveContext) *FileDataList {
	if md, ok := ctx.Metadata()[metadataName]; ok {
		if mdfl, ok := md.(*FileDataList); ok {
			return mdfl
		}
	}
	return &FileDataList{
		Fields: map[string]FileData{},
	}
}

func setMetadata(ctx debefix.ValueCallbackResolveContext, fileData FileData) {
	md := getMetadata(ctx)
	md.Fields[metadataName] = fileData
	ctx.SetMetadata(metadataName, md)
}
