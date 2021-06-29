package api

import (
	"net/http"
	"strconv"

	"github.com/ethereum/go-ethereum/common"
	"github.com/gin-gonic/gin"

	"prottohw/pkg/context"
	"prottohw/pkg/eth"
)

func MountEthRoutes(group *gin.RouterGroup, ethClient eth.Eth) error {
	handler := &ethHandler{
		ethClient: ethClient,
	}
	blocks := group.Group("/blocks")
	blocks.GET("/", handler.getBlocks)
	blocks.GET("/:hash", handler.getBlock)

	transation := group.Group("/transation")
	transation.GET("/:txHash", handler.getTransation)

	return nil
}

type ethHandler struct {
	ethClient eth.Eth
}

func (h *ethHandler) getBlocks(c *gin.Context) {
	val, _ := c.Get("ctx")
	ctx := val.(context.Context)

	limitStr := c.Query("limit")
	n, err := strconv.ParseInt(limitStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, errMsg(err))
		return
	}

	blocks, err := h.ethClient.GetBlocks(ctx, uint64(n))
	if err != nil {
		c.JSON(http.StatusBadRequest, errMsg(err))
		return
	}

	c.JSON(http.StatusOK, gin.H{"blocks": blocks})
}

func (h *ethHandler) getBlock(c *gin.Context) {
	val, _ := c.Get("ctx")
	ctx := val.(context.Context)

	hashHexStr := c.Param("hash")
	hash := common.HexToHash(hashHexStr)


	block, err := h.ethClient.GetBlock(ctx, hash)
	if err == eth.ErrNotFound {
		c.JSON(http.StatusNotFound, errMsg(err))
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, errMsg(err))
		return
	}

	c.JSON(http.StatusOK, block)
}

func (h *ethHandler) getTransation(c *gin.Context) {
	txHash := c.Param("txHash")
	c.String(http.StatusOK, txHash)
}
