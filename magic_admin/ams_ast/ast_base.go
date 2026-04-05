package ams_ast

import (
	"fmt"
	"go/ast"
	"go/token"
	"path/filepath"
	"strings"
)

type BaseAst struct {
}

// 添加导入声明

func (BaseAst) AddImport(f *ast.File, importPath string) {
	// 规范化导入路径，确保跨平台兼容性
	normalizedPath := filepath.ToSlash(strings.TrimSpace(importPath))

	for _, imp := range f.Imports {
		if strings.Trim(imp.Path.Value, `"`) == normalizedPath {
			return // 已经存在该导入
		}
	}

	// 创建新的导入声明
	newImport := &ast.ImportSpec{
		Path: &ast.BasicLit{
			Kind:  token.STRING,
			Value: fmt.Sprintf(`"%s"`, normalizedPath),
		},
	}

	// 添加到导入声明列表
	f.Imports = append(f.Imports, newImport)

	// 找到第一个 GenDecl(导入声明块)
	for _, decl := range f.Decls {
		if genDecl, ok := decl.(*ast.GenDecl); ok && genDecl.Tok == token.IMPORT {
			// 添加到现有的导入声明块
			genDecl.Specs = append(genDecl.Specs, newImport)
			return
		}
	}

	// 如果没有找到导入声明块，创建一个新的
	importDecl := &ast.GenDecl{
		Tok: token.IMPORT,
		Specs: []ast.Spec{
			newImport,
		},
	}

	// 将新的导入声明添加到文件开头
	newDecls := []ast.Decl{importDecl}
	f.Decls = append(newDecls, f.Decls...)
}

// 移除指定的导入

func (BaseAst) RemoveImport(f *ast.File, importPath string) {
	// 1. 从 f.Imports 中移除
	newImports := make([]*ast.ImportSpec, 0, len(f.Imports))
	for _, imp := range f.Imports {
		if strings.Trim(imp.Path.Value, `"`) != importPath {
			newImports = append(newImports, imp)
		}
	}
	f.Imports = newImports

	// 2. 从导入声明中移除
	for _, decl := range f.Decls {
		if genDecl, ok := decl.(*ast.GenDecl); ok && genDecl.Tok == token.IMPORT {
			newSpecs := make([]ast.Spec, 0, len(genDecl.Specs))
			for _, spec := range genDecl.Specs {
				if impSpec, ok := spec.(*ast.ImportSpec); ok {
					if strings.Trim(impSpec.Path.Value, `"`) != importPath {
						newSpecs = append(newSpecs, spec)
					}
				} else {
					newSpecs = append(newSpecs, spec)
				}
			}
			genDecl.Specs = newSpecs

			// 如果导入声明为空，可以移除整个声明
			if len(genDecl.Specs) == 0 {
				newDecls := make([]ast.Decl, 0, len(f.Decls))
				for _, d := range f.Decls {
					if d != decl {
						newDecls = append(newDecls, d)
					}
				}
				f.Decls = newDecls
			}
		}
	}
}

// 检查元素列表中是否已注入指定结果体 demo.RouterGroup{}

func (BaseAst) IsContainsModule(elts []ast.Expr, moduleName, structName string) bool {
	for _, elt := range elts {
		if compLit, ok := elt.(*ast.CompositeLit); ok {
			if selExpr, ok := compLit.Type.(*ast.SelectorExpr); ok {
				if xIdent, ok := selExpr.X.(*ast.Ident); ok && xIdent.Name == moduleName {
					if selExpr.Sel.Name == structName {
						return true
					}
				}
			}
		}
	}
	return false
}

// 检查是否为指定模块

func (BaseAst) IsTargetModule(expr ast.Expr, moduleName, structName string) bool {
	if compLit, ok := expr.(*ast.CompositeLit); ok {
		if selExpr, ok := compLit.Type.(*ast.SelectorExpr); ok {
			if xIdent, ok := selExpr.X.(*ast.Ident); ok && xIdent.Name == moduleName {
				if selExpr.Sel.Name == structName {
					return true
				}
			}
		}
	}
	return false
}

// 检查文件中指定结构体是否存在某属性

func (BaseAst) CheckFileStructHasFiled(file *ast.File, structName, filedName string) bool {
	// 检查是否已存在 UserServiceGroup
	hasService := false
	ast.Inspect(file, func(n ast.Node) bool {
		if compLit, ok := n.(*ast.CompositeLit); ok {
			if selExpr, ok := compLit.Type.(*ast.SelectorExpr); ok {
				if selExpr.Sel.Name == structName {
					for _, elt := range compLit.Elts {
						if kv, ok := elt.(*ast.KeyValueExpr); ok {
							if ident, ok := kv.Key.(*ast.Ident); ok && ident.Name == filedName {
								hasService = true
								return false
							}
						}
					}
				}
			}
		}
		return true
	})
	return hasService
}

// 检查文件中是否存在 struct 声明(类声明)

func (BaseAst) FindFileHasStruct(file *ast.File, className string) (*ast.TypeSpec, bool) {
	var classTypeSpec *ast.TypeSpec
	ast.Inspect(file, func(n ast.Node) bool {
		// 查找类型声明
		if typeSpec, ok := n.(*ast.TypeSpec); ok {
			//core.Log.Infof("查找文件中是否存在类型申明 typeSpec.Name.Name:%s", typeSpec.Name.Name)
			if typeSpec.Name.Name == className {
				classTypeSpec = typeSpec
				// 检查结构体类型
				return false
			}
		}
		return true
	})
	if classTypeSpec == nil {
		return nil, false
	}
	return classTypeSpec, true
}
