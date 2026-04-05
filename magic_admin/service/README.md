
# 模块功能介绍

## 服务实现层

### 系统配置基础代码 与 业务模块代码 实现具体的功能

#### 需要生成业务代码


```text

代码目录说明
service/system/ 系统管理基础代码存放目录
service/system/enter.go 系统服务入口 // 注入与回滚
service/system/*.go 系统服务实现 // 无需自动生成

service/biz_base/base.go 当前数据库公共方法 // 功能:none
service/biz/enter.go   数据库服务入口文件 // 功能：生成 注入 回滚
service/biz/user.go // 路径命名规则/ 数据库alisa/model.go  // 功能 生成 回滚

```