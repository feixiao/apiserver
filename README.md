# apiserver
demo for jupiter

### 安装脚手架
```
export GO111MODULE=on && go get -u github.com/swaggo/swag/cmd/swag
go get -u github.com/douyu/jupiter/tools/jupiter
```

#### 新建项目
```
jupiter new apiserver
```
#### 项目目录结构:
```go
build                           编译目录
cmd                             应用启动目录
config                          应用配置目录
internal
├─app                           应用目录
│  ├─engine                     
│  │  ├─engine.go               核心编排引擎(启动HTTP,GRPC,JOB等服务)
│  ├─grpc                       grpc服务实现目录
│  ├─handler                    控制器目录（接收用户请求）              
│  │  ├─user.go                 控制器文件
│  ├─model                      model目录（定义持久层结构体）
│  │  ├─db
│  │  │  ├─user.go
│  │  ├─init.go                 初始化全局数据库句柄
│  ├─service                    service层
│  │  ├─user                    模块
│  │  │  ├─impl  
│  │  │  │  ├─mysqlImpl.go      实现
│  │  │  ├─repository.go        service 接口
│  │  ├─init.go
pb                              proto文件
sql                             sql脚本
.gitignore
go.mod
Makefile
```
### 参考代码和文章
+ [apiserver_demos](https://github.com/feixiao/apiserver_demos)  基于Gin版本APIServer的代码实现。
+ [《Go API 开发实战》](https://cloud.tencent.com/developer/article/1427578) 基于Gin版本APIServer, 含有详细实现过程。

### 接口文档
```shell script
make run  
# 访问 http://localhost:8080/swagger/index.html
```

### 测试脚本
```shell
wrktest.sh
```

### 远程配置
+ [Juno部署](https://feixiao.github.io/2020/08/31/juno/)
   
### 关于接入Juno
+ 安装Juno
+ 启动apiserver，注意这边需要配置环境变量APP_NAME为apiserver，否则程序名字为main
+ 在页面上添加应用，应用名字很关键，这边是apiserver。
+ 关联apiserver和juno-agent
    + 这个时候juno-agent会去etcd获取apiserver的信息，然后写入promethus的配置路径(基于文件的动态获取抓取目标)
    + 这种配置也就说明juno-agent和prometheus要一对一配置
+ 顺利的情况下面我们可以在juno-admin上看到grafana的监控项了

### Todo
+ [x] 接入Juno
+ [ ] Repository接口实现
+ [ ] GRPC接口实现
+ [ ] 远程配置
+ [x] 服务注册
