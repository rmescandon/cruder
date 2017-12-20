// -*- Mode: Go; indent-tabs-mode: t -*-

/*
 * Copyright (C) 2017 Roberto Mier Escandon <rmescandon@gmail.com>
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License version 3 as
 * published by the Free Software Foundation.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 *
 */

package cruder

// This file contains the model construction by parsing source files.

import (
	"bytes"
	"errors"
	"fmt"
	"go/ast"
	"go/build"
	"go/parser"
	"go/token"
	"io"
	"log"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/golang/mock/mockgen/model"
)

// GetTypesMaps returns a map of types
func GetTypesMaps(filepath string) (map[string]map[string]string, error) {
	emptyMap := make(map[string]map[string]string)
	f, err := os.Open(filepath)
	if err != nil {
		return emptyMap, err
	}

	metafile, err := parse(f)
	if err != nil {
		return emptyMap, err
	}

	structsList, err := getStructs(metafile)
	if err != nil {
		return emptyMap, err
	}

	return decomposeStructs(structsList)
}

// Parse parses io reader to ast.File pointer
func parse(reader io.Reader) (*ast.File, error) {
	if reader == nil {
		return nil, errors.New("Reader is null")
	}

	var buf bytes.Buffer
	_, err := io.Copy(&buf, reader)
	if err != nil {
		return nil, err
	}

	buffer := buf.Bytes()

	fs := token.NewFileSet()
	return parser.ParseFile(fs, "", buffer, parser.Trace)
}

func getStructs(file *ast.File) (map[string]*ast.StructType, error) {
	structs := make(map[string]*ast.StructType)

	for _, decl := range file.Decls {
		switch decl.(type) {
		case *ast.GenDecl:
			genDecl := decl.(*ast.GenDecl)
			for _, spec := range genDecl.Specs {
				switch spec.(type) {
				case *ast.TypeSpec:
					typeSpec := spec.(*ast.TypeSpec)
					switch typeSpec.Type.(type) {
					case *ast.StructType:
						structName := typeSpec.Name.Name
						structType := typeSpec.Type.(*ast.StructType)
						structs[structName] = structType
					}
				}
			}
		}
	}

	return structs, nil
}

func decomposeStructs(structs map[string]*ast.StructType) (map[string]map[string]string, error) {
	structsMap := make(map[string]map[string]string)
	for structName := range structs {
		structMembers, err := decomposeStruct(structs[structName])
		if err != nil {
			return structsMap, err
		}
		structsMap[structName] = structMembers
	}
	return structsMap, nil
}

func decomposeStruct(structType *ast.StructType) (map[string]string, error) {
	fields := make(map[string]string)
	for _, field := range structType.Fields.List {
		if len(field.Names) == 2 {
			fields[field.Names[0].Name] = field.Names[1].Name
		}
	}
	return fields, nil
}

type fileParser struct {
	fileSet            *token.FileSet
	imports            map[string]string                        // package name => import path
	importedInterfaces map[string]map[string]*ast.InterfaceType // package (or "") => name => interface

	auxFiles      []*ast.File
	auxInterfaces map[string]map[string]*ast.InterfaceType // package (or "") => name => interface

	srcDir string
}

// ParseFile returns a file parsed to memory, ready to build the skeletom code
func ParseFile(source string) (*model.Package, error) {
	srcDir, err := filepath.Abs(filepath.Dir(source))
	if err != nil {
		return nil, fmt.Errorf("failed getting source directory: %v", err)
	}
	fs := token.NewFileSet()
	file, err := parser.ParseFile(fs, source, nil, 0)
	if err != nil {
		return nil, fmt.Errorf("failed parsing source file %v: %v", source, err)
	}

	p := &fileParser{
		fileSet:            fs,
		imports:            make(map[string]string),
		importedInterfaces: make(map[string]map[string]*ast.InterfaceType),
		auxInterfaces:      make(map[string]map[string]*ast.InterfaceType),
		srcDir:             srcDir,
	}

	pkg, err := p.parseFile(file)
	if err != nil {
		return nil, err
	}
	return pkg, nil
}

