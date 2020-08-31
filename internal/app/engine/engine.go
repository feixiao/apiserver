package engine

import (
	"apiserver/internal/app/grpc/greeter"
	"apiserver/internal/app/handler/sd"

	"apiserver/internal/app/handler/middleware"
	"apiserver/internal/app/handler/user"
	"net/http"
	"time"

	"github.com/douyu/jupiter"
	"github.com/douyu/jupiter/pkg/server/xgin"
	"github.com/douyu/jupiter/pkg/server/xgrpc"
	"github.com/douyu/jupiter/pkg/util/xgo"
	"github.com/douyu/jupiter/pkg/worker/xcron"
	"github.com/douyu/jupiter/pkg/xlog"
	"github.com/gin-gonic/gin"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
	"google.golang.org/grpc/examples/helloworld/helloworld"
)

type Engine struct {
	jupiter.Application
}

func NewEngine() *Engine {
	eng := &Engine{}
	if err := eng.Startup(
		xgo.ParallelWithError(
			eng.serveGRPC,
			eng.serveHTTP,
			eng.startJobs,
		),
	); err != nil {
		xlog.Panic("startup engine", xlog.Any("err", err))
	}
	return eng
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
	cron.Schedule(xcron.Every(time.Second*10), xcron.FuncJob(eng.execJob))
	return eng.Schedule(cron)
}

func (eng *Engine) execJob() error {
	xlog.Info("exec job", xlog.String("info", "print info"))
	xlog.Warn("exec job", xlog.String("warn", "print warning"))
	return nil
}
