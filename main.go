package main

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"prottohw/pkg/api"
	"prottohw/pkg/eth"
	"prottohw/pkg/log"
)

func main() {
	ethclient := eth.New("url")
	router := gin.Default()
	router.Use(api.InjectContext)
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, "healthy")
	})

	root := router.Group("/")
	if err := api.MountEthRoutes(root, ethclient); err != nil {
		log.Global().Error("api.MountEthRoutes failed")
		panic(err)
	}

	if err := router.Run(); err != nil {
		log.Global().Error("router.Run failed")
	}
}
