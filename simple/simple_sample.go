package main

import (
	"context"
	"fmt"
	"strings"

	"github.com/davecgh/go-spew/spew"
	"github.com/google/uuid"
	"github.com/rrgmc/debefix-db/v2/sql"
	"github.com/rrgmc/debefix-db/v2/sql/postgres"
	data2 "github.com/rrgmc/debefix-samples/v2/data"
	"github.com/rrgmc/debefix/v2"
)

func main() {
	ctx := context.Background()

	data := data2.Data()
	if data.Err() != nil {
		panic(data.Err())
	}

	// spew.Dump(data)

	err := debefix.ResolveCheck(ctx, data)
	if err != nil {
		panic(err)
	}

	// err = resolvePrint(ctx, data)
	// if err != nil {
	// 	panic(err)
	// }

	err = resolveSQL(ctx, data)
	if err != nil {
		panic(err)
	}
}

func resolvePrint(ctx context.Context, data *debefix.Data) error {
	_, err := debefix.Resolve(ctx, data, func(ctx context.Context, resolveInfo debefix.ResolveInfo, values debefix.ValuesMutable) error {
		fmt.Printf("%s %s %s\n", strings.Repeat("=", 10), resolveInfo.TableID.TableName(), strings.Repeat("=", 10))
		spew.Dump(values)

		resolved := map[string]any{}
		for fn, fv := range values.All {
			if fresolve, ok := fv.(debefix.ResolveValue); ok {
				rv, err := fresolve.ResolveValueParse(ctx, uuid.New())
				if err != nil {
					return fmt.Errorf("error parsing resolve value '%s': %w", fn, err)
				}
				values.Set(fn, rv)
				resolved[fn] = rv
			}
		}

		if len(resolved) > 0 {
			fmt.Println("---")
			spew.Dump(resolved)
		}

		return nil
	})
	return err
}

func resolveSQL(ctx context.Context, data *debefix.Data) error {
	_, err := debefix.Resolve(ctx, data, postgres.ResolveFunc(sql.NewDebugQueryInterface(nil)))
	return err
}
