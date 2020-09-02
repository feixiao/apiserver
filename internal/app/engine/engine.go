package engine

import (
	_ "apiserver/docs"
	"apiserver/internal/app/grpc/greeter"
	"apiserver/internal/app/handler/middleware"
	"apiserver/internal/app/handler/sd"
	"apiserver/internal/app/handler/user"
	"github.com/douyu/jupiter"
	"github.com/douyu/jupiter/pkg/conf"
	"github.com/douyu/jupiter/pkg/server/xgin"
	"github.com/douyu/jupiter/pkg/server/xgrpc"
	"github.com/douyu/jupiter/pkg/util/xgo"
	"github.com/douyu/jupiter/pkg/worker/xcron"
	"github.com/douyu/jupiter/pkg/xlog"
	"github.com/gin-gonic/gin"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
	"google.golang.org/grpc/examples/helloworld/helloworld"
	"net/http"
	"time"
)

type Engine struct {
	jupiter.Application
}

func NewEngine() *Engine {
	eng := &Engine{}
	if err := eng.Startup(
		xgo.ParallelWithError(
			eng.remoteConfigWatch,
			eng.serveGRPC,
			eng.serveHTTP,
			eng.startJobs,
		),
	); err != nil {
		xlog.Panic("startup engine", xlog.Any("err", err))
	}
	return eng
}

type People struct {
	Name string
}

func (eng *Engine) remoteConfigWatch() error {
	p := People{}
	conf.OnChange(func(config *conf.Configuration) {
		err := config.UnmarshalKey("people", &p)
		if err != nil {
			panic(err.Error())
		}
	})
	go func() {
		// 循环打印配置
		for {
			time.Sleep(10 * time.Second)
			xlog.Info("people info", xlog.String("name", p.Name), xlog.String("type", "structByFileWatch"))
		}
	}()
	return nil
}

func (eng *Engine) serveHTTP() error {
	server := xgin.StdConfig("http").WithLogger(xlog.DefaultLogger).Build()

	g := server
	// Middlewares.
	g.Use(gin.Recovery())
	g.Use(middleware.NoCache)
	g.Use(middleware.Options)
	g.Use(middleware.Secure)
	g.Use(middleware.Logging(),
		middleware.RequestId())
	// 404 Handler.
	g.NoRoute(func(c *gin.Context) {
		c.String(http.StatusNotFound, "The incorrect API route.")
	})

	// swagger api docs
	g.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// api for authentication functionalities
	g.POST("/login", user.Login)

	// The user handlers, requiring authentication
	u := g.Group("/v1/user")
	//u.Use(middleware.AuthMiddleware())
	{
		u.POST("", user.Create)
		u.DELETE("/:id", user.Delete)
		u.PUT("/:id", user.Update)
		u.GET("", user.List)
		u.GET("/:username", user.Get)
	}

	// The health check handlers
	svcd := g.Group("/sd")
	{
		svcd.GET("/health", sd.HealthCheck)
		svcd.GET("/disk", sd.DiskCheck)
		svcd.GET("/cpu", sd.CPUCheck)
		svcd.GET("/ram", sd.RAMCheck)
	}
	// server.GET("/jupiter", func(ctx echo.Context) error {
	// 	return ctx.JSON(200, "welcome to jupiter")
	// })
	// // Specify routing group
	// group := server.Group("/api")
	// group.GET("/user/:id", handler.GetUser)

	// //support proxy for http to grpc controller
	// g := greeter.Greeter{}
	// group2 := server.Group("/grpc")
	// group2.GET("/get", xecho.GRPCProxyWrapper(g.SayHello))
	// group2.POST("/post", xecho.GRPCProxyWrapper(g.SayHello))
	return eng.Serve(server)
}

func (eng *Engine) serveGRPC() error {
	server := xgrpc.StdConfig("grpc").Build()
	helloworld.RegisterGreeterServer(server.Server, new(greeter.Greeter))
	return eng.Serve(server)
}

func (eng *Engine) startJobs() error {
	cron := xcron.StdConfig("demo").Build()
	//cron.Schedule(xcron.Every(time.Second*10), xcron.FuncJob(eng.execJob))
	// https://blog.csdn.net/qq_37493556/article/details/105083396
	// spec := "*/5 * * * * ?" //cron表达式，每五秒一次
	//spec := "0 0 1 * * ?"  // 每天凌晨1点执行一次：0 0 1 * * ?
	spec := "0 31 16 * * ?"
	if _, err := cron.AddFunc(spec, eng.execJob); err != nil {
		xlog.Errorf("corn failed, err:%+v", err)
	}
	return eng.Schedule(cron)
}

func (eng *Engine) execJob() error {
	xlog.Info("exec job", xlog.String("info", "print info"))
	xlog.Warn("exec job", xlog.String("warn", "print warning"))
	return nil
}
