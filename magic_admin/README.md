# aicgold-admin-server

## 管理后台基础框架
```
本服务实现快速生成业务端数据库模型、路由、服务
1 生成查询+导出接口
2 可根据业务端模型变化而发生变化
3 对业务端数据只读 减少业务交互
```
## 项目目录功能说明
```
├── PLAN.md  -- 规划文件
├── README.md -- 使用说明
├── Taskfile.yaml -- task 命令配置
├── ams_ast  -- 代码自动生成模型模块
├── base -- 基础模块
│    ├── config -- 各服务配置文件
│    ├── app.go -- 项目配置
│       ├── core -- 基础核心实现
│       └── engine -- 基础web引擎
│           └── proxy.go  -- 代理请求实现
├── cmds  -- 命令行管理模块
├── conf.d -- 自动生成 配置文件
├── global -- DB全局服务
├── go.mod
├── go.sum
├── logs    -- 自动生成 日志存放文件夹
├── main.go -- 程序主入口
├── model   -- 模型层
├── router  -- 路由层
├── service -- 服务层
├── sign_test.go -- 测试文件
├── static -- 自动生成 静态文件存放路径
└── utils -- 公共工具类函数
```
#### task命令行工具安装
```shell
go install github.com/go-task/task/v3/cmd/task@latest
# 或者安装特定版本
go install github.com/go-task/task/v3/cmd/task@v3.30.1
task --version
task -l  # 列出所有可用的 task
```
### 使用说明
```
0 同步需要生成管理服务的业务端数据
1 初始化配置 使用命令 go run main.go config / task config（如果有安装task命令） 可使用命令进行安装
2 配置redis mysql 等配置文件
3 go run main.go migrate / task migrate 同步数据库
4 task run 启动服务
5.生成代码：go run main.go biz generate app all (包括模型 服务 路由)
  指定表名生成服务 go run main.go biz generate app magic_extra_stake_order
6.同步所有表结构：go run main.go biz model app all
7.同步指定表结构：go run main.go biz model app sys_config


```
