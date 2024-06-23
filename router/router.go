package router

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go-transactions-test/config"
	"go-transactions-test/dicontainer"
)

type IRouter interface {
	InitRoutes(container dicontainer.IDiContainer)
	GetEngine() *gin.Engine
}

type Router struct {
	engine *gin.Engine
}

func SetJSON(c *gin.Context) {
	c.Writer.Header().Set("Content-Type", "application/json")
	c.Next()
}
func NewRouter(ginMode string) IRouter {
	gin.SetMode(ginMode)
	router := &Router{
		engine: gin.Default(),
	}
	router.engine.Use(SetJSON)
	router.engine.Use(gin.Recovery())
	return router
}

func (r *Router) InitRoutes(container dicontainer.IDiContainer) {
	di := container.GetDiContainer()

	r.engine.POST("/accounts", di.TransactionServiceController.CreateAccount)
	r.engine.GET("/accounts/:account_id", di.TransactionServiceController.QueryAccount)
	r.engine.POST("/transactions", di.TransactionServiceController.SubmitTransaction)

	r.engine.GET("/swagger/*any", ginSwagger.CustomWrapHandler(config.NewGinSwaggerConfig(di.Config.SwaggerConfig), swaggerFiles.Handler))
}

func (r *Router) GetEngine() *gin.Engine {
	return r.engine
}
