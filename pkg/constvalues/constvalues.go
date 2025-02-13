package constvalues

import (
	"go/ast"
	"reflect"
	"sort"
	"strings"

	"golang.org/x/tools/go/packages"
)

var cache = make(map[string]map[string][]string)

func FreeCache() {
	if cache != nil {
		cache = nil
	}
}

func GetConstValuesByType(t any) []string {
	if t == nil {
		return nil
	}

	typ := getRootType(t)

	if !checkKind(typ) {
		return nil
	}

	var (
		packagPath = typ.PkgPath()
		typeName   = typ.Name()
	)

	if cache == nil {
		return nil
	}

	if _, ok := cache[packagPath][typeName]; !ok {
		pkg := loadPackage(packagPath)
		cache[packagPath] = getConstValueList(pkg, typeName)
		sort.Strings(cache[packagPath][typeName])
	}
	return cache[packagPath][typeName]
}

func checkKind(typ reflect.Type) bool {
	switch typ.Kind() {
	case reflect.String, reflect.Bool, reflect.Int64, reflect.Float64:
		return true
	default:
		return false
	}
}

func getRootType(t any) reflect.Type {
	typ := reflect.TypeOf(t)
	for typ.Kind() == reflect.Pointer {
		typ = typ.Elem()
	}
	return typ
}

func loadPackage(packagePath string) *packages.Package {
	pkg, err := packages.Load(&packages.Config{
		Mode: packages.NeedTypesInfo | packages.NeedSyntax,
	}, packagePath)
	if err != nil {
		return nil
	}
	return pkg[0]
}

func getConstValueList(pkg *packages.Package, typeName string) map[string][]string {
	var values = make(map[string][]string)
	for _, file := range pkg.Syntax {
		for _, decl := range file.Decls {
			if v, ok := decl.(*ast.GenDecl); ok {
				for _, spec := range v.Specs {
					switch vSpec := spec.(type) {
					case *ast.ValueSpec:
						switch vType := vSpec.Type.(type) {
						case *ast.Ident:
							if vType.Name != typeName {
								continue
							}
						default:
							continue
						}
						for _, value := range vSpec.Values {
							switch vExpr := value.(type) {
							case *ast.BasicLit:
								values[typeName] = append(values[typeName], strings.Trim(vExpr.Value, "\""))
							default:
							}
						}
					default:
					}
				}
			}
		}
	}
	return values
}
