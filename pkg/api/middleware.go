package api

import (
	"github.com/gin-gonic/gin"

	"porttohw/pkg/context"
)

func InjectContext(c *gin.Context) {
	c.Set("ctx", context.Background())
	c.Next()
}
