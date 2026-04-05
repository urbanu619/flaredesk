package ams_ast

import (
	"bytes"
	"fmt"
	"github.com/pkg/errors"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"text/template"
)

// 根据模板文件 生成 service文件

type ServiceRegister struct {
	TemplateFile      string `json:"TemplateFile"`      // 模块文件 resource/package/go_server/service/service.go.tpl
	TemplateEnterFile string `json:"TemplateEnterFile"` // 模块注册文件 resource/package/go_server/service/enter.go.tpl
	EnterFile         string `json:"enterFile"`         // 服务主入口文件 go_server/service/enter.go
	TargetEnterFile   string `json:"targetEnterFile"`   // 业务服务入口文件 go_server/service/alias/enter.go
	TargetFile        string `json:"targetFile"`        // 业务服务文件 go_server/service/alias/model.go
	ServicePkg        string `json:"servicePkg"`        // go_server/service/alias
	DbAlias           string `json:"dbAlias"`           // 数据库别名 aic_gold  模块名称
	DbName            string `json:"dbName"`            // 数据库名 aic_boost 数据库名
	TableName         string `json:"tableName"`         // 模型名称 demo_log 表名称
	ModelCap          string `json:"modelCap"`          // 模型名 首字母大写驼峰规则 DemoLog
	StructName        string `json:"structName"`        // 服务结构体名称 DemoLogService
	BaseAst
}

// alias 数据库别名 tableName 表名

func BuildServiceRegister(alias, dbName, tableName string) *ServiceRegister {
	register := &ServiceRegister{
		TemplateFile:      filepath.Join(template_root, server_service, "service.go.tpl"),
		TemplateEnterFile: filepath.Join(template_root, server_service, "enter.go.tpl"),
		DbAlias:           alias,
		DbName:            dbName,
		TableName:         tableName,
		ModelCap:          CapitalizeOrLower(toCamelCase(tableName)),
		StructName:        fmt.Sprintf("%s%s", CapitalizeOrLower(toCamelCase(tableName)), "Service"),
		BaseAst:           BaseAst{},
	}
	register.EnterFile = register.enterFile()
	register.TargetEnterFile = register.targetEnterFile(alias)
	register.TargetFile = register.targetFile(alias, tableName)
	register.ServicePkg = register.servicePkg(alias)
	return register
}

// 主入口服务名 BizGroup

func (s *ServiceRegister) enterServiceGroup() string {
	return fmt.Sprintf("%s%s", CapitalizeOrLower(toCamelCase(s.DbAlias)), "ServiceGroup")
}

// 服务注册入口文件 go_server/service/enter.go

func (s *ServiceRegister) enterFile() string {
	return filepath.Join(root, server_service, server_enter)
}

// 目标文件 go_server/service/biz/enter.go

func (s *ServiceRegister) targetEnterFile(alias string) string {
	return filepath.Join(root, server_service, alias, server_enter)
}

// 目标文件 go_server/service/biz/user.go

func (s *ServiceRegister) targetFile(alias, tableName string) string {
	return filepath.Join(root, server_service, alias, tableName+".go")
}

// 业务库基础服务 go_server/service/base

func (s *ServiceRegister) baseServicePkg() string {
	return filepath.Join(server_service, server_biz_base_service)
}

// 服务PKG go_server/service/alias

func (s *ServiceRegister) servicePkg(alias string) string {
	return filepath.Join(server_mod, server_service, alias)
}

// 初始化

func (s *ServiceRegister) Initialize() error {
	//core.Log.Infof("服务注入开始:%s-%s", s.DbAlias, s.TableName)
	var err error
	// 创建targetService
	if err = s.createTargetService(); err != nil {
		//core.Log.Infof("服务文件创建失败:%s", err.Error())
		return err
	}
	//core.Log.Infof("1 服务文件创建完成:service/%s/%s.go", s.DbAlias, s.TableName)
	// 创建模块入口文件 注入服务
	if err = s.inspectTargetServiceToGroup(); err != nil {
		//core.Log.Infof("服务注入失败:%s", err.Error())
		return err
	}
	//core.Log.Infof("2 服务注入完成:service/%s/enter.go", s.DbAlias)
	// 注入服务到总入口文件
	if err = s.inspectRootService(); err != nil {
		//core.Log.Infof("注入服务到总入口文件失败:%s", err.Error())
		return err
	}
	//core.Log.Infof("3 注入服务到总入口文件完成:service/enter.go")
	return nil
}

