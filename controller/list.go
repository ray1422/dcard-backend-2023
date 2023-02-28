package controller

import (
	"errors"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/ray1422/dcard-backend-2023/model"

	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// listHandler handles /list/<id>/<version>?next=[cursor]
func listHandler(c *gin.Context) {
	listIDStr := c.Param("id")
	listID, err := strconv.Atoi(listIDStr)
	if err != nil {
		c.Status(400)
		return
	}
	versionStr := c.Param("version")
	version, err := strconv.Atoi(versionStr)
	if err != nil {
		c.Status(400)
		return
	}
	cursorStr := c.Query("next")
	listNodes, cursor, err := model.GetListNodes(uint(listID), uint(version), cursorStr)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.Status(404)
			return
		}
		c.Status(500)
		logger.Default.Error(c, "%s", err)

	}
	articles := model.SerializeCursorPagedItems[model.ArticleSerializer](listNodes, cursor)
	c.JSON(200, articles)
}
