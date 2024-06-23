package router

import (
	"github.com/gin-gonic/gin"
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
}

func (r *Router) GetEngine() *gin.Engine {
	return r.engine
}