func (p *fileParser) errorf(pos token.Pos, format string, args ...interface{}) error {
	ps := p.fileSet.Position(pos)
	format = "%s:%d:%d: " + format
	args = append([]interface{}{ps.Filename, ps.Line, ps.Column}, args...)
	return fmt.Errorf(format, args...)
}

func (p *fileParser) parseFile(file *ast.File) (*model.Package, error) {
	var is []*model.Interface
	for ni := range iterInterfaces(file) {
		i, err := p.parseInterface(ni.name.String(), "", ni.it)
		if err != nil {
			return nil, err
		}
		is = append(is, i)
	}
	return &model.Package{
		Name:       file.Name.String(),
		Interfaces: is,
	}, nil
}

func (p *fileParser) parsePackage(path string) error {
	var pkgs map[string]*ast.Package
	if imp, err := build.Import(path, p.srcDir, build.FindOnly); err != nil {
		return err
	} else if pkgs, err = parser.ParseDir(p.fileSet, imp.Dir, nil, 0); err != nil {
		return err
	}
	for _, pkg := range pkgs {
		file := ast.MergePackageFiles(pkg, ast.FilterFuncDuplicates|ast.FilterUnassociatedComments|ast.FilterImportDuplicates)
		if _, ok := p.importedInterfaces[path]; !ok {
			p.importedInterfaces[path] = make(map[string]*ast.InterfaceType)
		}
		for ni := range iterInterfaces(file) {
			p.importedInterfaces[path][ni.name.Name] = ni.it
		}
		for pkgName, pkgPath := range importsOfFile(file) {
			if _, ok := p.imports[pkgName]; !ok {
				p.imports[pkgName] = pkgPath
			}
		}
	}
	return nil
}

func (p *fileParser) parseInterface(name, pkg string, it *ast.InterfaceType) (*model.Interface, error) {
	intf := &model.Interface{Name: name}
	for _, field := range it.Methods.List {
		switch v := field.Type.(type) {
		case *ast.FuncType:
			if nn := len(field.Names); nn != 1 {
				return nil, fmt.Errorf("expected one name for interface %v, got %d", intf.Name, nn)
			}
			m := &model.Method{
				Name: field.Names[0].String(),
			}
			var err error
			m.In, m.Variadic, m.Out, err = p.parseFunc(pkg, v)
			if err != nil {
				return nil, err
			}
			intf.Methods = append(intf.Methods, m)
		case *ast.Ident:
			// Embedded interface in this package.
			ei := p.auxInterfaces[pkg][v.String()]
			if ei == nil {
				if ei = p.importedInterfaces[pkg][v.String()]; ei == nil {
					return nil, p.errorf(v.Pos(), "unknown embedded interface %s", v.String())
				}
			}
			eintf, err := p.parseInterface(v.String(), pkg, ei)
			if err != nil {
				return nil, err
			}
			// Copy the methods.
			// TODO: apply shadowing rules.
			for _, m := range eintf.Methods {
				intf.Methods = append(intf.Methods, m)
			}
		case *ast.SelectorExpr:
			// Embedded interface in another package.
			fpkg, sel := v.X.(*ast.Ident).String(), v.Sel.String()
			epkg, ok := p.imports[fpkg]
			if !ok {
				return nil, p.errorf(v.X.Pos(), "unknown package %s", fpkg)
			}
			ei := p.auxInterfaces[fpkg][sel]
			if ei == nil {
				fpkg = epkg
				if _, ok = p.importedInterfaces[epkg]; !ok {
					if err := p.parsePackage(epkg); err != nil {
						return nil, p.errorf(v.Pos(), "could not parse package %s: %v", fpkg, err)
					}
				}
				if ei = p.importedInterfaces[epkg][sel]; ei == nil {
					return nil, p.errorf(v.Pos(), "unknown embedded interface %s.%s", fpkg, sel)
				}
			}
			eintf, err := p.parseInterface(sel, fpkg, ei)
			if err != nil {
				return nil, err
			}
			// Copy the methods.
			// TODO: apply shadowing rules.
			for _, m := range eintf.Methods {
				intf.Methods = append(intf.Methods, m)
			}
		default:
			return nil, fmt.Errorf("don't know how to mock method of type %T", field.Type)
		}
	}
	return intf, nil
}

