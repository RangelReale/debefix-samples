package main

import (
	"context"
	"errors"
	"fmt"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/davecgh/go-spew/spew"
	"github.com/google/uuid"
	"github.com/rrgmc/debefix"
	"github.com/rrgmc/debefix/db/sql"
	"github.com/rrgmc/debefix/db/sql/postgres"
)

func currentSourceDirectory() (string, error) {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		return "", errors.New("unable to get the current filename")
	}
	return filepath.Dir(filename), nil
}

func main() {
	curDir, err := currentSourceDirectory()
	if err != nil {
		panic(err)
	}

	data, err := debefix.Load(debefix.NewDirectoryFileProvider(filepath.Join(curDir, "data"),
		debefix.WithDirectoryAsTag()))
	if err != nil {
		panic(err)
	}

	// spew.Dump(data)

	resolveTags := []string{}

	err = debefix.ResolveCheck(data, debefix.WithResolveTags(resolveTags))
	if err != nil {
		panic(err)
	}

	// err = resolvePrint(data, resolveTags)
	// if err != nil {
	// 	panic(err)
	// }

	err = resolveSQL(context.Background(), data, resolveTags)
	if err != nil {
		panic(err)
	}
}

func resolvePrint(data *debefix.Data, resolveTags []string) error {
	_, err := debefix.Resolve(data, func(ctx debefix.ResolveContext, fields map[string]any) error {
		fmt.Printf("%s %s %s\n", strings.Repeat("=", 10), ctx.TableName(), strings.Repeat("=", 10))
		spew.Dump(fields)

		resolved := map[string]any{}
		for fn, fv := range fields {
			if fresolve, ok := fv.(debefix.ResolveValue); ok {
				switch fresolve.(type) {
				case *debefix.ResolveGenerate:
					ctx.ResolveField(fn, uuid.New())
					resolved[fn] = uuid.New()
				}
			}
		}

		if len(resolved) > 0 {
			fmt.Println("---")
			spew.Dump(resolved)
		}

		return nil
	}, debefix.WithResolveTags(resolveTags))
	return err
}

func resolveSQL(ctx context.Context, data *debefix.Data, resolveTags []string) error {
	_, err := postgres.Resolve(ctx, sql.NewDebugQueryInterface(nil), data, debefix.WithResolveTags(resolveTags))
	return err
}