// 创建目标服务文件

func (s *ServiceRegister) createTargetService() error {
	var files *template.Template
	files, err := template.ParseFiles(s.TemplateFile)
	if err != nil {
		return err
	}
	err = os.MkdirAll(filepath.Dir(s.TargetFile), os.ModePerm)
	if err != nil {
		return err
	}
	if FileExists(s.TargetFile) {
		err = os.Remove(s.TargetFile)
		if err != nil {
			return err
		}
	}
	var file *os.File
	file, err = os.Create(s.TargetFile)
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

// 注入子服务到服务组

func (s *ServiceRegister) inspectTargetServiceToGroup() error {
	var err error
	if err = s.checkTargetEnter(); err != nil {
		return err
	}
	// 创建 token.FileSet
	fSet := token.NewFileSet()
	// 解析文件
	file, err := parser.ParseFile(fSet, s.TargetEnterFile, nil, parser.ParseComments)
	if err != nil {
		return err
	}
	// 检查是否存在 ServiceGroup 的结构体声明
	classTypeSpec, isHasServiceGroup := s.FindFileHasStruct(file, fmt.Sprintf("ServiceGroup"))
	if !isHasServiceGroup {
		return errors.New("not found ServiceGroup")
	}
	// 检查并添加属性- 不带Key
	if err = s.checkAndInsertStructFiled(classTypeSpec, s.StructName); err != nil {
		return err
	}
	// 格式化并写回文件
	var buf bytes.Buffer
	if err = format.Node(&buf, fSet, file); err != nil {
		return err
	}
	return os.WriteFile(s.TargetEnterFile, buf.Bytes(), 0666)
}

// 查找结构中属性 如果不存在则增加该属性 属性类型与名称一致

func (s *ServiceRegister) checkAndInsertStructFiled(typeSpec *ast.TypeSpec, typeName string) error {
	hasField := false
	if structType, ok := typeSpec.Type.(*ast.StructType); ok {
		// 检查是否已有 AssetService 字段
		for _, field := range structType.Fields.List {
			if ident, ok := field.Type.(*ast.Ident); ok {
				if ident.Name == typeName {
					hasField = true
					break
				}
			}
			//field.Names
		}
	}
	if !hasField {
		// 添加 AssetService 字段
		structType, ok := typeSpec.Type.(*ast.StructType)
		if !ok {
			return errors.New("ServiceGroup is not a struct type")
		}
		// 创建新的字段
		newField := &ast.Field{
			Names: []*ast.Ident{
				{
					Name: "",
				},
			},
			Type: &ast.Ident{
				Name: typeName,
			},
		}

		// 添加到结构体字段列表
		structType.Fields.List = append(structType.Fields.List, newField)
	}
	return nil
}

// 检查与创建服务组入口文件

func (s *ServiceRegister) checkTargetEnter() error {
	if FileExists(s.TargetEnterFile) {
		return nil
	}
	var files *template.Template
	files, err := template.ParseFiles(s.TemplateEnterFile)
	if err != nil {
		return err
	}
	err = os.MkdirAll(filepath.Dir(s.TargetEnterFile), os.ModePerm)
	if err != nil {
		return err
	}
	var file *os.File
	file, err = os.Create(s.TargetEnterFile)
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

// 注入主服务 EnterFile

func (s *ServiceRegister) inspectRootService() error {
	// 创建 token.FileSet
	fSet := token.NewFileSet()
	// 解析文件
	file, err := parser.ParseFile(fSet, s.EnterFile, nil, parser.ParseComments)
	if err != nil {
		log.Fatal(err)
	}
	isHasServiceGroup := s.CheckFileStructHasFiled(file, "Group", fmt.Sprintf("%sServiceGroup", s.DbAlias))
	if isHasServiceGroup {
		return nil
	}
	// 添加导入声明 -- 需要兼容windows
	normalizedPkg := filepath.ToSlash(s.ServicePkg)

	s.AddImport(file, normalizedPkg)
	// 注入服务
	// 创建新的字段
	newField := &ast.Field{
		Names: []*ast.Ident{
			{
				Name: s.enterServiceGroup(),
			},
		},
		Type: &ast.Ident{
			Name: fmt.Sprintf("%s.ServiceGroup", s.DbAlias),
		},
	}

	classTypeSpec, isHasServiceGroup := s.FindFileHasStruct(file, fmt.Sprintf("Group"))
	if !isHasServiceGroup {
		return errors.New("not found ServiceGroup")
	}

	if err = s.addFieldToStruct(classTypeSpec, s.enterServiceGroup(), newField); err != nil {
		//core.Log.Infof("服务总入口文件:%s 注入:%+v", err.Error(), newField)
		return err
	}
	// 格式化并写回文件
	var buf bytes.Buffer
	if err = format.Node(&buf, fSet, file); err != nil {
		return err

	}
	// 保持原始文件的模式
	info, err := os.Stat(s.EnterFile)
	if err != nil {
		return err
	}
	if err = ioutil.WriteFile(s.EnterFile, buf.Bytes(), info.Mode()); err != nil {
		return err
	}

	return nil
}

func (s *ServiceRegister) addFieldToStruct(typeSpec *ast.TypeSpec, filedName string, newField *ast.Field) error {
	hasFiled := false
	if structType, ok := typeSpec.Type.(*ast.StructType); ok {
		for _, field := range structType.Fields.List {
			for _, name := range field.Names {
				//core.Log.Infof("已注入服务:%s 待注入服务：%s", name.Name, filedName)
				if name.Name == filedName {
					hasFiled = true
				}
			}
		}
	}
	if !hasFiled {
		structType, ok := typeSpec.Type.(*ast.StructType)
		if !ok {
			return errors.New("ServiceGroup is not a struct type")
		}
		// 添加到结构体字段列表
		structType.Fields.List = append(structType.Fields.List, newField)
	}
	return nil
}

// 回滚

func (s *ServiceRegister) RollBack() error {
	// delete targetService
	if FileExists(s.TargetFile) {
		err := os.Remove(s.TargetFile)
		if err != nil {
			return err
		}
	}
	if err := s.rollBackSubEnter(); err != nil {
		return err
	}
	// todo:回滚 总入口文件 -- 不需要回滚 模块一旦创建 不做回滚
	return nil
}

// 回滚子服务注入

func (s *ServiceRegister) rollBackSubEnter() error {
	// 创建 token.FileSet
	fSet := token.NewFileSet()
	// 解析文件
	file, err := parser.ParseFile(fSet, s.TargetEnterFile, nil, parser.ParseComments)
	if err != nil {
		return err
	}
	// 模块入口文件 回滚注入服务
	// 检查是否存在 ServiceGroup 的结构体声明
	//core.Log.Infof("检查文件:%s中是否存在 ServiceGroup 的结构体声明", s.TargetEnterFile)

	classTypeSpec, isHasServiceGroup := s.FindFileHasStruct(file, fmt.Sprintf("ServiceGroup"))
	if !isHasServiceGroup {
		return errors.New("not found ServiceGroup")
	}
	//core.Log.Infof("检查文件:%s中 存在 ServiceGroup", s.TargetEnterFile)

	if err := s.checkAndRemoveStructFiled(classTypeSpec, s.StructName); err != nil {
		return err
	}
	// 格式化并写回文件
	var buf bytes.Buffer
	if err = format.Node(&buf, fSet, file); err != nil {
		return err

	}
	// 保持原始文件的模式
	info, err := os.Stat(s.TargetEnterFile)
	if err != nil {
		return err
	}
	if err = ioutil.WriteFile(s.TargetEnterFile, buf.Bytes(), info.Mode()); err != nil {
		return err
	}

	// 回滚 总入口文件

	return nil
}

// 查找结构中属性 如果存在则删除该属性

func (s *ServiceRegister) checkAndRemoveStructFiled(typeSpec *ast.TypeSpec, typeName string) error {
	if structType, ok := typeSpec.Type.(*ast.StructType); ok {
		// 检查是否已有
		//core.Log.Infof("检查文件 ServiceGroup 中是否存在Field: %s", typeName)
		filteredFields := make([]*ast.Field, 0)
		for _, field := range structType.Fields.List {
			if ident, ok := field.Type.(*ast.Ident); ok {
				//core.Log.Infof("ident.Name:%s typeName:%s", ident.Name, typeName)
				if ident.Name == typeName {
					continue
				}
				// Skip this field (effectively removing it)
			}
			filteredFields = append(filteredFields, field)
		}
		//core.Log.Infof("len(filteredFields):%d len(structType.Fields.List):%d", len(filteredFields), len(structType.Fields.List))

		// If the length changed, we removed at least one field
		if len(filteredFields) < len(structType.Fields.List) {
			structType.Fields.List = filteredFields
			return nil
		}
	}
	return nil
}
