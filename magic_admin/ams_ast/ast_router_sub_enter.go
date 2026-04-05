package ams_ast

import (
	"bytes"
	"fmt"
	"github.com/pkg/errors"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"os"
	"path/filepath"
	"text/template"
)

// 模版自路由入口文件路径

func TemplateSubRouterEnterFile() string {
	return filepath.Join(template_root, "router/module/enter.go.tpl")
}

// 模版自路由文件路径

func TemplateSubRouterRegisterFile() string {
	return filepath.Join(template_root, "router/module/model.go.tpl")
}

// 通过模块名获取子路由入口文件

func ServerSubRouterEnterFile(module string) string {
	return filepath.Join(serverRootRouterPath, module, server_enter)
}

// 通过模块+模型名获取 子路由注册文件

func ServerSubRouterRegisterFile(module, model string) string {
	return filepath.Join(serverRootRouterPath, module, snakeString(model)+".go")
}

// 子路由注册器

type SubRouterRegister struct {
	TemplateEnterFile     string // `json:"templateEnterFile"` // 子路由入口文件
	TemplateModelFile     string // `json:"templateModelFile"` // 子路由模型文件
	TargetModuleEnterFile string // `json:"targetModuleFile"`  //  模块路由文件
	TargetModelFile       string // `json:"targetModelFile"`   //  模型路由文件
	Server                string // server mod 名 go_server
	Module                string // `json:"module" example:"保持与alias一致"`
	ModuleUpper           string // `json:"moduleUpper" example:"模块名大写首字母"`
	Model                 string // `json:"model" example:"模型名"` xx_yy
	ModelUpper            string // `json:"modelUpper" example:"模型名大写首字母"`
	RouterServiceName     string // xxYyXxYyService
	ValKey                string `json:"valKey"` // 默认值  priRouters
	BaseAst
}

func BuildSubRouterRegister(alias, model string) *SubRouterRegister {
	// aaBbAabb
	routerServiceName := fmt.Sprintf("%s%s", CapitalizeOrLower(toCamelCase(alias), false), CapitalizeOrLower(toCamelCase(model)))
	return &SubRouterRegister{
		TemplateEnterFile:     TemplateSubRouterEnterFile(),              // 入口模版文件
		TemplateModelFile:     TemplateSubRouterRegisterFile(),           // 模型模版文件
		TargetModuleEnterFile: ServerSubRouterEnterFile(alias),           // 目标模板文件
		TargetModelFile:       ServerSubRouterRegisterFile(alias, model), // 目标模型文件
		Server:                server_mod,                                // app服务模块
		Module:                alias,                                     // 保持别名不变                                                                      // xx_yy
		ModuleUpper:           CapitalizeOrLower(toCamelCase(alias)),     // 改为首字母大写 驼峰结构  AaBb                                                               // XxYy
		Model:                 snakeString(model),                        // xx_yy
		ModelUpper:            CapitalizeOrLower(toCamelCase(model)),     // XxYy
		RouterServiceName:     routerServiceName,                         // 路由服务名称 aaBbAabb
		ValKey:                "allRouters",
		BaseAst:               BaseAst{},
	}
}

func (s *SubRouterRegister) Inspect() error {
	var err error
	// 模版文件检查
	if !FileExists(s.TemplateEnterFile) {
		return errors.New("sub router enter template file is not exists")
	}
	if !FileExists(s.TemplateModelFile) {
		return errors.New("sub router model template file is not exists")
	}
	// 入口文件 -- 如果不存在则创建
	if !FileExists(s.TargetModuleEnterFile) {
		err = s.createEnterRouter()
		if err != nil {
			return err
		}
	}
	// 模型路由文件(仅支持创建与删除 不支持回滚)
	if !FileExists(s.TargetModelFile) {
		err = s.createModel()
		if err != nil {
			return err
		}
	}
	// 注入路由代码
	err = s.inspectRouter()
	if err != nil {
		return err
	}
	return nil
}

