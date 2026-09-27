package bitget

import (
	"go/ast"
	"go/parser"
	"go/token"
	"io/fs"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
	"testing"
)

// TestTimeFieldsDeclareFormat requires every exported time.Time / *time.Time
// struct field that takes part in JSON to declare its wire format with the
// `format` tag option (e.g. `json:"cTime,format:unixmilli"`). Without one,
// encoding/json/v2 falls back to RFC 3339, which no Bitget timestamp uses, so a
// missing tag would only surface as a decode error against the live API.
func TestTimeFieldsDeclareFormat(t *testing.T) {
	fset := token.NewFileSet()
	checked := 0
	err := filepath.WalkDir(".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			if name := d.Name(); path != "." && (strings.HasPrefix(name, ".") || name == "testdata") {
				return filepath.SkipDir
			}
			return nil
		}
		if !strings.HasSuffix(path, ".go") || strings.HasSuffix(path, "_test.go") {
			return nil
		}
		f, err := parser.ParseFile(fset, path, nil, parser.SkipObjectResolution)
		if err != nil {
			return err
		}
		ast.Inspect(f, func(n ast.Node) bool {
			st, ok := n.(*ast.StructType)
			if !ok {
				return true
			}
			for _, field := range st.Fields.List {
				if !isTimeType(field.Type) || len(field.Names) == 0 || !field.Names[0].IsExported() {
					continue
				}
				var tag string
				if field.Tag != nil {
					raw, _ := strconv.Unquote(field.Tag.Value)
					tag = reflect.StructTag(raw).Get("json")
				}
				if tag == "-" {
					continue
				}
				checked++
				opts := strings.Split(tag, ",")[1:]
				if len(opts) == 0 || !strings.HasPrefix(opts[len(opts)-1], "format:") {
					t.Errorf("%s: field %s has json tag %q without a trailing format option", fset.Position(field.Pos()), field.Names[0].Name, tag)
				}
			}
			return true
		})
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	if checked == 0 {
		t.Fatal("no time.Time fields found; is the test running from the module root?")
	}
}

func isTimeType(e ast.Expr) bool {
	if star, ok := e.(*ast.StarExpr); ok {
		e = star.X
	}
	sel, ok := e.(*ast.SelectorExpr)
	if !ok {
		return false
	}
	pkg, ok := sel.X.(*ast.Ident)
	return ok && pkg.Name == "time" && sel.Sel.Name == "Time"
}