func (p *fileParser) parseFunc(pkg string, f *ast.FuncType) (in []*model.Parameter, variadic *model.Parameter, out []*model.Parameter, err error) {
	if f.Params != nil {
		regParams := f.Params.List
		if isVariadic(f) {
			n := len(regParams)
			varParams := regParams[n-1:]
			regParams = regParams[:n-1]
			vp, err := p.parseFieldList(pkg, varParams)
			if err != nil {
				return nil, nil, nil, p.errorf(varParams[0].Pos(), "failed parsing variadic argument: %v", err)
			}
			variadic = vp[0]
		}
		in, err = p.parseFieldList(pkg, regParams)
		if err != nil {
			return nil, nil, nil, p.errorf(f.Pos(), "failed parsing arguments: %v", err)
		}
	}
	if f.Results != nil {
		out, err = p.parseFieldList(pkg, f.Results.List)
		if err != nil {
			return nil, nil, nil, p.errorf(f.Pos(), "failed parsing returns: %v", err)
		}
	}
	return
}

func (p *fileParser) parseFieldList(pkg string, fields []*ast.Field) ([]*model.Parameter, error) {
	nf := 0
	for _, f := range fields {
		nn := len(f.Names)
		if nn == 0 {
			nn = 1 // anonymous parameter
		}
		nf += nn
	}
	if nf == 0 {
		return nil, nil
	}
	ps := make([]*model.Parameter, nf)
	i := 0 // destination index
	for _, f := range fields {
		t, err := p.parseType(pkg, f.Type)
		if err != nil {
			return nil, err
		}

		if len(f.Names) == 0 {
			// anonymous arg
			ps[i] = &model.Parameter{Type: t}
			i++
			continue
		}
		for _, name := range f.Names {
			ps[i] = &model.Parameter{Name: name.Name, Type: t}
			i++
		}
	}
	return ps, nil
}