func (s *SubRouterRegister) Rollback() error {
	// 删除模型文件
	// 模型路由文件(仅支持创建与删除 不支持回滚)
	if FileExists(s.TargetModelFile) {
		err := os.Remove(s.TargetModelFile)
		if err != nil {
			return err
		}
	}
	// 子路由回滚
	err := s.rollbackRouter()
	if err != nil {
		return err
	}
	return nil
}

// 创建子路由入口文件

func (s *SubRouterRegister) createEnterRouter() error {

	var files *template.Template
	files, err := template.ParseFiles(s.TemplateEnterFile)
	if err != nil {
		return err
	}

	err = os.MkdirAll(filepath.Dir(s.TargetModuleEnterFile), os.ModePerm)
	if err != nil {
		return err
	}

	var file *os.File
	file, err = os.Create(s.TargetModuleEnterFile)
	if err != nil {
		return err
	}
	err = files.Execute(file, s)
	_ = file.Close()
	if err != nil {
		return err
	}
	return nil
}

// 创建模型文件

func (s *SubRouterRegister) createModel() error {
	var files *template.Template
	files, err := template.ParseFiles(s.TemplateModelFile)
	if err != nil {
		return err
	}
	//core.Log.Infof("子路由模板文件:%s", s.TemplateModelFile)

	err = os.MkdirAll(filepath.Dir(s.TargetModelFile), os.ModePerm)
	if err != nil {
		return err
	}
	//core.Log.Infof("子路由生成目标文件:%s", s.TargetModelFile)

	var file *os.File
	file, err = os.Create(s.TargetModelFile)
	if err != nil {
		return err
	}
	err = files.Execute(file, s)
	_ = file.Close()
	if err != nil {
		return err
	}
	return nil
}

// todo:子路由注入

func (s *SubRouterRegister) inspectRouter() error {
	//core.Log.Infof("开始注入子路由")
	//core.Log.Infof("子路由Model文件:%s", s.TargetModelFile)
	//core.Log.Infof("子路由入口文件:%s", s.TargetModuleEnterFile)
	// 创建 FileSet 用于跟踪文件位置信息
	fSet := token.NewFileSet()
	// 解析源代码为 AST
	node, err := parser.ParseFile(fSet, s.TargetModuleEnterFile, nil, parser.ParseComments)
	if err != nil {
		return err
	}
	// 检查是否已注入
	inspectIsExists := s.IsContainsSubRouter(node, "allRouters", fmt.Sprintf("%sRouter", s.ModelUpper))

	if inspectIsExists {
		//core.Log.Infof("子路由%s已注入:allRouters", s.ModelUpper)
		return nil
	}
	s.inspectRouterToVars(node, "allRouters", fmt.Sprintf("%sRouter", s.ModelUpper))
	var out []byte
	bf := bytes.NewBuffer(out)
	err = printer.Fprint(bf, fSet, node)
	if err != nil {
		return err
	}
	err = os.WriteFile(s.TargetModuleEnterFile, bf.Bytes(), 0666)
	if err != nil {
		return err
	}
	//core.Log.Infof("成功添加 %s{} 到 allRouters", s.ModelUpper)
	return nil
}

// todo: 子路由回滚

func (s *SubRouterRegister) rollbackRouter() error {
	// 创建 FileSet 用于跟踪文件位置信息
	fSet := token.NewFileSet()
	// 解析源代码为 AST
	node, err := parser.ParseFile(fSet, s.TargetModuleEnterFile, nil, parser.ParseComments)
	if err != nil {
		return err
	}
	// 检查是否已注入
	inspectIsExists := s.IsContainsSubRouter(node, "allRouters", fmt.Sprintf("%sRouter", s.ModelUpper))
	if !inspectIsExists {
		//core.Log.Infof("子路由%s已回滚:allRouters", s.ModelUpper)
		return nil
	}
	s.rollbackRouterFromVars(node, "allRouters", fmt.Sprintf("%sRouter", s.ModelUpper))
	var out []byte
	bf := bytes.NewBuffer(out)
	err = printer.Fprint(bf, fSet, node)
	if err != nil {
		return err
	}
	err = os.WriteFile(s.TargetModuleEnterFile, bf.Bytes(), 0666)
	if err != nil {
		return err
	}
	//core.Log.Infof("成功添从 allRouters 移除 %s{} ", fmt.Sprintf("%sRouter", s.ModelUpper))
	return nil
}

