# apiserver
demo for jupiter

### 安装脚手架
```
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


### 测试脚本
```shell
wrktest.sh
```

### Todo
+ [ ] Repository接口实现
+ [ ] GRPC接口实现
+ [ ] 远程配置
+ [ ] 服务注册
