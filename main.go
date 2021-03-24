// Copyright 2016 The OPA Authors.  All rights reserved.
// Use of this source code is governed by an Apache2
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"os"

	"github.com/open-policy-agent/opa/cmd"
    "github.com/open-policy-agent/opa/ast"
    "github.com/open-policy-agent/opa/rego"
    "github.com/open-policy-agent/opa/types"


    "github.com/mb0/wkt"
    "github.com/umahmood/haversine"
)

func main() {
    rego.RegisterBuiltin2(
        &rego.Function{
            Name: "distance",
            Decl: types.NewFunction(types.Args(types.S, types.S), types.N),
            Memoize: true,
        },
        func(bctx rego.BuiltinContext, a, b *ast.Term) (*ast.Term, error) {
            var as, bs string

            if err := ast.As(a.Value, &as); err != nil {
                return nil, err
            } else if ast.As(b.Value, &bs); err != nil {
                return nil, err
            }

            // var ap, bp wkt.Geom

            ag, err := wkt.Parse([]byte(as));
            if err != nil {
                return nil, err
            }

            bg, err := wkt.Parse([]byte(bs));
            if err != nil {
                return nil, err
            }

            ap := ag.(*wkt.Point)

            bp := bg.(*wkt.Point)

            // var ac, bc haversine.Coord

            ac := haversine.Coord{Lat: ap.Coord.Y, Lon: ap.Coord.X}
            bc := haversine.Coord{Lat: bp.Coord.Y, Lon: bp.Coord.X}
            _, km := haversine.Distance(ac, bc)

            return ast.FloatNumberTerm(km), nil
        },
    )

	if err := cmd.RootCommand.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

// Capabilities file generation:
//go:generate build/gen-run-go.sh internal/cmd/genopacapabilities/main.go capabilities.json

// WASM base binary generation:
//go:generate build/gen-run-go.sh internal/cmd/genopawasm/main.go -o internal/compiler/wasm/opa/opa.go internal/compiler/wasm/opa/opa.wasm  internal/compiler/wasm/opa/callgraph.csv
