// Copyright 2026 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package imports

import (
	"context"
	"strings"
	"testing"
)

type usedImportNamesSource struct {
	used map[PackageName]bool
}

func (s usedImportNamesSource) LoadPackageNames(context.Context, string, []ImportPath) (map[ImportPath]PackageName, error) {
	return nil, nil
}

func (s usedImportNamesSource) ResolveReferences(context.Context, string, References) ([]*Result, error) {
	return nil, nil
}

func (s usedImportNamesSource) UsedImportNames(context.Context, string) map[PackageName]bool {
	return s.used
}

func TestFixImportsKeepsExtensionImport(t *testing.T) {
	const src = `package main

import (
	"fmt"
	"linq"
)

func main() {
	_ = []int{1}.Where(n => n > 0)
	fmt.Println("ok")
}
`
	fixes, err := FixImports(context.Background(), "main.go", []byte(src), "", nil, usedImportNamesSource{
		used: map[PackageName]bool{"linq": true, "fmt": true},
	})
	if err != nil {
		t.Fatal(err)
	}
	for _, fix := range fixes {
		if fix.FixType == DeleteImport && fix.StmtInfo.ImportPath == "linq" {
			t.Fatalf("organize imports deleted linq: %#v", fix)
		}
	}
	out, err := ApplyFixes(fixes, "main.go", []byte(src), &Options{Comments: true}, 0)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(out), `"linq"`) {
		t.Fatalf("linq import missing after fix:\n%s", out)
	}
}