func (p *fileParser) parseType(pkg string, typ ast.Expr) (model.Type, error) {
	switch v := typ.(type) {
	case *ast.ArrayType:
		ln := -1
		if v.Len != nil {
			x, err := strconv.Atoi(v.Len.(*ast.BasicLit).Value)
			if err != nil {
				return nil, p.errorf(v.Len.Pos(), "bad array size: %v", err)
			}
			ln = x
		}
		t, err := p.parseType(pkg, v.Elt)
		if err != nil {
			return nil, err
		}
		return &model.ArrayType{Len: ln, Type: t}, nil
	case *ast.ChanType:
		t, err := p.parseType(pkg, v.Value)
		if err != nil {
			return nil, err
		}
		var dir model.ChanDir
		if v.Dir == ast.SEND {
			dir = model.SendDir
		}
		if v.Dir == ast.RECV {
			dir = model.RecvDir
		}
		return &model.ChanType{Dir: dir, Type: t}, nil
	case *ast.Ellipsis:
		// assume we're parsing a variadic argument
		return p.parseType(pkg, v.Elt)
	case *ast.FuncType:
		in, variadic, out, err := p.parseFunc(pkg, v)
		if err != nil {
			return nil, err
		}
		return &model.FuncType{In: in, Out: out, Variadic: variadic}, nil
	case *ast.Ident:
		if v.IsExported() {
			// `pkg` may be an aliased imported pkg
			// if so, patch the import w/ the fully qualified import
			maybeImportedPkg, ok := p.imports[pkg]
			if ok {
				pkg = maybeImportedPkg
			}
			// assume type in this package
			return &model.NamedType{Package: pkg, Type: v.Name}, nil
		}
		// assume predeclared type
		return model.PredeclaredType(v.Name), nil
	case *ast.InterfaceType:
		if v.Methods != nil && len(v.Methods.List) > 0 {
			return nil, p.errorf(v.Pos(), "can't handle non-empty unnamed interface types")
		}
		return model.PredeclaredType("interface{}"), nil
	case *ast.MapType:
		key, err := p.parseType(pkg, v.Key)
		if err != nil {
			return nil, err
		}
		value, err := p.parseType(pkg, v.Value)
		if err != nil {
			return nil, err
		}
		return &model.MapType{Key: key, Value: value}, nil
	case *ast.SelectorExpr:
		pkgName := v.X.(*ast.Ident).String()
		pkg, ok := p.imports[pkgName]
		if !ok {
			return nil, p.errorf(v.Pos(), "unknown package %q", pkgName)
		}
		return &model.NamedType{Package: pkg, Type: v.Sel.String()}, nil
	case *ast.StarExpr:
		t, err := p.parseType(pkg, v.X)
		if err != nil {
			return nil, err
		}
		return &model.PointerType{Type: t}, nil
	case *ast.StructType:
		if v.Fields != nil && len(v.Fields.List) > 0 {
			return nil, p.errorf(v.Pos(), "can't handle non-empty unnamed struct types")
		}
		return model.PredeclaredType("struct{}"), nil
	}

	return nil, fmt.Errorf("don't know how to parse type %T", typ)
}

// importsOfFile returns a map of package name to import path
// of the imports in file.
func importsOfFile(file *ast.File) map[string]string {
	/* We have to make guesses about some imports, because imports are not required
	 * to have names. Named imports are always certain. Unnamed imports are guessed
	 * to have a name of the last path component; if the last path component has dots,
	 * the first dot-delimited field is used as the name.
	 */

	m := make(map[string]string)
	for _, is := range file.Imports {
		var pkg string
		importPath := is.Path.Value[1 : len(is.Path.Value)-1] // remove quotes

		if is.Name != nil {
			if is.Name.Name == "_" {
				continue
			}
			pkg = removeDot(is.Name.Name)
		} else {
			_, last := path.Split(importPath)
			pkg = strings.SplitN(last, ".", 2)[0]
		}
		if _, ok := m[pkg]; ok {
			log.Fatalf("imported package collision: %q imported twice", pkg)
		}
		m[pkg] = importPath
	}
	return m
}

type namedInterface struct {
	name *ast.Ident
	it   *ast.InterfaceType
}

// Create an iterator over all interfaces in file.
func iterInterfaces(file *ast.File) <-chan namedInterface {
	ch := make(chan namedInterface)
	go func() {
		for _, decl := range file.Decls {
			gd, ok := decl.(*ast.GenDecl)
			if !ok || gd.Tok != token.TYPE {
				continue
			}
			for _, spec := range gd.Specs {
				ts, ok := spec.(*ast.TypeSpec)
				if !ok {
					continue
				}
				it, ok := ts.Type.(*ast.InterfaceType)
				if !ok {
					continue
				}

				ch <- namedInterface{ts.Name, it}
			}
		}
		close(ch)
	}()
	return ch
}

// isVariadic returns whether the function is variadic.
func isVariadic(f *ast.FuncType) bool {
	nargs := len(f.Params.List)
	if nargs == 0 {
		return false
	}
	_, ok := f.Params.List[nargs-1].Type.(*ast.Ellipsis)
	return ok
}

func removeDot(s string) string {
	if len(s) > 0 && s[len(s)-1] == '.' {
		return s[0 : len(s)-1]
	}
	return s
}
