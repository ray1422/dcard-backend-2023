package controller

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/ray1422/dcard-backend-2023/model"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// headHandler returns the head of list /list/<key>.
// due to the conflicts of the routes, the param is named `id` but it's typed string and should pass the `key` of the list to it.
func headHandler(c *gin.Context) {
	listKey, exists := c.Params.Get("id")
	if !exists {
		c.Status(400)
		return
	}
	head, err := model.FindListByKey(listKey)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.Status(404)
			return
		}
		c.Status(500)
		logger.Default.Error(c, "%s", err)
		return
	}

	c.JSON(200, head.Serialize())
}