// 检查元素列表中是否已注入指定结果体
// 检查变了列表 allRouters 中是否存在 XxYyRouterGroup{}
// routerVal:allRouters

func (s *SubRouterRegister) IsContainsSubRouter(node ast.Node, routerVal, subRouterName string) bool {
	// 2. 检查 UserRouter{} 是否存在
	routerExists := false
	ast.Inspect(node, func(n ast.Node) bool {
		// 查找 allRouters 变量声明
		if valueSpec, ok := n.(*ast.ValueSpec); ok && len(valueSpec.Names) > 0 && valueSpec.Names[0].Name == routerVal {
			if compLit, ok := valueSpec.Values[0].(*ast.CompositeLit); ok {
				for _, elt := range compLit.Elts {
					switch expr := elt.(type) {
					case *ast.CompositeLit:
						// 处理 AssetRouter{} 情况
						if selExpr, ok := expr.Type.(*ast.SelectorExpr); ok {
							if ident, ok := selExpr.X.(*ast.Ident); ok && ident.Name == subRouterName {
								routerExists = true
								return false
							}
						} else if ident, ok := expr.Type.(*ast.Ident); ok && ident.Name == subRouterName {
							routerExists = true
							return false
						}
					case *ast.SelectorExpr:
						// 处理 AssetRouter 情况 (没有花括号)
						if ident, ok := expr.X.(*ast.Ident); ok && ident.Name == subRouterName {
							routerExists = true
							return false
						}
					case *ast.Ident:
						// 处理 AssetRouter 情况 (单独标识符)
						if expr.Name == subRouterName {
							routerExists = true
							return false
						}
					}
				}
			}
		}
		return true
	})

	return routerExists
}

// 注入元素到列表

func (s *SubRouterRegister) inspectRouterToVars(node ast.Node, routerVal, subRouterName string) {
	ast.Inspect(node, func(n ast.Node) bool {
		if valueSpec, ok := n.(*ast.ValueSpec); ok && len(valueSpec.Names) > 0 && valueSpec.Names[0].Name == routerVal {
			if compLit, ok := valueSpec.Values[0].(*ast.CompositeLit); ok {
				// 创建新的 AssetRouter{} 元素
				newElement := &ast.CompositeLit{
					Type: &ast.Ident{Name: subRouterName},
				}
				// 添加到切片中
				compLit.Elts = append(compLit.Elts, newElement)
				return false
			}
		}
		return true
	})
}

// 从列表移除元素

func (s *SubRouterRegister) rollbackRouterFromVars(node ast.Node, routerVal, subRouterName string) {
	ast.Inspect(node, func(n ast.Node) bool {
		if valueSpec, ok := n.(*ast.ValueSpec); ok && len(valueSpec.Names) > 0 && valueSpec.Names[0].Name == routerVal {
			if compLit, ok := valueSpec.Values[0].(*ast.CompositeLit); ok {
				// Filter out elements with the specified subRouterName
				filteredEls := make([]ast.Expr, 0, len(compLit.Elts))
				for _, elt := range compLit.Elts {
					if cl, ok := elt.(*ast.CompositeLit); ok {
						if ident, ok := cl.Type.(*ast.Ident); ok {
							//core.Log.Infof("ident.name:%s", ident.Name)
							//core.Log.Infof("subRouterName:%s", subRouterName)
							if ident.Name == subRouterName {
								continue
							}
							// Skip this element (remove it)
						}
					}
					filteredEls = append(filteredEls, elt)
				}
				// Update the elements if any were removed
				if len(filteredEls) < len(compLit.Elts) {
					compLit.Elts = filteredEls
				}
				return false
			}
		}
		return true
	})

}
