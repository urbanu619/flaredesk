package system

import (
	"github.com/gin-gonic/gin"
	"go_server/ams_ast"
	"go_server/base/core"
	"go_server/global"
	"go_server/model/common/response"
	"go_server/service/base"
	"go_server/utils"
	"gorm.io/gorm"
)

type AutoService struct {
	base.SysCommonService
}

var isOpenAllTable bool = true // 是否开启全量同步

// 仅全量同步模型

func (s *AutoService) ModelAuto(c *gin.Context) {
	type AutoCode struct {
		TableName string `json:"tableName" example:"表名"`
		DbAlias   string `json:"dbAlias" example:"业务数据库"`
	}
	var req AutoCode
	var err error
	defer func() {
		if err != nil {
			response.Resp(c, err.Error())
		}
	}()
	err = c.ShouldBindJSON(&req)
	if err != nil {
		return
	}
	err = utils.Verify(req, utils.AutoCodeVerify)
	if err != nil {
		return
	}
	if req.DbAlias == "" {
		response.Resp(c, "请填写别名")
		return
	}
	// 不开放一次生成全库
	if !isOpenAllTable && req.TableName == "" {
		response.Resp(c, "未开放全库生成")
		return
	}
	err = s.ModelAutoForCmd(req.DbAlias, req.TableName)
	if err != nil {
		return
	}
	response.Resp(c)
}

func (s *AutoService) ModelAutoForCmd(alias, tableName string) error {
	db, err := global.BizDBByAlias(alias)
	if err != nil {
		return err
	}
	dbAlias := alias
	dbName := db.Migrator().CurrentDatabase()
	tableNames := []string{}
	if tableName == "" {
		tbs, err := s.GetTables(dbAlias)
		if err != nil {
			return err
		}
		for _, tb := range tbs {
			tableNames = append(tableNames, tb.TableName)
		}
	} else {
		tableNames = append(tableNames, tableName)
	}
	for _, tbn := range tableNames {
		// 模型构建
		core.Log.Infof("模型:%s构建", tbn)
		err = ams_ast.BuildModelRegister(dbAlias, dbName, tbn).Initialize(db)
		if err != nil {
			return err
		}
	}
	return nil
}

// 通过表生成前后端代码

func (s *AutoService) ServerCode(c *gin.Context) {
	type AutoCode struct {
		TableName string `json:"tableName" example:"表名"`
		DbAlias   string `json:"dbAlias" example:"业务数据库"`
	}
	var req AutoCode
	var err error
	defer func() {
		if err != nil {
			response.Resp(c, err.Error())
		}
	}()
	err = c.ShouldBindJSON(&req)
	if err != nil {
		return
	}
	err = utils.Verify(req, utils.AutoCodeVerify)
	if err != nil {
		return
	}
	if req.DbAlias == "" {
		response.Resp(c, "请填写别名")
		return
	}
	// 不开放一次生成全库
	if !isOpenAllTable && req.TableName == "" {
		response.Resp(c, "未开放全库生成")
		return
	}
	db, err := global.BizDBByAlias(req.DbAlias)
	if err != nil {
		return
	}
	dbAlias := req.DbAlias

	dbName := db.Migrator().CurrentDatabase()
	tableNames := []string{}
	if req.TableName == "" {
		tbs, err := s.GetTables(dbAlias)
		if err != nil {
			return
		}
		for _, tb := range tbs {
			tableNames = append(tableNames, tb.TableName)
		}
	} else {
		tableNames = append(tableNames, req.TableName)
	}
	for _, tableName := range tableNames {
		err = s.autoCodeWithTable(db, req.DbAlias, dbName, tableName)
		if err != nil {
			return
		}
	}
	response.Resp(c)
}

