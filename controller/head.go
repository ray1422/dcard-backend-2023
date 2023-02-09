package controller

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/ray1422/dcard-backend-2023/model"
	"github.com/ray1422/dcard-backend-2023/utils/db"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// headHandler returns the head of list /list/<id>
func headHandler(c *gin.Context) {
	listKey, exists := c.Params.Get("id")
	if !exists {
		c.Status(400)
		return
	}
	head, err := head(listKey)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.Status(404)
			return
		}
		c.Status(500)
		logger.Default.Error(c, "%s", err)
		return
	}

	c.JSON(200, head)
}

func head(key string) (*model.List, error) {
	list := model.List{}
	err := db.GormDB().Take(&list, "key = ?", key).Error
	if err != nil {
		return nil, err
	}

	return &list, nil

}
