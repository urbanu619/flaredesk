package ams_ast

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"os"
	"path/filepath"
)

// 根路由注册器

type RootRouterRegister struct {
	// demo : /Users/***/app_manage_system/go_server/router/enter.go
	TargetFile  string   `json:"targetFile"`  // 目标文件-固定-构建时使用绝对路径 go_server/router/enter.go
	ModuleNames []string `json:"moduleNames"` // 模块名
	StructName  string   `json:"structName"`  // 结构体名称 -- RouterGroup
	ValKey      string   `json:"valKey"`      // 默认值  priRouters
	BaseAst
}

// 服务根路由入口文件地址

func ServerRootRouterEnterFile() string {
	return serverRootRouterEnterFile
}

// 支持同时注册多个模块 moduleNames

func BuildRootRouterRegister(alias []string) *RootRouterRegister {
	return &RootRouterRegister{
		TargetFile:  serverRootRouterEnterFile,
		ModuleNames: alias,
		StructName:  "RouterGroup",
		ValKey:      "priRouters",
	}
}

func (r *RootRouterRegister) routerModulePkg(module string) string {
	return fmt.Sprintf("go_server/router/%s", module)
}

// 总路由注册

func (r *RootRouterRegister) Inspect() error {
	// 读取目标文件
	src, err := os.ReadFile(r.TargetFile)
	if err != nil {
		return err
	}
	// 创建 FileSet 用于跟踪文件位置信息
	fSet := token.NewFileSet()
	// 解析源代码为 AST
	f, err := parser.ParseFile(fSet, "", src, 0)
	if err != nil {
		return err
	}
	ast.Inspect(f, func(n ast.Node) bool {
		if decl, ok := n.(*ast.GenDecl); ok && decl.Tok == token.VAR {
			for _, spec := range decl.Specs {
				if valueSpec, ok := spec.(*ast.ValueSpec); ok {
					for i, name := range valueSpec.Names {
						if name.Name == r.ValKey && len(valueSpec.Values) > i {
							if compLit, ok := valueSpec.Values[i].(*ast.CompositeLit); ok {
								// 元素确定
								if !r.isGlobalContextInterface(compLit) {
									return false
								}
								for _, moduleName := range r.ModuleNames {
									// 添加新的元素 ModelName.RouterGroup{}
									newElement := &ast.CompositeLit{
										Type: &ast.SelectorExpr{
											X:   ast.NewIdent(moduleName),
											Sel: ast.NewIdent(r.StructName),
										},
									}
									// 检查是否已注入
									if !r.IsContainsModule(compLit.Elts, moduleName, r.StructName) {
										compLit.Elts = append(compLit.Elts, newElement)
									}
									// 确保导入了 模块对应的包
									normalizedPkg := filepath.ToSlash(moduleName)
									r.AddImport(f, r.routerModulePkg(normalizedPkg))
								}
							}
						}
					}
				}
			}
		}
		return true
	})
	var out []byte
	bf := bytes.NewBuffer(out)
	err = printer.Fprint(bf, fSet, f)
	if err != nil {
		return err
	}
	return os.WriteFile(r.TargetFile, bf.Bytes(), 0666)
}

// 可做进一步检查检查类型是否是 []global.ContextInterface

func (r *RootRouterRegister) isGlobalContextInterface(compLit *ast.CompositeLit) bool {
	if arrType, ok := compLit.Type.(*ast.ArrayType); ok {
		if selExpr, ok := arrType.Elt.(*ast.SelectorExpr); ok {
			if xIdent, ok := selExpr.X.(*ast.Ident); ok && xIdent.Name == "global" {
				if selExpr.Sel.Name == "ContextInterface" {
					// 确认是我们要找的类型
					return true
				}
			}
		}
	}
	return false
}

// 总路由回滚 -- 仅支持单一模块回滚移除

func (r *RootRouterRegister) RollbackRootRouter() error {
	// 读取目标文件
	src, err := os.ReadFile(r.TargetFile)
	if err != nil {
		return err
	}
	// 创建 FileSet 用于跟踪文件位置信息
	fSet := token.NewFileSet()
	// 解析源代码为 AST
	f, err := parser.ParseFile(fSet, "", src, 0)
	if err != nil {
		return err
	}
	ast.Inspect(f, func(n ast.Node) bool {
		if decl, ok := n.(*ast.GenDecl); ok && decl.Tok == token.VAR {
			for _, spec := range decl.Specs {
				if valueSpec, ok := spec.(*ast.ValueSpec); ok {
					for i, name := range valueSpec.Names {
						if name.Name == r.ValKey && len(valueSpec.Values) > i {
							if compLit, ok := valueSpec.Values[i].(*ast.CompositeLit); ok {
								// 元素确定
								if !r.isGlobalContextInterface(compLit) {
									return false
								}
								newEls := make([]ast.Expr, 0, len(compLit.Elts))
								removeImport := false
								for _, elt := range compLit.Elts {
									if !r.IsTargetModule(elt, r.ModuleNames[0], r.StructName) {
										newEls = append(newEls, elt)
									} else {
										removeImport = true
									}
									if removeImport {
										// 确保删除模块对应的包
										r.RemoveImport(f, r.routerModulePkg(r.ModuleNames[0]))
									}
								}
								compLit.Elts = newEls

							}
						}
					}
				}
			}
		}
		return true
	})
	var out []byte
	bf := bytes.NewBuffer(out)
	err = printer.Fprint(bf, fSet, f)
	if err != nil {
		return err
	}
	err = os.Remove(r.TargetFile)
	if err != nil {
		return err
	}
	return os.WriteFile(r.TargetFile, bf.Bytes(), 0666)
}