func (s *AutoService) AutoServerCodeWithAlias(alias, tableName string) error {
	db, err := global.BizDBByAlias(alias)
	if err != nil {
		return err
	}
	dbAlias := alias

	dbName := db.Migrator().CurrentDatabase()
	tableNames := []string{}
	if tableName == "" {
		tbs, err := s.GetTables(dbAlias)
		if err != nil {
			return err

		}
		for _, tb := range tbs {
			tableNames = append(tableNames, tb.TableName)
		}
	} else {
		tableNames = append(tableNames, tableName)
	}
	for _, tablename := range tableNames {
		err = s.autoCodeWithTable(db, dbAlias, dbName, tablename)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *AutoService) autoCodeWithTable(db *gorm.DB, alias, dbName, tableName string) (err error) {
	// 模型构建
	err = ams_ast.BuildModelRegister(alias, dbName, tableName).Initialize(db)
	if err != nil {
		return
	}
	// 服务构建
	err = ams_ast.BuildServiceRegister(alias, dbName, tableName).Initialize()
	if err != nil {
		return
	}
	// 注入总路由
	err = ams_ast.BuildRootRouterRegister([]string{alias}).Inspect()
	if err != nil {
		return
	}
	// 注入子路由
	err = ams_ast.BuildSubRouterRegister(alias, tableName).Inspect()
	if err != nil {
		return
	}
	return
}

func (s *AutoService) Rollback(c *gin.Context) {
	type AutoCode struct {
		DbAlias   string `json:"dbAlias" example:"业务数据库"`
		TableName string `json:"tableName" example:"表名"`
	}
	var req AutoCode
	var err error
	defer func() {
		if err != nil {
			response.Resp(c, err.Error())
		}
	}()
	err = c.ShouldBindJSON(&req)
	if err != nil {
		return
	}
	err = utils.Verify(req, utils.AutoCodeVerify)
	if err != nil {
		return
	}

	if req.DbAlias == "" {
		response.Resp(c, "请填写别名")
		return
	}
	// 不开放一次生成全库
	if !isOpenAllTable && req.TableName == "" {
		response.Resp(c, "未开放全库回滚")
		return
	}
	db, err := global.BizDBByAlias(req.DbAlias)
	if err != nil {
		return
	}
	dbAlias := req.DbAlias
	dbName := db.Migrator().CurrentDatabase()
	tableNames := []string{}
	if req.TableName == "" {
		tbs, err := s.GetTables(dbAlias)
		if err != nil {
			return
		}
		for _, tb := range tbs {
			tableNames = append(tableNames, tb.TableName)
		}
	} else {
		tableNames = append(tableNames, req.TableName)
	}
	for _, tableName := range tableNames {
		err = s.rollbackTable(dbAlias, dbName, tableName)
		if err != nil {
			return
		}
	}
	response.Resp(c)
}

func (s *AutoService) RollbackWithAlias(dbAlias, tableName string) error {
	db, err := global.BizDBByAlias(dbAlias)
	if err != nil {
		return err
	}
	dbName := db.Migrator().CurrentDatabase()
	tableNames := []string{}
	if tableName == "" {
		tbs, err := s.GetTables(dbAlias)
		if err != nil {
			return err
		}
		for _, tb := range tbs {
			tableNames = append(tableNames, tb.TableName)
		}
	} else {
		tableNames = append(tableNames, tableName)
	}
	for _, tb := range tableNames {
		err = s.rollbackTable(dbAlias, dbName, tb)
		if err != nil {
			return err
		}
	}
	return nil

}

func (s *AutoService) rollbackTable(dbAlias, dbName, tableName string) (err error) {
	// 模型构建
	err = ams_ast.BuildModelRegister(dbAlias, dbName, tableName).RollbackModel()
	if err != nil {
		return
	}
	// 服务构建
	err = ams_ast.BuildServiceRegister(dbAlias, dbName, tableName).RollBack() // todo -- 待验证
	if err != nil {
		return
	}
	// 注入总路由
	err = ams_ast.BuildRootRouterRegister([]string{dbAlias}).RollbackRootRouter()
	if err != nil {
		return
	}
	// 注入子路由
	err = ams_ast.BuildSubRouterRegister(dbAlias, tableName).Rollback() // todo
	if err != nil {
		return
	}
	return
}
