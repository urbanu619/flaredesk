
# 模块功能介绍

### 路由层

#### 从api接口层获取所有暴露的接口并注册进路由

### 路由层代码自动生成需写入信息
```text
1 router/enter.go  priRouters 增加路由组 -- 业务管理仅开放私有路由注册

2 router/model_name/*.go 增加模组路由信息

3 router/model_name/enter.go 增加注册信息

```

### 对应的模块路由 增加表所对应的服务