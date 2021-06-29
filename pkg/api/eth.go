package api

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"prottohw/pkg/eth"
)

func MountEthRoutes(group *gin.RouterGroup, ethClient eth.Eth) error {
	handler := &ethHandler{
		ethClient: ethClient,
	}
	blocks := group.Group("/blocks")
	blocks.GET("/", handler.getBlocks)
	blocks.GET("/:id", handler.getBlocks)

	transation := group.Group("/transation")
	transation.GET("/:txHash", handler.getTransation)

	return nil
}

type ethHandler struct {
	ethClient eth.Eth
}

func (h *ethHandler) getBlocks(c *gin.Context) {
	n := c.Query("limit")
	c.String(http.StatusOK, n)
}

func (h *ethHandler) getBlock(c *gin.Context) {
	id := c.Param("id")
	c.String(http.StatusOK, id)
}

func (h *ethHandler) getTransation(c *gin.Context) {
	txHash := c.Param("txHash")
	c.String(http.StatusOK, txHash)
}
